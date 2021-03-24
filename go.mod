module github.com/buidl-labs/miner-marketplace-backend

go 1.15

require (
	github.com/99designs/gqlgen v0.13.0
	github.com/buidl-labs/filecoin-chain-indexer v0.0.0-20210324200839-a1f798dd8e88
	github.com/go-chi/chi v3.3.2+incompatible
	github.com/go-pg/pg/v10 v10.7.5
	github.com/lib/pq v1.9.0
	github.com/rs/cors v1.7.0
	github.com/vektah/gqlparser/v2 v2.1.0
)

replace (
	github.com/filecoin-project/fil-blst => ./extern/fil-blst
	github.com/filecoin-project/filecoin-ffi => ./extern/filecoin-ffi-stub
	github.com/supranational/blst => ./extern/fil-blst/blst
)

// Supports go-ipld-prime v7
// TODO: remove once https://github.com/filecoin-project/statediff/pull/155 is merged
replace github.com/filecoin-project/statediff => github.com/filecoin-project/statediff v0.0.19-0.20210225063407-9e38aa4b7ede

// Supports go-ipld-prime v7
// TODO: remove once https://github.com/filecoin-project/go-hamt-ipld/pull/70 is merged
replace github.com/filecoin-project/go-hamt-ipld/v2 => github.com/willscott/go-hamt-ipld/v2 v2.0.1-0.20210225034344-6d6dfa9b3960
