package httpclient

import (
	"context"
	"crypto/tls"
	"io"
	"net/http"

	"golang.org/x/net/http2"
)

type (
	defaultClient struct {
		client *http.Client
		req    *http.Request
	}
)

func newClient() *http.Client {
	tlsConfig := &tls.Config{
		InsecureSkipVerify: true,
	}
	client := http.DefaultClient
	client.Transport = http.DefaultTransport
	if v, ok := client.Transport.(*http.Transport); ok {
		http2.ConfigureTransport(v)
		v.TLSClientConfig = tlsConfig
	}
	return client
}

func getRequest(ctx context.Context, method, url string, body io.Reader, header http.Header) *http.Request {
	req, _ := http.NewRequestWithContext(ctx, method, url, body)
	req.Header = make(http.Header, len(header))
	for k, s := range header {
		req.Header[k] = append([]string(nil), s...)
	}
	return req
}

func newDefaultClient(ctx context.Context, method, url string, body io.Reader, header http.Header) IHttpClient {
	return &defaultClient{
		client: newClient(),
		req:    getRequest(ctx, method, url, body, header),
	}
}

func (c *defaultClient) GetClient() *http.Client {
	return c.client
}

func (c *defaultClient) GetRequest() *http.Request {
	return c.req
}

func (c *defaultClient) Do() (*http.Response, error) {
	return c.client.Do(c.req)
}
