package pingpong

import (
	"encoding/json"
	"fmt"
	"github.com/micro/go-micro/broker"
	"github.com/micro/go-plugins/broker/nats"
	gonats "github.com/nats-io/go-nats"
	"github.com/satori/go.uuid"
	logger "github.com/sirupsen/logrus"
	"time"
)

type Reply struct {
	AllowReconnect      bool
	MaxReconnect        int
	ConnectTime         time.Duration
	ReconnectTimeout    time.Duration
	PingIntervalTimeout time.Duration
	UUID                string
	Name                string
	Broker              broker.Broker
	NATSUrl             string
}

func NewReply() *Reply {
	return &Reply{
		AllowReconnect:      true,
		MaxReconnect:        10,
		ConnectTime:         3 * time.Second,
		ReconnectTimeout:    5 * time.Second,
		PingIntervalTimeout: 10 * time.Second,
		NATSUrl:             "127.0.0.1:4222",
		UUID:                uuid.NewV4().String(),
		Name:                "pong",
	}
}

// CBGetPorts function types
type CBGetPorts func(input string) (string, error)

func (c *Reply) Init() error {

	option := gonats.Options{
		AllowReconnect: c.AllowReconnect,
		MaxReconnect:   c.MaxReconnect,
		Timeout:        c.ConnectTime,
		ReconnectWait:  c.ReconnectTimeout,
		PingInterval:   c.PingIntervalTimeout,
		DisconnectedCB: func(nc *gonats.Conn) {
			logger.Warn("Disconnect")
		},
		ReconnectedCB: func(nc *gonats.Conn) {
			logger.Warn("Reconnect")
		},
	}
	c.Broker = nats.NewBroker(broker.Addrs(c.NATSUrl), nats.Options(option))
	if err := c.Broker.Init(); err != nil {
		logger.Errorf("Broker Init error: %v", err)
	}
	if err := c.Broker.Connect(); err != nil {
		logger.Errorf("Broker Connect error: %v", err)
	}

	if err := c.PongHandler(); err != nil {
		logger.Errorf("PingPong Alive Agent error: %v", err)
	}
	return nil
}

func (c *Reply) PongHandler() error {

	_, err := c.Broker.Subscribe("request_message", func(p broker.Publication) error {
		logger.Debugf("message: %v", p.Message().Body)

		request := Message{
			UUID: uuid.NewV4().String(),
			From: "pong",
			To:   "ping",
			Time: time.Now().String(),
		}
		reqData, err := json.Marshal(request)
		if err != nil {
			fmt.Println(err.Error())
			return err
		}

		err = c.Broker.Publish("receive_message", &broker.Message{Body: reqData})
		return err
	})
	if err != nil {
		return err
	}

	_, err = c.Broker.Subscribe("request_alive_info", func(p broker.Publication) error {
		fmt.Printf("message: %s\n", p.Message().Body)
		reply := AliveInfoReply{}
		reply.UUID = c.UUID
		reply.Name = c.Name

		data, err := json.Marshal(reply)
		if err != nil {
			fmt.Println(err.Error())
			return err
		}
		err = c.Broker.Publish("reply_alive_info", &broker.Message{Body: data})
		return err
	})
	return err
}
