package main

import (
	"log"
	"net/http"

	"github.com/ankush-web-eng/brolang/api/handler"
)

func main() {
	http.HandleFunc("/compile", handler.CompilerHandler)

	log.Println("Server starting on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
