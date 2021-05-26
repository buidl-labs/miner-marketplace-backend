package main

import (
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/buidl-labs/miner-marketplace-backend/graph"
	"github.com/buidl-labs/miner-marketplace-backend/graph/generated"
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

	router := chi.NewRouter()

	router.Use(cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowCredentials: true,
		Debug:            true,
	}).Handler)

	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &graph.Resolver{
		DB: newDB,
	}}))

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
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
