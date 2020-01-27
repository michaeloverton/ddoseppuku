package main

import (
	log "github.com/sirupsen/logrus"

	"github.com/michaeloverton/ddos-laser/internal/env"
)

func main() {
	// Load environment.
	_, err := env.LoadTargetEnv()
	if err != nil {
		log.Fatal("error loading environment: ", err.Error())
	}
}
