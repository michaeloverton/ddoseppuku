package server

import (
	"encoding/json"
	"net/http"

	"github.com/michaeloverton/ddoseppuku/internal/attack"
	"github.com/michaeloverton/ddoseppuku/internal/message"

	"github.com/sirupsen/logrus"
)

func (s *Server) Routes() {
	s.router.Methods("POST").Path("/attack").HandlerFunc(s.postAttack)
	s.router.Methods("GET").Path("/attack").HandlerFunc(s.getAttack)
	s.router.Methods("GET").Path("/ceasefire").HandlerFunc(s.quit)
}

func (s *Server) postAttack(res http.ResponseWriter, req *http.Request) {
	// Decode response.
	var a attack.Attack
	err := json.NewDecoder(req.Body).Decode(&a)
	if err != nil {
		logrus.Error("could not decode request body", err)
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	logrus.Info("sentinel attack: ", a)

	// Ensure URL is valid.
	if a.URL == "" {
		logrus.Error("URL is required", err)
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Publish the message to the attack topic.
	err = s.redisClient.Publish("topic", a).Err()
	if err != nil {
		logrus.Error("failed to publish to topic: ", err)
		res.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (s *Server) getAttack(res http.ResponseWriter, req *http.Request) {
	// Get the attack URL from the query params.
	attackURL := req.URL.Query().Get("url")

	// Ensure URL is present.
	if attackURL == "" {
		logrus.Error("URL is required")
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Form the attack message.
	msg := message.Message{
		URL:  attackURL,
		Type: "GET",
	}

	// Publish the message to the attack topic.
	err := s.redisClient.Publish("topic", msg).Err()
	if err != nil {
		logrus.Error("failed to publish to topic: ", err)
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
