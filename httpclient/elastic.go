package httpclient

import (
	"context"
	"io"
	"net/http"

	"go.elastic.co/apm/module/apmhttp"
)

type (
	elasticClient struct {
		*defaultClient
	}
)

func newElasticClient(ctx context.Context, method, url string, body io.Reader, header http.Header) IHttpClient {
	client := newClient()
	client = apmhttp.WrapClient(client)
	req := getRequest(ctx, method, url, body, header)
	return &elasticClient{
		defaultClient: &defaultClient{
			client: client,
			req:    req,
		},
	}
}
