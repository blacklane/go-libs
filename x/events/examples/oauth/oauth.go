package oauth

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
)

type Cfg struct {
	ClientID     string `env:"OAUTH_CLIENT_ID,required"`
	ClientSecret string `env:"OAUTH_CLIENT_SECRET,required"`
	TokenURL     string `env:"OAUTH_TOKEN_ENDPOINT_URI,required"`
	KafkaServer  string `env:"KAFKA_BOOTSTRAP_SERVERS,required"`
	KafkaGroupID string `env:"KAFKA_GROUP_ID,required"`
	Topic        string `env:"KAFKA_TOPIC" envDefault:"test-go-libs-events"`
}

type source struct {
	ClientID     string
	ClientSecret string
	URL          string

	refreshTTL time.Duration
	httpClient http.Client
	token      *oauth2.Token
	Config     clientcredentials.Config
}

func NewTokenSource(clientID,
	clientSecret,
	url string,
	refreshTTL time.Duration,
	httpClient http.Client) oauth2.TokenSource {

	return source{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		URL:          url,

		httpClient: httpClient,
		refreshTTL: refreshTTL,
		Config: clientcredentials.Config{
			ClientID:     clientID,
			ClientSecret: clientSecret,
			TokenURL:     url,
			AuthStyle:    oauth2.AuthStyleInParams,
		},
	}
}

func (t source) Token() (*oauth2.Token, error) {
	ctx := context.WithValue(context.Background(), oauth2.HTTPClient, &t.httpClient)

	if t.token == nil || t.token.Expiry.Sub(time.Now()) < t.refreshTTL {
		token, err := t.Config.Token(ctx)
		if err != nil {
			return nil, fmt.Errorf("could not retirve a new token: %w", err)
		}

		t.token = token
	}

	return t.token, nil
}

type Token struct {
	AccessToken      string `json:"access_token"`
	ExpiresIn        int    `json:"expires_in"`
	RefreshExpiresIn int    `json:"refresh_expires_in"`
	RefreshToken     string `json:"refresh_token"`
	TokenType        string `json:"token_type"`
	NotBeforePolicy  int    `json:"not-before-policy"`
	SessionState     string `json:"session_state"`
	Scope            string `json:"scope"`
}
