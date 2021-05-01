package httpclient

import (
	"context"
	"io"
	"net/http"

	"github.com/organization-service/goorg/v2/internal"
)

type (
	IHttpClient interface {
		GetClient() *http.Client
		GetRequest() *http.Request
		Do() (res *http.Response, err error)
	}
)

func NewClient(ctx context.Context, method, url string, body io.Reader, header http.Header) IHttpClient {
	switch internal.GetApmName() {
	case internal.Elastic:
		return newElasticClient(ctx, method, url, body, header)
	case internal.Newrelic:
		return newNewrelicClient(ctx, method, url, body, header)
	default:
		return newDefaultClient(ctx, method, url, body, header)
	}
}
