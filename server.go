package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/go-pg/pg/v10"
	// pq postgresql driver
	_ "github.com/lib/pq"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/buidl-labs/miner-marketplace-backend/graph"
	"github.com/buidl-labs/miner-marketplace-backend/graph/generated"
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
	srv := handler.NewDefaultServer(
		generated.NewExecutableSchema(generated.Config{
			Resolvers: &graph.Resolver{
				DB:   myNewDB,
				PQDB: myNewPqDB,
			},
		}),
	)

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func NewDB() (*pg.DB, error) {
	// dburl postgres://rajdeep@localhost/filecoinminermarketplace?sslmode=disable
	// db := pg.Connect(&pg.Options{
	// 	Addr:     ":5432",
	// 	User:     "rajdeep",
	// 	Database: "filecoinminermarketplace",
	// })
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
