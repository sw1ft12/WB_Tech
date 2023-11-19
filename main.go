package main

import (
	"awesomeProject/cache"
	"awesomeProject/http_server"
	"awesomeProject/postgresql"
	"awesomeProject/repository"
	"context"
	"encoding/json"
	"github.com/nats-io/stan.go"
	"log"
)

func main() {
	ctx := context.Background()
	postgreSQLConn, err := postgresql.NewConnection(ctx, 5,
		postgresql.Parameters{Username: "postgres", Password: "postgres", Host: "localhost",
			Port: "5432", Database: "postgres"})

	if err != nil {
		log.Fatalf("failed to connect to Postgres: %v", err)
	}

	rep := repository.Repository{Conn: postgreSQLConn}

	storage := cache.RestoreCache(ctx, &rep)

	sc, err := stan.Connect("my_cluster", "sw1ft")
	if err != nil {
		log.Fatalf("failed to connect to NATS server: %v", err)
	}

	var sub stan.Subscription
	sub, err = sc.Subscribe("foo", func(msg *stan.Msg) {
		var dt repository.Order

		if err = json.Unmarshal(msg.Data, &dt); err != nil {
			log.Fatalf("Cannot unmarshal: %v", err)
		}

		if err = rep.Create(ctx, dt); err != nil {
			log.Fatalf("Bad SQL query: %v", err)
		}

		storage.Set(dt.Id, dt, 0)
	})

	if err != nil {
		log.Fatalf("Failed to subscribe: %v", err)
	}

	handler := http_server.Handler{Storage: storage, Sc: sc}
	router := handler.InitRoutes()
	err = router.Run("localhost:8000")

	if err != nil {
		log.Fatalf("Cannot run http-server: %v", err)
	}

	if err = sub.Unsubscribe(); err != nil {
		log.Fatalf("Failed to unsubscribe: %v", err)
	}

	if err = sc.Close(); err != nil {
		log.Fatalf("Failed to close: %v", err)
	}

	postgreSQLConn.Close()
}
