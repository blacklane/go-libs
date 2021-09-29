package camunda

import "net/http"

// This is sould be removed. There is no need for that, if you need to receive
// an http.Client, just take it in. Besides it's only used internally. No need ot export that.
// But please remove it.
type HttpClient interface {
	Do(r *http.Request) (*http.Response, error)
}

type BasicAuthCredentials struct {
	User     string
	Password string
}
