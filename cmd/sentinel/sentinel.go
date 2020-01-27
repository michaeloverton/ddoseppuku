package main

import (
	"net/http"

	"github.com/michaeloverton/ddos-laser/internal/env"
	log "github.com/sirupsen/logrus"
)

func main() {
	// Load environment.
	env, err := env.LoadSentinelEnv()
	if err != nil {
		log.Fatal("error loading environment: ", err.Error())
	}

	// Set up handler function.
	http.HandleFunc("/", nil)

	// Serve.
	log.Infof("starting server on port: %s", env.Port)
	err = http.ListenAndServe(":"+env.Port, nil)
	if err != nil {
		log.Fatal("failed to serve", err)
	}

}
