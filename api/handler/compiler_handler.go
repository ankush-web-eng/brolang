package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/ankush-web-eng/brolang/evaluator"
	"github.com/ankush-web-eng/brolang/lexer"
	"github.com/ankush-web-eng/brolang/object"
	"github.com/ankush-web-eng/brolang/parser"
)

var GlobalEnv *object.Environment

// SetGlobalEnvironment sets the global environment.
func SetGlobalEnvironment(env *object.Environment) {
	GlobalEnv = env
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

	// Break the code into small parts and parse it
	l := lexer.New(req.Code)
	p := parser.New(l)
	program := p.ParseProgram()

	// if there are any errors while parsing the code, return the error
	var customErrors strings.Builder
	if len(p.Errors()) > 0 {
		fmt.Printf("Parser has %v error:\n", p.Errors())
		for _, value := range p.Errors() {
			customErrors.WriteString(value)
			customErrors.WriteString(" ")
		}
		response := CompileResponse{
			Error: customErrors.String(),
		}
		json.NewEncoder(w).Encode(response)
		return
	}

	// Initialize a global environment to hanydle variables
	env := object.NewEnvironment()

	// Get the evaluated code and return to the client
	result := evaluator.Eval(program, env)

	// env.OutputBuilder.WriteString(result.Inspect())

	response := CompileResponse{
		Result: env.OutputBuilder.String(), // Use accumulated output
	}

	if result.Type() == object.ERROR_OBJ {
		response.Error = result.Inspect()
	}

	json.NewEncoder(w).Encode(response)
}
