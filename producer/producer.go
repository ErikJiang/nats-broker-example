package main

import (
	"fmt"
	"github.com/micro/go-plugins/broker/nats"
	"log"
	"time"

	"github.com/micro/go-micro/broker"
	"github.com/micro/go-micro/cmd"
)

var (
	topic = "go.micro.topic.foo"
	natsBroker broker.Broker
)

func pub() {
	tick := time.NewTicker(time.Second)
	i := 0
	for t := range tick.C {
		msg := &broker.Message{
			Header: map[string]string{
				"id": fmt.Sprintf("%d", i),
			},
			Body: []byte(fmt.Sprintf("%d: %s", i, t.String())),
		}
		if err := natsBroker.Publish(topic, msg); err != nil {
			log.Printf("[pub] failed: %v", err)
		} else {
			fmt.Println("[pub] pubbed message:", string(msg.Body))
		}
		i++
	}
}

func main() {
	_ = cmd.Init()
	natsBroker = nats.NewBroker(broker.Addrs("127.0.0.1:4222"))
	if err := natsBroker.Init(); err != nil {
		log.Fatalf("Broker Init error: %v", err)
	}

	if err := natsBroker.Connect(); err != nil {
		log.Fatalf("Broker Connect error: %v", err)
	}

	pub()
}
