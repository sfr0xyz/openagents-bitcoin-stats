package main

import (
	"encoding/json"
	"strings"

	"github.com/extism/go-pdk"
)

//export run
func run() int32 {
	keywords := map[string]func() (interface{}, error){
		"market":      GetMarketData,
		"latestBlock": GetLatestBlockInfo,
		"mining":      GetMiningData,
		"fees":        GetFeeRecommendation,
		"mempool":     GetMempoolStats,
		"lightning":   GetLightningStats,
		"nodes":       GetBitcoinNodeStats,
	}

	outputData := make(map[string]interface{})

	input := strings.ToLower(pdk.InputString())

	includeData := ""

	for kw := range keywords {
		if strings.Contains(input, "-"+kw) {
			continue
		}
		if strings.Contains(input, kw) {
			includeData = input
			break
		}

		includeData = includeData + " " + kw
	}

	for kw, function := range keywords {
		if strings.Contains(includeData, kw) && !strings.Contains(includeData, "-"+kw) {
			data, err := function()
			if err != nil {
				pdk.SetErrorString("Failed to get " + kw + " data: " + err.Error())
				return 1
			}

			outputData[kw] = data
		}
	}

	output, err := json.Marshal(outputData)
	if err != nil {
		pdk.SetErrorString("Failed to marshal the result.")
		return 1
	}

	pdk.OutputString(string(output))

	return 0
}

func main() {}
