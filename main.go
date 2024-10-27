package main

import (
	"net/http"

	"github.com/ankush-web-eng/brolang/api/handler"
	"github.com/ankush-web-eng/brolang/object"
)

func main() {
	env := object.NewEnvironment()
	handler.SetGlobalEnvironment(env)
	http.HandleFunc("/compile", handler.CompilerHandler)
	http.ListenAndServe(":8080", nil)
}
