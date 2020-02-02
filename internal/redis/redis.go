package redis

import (
	"github.com/go-redis/redis"
)

// Redis wraps the Redis client.
type Redis struct {
	Client *redis.Client
}

// New creates a new Redis client.
func New(addr string) (*Redis, error) {
	// Set up client.
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	// Test connection to client.
	_, err := client.Ping().Result()
	if err != nil {
		return nil, err
	}

	// Return client.
	r := Redis{
		Client: client,
	}
	return &r, nil
}

func (c Redis) Publish(topic, URL string) error {
	err := c.Client.Publish(topic, URL).Err()
	if err != nil {
		return err
	}

	return nil
}

// func (c Redis) Receive() (interface{}, error) {
// 	sub := c.Client.Subscribe(c.LaserTopic)
// 	msg, err := sub.ReceiveMessage()
// 	// test reception of topic
// 	// msg, err := sub.Receive()
// 	if err != nil {
// 		return nil, err
// 	}

// 	return msg, nil
// }
