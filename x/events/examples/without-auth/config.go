package config

//revive:disable:exported Obvious config struct.
type Config struct {
	KafkaServer  string `env:"KAFKA_BOOTSTRAP_SERVERS,required"`
	KafkaGroupID string `env:"KAFKA_GROUP_ID,required"`
	Topic        string `env:"KAFKA_TOPIC" envDefault:"test-go-libs-events"`
}
