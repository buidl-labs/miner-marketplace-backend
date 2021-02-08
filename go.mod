module github.com/buidl-labs/miner-marketplace-backend

go 1.15

require (
	github.com/99designs/gqlgen v0.13.0
	github.com/buidl-labs/filecoin-chain-indexer v0.0.0-20210208150437-4352a908186e
	github.com/go-pg/pg/v10 v10.7.5
	github.com/vektah/gqlparser/v2 v2.1.0
	golang.org/x/exp v0.0.0-20200924195034-c827fd4f18b9 // indirect
	honnef.co/go/tools v0.0.1-2020.1.4 // indirect
)

replace (
	github.com/filecoin-project/fil-blst => ./extern/fil-blst
	github.com/filecoin-project/filecoin-ffi => ./extern/filecoin-ffi-stub
	github.com/supranational/blst => ./extern/fil-blst/blst
)
