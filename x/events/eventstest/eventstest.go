package eventstest

//go:generate go run github.com/golang/mock/mockgen -source=../producer.go -destination=./mock_producer.go -package=eventstest Producer
