package main

import (
	"fmt"
	"github.com/micro/go-plugins/broker/nats"
	"log"

	"github.com/micro/go-micro/broker"
	"github.com/micro/go-micro/cmd"
)

var (
	topic = "go.micro.topic.foo"
	natsBroker broker.Broker
)

// Example of a shared subscription which receives a subset of messages
func sharedSub() {
	_, err := natsBroker.Subscribe(topic, func(p broker.Publication) error {
		fmt.Println("[sub] received message:", string(p.Message().Body), "header", p.Message().Header)
		return nil
	}, broker.Queue("consumer"))
	if err != nil {
		fmt.Println(err)
	}
}

// Example of a subscription which receives all the messages
func sub() {
	_, err := natsBroker.Subscribe(topic, func(p broker.Publication) error {
		fmt.Println("[sub] received message:", string(p.Message().Body), "header", p.Message().Header)
		return nil
	})
	if err != nil {
		fmt.Println(err)
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

	sub()
	// sharedSub()
	select {}
}
