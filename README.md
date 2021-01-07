# edat-stan - Streaming NATS for edat

## Installation

    go get -u github.com/stackus/edat-stan

## Usage Example

    import "github.com/stack/edat-stan"

    conn, _ := stan.Connect(clusterID, clientID, options)

    // Create a consumer and use it in a message subscriber
    consumer := edatstan.NewConsumer(conn, groupID)
    subscriber := msg.NewSubscriber(consumer)

    // Create a producer and use it in a message publisher
    producer := edatstan.NewProducer(conn)
    publisher := msg.NewPublisher(producer)


## Prerequisites

Go 1.15

## Features

- Message Consumer `NewConsumer(stan.Conn, groupID, ...options)`
- Message Producer `NewProducer(stan.Conn, ...options)`

## TODOs

- Documentation
- Tests, tests, and more tests

## Contributing
Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.

Please make sure to update tests as appropriate.

## License

MIT
