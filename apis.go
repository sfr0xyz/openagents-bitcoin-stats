package main

import (
	"encoding/json"
	"errors"
	"math"
	"strconv"
	"time"

	"github.com/extism/go-pdk"
)

const mempoolspaceAPI string = "https://mempool.space/api"
const bitnodesAPI string = "https://bitnodes.io/api/v1"
const coincapAPI string = "https://api.coincap.io/v2"

const satsPerBTC = 100000000
const bytesPerMegabyte = 1000000

type MarketData struct {
	Supply                float64 `json:"supply"`
	SupplyPercent         float64 `json:"supplyPercent"`
	Price                 float64 `json:"price"`
	PriceChangePercent24h float64 `json:"priceChange24hPercent"`
	MoscowTime            int64   `json:"moscowTime"`
	MarketCap             float64 `json:"marketCap"`
}

type CoincapResponse struct {
	Data struct {
		Supply            string `json:"supply"`
		MaxSupply         string `json:"maxSupply"`
		MarketCapUsd      string `json:"marketCapUsd"`
		PriceUsd          string `json:"priceUsd"`
		ChangePercent24Hr string `json:"changePercent24Hr"`
	} `json:"data"`
	Timestamp int64 `json:"timestamp"`
}

func GetMarketData() (interface{}, error) {
	const endpoint string = "/assets/bitcoin"

	req := pdk.NewHTTPRequest(pdk.MethodGet, coincapAPI+endpoint)
	res := req.Send()

	if res.Status() != 200 {
		errMsg := "CoinCap.IO API request returned status: " + strconv.Itoa(int(res.Status()))
		return MarketData{}, errors.New(errMsg)
	}

	data := CoincapResponse{}
	err := json.Unmarshal(res.Body(), &data)
	if err != nil {
		return MarketData{}, err
	}

	sup, err := strconv.ParseFloat(data.Data.Supply, 64)
	if err != nil {
		return MarketData{}, err
	}
	maxSup, err := strconv.ParseFloat(data.Data.MaxSupply, 64)
	if err != nil {
		return MarketData{}, err
	}
	priceUSD, err := strconv.ParseFloat(data.Data.PriceUsd, 64)
	if err != nil {
		return MarketData{}, err
	}
	changePercent24h, err := strconv.ParseFloat(data.Data.ChangePercent24Hr, 64)
	if err != nil {
		return MarketData{}, err
	}
	marketCapUSD, err := strconv.ParseFloat(data.Data.MarketCapUsd, 64)
	if err != nil {
		return MarketData{}, err
	}

	var marketData MarketData
	marketData.Supply = float64(sup)
	marketData.SupplyPercent = roundDec((sup/maxSup)*100, 2)
	marketData.Price = roundDec(float64(priceUSD), 2)
	marketData.PriceChangePercent24h = roundDec(changePercent24h, 2)
	marketData.MoscowTime = int64(float64(satsPerBTC) / priceUSD)
	marketData.MarketCap = roundDec(float64(marketCapUSD), 2)

	return marketData, nil
}

type DiffAdjData struct {
	DifficultyChangePercent float64
	RemainingBlocks         int
	EstimatedRetargetDate   string
}

type DiffAdjResponse struct {
	DifficultyChange      float64 `json:"difficultyChange"`
	EstimatedRetargetDate int64   `json:"estimatedRetargetDate"`
	RemainingBlocks       int     `json:"remainingBlocks"`
}

func getDifficultyAdjustmentData() (DiffAdjData, error) {
	const endpoint string = "/v1/difficulty-adjustment"

	req := pdk.NewHTTPRequest(pdk.MethodGet, mempoolspaceAPI+endpoint)
	res := req.Send()

	if res.Status() != 200 {
		errMsg := "Mempool.Space API request returned status: " + strconv.Itoa(int(res.Status()))
		return DiffAdjData{}, errors.New(errMsg)
	}

	data := DiffAdjResponse{}
	err := json.Unmarshal(res.Body(), &data)
	if err != nil {
		return DiffAdjData{}, err
	}

	var diffAdj DiffAdjData
	diffAdj.DifficultyChangePercent = roundDec(data.DifficultyChange, 2)
	diffAdj.RemainingBlocks = data.RemainingBlocks
	diffAdj.EstimatedRetargetDate = time.UnixMilli(data.EstimatedRetargetDate).Format(time.RFC1123)

	return diffAdj, nil
}

type HRAndDiff struct {
	Hashrate   float64
	Difficulty float64
}

type HRResponse struct {
	CurrentHashrate   float64 `json:"currentHashrate"`
	CurrentDifficulty float64 `json:"currentDifficulty"`
}

func getHashrateAndDifficulty() (HRAndDiff, error) {
	const endpoint string = "/v1/mining/hashrate/1s"

	req := pdk.NewHTTPRequest(pdk.MethodGet, mempoolspaceAPI+endpoint)
	res := req.Send()

	if res.Status() != 200 {
		errMsg := "Mempool.Space API request returned status: " + strconv.Itoa(int(res.Status()))
		return HRAndDiff{}, errors.New(errMsg)
	}

	var data HRResponse
	err := json.Unmarshal(res.Body(), &data)
	if err != nil {
		return HRAndDiff{}, err
	}

	var hashrateAndDifficulty HRAndDiff
	hashrateAndDifficulty.Hashrate = data.CurrentHashrate
	hashrateAndDifficulty.Difficulty = data.CurrentDifficulty

	return hashrateAndDifficulty, nil
}

type MiningData struct {
	Hashrate                 float64 `json:"hashrate"`
	Difficulty               float64 `json:"difficulty"`
	RetargetDifficultyChange float64 `json:"retargetDifficultyChangePercent"`
	RetargetRemainingBlocks  int     `json:"retargetRemainingBlocks"`
	RetargetEstimatedDate    string  `json:"retargetEstimatedDate"`
}

func GetMiningData() (interface{}, error) {
	var err error

	var da DiffAdjData
	da, err = getDifficultyAdjustmentData()
	if err != nil {
		return MiningData{}, err
	}

	var hd HRAndDiff
	hd, err = getHashrateAndDifficulty()
	if err != nil {
		return MiningData{}, err
	}

	var mining MiningData
	mining.Hashrate = hd.Hashrate
	mining.Difficulty = hd.Difficulty
	mining.RetargetDifficultyChange = da.DifficultyChangePercent
	mining.RetargetRemainingBlocks = da.RemainingBlocks
	mining.RetargetEstimatedDate = da.EstimatedRetargetDate

	return mining, nil
}

type MempoolStats struct {
	UnconfirmedTXs int64   `json:"unconfirmedTXs"`
	VSize          float64 `json:"vSize"`
	PendingFees    float64 `json:"pendingFees"`
	BlocksToClear  int     `json:"blocksToClear"`
}

type MempoolResponse struct {
	Count    int64 `json:"count"`
	Vsize    int64 `json:"vsize"`
	TotalFee int64 `json:"total_fee"`
}

func GetMempoolStats() (interface{}, error) {
	const endpoint string = "/mempool"
	const maxBlockVSize int64 = 1000000 // maxBlockWeight = 4_000_000 (vSize = weight / 4)

	req := pdk.NewHTTPRequest(pdk.MethodGet, mempoolspaceAPI+endpoint)
	res := req.Send()

	if res.Status() != 200 {
		errMsg := "Mempool.Space API request returned status: " + strconv.Itoa(int(res.Status()))
		return MempoolStats{}, errors.New(errMsg)
	}

	var data MempoolResponse
	err := json.Unmarshal(res.Body(), &data)
	if err != nil {
		return MempoolStats{}, err
	}

	var mempoolStats MempoolStats
	mempoolStats.UnconfirmedTXs = data.Count
	mempoolStats.VSize = roundDec(float64(data.Vsize)/bytesPerMegabyte, 2)
	mempoolStats.PendingFees = roundDec(float64(data.TotalFee)/satsPerBTC, 3)
	mempoolStats.BlocksToClear = int((data.Vsize / maxBlockVSize) + 1)

	return mempoolStats, nil
}

type NodeStats struct {
	TotalNodes      int     `json:"totalNodes"`
	TorNodes        int     `json:"torNodes"`
	TorNodesPercent float64 `json:"torNodesPercent"`
}

type BitnodesResponse struct {
	Timestamp    int64                    `json:"timestamp"`
	TotalNodes   int                      `json:"total_nodes"`
	LatestHeight int                      `json:"latest_height"`
	Nodes        map[string][]interface{} `json:"nodes"`
}

func GetBitcoinNodeStats() (interface{}, error) {
	const endpoint string = "/snapshots/latest"

	req := pdk.NewHTTPRequest(pdk.MethodGet, bitnodesAPI+endpoint)
	res := req.Send()

	if res.Status() != 200 {
		errMsg := "BitNodes.IO API request returned status: " + strconv.Itoa(int(res.Status()))
		return NodeStats{}, errors.New(errMsg)
	}

	var data BitnodesResponse
	err := json.Unmarshal(res.Body(), &data)
	if err != nil {
		return NodeStats{}, err
	}

	var nodeStats NodeStats
	nodeStats.TotalNodes = data.TotalNodes

	torNodes := 0
	for _, details := range data.Nodes {
		if details[11] == "TOR" {
			torNodes++
		}
	}

	nodeStats.TorNodes = torNodes
	nodeStats.TorNodesPercent = roundDec((float64(torNodes)/float64(data.TotalNodes))*100, 2)

	return nodeStats, nil
}

type Fees struct {
	Fastest  int `json:"fastest"`
	HalfHour int `json:"halfHour"`
	Hour     int `json:"hour"`
	Economy  int `json:"economy"`
	Minimum  int `json:"minimum"`
}

type FeeResponse struct {
	FastestFee  int `json:"fastestFee"`
	HalfHourFee int `json:"halfHourFee"`
	HourFee     int `json:"hourFee"`
	EconomyFee  int `json:"economyFee"`
	MinimumFee  int `json:"minimumFee"`
}

func GetFeeRecommendation() (interface{}, error) {
	const endpoint string = "/v1/fees/recommended"

	req := pdk.NewHTTPRequest(pdk.MethodGet, mempoolspaceAPI+endpoint)
	res := req.Send()

	if res.Status() != 200 {
		errMsg := "Mempool.Space API request returned status: " + strconv.Itoa(int(res.Status()))
		return Fees{}, errors.New(errMsg)
	}

	var data FeeResponse
	err := json.Unmarshal(res.Body(), &data)
	if err != nil {
		return Fees{}, err
	}

	var fees Fees
	fees.Fastest = data.FastestFee
	fees.HalfHour = data.HalfHourFee
	fees.Hour = data.HourFee
	fees.Economy = data.EconomyFee
	fees.Minimum = data.MinimumFee

	return fees, nil
}

type LNStats struct {
	TotalNodes             int     `json:"totalNodes"`
	TorNodes               int     `json:"torNodes"`
	ClearnetNodes          int     `json:"clearnetNodes"`
	ClearnetTorNodes       int     `json:"clearnetTorNodes"`
	UnannouncedNodes       int     `json:"unannouncedNodes"`
	Channels               int     `json:"channels"`
	TotalCapacity          float64 `json:"totalCapacity"`
	AverageChannelCapacity float64 `json:"averageChannelCapacity"`
	MedianChannelCapacity  float64 `json:"medianChannelCapacity"`
}

type LNResponse struct {
	Latest struct {
		ChannelCount     int   `json:"channel_count"`
		NodeCount        int   `json:"node_count"`
		TotalCapacity    int64 `json:"total_capacity"`
		TorNodes         int   `json:"tor_nodes"`
		ClearnetNodes    int   `json:"clearnet_nodes"`
		UnannouncedNodes int   `json:"unannounced_nodes"`
		AverageCapacity  int64 `json:"avg_capacity"`
		MedianCapacity   int64 `json:"med_capacity"`
		ClearnetTorNodes int   `json:"clearnet_tor_nodes"`
	} `json:"latest"`
}

func GetLightningStats() (interface{}, error) {
	const endpoint string = "/v1/lightning/statistics/latest"

	req := pdk.NewHTTPRequest(pdk.MethodGet, mempoolspaceAPI+endpoint)
	res := req.Send()

	if res.Status() != 200 {
		errMsg := "Mempool.Space API request returned status: " + strconv.Itoa(int(res.Status()))
		return LNStats{}, errors.New(errMsg)
	}

	var data LNResponse
	err := json.Unmarshal(res.Body(), &data)
	if err != nil {
		return LNStats{}, err
	}

	var lnStats LNStats
	lnStats.TotalNodes = data.Latest.NodeCount
	lnStats.TorNodes = data.Latest.TorNodes
	lnStats.ClearnetNodes = data.Latest.ClearnetNodes
	lnStats.ClearnetTorNodes = data.Latest.ClearnetTorNodes
	lnStats.UnannouncedNodes = data.Latest.UnannouncedNodes
	lnStats.Channels = data.Latest.ChannelCount
	lnStats.TotalCapacity = roundDec(float64(data.Latest.TotalCapacity)/satsPerBTC, 3)
	lnStats.AverageChannelCapacity = roundDec(float64(data.Latest.AverageCapacity)/satsPerBTC, 3)
	lnStats.MedianChannelCapacity = roundDec(float64(data.Latest.MedianCapacity)/satsPerBTC, 3)

	return lnStats, nil
}

type BlockInfo struct {
	Height    int     `json:"height"`
	Timestamp string  `json:"timestamp"`
	TXs       int64   `json:"transactions"`
	Size      float64 `json:"size"`
	Reward    float64 `json:"totalReward"`
	TotalFees float64 `json:"totalFees"`
	FeeRate   float64 `json:"medianFeeRate"`
	Miner     string  `json:"miner"`
}

type BlockResponse struct {
	Height    int   `json:"height"`
	Timestamp int64 `json:"timestamp"`
	TXCount   int64 `json:"tx_count"`
	Size      int64 `json:"size"`
	Extras    struct {
		Reward    int64   `json:"reward"`
		TotalFees int64   `json:"totalFees"`
		FeeRate   float64 `json:"medianFee"`
		Pool      struct {
			Name string `json:"name"`
		} `json:"pool"`
	} `json:"extras"`
}

func GetLatestBlockInfo() (interface{}, error) {
	height, err := getTipHeight()
	if err != nil {
		return BlockInfo{}, err
	}

	endpoint := "/v1/blocks/" + string(height)

	req := pdk.NewHTTPRequest(pdk.MethodGet, mempoolspaceAPI+endpoint)
	res := req.Send()

	if res.Status() != 200 {
		errMsg := "Mempool.Space API request returned status: " + strconv.Itoa(int(res.Status()))
		return BlockInfo{}, errors.New(errMsg)
	}

	var dataArray []BlockResponse
	err = json.Unmarshal(res.Body(), &dataArray)
	if err != nil {
		return BlockInfo{}, err
	}

	data := dataArray[0]

	var block BlockInfo
	block.Height = data.Height
	block.Timestamp = time.Unix(data.Timestamp, 0).Format(time.RFC1123)
	block.TXs = data.TXCount
	block.Size = roundDec(float64(data.Size)/bytesPerMegabyte, 2)
	block.Reward = roundDec(float64(data.Extras.Reward)/satsPerBTC, 3)
	block.TotalFees = roundDec(float64(data.Extras.TotalFees)/satsPerBTC, 3)
	block.FeeRate = roundDec(data.Extras.FeeRate, 1)
	block.Miner = data.Extras.Pool.Name

	return block, nil
}

func getTipHeight() (string, error) {
	const endpoint string = "/blocks/tip/height"

	req := pdk.NewHTTPRequest(pdk.MethodGet, mempoolspaceAPI+endpoint)
	res := req.Send()

	if res.Status() != 200 {
		errMsg := "Mempool.Space API request returned status: " + strconv.Itoa(int(res.Status()))
		return "", errors.New(errMsg)
	}

	height := string(res.Body())

	return height, nil
}

func roundDec(x float64, d int) float64 {
	fac := math.Pow(10, float64(d))
	return float64(math.Round(x*fac) / fac)
}
