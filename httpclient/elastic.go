package httpclient

import (
	"go.elastic.co/apm/module/apmhttp"
)

type (
	elasticClient struct {
		*defaultClient
	}
)

func newElasticClient() IHttpClient {
	client := newClient()
	client = apmhttp.WrapClient(client)
	return &elasticClient{
		defaultClient: &defaultClient{
			client: client,
		},
	}
}
