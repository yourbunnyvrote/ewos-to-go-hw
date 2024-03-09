package middleware

import (
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/middleware"
)

func Logger() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
			rw := &ResponseWriter{WrapResponseWriter: ww}

			defer func() {
				status := rw.Status()

				var levelStr string

				switch {
				case status >= http.StatusInternalServerError:
					levelStr = "SERVER ERROR"
				case status >= http.StatusBadRequest:
					levelStr = "CLIENT ERROR"
				case status >= http.StatusMultipleChoices:
					levelStr = "REDIRECTION"
				case status >= http.StatusOK:
					levelStr = "SUCCESS"
				case status >= http.StatusContinue:
					levelStr = "INFO"
				default:
					levelStr = "UNKNOWN"
				}

				log.Printf("[%s] %s - %s %s - %d: %s\n", levelStr, time.Now().Format(time.RFC3339), r.Method, r.URL.Path, status, rw.message)
			}()

			next.ServeHTTP(rw, r)
		})
	}
}

type ResponseWriter struct {
	middleware.WrapResponseWriter
	message string
}

func (rw *ResponseWriter) Write(p []byte) (int, error) {
	rw.message = string(p)
	return rw.WrapResponseWriter.Write(p)
}
