package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/buidl-labs/filecoin-chain-indexer/config"
	"github.com/buidl-labs/filecoin-chain-indexer/lens/lotus"
	"github.com/buidl-labs/miner-marketplace-backend/graph"
	"github.com/buidl-labs/miner-marketplace-backend/graph/generated"
	"github.com/buidl-labs/miner-marketplace-backend/service"
	"github.com/go-chi/chi"
	"github.com/go-pg/pg/v10"
	"github.com/rs/cors"
)

const defaultPort = "8080"

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	newDB, err := NewDB()
	if err != nil {
		log.Fatal("connecting to db: ", err)
	}

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

	var command string
	flag.StringVar(&command, "cmd", "", "Command to run")
	flag.Parse()

	fmt.Println("command selected:", command)
	if command == "psdm" {
		service.PublishStorageDealsMessages(newDB, node)
	} else if command == "wbmm" {
		service.WithdrawBalanceMarketMessages(newDB, node)
	} else if command == "abm" {
		service.AddBalanceMessages(newDB, node)
	} else if command == "mm" {
		service.MinerPageMessages(newDB, node)
	} else if command == "dtc" {
		service.StorageDeals(newDB, node)
	} else if command == "ams" {
		service.AddressMessages(newDB, node)
	} else if command == "idx" {
		service.Indexer(newDB, node)
	} else {
		// go service.Indexer(newDB, node)

		router := chi.NewRouter()

		router.Use(cors.New(cors.Options{
			AllowedOrigins:   []string{"https://datastation.app", "https://www.datastation.app", "https://filecoin-miner-marketplace.onrender.com"},
			AllowCredentials: true,
			AllowedMethods:   []string{"GET", "POST", "PUT", "OPTIONS"},
			AllowedHeaders:   []string{"*"},
			Debug:            true,
		}).Handler)
		// router.Use(cors.AllowAll().Handler)

		srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &graph.Resolver{
			DB:      newDB,
			LensAPI: node,
		}}))

		router.Handle("/", playground.Handler("GraphQL playground", "/query"))
		router.Handle("/query", srv)

		log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
		log.Fatal(http.ListenAndServe(":"+port, router))
	}
}

func NewDB() (*pg.DB, error) {
	opt, err := pg.ParseURL(os.Getenv("DB"))
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	db := pg.Connect(opt)
	// db = db.WithTimeout(time.Second * 20)
	return db, nil
}
