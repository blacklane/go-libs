package camunda

import (
	"net/http"
)

type HttpClient interface {
	Do(r *http.Request) (*http.Response, error)
}

type BasicAuthCredentials struct {
	User     string
	Password string
}
