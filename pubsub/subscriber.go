package pubsub

import (
	"log"
	"net"
	"reflect"

	"github.com/go-redis/redis"
)

// Subscriber ..
type Subscriber struct {
	pubsub   *redis.Client
	channel  string
	callback processFunc
}

type processFunc func(string, string)

// NewSubscriber ..
func NewSubscriber(channel string, fn processFunc) (*Subscriber, error) {
	var err error
	// TODO Timeout param?

	s := Subscriber{
		pubsub:   Service.client,
		channel:  channel,
		callback: fn,
	}

	// Subscribe to the channel
	pb, err := s.subscribe()
	if err != nil {
		return nil, err
	}
	// Listen for messages
	go s.listen(pb)

	return &s, nil
}

func (s *Subscriber) subscribe() (*redis.PubSub, error) {

	pb := s.pubsub.Subscribe(s.channel)
	return pb, nil
}

func (s *Subscriber) listen(pb *redis.PubSub) error {
	var channel string
	var payload string

	for {
		msg, err := pb.Receive()
		if err != nil {
			if reflect.TypeOf(err) == reflect.TypeOf(&net.OpError{}) && reflect.TypeOf(err.(*net.OpError).Err).String() == "*net.timeoutError" {
				// Timeout, ignore
				continue
			}
			// Actual error
			log.Print("Error in ReceiveTimeout()", err)
		}

		channel = ""
		payload = ""

		switch m := msg.(type) {
		case *redis.Subscription:
			log.Printf("Subscription Message: %v to channel '%v'. %v total subscriptions.", m.Kind, m.Channel, m.Count)
			continue
		case *redis.Message:
			channel = m.Channel
			payload = m.Payload

			// Process the message
			go s.callback(channel, payload)
		}
	}
}
