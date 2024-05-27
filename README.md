# Bitcoin Stats Plugin for OpenAgents

Compile it:

```sh
tinygo build -o btcstats.wasm -target wasi .
```

Run it:

```sh
extism call btcstats.wasm run --input '<YOUR INPUT>' --wasi --allow-host '*'
```

## Getting started

- market
  - price
  - price change 24h
  - moscow time
  - marketcap
  - supply
  - supply percent
- latestBlock
  - height
  - timestamp
  - tx count
  - size
  - reward
  - total fees
  - median fee rate
  - miner / miner pool
- mining
  - current hashrate
  - current difficulty
  - difficulty adjustment
    - chnage percent
    - remaining blocks
    - estimated retarget date
- fees
  - fastest
  - half hour
  - hour
  - economy
  - minimum
- mempool stats
  - unconfirmed txs
  - total vsize
  - total fees
- lightning stats
  - total nodes
  - tor nodes
  - clearnet nodes
  - clearnet+tor nodes
  - unannounced nodes
  - num channels
  - total capacity
  - average channel capacity
  - median channel capacity
- nodes stats
  - total nodes
  - tor nodes
  - tor nodes percent
