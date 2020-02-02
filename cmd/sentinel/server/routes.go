package server

import (
	"encoding/json"
	"net/http"

	"github.com/sirupsen/logrus"

	"github.com/gorilla/mux"
)

// // NewRouter returns the sentinel's router.
// func NewRouter() *mux.Router {
// 	router := mux.NewRouter()
// 	return router
// }

func (s *Server) Routes() {
	s.router.Methods("GET").Path("/attack/{URL}").HandlerFunc(s.getAttack)
	s.router.Methods("POST").Path("/attack").HandlerFunc(s.postAttack)
}

func (s *Server) getAttack(res http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	URL := vars["URL"]

	// Publish the message to the attack topic.
	err := s.redisClient.Publish("topic", URL).Err()
	if err != nil {
		logrus.Info("failed to publish to topic", err)
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

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
