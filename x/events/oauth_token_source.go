package events

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
)

// type Token struct {
// 	AccessToken      string `json:"access_token"`
// 	ExpiresIn        int    `json:"expires_in"`
// 	RefreshExpiresIn int    `json:"refresh_expires_in"`
// 	RefreshToken     string `json:"refresh_token"`
// 	TokenType        string `json:"token_type"`
// 	NotBeforePolicy  int    `json:"not-before-policy"`
// 	SessionState     string `json:"session_state"`
// 	Scope            string `json:"scope"`
// }

type source struct {
	clientID     string
	clientSecret string
	url          string

	refreshBefore time.Duration
	httpClient    http.Client
	token         *oauth2.Token
	config        clientcredentials.Config
}

// NewTokenSource returns a simple oauth2.TokenSource implementation.
// It refreshes the token refreshBefore the token expiration.
func NewTokenSource(
	clientID,
	clientSecret,
	url string,
	refreshBefore time.Duration,
	httpClient http.Client) oauth2.TokenSource {

	return source{
		clientID:     clientID,
		clientSecret: clientSecret,
		url:          url,

		httpClient:    httpClient,
		refreshBefore: refreshBefore,
		config: clientcredentials.Config{
			ClientID:     clientID,
			ClientSecret: clientSecret,
			TokenURL:     url,
			AuthStyle:    oauth2.AuthStyleInParams,
		},
	}
}

func (t source) Token() (*oauth2.Token, error) {
	ctx := context.WithValue(
		context.TODO(), oauth2.HTTPClient, &t.httpClient)

	if t.token == nil || t.token.Expiry.Sub(time.Now()) < t.refreshBefore {
		token, err := t.config.Token(ctx)
		if err != nil {
			return nil, fmt.Errorf("could not retirve a new token: %w", err)
		}

		t.token = token
	}

	return t.token, nil
}
