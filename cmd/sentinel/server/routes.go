package server

import (
	"encoding/json"
	"net/http"

	"github.com/sirupsen/logrus"
)

func (s *Server) Routes() {
	s.router.Methods("POST").Path("/attack").HandlerFunc(s.postAttack)
	s.router.Methods("GET").Path("/ceasefire").HandlerFunc(s.quit)
}

type attack struct {
	URL string `json:"url"`
}

func (s *Server) postAttack(res http.ResponseWriter, req *http.Request) {
	// Decode response.
	var a attack
	err := json.NewDecoder(req.Body).Decode(&a)
	if err != nil {
		logrus.Error("could not decode request body", err)
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Ensure URL is valid.
	if a.URL == "" {
		logrus.Error("URL is required", err)
		res.WriteHeader(http.StatusInternalServerError)
	}

	// Publish the message to the attack topic.
	err = s.redisClient.Publish("topic", a.URL).Err()
	if err != nil {
		logrus.Error("failed to publish to topic", err)
		res.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (s *Server) quit(res http.ResponseWriter, req *http.Request) {
	// Publish the message to the quit topic.
	err := s.redisClient.Publish("quit", "cease fire").Err()
	if err != nil {
		logrus.Error("failed to publish to topic", err)
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

}
