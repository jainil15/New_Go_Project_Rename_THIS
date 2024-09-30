package main

import (
	"fmt"
	"kitchen/pkg/config"
	"kitchen/pkg/response"
	"log"
	"net/http"
	"os"
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
	router.HandleFunc("GET /favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		im, err := os.ReadFile("web/img/doomguy.ico")
		if err != nil {
			log.Println("Error ->", err)
			response.ErrorHandler(
				w,
				&response.Error{
					StatusCode: http.StatusNotFound,
					Error:      err,
					Message:    "Favicon not found",
				},
			)
			return
		}
		w.Header().Add("Content-type", "image/x-icon")
		w.WriteHeader(http.StatusOK)
		w.Write(im)
		return
	})
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		response.ErrorHandler(
			w,
			&response.Error{
				StatusCode: http.StatusNotFound,
				Message:    "Route not found",
			},
		)
		return
	})
	port := config.Env.Port
	server := http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: LoggingMiddleware(router),
	}
	log.Println("Server starting on port", port)
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Server Crashed")
	}
}
