package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/go-redis/redis"
	"github.com/michaeloverton/ddos-laser/cmd/sentinel/server"
	"github.com/michaeloverton/ddos-laser/internal/env"
	log "github.com/sirupsen/logrus"
)

func main() {
	// Load environment.
	env, err := env.LoadSentinelEnv()
	if err != nil {
		log.Fatal("error loading environment: ", err.Error())
	}

	// Initialize Redis client.
	// redisClient, err := redis.New(env.RedisAddress)
	// if err != nil {
	// 	log.Fatal("failed to initialize redis client: ", err)
	// }

	// // test publishing to a redis topic
	// err = redisClient.Publish("target", "target.com")
	// if err != nil {
	// 	log.Fatal("failed to publish to topic")
	// }

	pubClient := redis.NewClient(&redis.Options{
		Addr:     env.RedisAddress,
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	// Test connection to client.
	_, err = pubClient.Ping().Result()
	if err != nil {
		log.Fatal("failed to connect to redis", err)
	}

	// Set up server and routes.
	s := server.NewServer(pubClient)
	s.Routes()

	// Start the server
	server := &http.Server{
		Addr:    ":" + env.Port,
		Handler: s.Router(),
	}
	go func() {
		log.Info("sentinel serving on: ", env.Port)
		if err := server.ListenAndServe(); err != nil {
			log.Fatal("server failure", err)
		}
	}()

	// Allow interrupt signal to gracefully shutdown with a 5-second timeout.
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err = server.Shutdown(ctx)
	if err != nil {
		panic(err)
	} else {
		log.Info("gracefully shutting down")
	}

}
