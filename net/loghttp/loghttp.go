package loghttp

import (
	"fmt"
	"net/http"
)

type statusWriter struct {
	http.ResponseWriter
	status int
}

func (w *statusWriter) WriteHeader(status int) {
	w.status = status
	w.ResponseWriter.WriteHeader(status)
}

// Log common feild of HTTP like nginx access log.
// Example: http.ListenAndServe(":80", LogHandler(http.DefaultServeMux))
func LogHandler(handler http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ww := &statusWriter{ResponseWriter: w}
		handler.ServeHTTP(ww, r)
		status := ww.status
		if status == 0 {
			status = 200
		}
		fmt.Println(r.RemoteAddr, r.Method, r.URL.Path, status,
			r.Header.Get("User-Agent"))
	}
}
