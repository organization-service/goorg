package httpclient

import (
	"context"
	"io"
	"net/http"

	newrelic "github.com/newrelic/go-agent"
)

type (
	newrelicClient struct {
		*defaultClient
	}
)

type roundTripperFunc func(*http.Request) (*http.Response, error)

func (f roundTripperFunc) RoundTrip(r *http.Request) (*http.Response, error) { return f(r) }

func newNewrelicClient(ctx context.Context, method, url string, body io.Reader, header http.Header) IHttpClient {
	client := newClient()
	client.Transport = newrelicRoundTripper(client.Transport)
	req := getRequest(ctx, method, url, body, header)
	return &newrelicClient{
		defaultClient: &defaultClient{
			client: newClient(),
			req:    req,
		},
	}
}

func newrelicRoundTripper(original http.RoundTripper) http.RoundTripper {
	return roundTripperFunc(func(r *http.Request) (*http.Response, error) {
		txn := newrelic.FromContext(r.Context())
		segment := newrelic.StartExternalSegment(txn, r)
		if nil == original {
			original = http.DefaultTransport
		}
		response, err := original.RoundTrip(r)
		segment.Response = response
		segment.End()

		return response, err
	})
}
