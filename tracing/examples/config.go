package examples

import (
	"encoding/json"
	"os"

	"github.com/blacklane/go-libs/logger"
	"github.com/caarlos0/env"
)

//revive:disable // Default and obvious config package.
type Config struct {
	ServiceName string

	KafkaServer  string `env:"KAFKA_BOOTSTRAP_SERVERS" envDefault:"127.0.0.1:9092"`
	KafkaGroupID string `env:"KAFKA_GROUP_ID" envDefault:"tracing-example"`
	Topic        string `env:"KAFKA_TOPIC" envDefault:"tracing-example"`
	TracerHost   string `env:"TRACER_HOST" envDefault:"localhost:55680"`

	Log logger.Logger `json:"-"`
}

func ParseConfig(serviceName string) Config {
	var cfg = Config{
		Log: logger.New(logger.ConsoleWriter{Out: os.Stdout}, serviceName)}

	if err := env.Parse(&cfg); err != nil {
		cfg.Log.Panic().Err(err).Msg("failed to load environment variables")
	}

	cfg.ServiceName = serviceName

	bs, _ := json.MarshalIndent(cfg, "", "  ")
	cfg.Log.Debug().Msgf("config: %s", bs)

	return cfg
}
