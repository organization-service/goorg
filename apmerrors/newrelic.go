package apmerrors

import (
	"context"

	newrelic "github.com/newrelic/go-agent"
)

func newrelicError(c context.Context, err error) {
	tx := newrelic.FromContext(c)
	tx.NoticeError(err)
}
