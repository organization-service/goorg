// This file was generated by goorg/cli/swagger at {{ .Timestamp }}

package handler

import (
	"bytes"
	"net/http"
	"strings"

	"gopkg.in/yaml.v2"
)

type (
	responseWriter struct {
		http.ResponseWriter
		buf bytes.Buffer
	}
)

func (w *responseWriter) write() (int, error) {
	return w.ResponseWriter.Write(w.buf.Bytes())
}

func (w *responseWriter) Write(buf []byte) (int, error) {
	return w.buf.Write(buf)
}

func (w *responseWriter) WriteHeader(statusCode int) {
	w.ResponseWriter.WriteHeader(statusCode)
}

func wrapResponseWriter(w http.ResponseWriter) http.ResponseWriter {
	rw := responseWriter{
		ResponseWriter: w,
		buf:            bytes.Buffer{},
	}

	return &rw
}

func SwaggerHandler(h http.HandlerFunc) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		accept := r.Header.Get("Accept")
		w := wrapResponseWriter(rw)
		h(w, r)
		ww := w.(*responseWriter)
		if strings.Contains(accept, "application/yaml") {
			mp := map[string]interface{}{}
			buf := ww.buf.Bytes()
			yaml.Unmarshal(buf, &mp)
			buf, _ = yaml.Marshal(mp)
			ww.ResponseWriter.Header().Set("Content-Type", "application/yaml")
			ww.ResponseWriter.Write(buf)
		} else {
			ww.write()
		}
	}
}
