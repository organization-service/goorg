package httpclient

import (
	"context"
	"io"
	"net/http"

	"github.com/organization-service/goorg/internal"
)

type (
	IHttpClient interface {
		GetClient() *http.Client
		Request(ctx context.Context, method, url string, body io.Reader, header http.Header) (res *http.Response, err error)
	}
)

func NewClient() IHttpClient {
	switch internal.GetApmName() {
	case internal.Elastic:
		return newElasticClient()
	case internal.Newrelic:
		return newNewrelicClient()
	default:
		return newDefaultClient()
	}
}
