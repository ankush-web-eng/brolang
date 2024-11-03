package handler

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"testing"
)

func TestCompilerHandler(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			`bhai_sun x = 5;
            bol_bhai(x);`,
			"5\n",
		},
		{
			`bhai_sun arr = [1, 2, 3];
            bol_bhai(arr[0]);`,
			"1\n",
		},
	}

	for _, tt := range tests {
		req := CompileRequest{
			Code: tt.input,
		}
		reqBody, _ := json.Marshal(req)

		w := httptest.NewRecorder()
		r := httptest.NewRequest("POST", "/compile", bytes.NewBuffer(reqBody))
		r.Header.Set("Content-Type", "application/json")

		CompilerHandler(w, r)

		var resp CompileResponse
		json.NewDecoder(w.Body).Decode(&resp)

		if resp.Result != tt.expected {
			t.Errorf("expected=%q, got=%q", tt.expected, resp.Result)
		}
	}
}
