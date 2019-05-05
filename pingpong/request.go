package pingpong

import (
	"encoding/json"
	"fmt"
	"github.com/micro/go-micro/broker"
	"github.com/micro/go-plugins/broker/nats"
	gonats "github.com/nats-io/go-nats"
	"github.com/satori/go.uuid"
	logger "github.com/sirupsen/logrus"
	"sync"
	"time"
)

type Metadata struct {
	UUID string
	Name string
	Time time.Time
}

type Request struct {
	AllowReconnect      bool
	MaxReconnect        int
	ConnectTime         time.Duration
	ReconnectTimeout    time.Duration
	PingIntervalTimeout time.Duration
	Broker              broker.Broker
	NATSUrl             string
	Mutex               *sync.RWMutex
	AliveData           map[string]*Metadata
}

func NewRequest() *Request {
	return &Request{
		AllowReconnect:      true,
		MaxReconnect:        10,
		ConnectTime:         3 * time.Second,
		ReconnectTimeout:    5 * time.Second,
		PingIntervalTimeout: 10 * time.Second,
		NATSUrl:             "127.0.0.1:4222",
		Mutex:               &sync.RWMutex{},
		AliveData:           make(map[string]*Metadata),
	}
}

func (c *Request) Init() error {

	option := gonats.Options{
		Timeout:        c.ConnectTime,
		AllowReconnect: c.AllowReconnect,
		MaxReconnect:   c.MaxReconnect,
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

	if err := c.PingHandler(); err != nil {
		logger.Errorf("PingPong Alive Agent error: %v", err)
	}
	return nil
}

func (c *Request) PingHandler() error {
	_, err := c.Broker.Subscribe("reply_alive_info", func(p broker.Publication) error {
		fmt.Println("[sub] received message:", string(p.Message().Body), "header", p.Message().Header)

		reply := AliveInfoReply{}
		err := json.Unmarshal(p.Message().Body, &reply)
		if err != nil {
			logger.Errorf(fmt.Sprintf("reply_alive_agent json error: %s", err.Error()))
			return err
		}

		meta := Metadata{}
		meta.UUID = reply.UUID
		meta.Name = reply.Name
		meta.Time = time.Now()

		c.Mutex.Lock()
		defer c.Mutex.Unlock()
		c.AliveData[reply.UUID] = &meta
		return nil
	})
	if err != nil {
		logger.Errorf("PingHandler error: %v", err)
		return err
	}

	// 2秒发送一次通知所有的在线主机.
	go func() {
		count := 1
		tick := time.NewTicker(2 * time.Second)
		for range tick.C {
			msg := &broker.Message{
				Body: []byte(fmt.Sprintf("%d", count)),
			}
			logger.Infof("msg: %s", msg.Body)
			err := c.Broker.Publish("request_alive_info", msg)
			if err != nil {
				logger.Errorf(fmt.Sprintf("nats ping_pong: %s", err.Error()))
				// return
			}
			count++
		}
	}()
	return nil
}

func (c *Request) Test() error {
	// 请求
	request := Message{
		UUID: uuid.NewV4().String(),
		From: "ping",
		To:   "pong",
		Time: time.Now().String(),
	}
	reqData, err := json.Marshal(request)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	err = c.Broker.Publish("request_message", &broker.Message{Body: reqData})
	if err != nil {
		return err
	}

	// 应答
	sub, err := c.Broker.Subscribe("receive_message", func(p broker.Publication) error {
		reply := Message{}
		err := json.Unmarshal(p.Message().Body, &reply)
		if err != nil {
			logger.Errorf("receive_ports() json error: %v", err)
			return err
		}

		logger.Debug("message: %v", reply)
		return nil
	})
	if err != nil {
		return err
	}
	defer sub.Unsubscribe()

	return nil
}
