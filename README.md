# nats-broker-example

使用 [go-micro/broker/nats](https://godoc.org/github.com/micro/go-micro/broker#Broker) 完成消息的发布与订阅

### 1. 安装并运行 NATS 服务:
``` sh
$ go get github.com/nats-io/gnatsd
$ gnatsd
```

### 2. 运行 producer
``` sh
$ go run producer/producer.go
```

### 3. 运行 consumer
``` sh
$ go run consumer/consumer.go
```
