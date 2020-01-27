package main

// import (
// 	"flag"
// 	"io/ioutil"
// 	"net/http"
// 	"os"
// 	"os/signal"
// 	"syscall"

// 	log "github.com/sirupsen/logrus"
// )

// // Temporary constant till we solve the error:
// // socket: too many open files
// // possible reference:
// // http://craigwickesser.com/2015/01/golang-http-to-many-open-files/
// // https://mtyurt.net/post/docker-how-to-increase-number-of-open-files-limit.html
// const maxRequests = 100

// func main() {
// 	// Parse URL from command line.
// 	var URL string
// 	flag.StringVar(&URL, "url", "", "target url")
// 	// Parse number of requests from command line.
// 	var count int
// 	flag.IntVar(&count, "count", 10, "number of requests")

// 	flag.Parse()

// 	if URL == "" {
// 		log.Fatal("target required")
// 	}

// 	log.Info("firing laser at: ", URL)

// 	var client http.Client

// 	// Set up quit channel.
// 	quit := make(chan bool)

// 	// Fire a goroutine that waits for OS signals. Upon receiving signal, quit channel will receive.
// 	// This will allow the program to terminate.
// 	go func(quit chan bool) {
// 		signals := make(chan os.Signal)
// 		defer close(signals)

// 		signal.Notify(signals, syscall.SIGQUIT, syscall.SIGTERM, os.Interrupt)
// 		defer signal.Stop(signals)

// 		signal := <-signals
// 		log.Infof("signal received: %s \n", signal.String())
// 		quit <- true
// 	}(quit)

// 	// Make requests until we receive a stop signal.
// 	conns := 0
// 	for {
// 		select {
// 		case <-quit:
// 			log.Info("quitting by command")
// 			return
// 		default:
// 			if conns < maxRequests {
// 				conns++
// 				go makeRequest(client, URL)
// 			}
// 		}
// 	}
// }

// func makeRequest(c http.Client, URL string) {
// 	// Create the request to the target.
// 	req, err := http.NewRequest("GET", URL, nil)
// 	if err != nil {
// 		log.Fatalf("failed to create request: %s", err)
// 	}

// 	// Make the request.
// 	res, err := c.Do(req)
// 	if err != nil {
// 		log.Fatalf("request failed: %s", err)
// 	}
// 	defer res.Body.Close()

// 	// Decode the body for now.
// 	body, err := ioutil.ReadAll(res.Body)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	// Log body.
// 	log.Info("response code:", res.StatusCode)
// 	log.Info(string(body))
// }
