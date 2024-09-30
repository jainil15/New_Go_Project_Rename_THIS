package main

import (
	"kitchen/pkg/response"
	"log"
	"net/http"
	"time"
)

type wrappedResponseWriter struct {
	http.ResponseWriter
	StatusCode int
}

func (w *wrappedResponseWriter) WriteHeader(code int) {
	w.StatusCode = code
	w.ResponseWriter.WriteHeader(code)
}

func LoggingMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		ww := &wrappedResponseWriter{w, http.StatusOK}
		h.ServeHTTP(ww, r)
		log.Printf("%v %v %v %v", ww.StatusCode, r.Method, r.RequestURI, time.Since(start))
	})
}

func main() {
	router := http.NewServeMux()
	router.HandleFunc(
		"GET /health",
		func(w http.ResponseWriter, r *http.Request) {
			response.ResponseHandler(w, &response.Response{Message: "Server Running"})
			return
		},
	)
	server := http.Server{
		Addr:    ":8080",
		Handler: LoggingMiddleware(router),
	}
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Server Crashed")
	}
}
