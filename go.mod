module nats-broker-example

go 1.12

replace (
	github.com/golang/lint v0.0.0-20190313153728-d0100b6bd8b3 => github.com/golang/lint v0.0.0-20190409202823-5614ed5bae6fb75893070bdc0996a68765fdd275
	github.com/testcontainers/testcontainer-go => github.com/testcontainers/testcontainers-go v0.0.0-20190108154635-47c0da630f72
)

require (
	github.com/armon/circbuf v0.0.0-20190214190532-5111143e8da2 // indirect
	github.com/google/uuid v1.1.1
	github.com/mattn/go-colorable v0.1.1 // indirect
	github.com/micro/go-micro v1.1.0
	github.com/micro/go-plugins v1.1.0
	github.com/nats-io/go-nats v1.7.2
	github.com/satori/go.uuid v1.2.0
	github.com/sirupsen/logrus v1.4.1
)
