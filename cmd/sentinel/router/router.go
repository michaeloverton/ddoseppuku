package router

import (
	"net/http"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

// NewRouter returns the sentinel's router.
func NewRouter() *mux.Router {
	router := mux.NewRouter()
	router.Methods("GET").Path("/attack/{URL}").HandlerFunc(getAttack)
	return router
}

func getAttack(res http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	URL := vars["URL"]

	log.Infof("attacking: %s", URL)
}
