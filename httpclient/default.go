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
	}
)

func newClient() *http.Client {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}
	http2.ConfigureTransport(tr)
	return &http.Client{
		Transport: tr,
	}
}

func newDefaultClient() IHttpClient {
	return &defaultClient{
		client: newClient(),
	}
}

func (c *defaultClient) GetClient() *http.Client {
	return c.client
}

func (c *defaultClient) Request(ctx context.Context, method, url string, body io.Reader, header http.Header) (*http.Response, error) {
	req, _ := http.NewRequest(method, url, body)
	req.Header = header
	*req = *req.WithContext(ctx)
	return c.client.Do(req)
}
