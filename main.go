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
		origin := r.Header.Get("Origin")
		if origin == "http://localhost:3000" || origin == "https://brolang.ankushsingh.tech" {
			w.Header().Set("Access-Control-Allow-Origin", origin)
		}
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	}
}

// package main

// import (
// 	"log"
// 	"net/http"

// 	"github.com/ankush-web-eng/brolang/api/handler"
// 	"github.com/ankush-web-eng/brolang/object"
// )

// func main() {
// 	env := object.NewEnvironment()
// 	handler.SetGlobalEnvironment(env)

// 	// Initialize Redis client for Pub/Sub
// 	handler.InitializeRedisClient()

// 	http.HandleFunc("/compile", corsMiddleware(handler.CompilerHandler))

// 	log.Println("Server is starting on port 8080...")
// 	if err := http.ListenAndServe(":8080", nil); err != nil {
// 		log.Fatalf("Failed to start server: %v", err)
// 	}
// }

// func corsMiddleware(next http.HandlerFunc) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		origin := r.Header.Get("Origin")
// 		if origin == "http://localhost:3000" || origin == "https://brolang.ankushsingh.tech" {
// 			w.Header().Set("Access-Control-Allow-Origin", origin)
// 		}
// 		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
// 		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

// 		if r.Method == "OPTIONS" {
// 			w.WriteHeader(http.StatusOK)
// 			return
// 		}

// 		next.ServeHTTP(w, r)
// 	}
// }
