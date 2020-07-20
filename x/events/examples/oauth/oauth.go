package oauth

type Cfg struct {
	ClientID     string `env:"OAUTH_CLIENT_ID,required"`
	ClientSecret string `env:"OAUTH_CLIENT_SECRET,required"`
	TokenURL     string `env:"OAUTH_TOKEN_ENDPOINT_URI,required"`
	KafkaServer  string `env:"KAFKA_BOOTSTRAP_SERVERS,required"`
	KafkaGroupID string `env:"KAFKA_GROUP_ID,required"`
	Topic        string `env:"KAFKA_TOPIC" envDefault:"test-go-libs-events"`
}
