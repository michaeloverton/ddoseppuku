package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/michaeloverton/ddos-laser/cmd/target/router"
	"github.com/michaeloverton/ddos-laser/internal/env"
)

func main() {
	// Load environment.
	env, err := env.LoadTargetEnv()
	if err != nil {
		log.Fatal("error loading environment: ", err.Error())
	}

	// Set up the router.
	router := router.NewRouter()

	// Start the server
	s := &http.Server{
		Addr:    ":" + env.Port,
		Handler: router,
	}
	go func() {
		log.Info("target serving on: ", env.Port)
		if err := s.ListenAndServe(); err != nil {
			log.Fatal("server failure", err)
		}
	}()

	// Allow interrupt signal to gracefully shutdown with a 5-second timeout.
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err = s.Shutdown(ctx)
	if err != nil {
		panic(err)
	} else {
		log.Info("gracefully shutting down")
	}
}
