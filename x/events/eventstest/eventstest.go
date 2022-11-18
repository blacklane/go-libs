package eventstest

//go:generate go run -mod=mod github.com/golang/mock/mockgen -source=../producer.go -destination=./mock_producer.go -package=eventstest Producer
