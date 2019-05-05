package main

import (
	logger "github.com/sirupsen/logrus"
	"nats-broker-example/pingpong"
)

func main() {
	natsReq := pingpong.NewRequest()
	err := natsReq.Init()
	if err != nil {
		logger.Errorf("Init err: %v", err)
		return
	}
	select {}
}
