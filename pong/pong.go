package main

import (
	logger "github.com/sirupsen/logrus"
	"nats-broker-example/pingpong"
)

func main() {
	natsRep := pingpong.NewReply()
	err := natsRep.Init()
	if err != nil {
		logger.Errorf("Init err: %v", err)
		return
	}
	select {}
}
