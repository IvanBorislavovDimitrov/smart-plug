package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	_ "github.com/lib/pq"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/IvanBorislavovDimitrov/smart-charger/graph"
	"github.com/IvanBorislavovDimitrov/smart-charger/service"
)

const defaultPort = "8081"
const connStr = "postgresql://postgres:123@127.0.0.1:5432/smart_plug?sslmode=disable"

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	resolver := graph.NewResolver(service.NewPlugService(db))
	srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: resolver}))

	http.Handle("/", playground.Handler("GraphQL playground", "/graphql/plugs"))
	http.Handle("/graphql/plugs", srv)

	log.Printf("connect to http://localhost:%s/ for GraphQL server", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
