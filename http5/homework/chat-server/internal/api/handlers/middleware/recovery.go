package middleware

import (
	"log"
	"net/http"
	"runtime/debug"

	"github.com/ew0s/ewos-to-go-hw/pkg/httputils/baseresponse"
)

type RecoveredError struct {
	message string
}

func (e RecoveredError) Error() string {
	return e.message
}

func Recovery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				debug.PrintStack()

				errMsg, ok := err.(string)
				if !ok {
					errMsg = "not string type panic"
				}

				re := RecoveredError{errMsg}

				log.Printf("Recovered from panic: %v\n", re)

				baseresponse.RenderErr(w, r, http.StatusInternalServerError, re)
			}
		}()

		next.ServeHTTP(w, r)
	})
}
