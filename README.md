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

### Run localhost

```
./miner-marketplace-backend
```
