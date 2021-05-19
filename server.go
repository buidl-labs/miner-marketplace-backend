package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/buidl-labs/filecoin-chain-indexer/config"
	"github.com/buidl-labs/filecoin-chain-indexer/lens/lotus"
	"github.com/go-pg/pg/v10"
	_ "github.com/lib/pq"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/buidl-labs/miner-marketplace-backend/graph"
	"github.com/buidl-labs/miner-marketplace-backend/graph/generated"
	"github.com/go-chi/chi"
	"github.com/rs/cors"
)

const defaultPort = "8080"

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	myNewDB, err := NewDB()
	if err != nil {
		log.Fatal("connecting to db: ", err)
	}
	myNewPqDB, err := NewPqDB()
	if err != nil {
		log.Fatal("connecting to db: ", err)
	}

	router := chi.NewRouter()

	router.Use(cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowCredentials: true,
		Debug:            true,
	}).Handler)

	lensOpener, lensCloser, err := lotus.NewAPIOpener(config.Config{
		FullNodeAPIInfo: os.Getenv("FULLNODE_API_INFO"),
		CacheSize:       1,
	}, context.Background())
	if err != nil {
		log.Fatal("creating lotus API opener: ", err)
	}
	defer func() {
		lensCloser()
	}()
	node, closer, err := lensOpener.Open(context.Background())
	if err != nil {
		log.Fatal("opening lotus lens API: ", err)
	}
	defer closer()

	srv := handler.NewDefaultServer(
		generated.NewExecutableSchema(generated.Config{
			Resolvers: &graph.Resolver{
				DB:      myNewDB,
				PQDB:    myNewPqDB,
				LensAPI: node,
			},
		}),
	)

	router.Handle("/", playground.Handler("GraphQL playground", "/query"))
	router.Handle("/query", srv)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}

func NewDB() (*pg.DB, error) {
	opt, err := pg.ParseURL(os.Getenv("DB"))
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	db := pg.Connect(opt)
	return db, nil
}

func NewPqDB() (*sql.DB, error) {
	db, err := sql.Open("postgres", os.Getenv("DB"))
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	return db, nil
}
