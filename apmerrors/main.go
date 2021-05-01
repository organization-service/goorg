package apmerrors

import (
	"context"

	"github.com/organization-service/goorg/v2/internal"
)

func SendError(c context.Context, err error) {
	switch internal.GetApmName() {
	case internal.Elastic:
		elasticError(c, err)
	case internal.Newrelic:
		newrelicError(c, err)
	}
}
