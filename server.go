package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-co-op/gocron"
	"github.com/jackc/pgx/v5"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/IvanBorislavovDimitrov/smart-charger/graph"
	"github.com/IvanBorislavovDimitrov/smart-charger/scheduler"
	"github.com/IvanBorislavovDimitrov/smart-charger/service"
)

const defaultPort = "8081"
const connStr = "postgresql://postgres:123@127.0.0.1:5432/smart_plug?sslmode=disable"

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}
	conn, err := pgx.Connect(context.Background(), connStr)
	if err != nil {
		panic(err)
	}
	defer conn.Close(context.Background())
	plugService := service.NewPlugService(conn)
	resolver := graph.NewResolver(plugService)
	srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: resolver}))

	http.Handle("/", playground.Handler("GraphQL playground", "/graphql/plugs"))
	http.Handle("/graphql/plugs", srv)

	log.Printf("connect to http://localhost:%s/ for GraphQL server", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func startPowerScheduler(plugService *service.PlugService) {
	powerScheduler := scheduler.NewPowerScheduler(plugService)
	s := gocron.NewScheduler(time.Local)
	log.Println("Scheduler was configured for every 12 hours")
	s.Every(1).Minutes().Do(powerScheduler.ReconcilePlugsStates)
}
