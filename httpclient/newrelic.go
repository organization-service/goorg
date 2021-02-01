package httpclient

import (
	"github.com/newrelic/go-agent/v3/newrelic"
)

type (
	newrelicClient struct {
		*defaultClient
	}
)

func newNewrelicClient() IHttpClient {
	client := newClient()
	client.Transport = newrelic.NewRoundTripper(client.Transport)
	return &newrelicClient{
		defaultClient: &defaultClient{
			client: client,
		},
	}
}
