package dhl

import (
	"net/http"
)

type myRoundTripper struct {
	roundTripper http.RoundTripper
	apiKey       string
}

func (mrt myRoundTripper) RoundTrip(r *http.Request) (*http.Response, error) {
	r.Header.Add("DHL-API-Key", mrt.apiKey)
	return mrt.roundTripper.RoundTrip(r)
}

// NewClient is creating a http.Client with authentication
func NewClient(apiKey string) http.Client {
	return http.Client{
		Transport: myRoundTripper{
			roundTripper: http.DefaultTransport,
			apiKey:       apiKey,
		},
	}
}
