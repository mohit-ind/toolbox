package middlewares

// The logging middleware is heavily based on this article
// https://blog.questionable.services/article/guide-logging-middleware-go/

import (
	"fmt"
	"net/http"
	"time"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	logger "github.com/toolboxlogger"
)

// responseWriter is a minimal wrapper for http.ResponseWriter that allows the
// HTTP status code to be captured for logging.
type responseWriter struct {
	http.ResponseWriter
	status      int
	wroteHeader bool
}

// wrapResponseWriter wraps a basic ResponseWriter with the responseWriter wrapper struct
func wrapResponseWriter(w http.ResponseWriter) *responseWriter {
	return &responseWriter{
		ResponseWriter: w,
		status:         200,
	}
}

// WriteHeader registers the supplied HTTP status code, and passes it to the underlying ResponseWriter
// status code only can be written once, any further call will have no affect
func (rw *responseWriter) WriteHeader(code int) {
	if rw.wroteHeader {
		return
	}

	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
	rw.wroteHeader = true
}

// LoggingMiddleware logs the incoming HTTP request & its duration onto the supplied logger when the response is sent.
// this middleware also recovers from any panic in the middleware chain, and turns it into an error log entry
func LoggingMiddleware(logger *logger.Logger) func(http.Handler) http.Handler {

	return func(next http.Handler) http.Handler {

		fn := func(w http.ResponseWriter, r *http.Request) {
			// register the time when the request arrived
			start := time.Now()
			// defer the recovery from any panic on the middleware chain
			defer func() {
				if err := recover(); err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					logger.WithError(errors.New(fmt.Sprintf("%+v", err))).Error("Recovered from panic")
				}
			}()
			// wrap the original ResponseWriter
			wrapped := wrapResponseWriter(w)
			// call next middleware
			next.ServeHTTP(wrapped, r)
			// when the call to next returns, we log out the request and the elapsed time
			fields := logrus.Fields{
				"status":   wrapped.status,
				"method":   r.Method,
				"path":     r.URL.EscapedPath(),
				"duration": time.Since(start),
			}
			// depending on status log info or warn
			if wrapped.status >= 200 && wrapped.status <= 399 {
				logger.WithFields(fields).Info("Request")
			} else {
				logger.WithFields(fields).Warn("Request")
			}
		}

		return http.HandlerFunc(fn)
	}
}
