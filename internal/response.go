package internal

import (
	"encoding/json"
	"io"
	"net/http"
)

type ErrorBody struct {
	Message string `json:"message"`
}

func NewErrorBody(message string) *ErrorBody {
	return &ErrorBody{
		Message: message,
	}
}

func (e *ErrorBody) ToJson() []byte {
	if buf, err := json.Marshal(e); err != nil {
		return []byte(e.Message)
	} else {
		return buf
	}
}

func (e *ErrorBody) Write(w io.Writer) (int, error) {
	strJson := e.ToJson()
	return w.Write(strJson)
}

func OnError(w http.ResponseWriter, r *http.Request, err string) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(http.StatusUnauthorized)
	body := NewErrorBody(err)
	body.Write(w)
}
