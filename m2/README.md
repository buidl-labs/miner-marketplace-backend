# miner-marketplace-backend

GraphQL API for Filecoin Miner Marketplace

## Getting Started

- Clone the repo: `git clone https://github.com/buidl-labs/miner-marketplace-backend` and `cd` into it.
- Build: `go build`
- Set environment variables ([sample](.env.sample))
  - `DB` is the uri of the postgres db.
  - `ADDR_CHANGES_URL` is the url where the [address changes file](https://github.com/buidl-labs/filecoin-chain-indexer/blob/main/servedir/addrChanges.json) is served (typically by the indexer). The default value provided in sample will work if any of the [indexer tasks](https://github.com/buidl-labs/filecoin-chain-indexer#list-of-tasks) are running.
- Run: `./miner-marketplace-backend`
