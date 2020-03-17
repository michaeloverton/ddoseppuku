package server

import (
	"encoding/json"
	"net/http"

	"github.com/michaeloverton/ddoseppuku/internal/message"
	"github.com/sirupsen/logrus"
)

func (s *Server) Routes() {
	s.router.Methods("POST").Path("/attack").HandlerFunc(s.postAttack)
}

func (s *Server) postAttack(res http.ResponseWriter, req *http.Request) {
	// Decode response.
	var m message.Message
	err := json.NewDecoder(req.Body).Decode(&m)
	if err != nil {
		logrus.Errorf("could not decode request body: %s", err)
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Ensure URL is present.
	if m.URL == "" {
		logrus.Error("URL is required")
		res.WriteHeader(http.StatusBadRequest)
	}

	// Enforce that request is GET or POST.
	if m.Method != http.MethodGet && m.Method != http.MethodPost {
		logrus.Error("method must be GET or POST")
		res.WriteHeader(http.StatusBadRequest)
	}

	// If method is POST, require body.
	if m.Method == http.MethodPost && len(m.Body) == 0 {
		logrus.Error("if method is POST, body is required")
		res.WriteHeader(http.StatusBadRequest)
	}

	// Publish the message to the attack topic.
	err = s.redisClient.Publish("attack", m).Err()
	if err != nil {
		logrus.Error("failed to publish to topic", err)
		res.WriteHeader(http.StatusInternalServerError)
		return
	}
}
