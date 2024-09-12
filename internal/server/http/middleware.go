package internalhttp

import (
	"net/http"
)

type statusWriter struct {
	http.ResponseWriter
	status int
}

func (handler *statusWriter) WriteHeader(code int) {
	handler.status = code
	handler.ResponseWriter.WriteHeader(code)
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// startTime := time.Now()
		httpWriter := statusWriter{w, 200}
		next.ServeHTTP(&httpWriter, r)
		/*log.Println(r.RemoteAddr, startTime.String(),
		r.Method, r.URL.Path, r.Proto, httpWriter.status,
		time.Since(startTime), r.UserAgent())*/
	})
}
