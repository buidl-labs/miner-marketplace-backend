# miner-marketplace-backend

GraphQL API for Filecoin Miner Marketplace

## Getting Started

### Build

```
go build
```

### Update env vars

Sample present in [.env.sample](.env.sample)

### Run db migrations

```
goose -dir ./db/migration postgres "user=username dbname=fmmdbname sslmode=disable" up
```

### Run server

```
./miner-marketplace-backend
```

## Data Indexing Jobs

Set flag `--cmd` to invoke data indexing commands

| Command | Description                                                                                                                                    |
| ------- | ---------------------------------------------------------------------------------------------------------------------------------------------- |
| `idx`   | Index all data periodically.                                                                                                                   |
| `mm`    | Miner actor messages ([sample](https://filfox.info/api/v1/address/f0115238/messages?pageSize=100&page=0)).                                     |
| `dtc`   | Storage deals ([sample](https://filfox.info/api/v1/deal/list?pageSize=100&page=0)).                                                            |
| `psdm`  | PublishStorageDeals messages ([sample](https://filfox.info/api/v1/message/list?pageSize=100&page=0&method=PublishStorageDeals)).               |
| `wbmm`  | WithdrawBalance (market) messages ([sample](<https://filfox.info/api/v1/message/list?pageSize=100&page=0&method=WithdrawBalance%20(market)>)). |
