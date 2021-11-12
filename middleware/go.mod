module github.com/blacklane/go-libs/middleware

go 1.17

require (
	github.com/blacklane/go-libs/logger v0.5.2-0.20211110123002-5f1d885ee0d0
	github.com/blacklane/go-libs/otel v0.0.0-20211110143904-be9482540221
	github.com/blacklane/go-libs/tracking v0.2.1
	github.com/blacklane/go-libs/x/events v0.1.1-0.20211110143534-f824e5cfdd98
	github.com/google/go-cmp v0.5.6
	github.com/google/uuid v1.3.0
)

require (
	go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp v0.26.1 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc v1.1.0 // indirect
	go.opentelemetry.io/otel/exporters/stdout/stdouttrace v1.1.0 // indirect
	go.opentelemetry.io/proto/otlp v0.10.0 // indirect
	golang.org/x/sys v0.0.0-20211109184856-51b60fd695b3 // indirect
	golang.org/x/text v0.3.7 // indirect
	google.golang.org/genproto v0.0.0-20211104193956-4c6863e31247 // indirect
	google.golang.org/grpc v1.42.0 // indirect
)
