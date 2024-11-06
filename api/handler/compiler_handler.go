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

	if req.Code == "" {
		response := CompileResponse{
			Error: "Kuchh likh to sahi be!",
		}
		json.NewEncoder(w).Encode(response)
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

// package handler

// import (
// 	"context"
// 	"encoding/json"
// 	"log"
// 	"net/http"
// 	"strings"

// 	"github.com/ankush-web-eng/brolang/evaluator"
// 	"github.com/ankush-web-eng/brolang/lexer"
// 	"github.com/ankush-web-eng/brolang/object"
// 	"github.com/ankush-web-eng/brolang/parser"
// 	"github.com/go-redis/redis/v8"
// )

// var GlobalEnv *object.Environment
// var redisClient *redis.Client
// var ctx = context.Background()

// // SetGlobalEnvironment initializes the global environment for evaluation
// func SetGlobalEnvironment(env *object.Environment) {
// 	GlobalEnv = env
// }

// // InitializeRedisClient sets up the Redis client connection
// func InitializeRedisClient() {
// 	redisClient = redis.NewClient(&redis.Options{
// 		Addr: "localhost:6379",
// 	})
// 	if err := redisClient.Ping(ctx).Err(); err != nil {
// 		log.Fatalf("Failed to connect to Redis: %v", err)
// 	}
// }

// type CompileRequest struct {
// 	Code      string `json:"code"`
// 	RequestID string `json:"requestId"`
// }

// type CompileResponse struct {
// 	Result string `json:"result,omitempty"`
// 	Error  string `json:"error,omitempty"`
// }

// type HTTPError struct {
// 	Status  int
// 	Message string
// }

// func (e *HTTPError) Error() string {
// 	return e.Message
// }

// func validateRequest(r *http.Request) (*CompileRequest, error) {
// 	if r.Method != http.MethodPost {
// 		return nil, &HTTPError{
// 			Status:  http.StatusMethodNotAllowed,
// 			Message: "Method not allowed",
// 		}
// 	}

// 	var req CompileRequest
// 	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
// 		return nil, &HTTPError{
// 			Status:  http.StatusBadRequest,
// 			Message: "Invalid request body",
// 		}
// 	}

// 	if req.Code == "" {
// 		return nil, &HTTPError{
// 			Status:  http.StatusBadRequest,
// 			Message: "Kuchh likh to sahi be!",
// 		}
// 	}

// 	if req.RequestID == "" {
// 		return nil, &HTTPError{
// 			Status:  http.StatusBadRequest,
// 			Message: "Missing requestId",
// 		}
// 	}

// 	return &req, nil
// }

// // compileCode compiles the given code and publishes the result to Redis
// func compileCode(code, requestId string, env *object.Environment) {
// 	l := lexer.New(code)
// 	p := parser.New(l)
// 	program := p.ParseProgram()

// 	if len(p.Errors()) > 0 {
// 		var errorMsg strings.Builder
// 		for _, msg := range p.Errors() {
// 			errorMsg.WriteString(msg + " ")
// 		}

// 		// Publish error to Redis channel identified by requestId
// 		publishToRedis(requestId, CompileResponse{Error: errorMsg.String()})
// 		return
// 	}

// 	result := evaluator.Eval(program, env)
// 	if result != nil && result.Type() == object.ERROR_OBJ {
// 		// Publish runtime error to Redis
// 		publishToRedis(requestId, CompileResponse{Error: result.Inspect()})
// 		return
// 	}

// 	// Publish successful result to Redis
// 	publishToRedis(requestId, CompileResponse{Result: env.OutputBuilder.String()})
// }

// // publishToRedis sends the response to the Redis Pub/Sub channel
// func publishToRedis(channel string, response CompileResponse) {
// 	data, err := json.Marshal(response)
// 	if err != nil {
// 		log.Printf("Failed to marshal response: %v", err)
// 		return
// 	}

// 	if err := redisClient.Publish(ctx, channel, data).Err(); err != nil {
// 		log.Printf("Failed to publish to Redis: %v", err)
// 	}
// }

// // CompilerHandler is the HTTP handler for the compile endpoint
// func CompilerHandler(w http.ResponseWriter, r *http.Request) {
// 	w.Header().Set("Content-Type", "application/json")

// 	req, err := validateRequest(r)
// 	if err != nil {
// 		if httpErr, ok := err.(*HTTPError); ok {
// 			http.Error(w, httpErr.Message, httpErr.Status)
// 			return
// 		}
// 		http.Error(w, "Internal server error", http.StatusInternalServerError)
// 		return
// 	}

// 	// Call compileCode in a separate goroutine to allow asynchronous processing
// 	go compileCode(req.Code, req.RequestID, GlobalEnv)

// 	// Respond immediately to the HTTP request indicating successful submission
// 	w.WriteHeader(http.StatusAccepted)
// 	_ = json.NewEncoder(w).Encode(map[string]string{"status": "Code submitted successfully"})
// }
