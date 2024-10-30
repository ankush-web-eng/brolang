package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/ankush-web-eng/brolang/evaluator"
	"github.com/ankush-web-eng/brolang/lexer"
	"github.com/ankush-web-eng/brolang/object"
	"github.com/ankush-web-eng/brolang/parser"
)

var globalEnv *object.Environment

// SetGlobalEnvironment sets the global environment.
func SetGlobalEnvironment(env *object.Environment) {
	globalEnv = env
}

type CompileRequest struct {
	Code string `json:"code"`
}

type CompileResponse struct {
	Result string `json:"result"`
	Error  string `json:"error,omitempty"`
}

func CompilerHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req CompileRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	l := lexer.New(req.Code)
	p := parser.New(l)
	program := p.ParseProgram()

	if len(p.Errors()) > 0 {
		fmt.Printf("Parser has %v error:\n", p.Errors())
		response := CompileResponse{
			Error: "Bhaap ko bhej, tere bas ki nahi hai",
		}
		json.NewEncoder(w).Encode(response)
		return
	}

	// Initialize a global environment
	env := object.NewEnvironment()

	result := evaluator.Eval(program, env)
	if result.Type() == "ERROR" {
		response := CompileResponse{
			Error: result.Inspect(),
		}
		json.NewEncoder(w).Encode(response)
		return
	}

	response := CompileResponse{
		Result: result.Inspect(),
	}
	json.NewEncoder(w).Encode(response)
}
