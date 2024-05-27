# Bitcoin Stats Plugin for OpenAgents

This is a plugin for [OpenAgents](https://openagents.com) that gives you the latest statistics about the Bitcoin network, its mempool and the Lightning Network.

It is inspired by Clark Moody's [Bitcoin Dashboard](https://bitcoin.clarkmoody.com/dashboard) and the [TimeChainCalendar](https://timechaincalendar.com).

The project uses the [Extism Framework](https://extism.org/docs/overview/), in particular its [Go PDK](https://github.com/extism/go-pdk), and the REST APIs of [mempool.space](https://mempool.space/docs/api/rest), [bitnodes.io](https://bitnodes.io/api) and [coincap.io](https://docs.coincap.io/).

## Installation

Make sure you have both [Go](https://go.dev/doc/install) and [TinyGo](https://tinygo.org/getting-started/install/).

Clone this repo:

```sh
git clone https://github.com/sfr0xyz/openagents-bitcoin-stats.git
```

Build the plugin with TinyGo:

```sh
tinygo build -o btcstats.wasm -target wasi .
```

## Usage

Ensure that you have installed the [Extism CLI](https://github.com/extism/cli?tab=readme-ov-file#extism-cli) and downloaded the `btcstats.wasm` file.

You can call the plugin with the Extism CLI:

```sh
extism call btcstats.wasm run --input '<YOUR INPUT>' --wasi --allow-host '*'
```

Replace `<YOUR INPUT>` with a list of statistics you are interested in.

Available statistics (see detailed descriptions [below](#description-of-the-values))

- `market`: Bitcoin market data
- `latestBlock`: Information about the latest block
- `mining`: Mining data like the current hashrate and difficulty
- `fees`: Recommended fees based on the current mempool
- `mempool`: Mempool stats
- `lightning`: Lightning Network stats
- `nodes`: Bitcoin nodes stats

You can include more than one of the above at once.

If you leave the field empty (`''`), or include none of the above, all stats will be requested.

With the prefix `-`, e.g. `-nodes`, you can exclude stats, i.e. "I want all stats except `nodes`".

> Note: If you request the `nodes` statistic, it will take a while, a few seconds, for the result to be displayed.

### Examples

> Note: The resulting JSON is prettified here for better readability.

Get **all** stats:

```plain
$ extism call btcstats.wasm run --input '' --wasi --allow-host '*'

{
  "fees": {
    "fastest": 9,
    "halfHour": 9,
    "hour": 9,
    "economy": 6,
    "minimum": 3
  },
  "latestBlock": {
    "height": 845325,
    "timestamp": "Mon, 27 May 2024 01:34:54 UTC",
    "transactions": 5692,
    "size": 1.55,
    "totalReward": 3.272,
    "totalFees": 0.147,
    "medianFeeRate": 8.1,
    "miner": "MARA Pool"
  },
  "lightning": {
    "totalNodes": 12836,
    "torNodes": 8930,
    "clearnetNodes": 1700,
    "clearnetTorNodes": 1360,
    "unannouncedNodes": 846,
    "channels": 50872,
    "totalCapacity": 4980.822,
    "averageChannelCapacity": 0.098,
    "medianChannelCapacity": 0.02
  },
  "market": {
    "supply": 19699693,
    "supplyPercent": 93.81,
    "price": 69020.61,
    "priceChange24hPercent": -0.32,
    "moscowTime": 1448,
    "marketCap": 1359684729700.71
  },
  "mempool": {
    "unconfirmedTXs": 170669,
    "vSize": 180.66,
    "pendingFees": 5.257,
    "blocksToClear": 181
  },
  "mining": {
    "hashrate": 677612693053317300000,
    "difficulty": 84381461788831.34,
    "retargetDifficultyChangePercent": 13.39,
    "retargetRemainingBlocks": 1395,
    "retargetEstimatedDate": "Tue, 04 Jun 2024 15:07:02 UTC"
  },
  "nodes": {
    "totalNodes": 18037,
    "torNodes": 10942,
    "torNodesPercent": 60.66
  }
}
```

Get **all** stats except `nodes`:

```plain
$ extism call btcstats.wasm run --input '-nodes' --wasi --allow-host '*'

{
  "fees": {
    "fastest": 9,
    "halfHour": 9,
    "hour": 9,
    "economy": 6,
    "minimum": 3
  },
  "latestBlock": {
    "height": 845325,
    "timestamp": "Mon, 27 May 2024 01:34:54 UTC",
    "transactions": 5692,
    "size": 1.55,
    "totalReward": 3.272,
    "totalFees": 0.147,
    "medianFeeRate": 8.1,
    "miner": "MARA Pool"
  },
  "lightning": {
    "totalNodes": 12836,
    "torNodes": 8930,
    "clearnetNodes": 1700,
    "clearnetTorNodes": 1360,
    "unannouncedNodes": 846,
    "channels": 50872,
    "totalCapacity": 4980.822,
    "averageChannelCapacity": 0.098,
    "medianChannelCapacity": 0.02
  },
  "market": {
    "supply": 19699693,
    "supplyPercent": 93.81,
    "price": 69020.61,
    "priceChange24hPercent": -0.32,
    "moscowTime": 1448,
    "marketCap": 1359684729700.71
  },
  "mempool": {
    "unconfirmedTXs": 170669,
    "vSize": 180.66,
    "pendingFees": 5.257,
    "blocksToClear": 181
  },
  "mining": {
    "hashrate": 677612693053317300000,
    "difficulty": 84381461788831.34,
    "retargetDifficultyChangePercent": 13.39,
    "retargetRemainingBlocks": 1395,
    "retargetEstimatedDate": "Tue, 04 Jun 2024 15:07:02 UTC"
  }
}
```

Get `latestBlock` and `mempool` stats:

```plain
$ extism call btcstats.wasm run --input 'latestBlock mempool' --wasi --allow-host '*'

{
  "latestBlock": {
    "height": 845325,
    "timestamp": "Mon, 27 May 2024 01:34:54 UTC",
    "transactions": 5692,
    "size": 1.55,
    "totalReward": 3.272,
    "totalFees": 0.147,
    "medianFeeRate": 8.1,
    "miner": "MARA Pool"
  },
  "mempool": {
    "unconfirmedTXs": 170669,
    "vSize": 180.66,
    "pendingFees": 5.257,
    "blocksToClear": 181
  },
  "mining": {
    "hashrate": 677612693053317300000,
    "difficulty": 84381461788831.34,
    "retargetDifficultyChangePercent": 13.39,
    "retargetRemainingBlocks": 1395,
    "retargetEstimatedDate": "Tue, 04 Jun 2024 15:07:02 UTC"
  }
}
```

## Description of the values

- `market`
  - `price`: current US Dollar price for 1 bitcoin
  - `priceChange24hPercent`: US Dollar price change withng the last 24 hour in percent
  - `moscowTime`: current Moscow time, i.e. SATS/USD
  - `supply`: current amount of available bitcoin
  - `supplyPercent`: percentage of current supply from the maximum of 21'000'000 bitcoin
  - `marketCap`: current market capitalisation (price x supply)
- `latestBlock`
  - `height`: block height
  - `timestamp`: timestamp string in the format dd mmm YYYY HH:mm:ss Z
  - `transactions`: number of transactions in the latest block
  - `size`: size in MB
  - `reward`: total block reward, i.e. subsidy + total fees, in BTC
  - `totalFees`: total amount of fees from the transactions in BTC
  - `medianFeeRate`: median fee rate paid by the included transactions in SATS/vB
  - `miner`: name of the miner's mining pool
- `mining`
  - `hashrate`: current hashrate
  - `difficulty`" current difficulty
  - `retargetDifficultyChangePercent`: difficulty change at the next difficulty adjustment in percent
  - `retargetRemainingBlocks`: remaining blocks until difficulty adjustment
  - `retargetEstimatedDate`: estimated date of the difficulty adjustment (dd mmm YYYY HH:mm:ss Z)
- `fees`
  - `fastest`
  - `halfHour`
  - `hour`
  - `economy`
  - `minimum`
- `mempool`
  - `unconfirmedTXs`: current number of transactions in the mempool
  - `vSize`: virtual size of all transaction in the mempool
  - `pendingFees`: total amount of fees from the transaction in the mempool in BTC
  - `blocksToClear`: estimated number of blocks until all transaction in the mempool are included in a block
- `lightning`
  - `totalNodes`: total number of Lightning nodes
  - `torNodes`: number Lightning nodes running behind TOR
  - `clearnetNodes`: number of Lightning nodes running on clearnet
  - `clearnetTorNodes`: number of Lightning nodes that run on clearnet and behind TOR
  - `unannouncedNodes`: number of unannounced Lightning nodes
  - `channels`: number of announced Lightning channels
  - `totalCapacity`: total amount of (visible) BTC in Lightning channels
  - `averageChannelCapacity`: average amount of BTC per channel
  - `medianChannelCapacity`: median amount of BTC per channel
- `nodes`
  - `totalNodes`: total number of Bitcoin nodes
  - `torNodes`: number of Bitcoin nodes behind TOR
  - `torNodesPercent`: percentage of nodes behind TOR from the total number
