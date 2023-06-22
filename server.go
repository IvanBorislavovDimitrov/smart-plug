package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/IvanBorislavovDimitrov/smart-charger/graph"
	"github.com/IvanBorislavovDimitrov/smart-charger/scheduler"
	"github.com/IvanBorislavovDimitrov/smart-charger/service"
	"github.com/go-co-op/gocron"
	"github.com/jackc/pgx/v5/pgxpool"
)

const defaultPort = "8081"
const defaultConnStr = "postgresql://postgres:123@localhost:5432/smart_plug?sslmode=disable"

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}
	connStr := os.Getenv("CONN_STR")
	if connStr == "" {
		connStr = defaultConnStr
	}
	fmt.Println("DB connection string: " + connStr)
	pool, err := pgxpool.New(context.Background(), connStr)
	if err != nil {
		panic(err)
	}
	defer pool.Close()
	plugService := service.NewPlugService(pool)
	resolver := graph.NewResolver(plugService)
	srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: resolver}))
	startPowerScheduler(plugService)
	http.Handle("/", playground.Handler("GraphQL playground", "/graphql/plugs"))
	http.Handle("/graphql/plugs", srv)

	log.Printf("connect to http://localhost:%s/ for GraphQL server", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func startPowerScheduler(plugService *service.PlugService) {
	powerScheduler := scheduler.NewPowerScheduler(plugService)
	s := gocron.NewScheduler(time.Local)
	log.Println("Scheduler was configured to start every 12 hours")
	s.Cron("*/15 * * * *").Do(powerScheduler.ReconcilePlugsStates)
	s.Cron("1 */12 * * *").Do(powerScheduler.TurnOnPlugs)
	s.StartAsync()
}
