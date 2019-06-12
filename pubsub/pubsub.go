package pubsub

import (
	"encoding/json"

	"github.com/go-redis/redis"
)

// PubSub .
type PubSub struct {
	client *redis.Client
}

// Service ..
var Service *PubSub

func init() {
	var client *redis.Client
	client = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
		PoolSize: 10,
	})
	Service = &PubSub{client}
}

// PublishString ..
func (ps *PubSub) PublishString(channel, message string) *redis.IntCmd {
	return ps.client.Publish(channel, message)
}

// Publish ..
func (ps *PubSub) Publish(channel string, message interface{}) *redis.IntCmd {
	jsonBytes, err := json.Marshal(message)
	if err != nil {
		panic(err)
	}
	messageString := string(jsonBytes)
	return ps.client.Publish(channel, messageString)
}
