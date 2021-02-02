package apmerrors

import (
	"context"

	"go.elastic.co/apm"
)

func elasticError(c context.Context, err error) {
	e := apm.CaptureError(c, err)
	e.Send()
}
