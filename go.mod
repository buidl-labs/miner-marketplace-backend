module github.com/buidl-labs/miner-marketplace-backend

go 1.16

require (
	github.com/99designs/gqlgen v0.13.0
	github.com/buidl-labs/filecoin-chain-indexer v0.0.0-20210530115401-579e482bb5d3
	github.com/filecoin-project/go-address v0.0.5
	github.com/filecoin-project/go-state-types v0.1.0
	github.com/filecoin-project/lotus v1.9.0
	github.com/filecoin-project/specs-actors/v4 v4.0.1
	github.com/go-chi/chi v4.1.2+incompatible
	github.com/go-pg/pg/v10 v10.9.3
	github.com/onsi/ginkgo v1.15.0 // indirect
	github.com/onsi/gomega v1.10.5 // indirect
	github.com/rs/cors v1.7.0
	github.com/vektah/gqlparser/v2 v2.1.0
)

replace (
	github.com/filecoin-project/fil-blst => ./extern/fil-blst
	github.com/filecoin-project/filecoin-ffi => ./extern/filecoin-ffi-stub
	github.com/supranational/blst => ./extern/fil-blst/blst
)

// // Supports go-ipld-prime v7
// // TODO: remove once https://github.com/filecoin-project/statediff/pull/155 is merged
// replace github.com/filecoin-project/statediff => github.com/filecoin-project/statediff v0.0.19-0.20210225063407-9e38aa4b7ede

// Supports go-ipld-prime v7
// TODO: remove once https://github.com/filecoin-project/go-hamt-ipld/pull/70 is merged
replace github.com/filecoin-project/go-hamt-ipld/v2 => github.com/willscott/go-hamt-ipld/v2 v2.0.1-0.20210225034344-6d6dfa9b3960
