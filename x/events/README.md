events
-----------
An event abstraction library with a Kafka implementation.

You can find examples on [examples](examples) folder.

## Usage

### Consumer

```go
	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"group.id":           "consumer-example" + strconv.Itoa(int(time.Now().Unix())),
		"bootstrap.servers":  "localhost:9092",
		"session.timeout.ms": 6000,
		"auto.offset.reset":  "earliest",
	})
	if err != nil {
		log.Panicf("could not create kafka consumer: %v", err)
	}

	if err := consumer.Subscribe(topic, nil); err != nil {
		log.Panicf("failed to subscribe to Kafka topic %s: %v", topic, err)
	}

	c := events.NewKafkaConsumer(consumer, events.HandlerFunc(
		func(ctx context.Context, e events.Event) error {
			log.Printf("consumed event: %s", e.Payload)
			return nil
		}))

	c.Run(time.Second)
```

### Producer

```go
	errHandler := events.ErrorHandlerFunc(func(event events.Event, err error) {
		log.Panicf("failed to deliver the event %s: %v", string(event.Payload), err)
	})

	p, err := events.NewKafkaProducer(
		&kafka.ConfigMap{
			"bootstrap.servers":  "localhost:9092",
			"message.timeout.ms": "1000"},
		errHandler)
	if err != nil {
		log.Panicf("%v", err)
	}

	// handle failed deliveries
	_ = p.HandleMessages()
	defer p.Shutdown(context.Background())

	e := events.Event{Payload: []byte("Hello, Gophers")}
	err = p.Send(
		e,
		"events-example-topic")
	if err != nil {
		log.Panicf("error sending the event %s: %v", e, err)
	}
```

## Development

We provide a [`docker-compose`](docker-compose.yml) to spin up all the needed
dependencies and run the tests as well as `.env_local` with the needed 
environment variables. 

To spin up a kafka cluster and connect from you local machine, run:
```shell script
make compose-kafka
```
 - kafka is reachable on `localhost:9092`
 - zookeeper is reachable on `localhost:2181`

you might run `make clean` to stop and remove the containers created by docker-compose  

## Tests

 - Pass the build tag `integration` to `go test` to run integrations tests
 - Use the `-short` to skip slow tests. As some tests need to publish and 
 consume messages they sleep for a few seconds.
 - use `.env_compose` if running the tests through dockercompose

### Docker compose
```shell script
make compose-tests
```

### Local
```shell script
make test
```


## IDEs

As build tags are used to separate integration tests from unit tests make sure
to set up your IDE to include the files with the `integration` build flag 

### GoLand

To set it up to always use custom build flags:

- Go to `Preferences > Go > Build Tags & Vendoring` and fill in `Custom tags`

![Build Tags & Vendoring](doc/imgs/goBuildTagsConfig.png)

- Then on `Run/Debug Configurations` set `Templates > Go Test` 
to _Use all custom build tags_

![Build Tags & Vendoring](doc/imgs/goBuildTagsConfig.png)

![Open Run/Debug Configurations](doc/imgs/openRunDebugConfig.png)

![Set Run/Debug Configurations](doc/imgs/runDebugConfig.png)
