package main

import (
	"net/http"

	"github.com/ankush-web-eng/brolang/api/handler"
	"github.com/ankush-web-eng/brolang/object"
)

func main() {
	env := object.NewEnvironment()
	handler.SetGlobalEnvironment(env)

	http.HandleFunc("/compile", corsMiddleware(handler.CompilerHandler))
	http.ListenAndServe(":8080", nil)
}

func corsMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	}
}
