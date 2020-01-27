package router

import (
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/michaeloverton/ddos-laser/internal/server"
	log "github.com/sirupsen/logrus"
)

// NewRouter returns the targets's router.
func NewRouter() *mux.Router {
	router := mux.NewRouter()
	router.Methods("GET").Path("/thrash").HandlerFunc(thrash)
	return router
}

type targetResponse struct {
	Text string `json:"text"`
}

func thrash(res http.ResponseWriter, req *http.Request) {
	log.Info("responding to request from: ", req.Host)
	resp := targetResponse{
		Text: "whatever",
	}
	time.Sleep(time.Second * 5)
	server.Serve(res, resp, 200)
}
