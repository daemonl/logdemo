package main

import (
	"context"
	"net/http"

	"github.com/anz-bank/pkg/log"
)

func main() {

	// Intention:
	// - Create some application wide fields
	// - Set up the logger once
	// - Include the logger / fields in context of HTTP Requests
	// - Log some basic request info in middleware

	logger := log.NewStandardLogger()

	ctx := context.Background()
	ctx = log.WithLogger(logger).With("application", "LogDemo").With("version", "1.2").Onto(ctx)
	log.From(ctx).Debug(ctx, "Hello")

	var handler http.Handler = http.HandlerFunc(helloHandler)
	handler = reqLogMiddleware(handler)
	handler = logContextMiddleware(logger)(handler)

	http.ListenAndServe(":8080", handler)
}

func logContextMiddleware(logger log.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
			ctx := req.Context()
			ctx = log.WithLogger(logger).Onto(ctx)
			req = req.WithContext(ctx)
			next.ServeHTTP(res, req)
		})
	}
}

func reqLogMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		log.With("path", req.URL.Path).From(req.Context()).Info("Request")
		next.ServeHTTP(res, req)
	})
}

func helloHandler(res http.ResponseWriter, req *http.Request) {
	log.With("some", "field").From(req.Context()).Debug("Doing a thing now")
	res.Write([]byte("OK\n"))
}
