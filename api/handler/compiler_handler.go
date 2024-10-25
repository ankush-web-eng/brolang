package handler

import (
	"encoding/json"
	"net/http"

	"github.com/ankush-web-eng/brolang/evaluator"
	"github.com/ankush-web-eng/brolang/lexer"
	"github.com/ankush-web-eng/brolang/parser"
)

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
		response := CompileResponse{
			Error: "bhadwe kya coder banega tu, semicolon bhool gaya!!!",
		}
		json.NewEncoder(w).Encode(response)
		return
	}

	result := evaluator.Eval(program)
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
