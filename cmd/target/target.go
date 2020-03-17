package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"

	"github.com/michaeloverton/ddoseppuku/internal/env"
)

func main() {
	// Load environment.
	env, err := env.LoadTargetEnv()
	if err != nil {
		log.Fatal("error loading environment: ", err.Error())
	}

	// Set up endpoints.
	http.HandleFunc("/health", health)
	http.HandleFunc("/thrash", taskHandler(env))
	http.HandleFunc("/login", loginHandler(env))

	// Serve.
	log.Infof("target serving on: %s", env.Port)
	err = http.ListenAndServe(":"+env.Port, nil)
	if err != nil {
		log.Fatal("failed to serve", err)
	}
}

// health is the health check endpoint.
func health(res http.ResponseWriter, req *http.Request) {
	res.WriteHeader(http.StatusOK)
}

// taskHandler performs a mock task of configurable intensity.
func taskHandler(env *env.TargetEnv) http.HandlerFunc {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		// Perform the task. Intensity set via env var.
		start := time.Now()
		for i := 0; i < env.TaskIntensity; i++ {
			_ = reverse(positiveInfinity)
		}
		elapsed := time.Since(start)
		log.Info("negating infinity took: ", elapsed)

		res.WriteHeader(http.StatusOK)
	})
}

type Login struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// loginHandler receives a request body and also performs a mock task of configurable intensity.
func loginHandler(env *env.TargetEnv) http.HandlerFunc {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		var l Login
		err := json.NewDecoder(req.Body).Decode(&l)
		if err != nil {
			logrus.Error("could not decode request body", err)
			res.WriteHeader(http.StatusInternalServerError)
			return
		}

		// Perform the task. Intensity set via env var.
		start := time.Now()
		for i := 0; i < env.TaskIntensity; i++ {
			_ = reverse(positiveInfinity)
		}
		elapsed := time.Since(start)
		log.Info("negating infinity took: ", elapsed)

		res.WriteHeader(http.StatusOK)
	})
}

func reverse(s string) string {
	chars := []rune(s)
	for i, j := 0, len(chars)-1; i < j; i, j = i+1, j-1 {
		chars[i], chars[j] = chars[j], chars[i]
	}
	return string(chars)
}
