package evaluator

import (
	"testing"

	helper "github.com/ankush-web-eng/brolang/helpers"
	"github.com/ankush-web-eng/brolang/lexer"
	"github.com/ankush-web-eng/brolang/object"
	"github.com/ankush-web-eng/brolang/parser"
)

func testEval(input string) object.Object {
	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()
	env := object.NewEnvironment()

	return Eval(program, env)
}

func TestEvalIntegerExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"5", 5},
		{"10", 10},
		{"5 + 5", 10},
		{"5 - 3", 2},
		{"2 * 2", 4},
		{"4 / 2", 2},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		helper.TestIntegerObject(t, evaluated, tt.expected)
	}
}
