package router

import (
	"io/ioutil"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/michaeloverton/ddos-laser/cmd/target/task"
	"github.com/michaeloverton/ddos-laser/internal/server"
	log "github.com/sirupsen/logrus"
)

// NewRouter returns the targets's router.
func NewRouter() *mux.Router {
	router := mux.NewRouter()
	router.Methods("GET").Path("/thrash").HandlerFunc(thrash)
	return router
}

// type targetResponse struct {
// 	Duration float64 `json:"duration"`
// }

func thrash(res http.ResponseWriter, req *http.Request) {
	log.Info("responding to request from: ", req.Host)

	start := time.Now()
	negativeInfinity := task.Reverse(task.PositiveInfinity)
	elapsed := time.Since(start)
	log.Info("negating infinity took: ", elapsed)

	output := []byte(negativeInfinity)
	err := ioutil.WriteFile("/tmp/negative", output, 0644)
	if err != nil {
		panic(err)
	}

	// resp := targetResponse{
	// 	Duration: elapsed.Seconds(),
	// }
	// time.Sleep(time.Second * 5)

	server.Serve(res, nil, 200)
}
