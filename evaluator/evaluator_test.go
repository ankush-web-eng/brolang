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

func TestLoopControlStatements(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		// Test break in while loop
		{
			`
            bhai_sun x = 0;
            jaha_tak (x < 5) {
                bol_bhai(x);
                agar (x == 2) {
                    bas_kar_bhai;
                }
                x = x + 1;
            }
            `,
			"0\n1\n2\n",
		},
		// Test continue in while loop
		{
			`
            bhai_sun x = 0;
            jaha_tak (x < 3) {
                x = x + 1;
                agar (x == 2) {
                    aage_bhad;
                }
                bol_bhai(x);
            }
            `,
			"1\n3\n",
		},
		// Test break in for loop
		{
			`
            chal_bhai (bhai_sun i = 0; i < 5; i = i + 1) {
                bol_bhai(i);
                agar (i == 2) {
                    bas_kar_bhai;
                }
            }
            `,
			"0\n1\n2\n",
		},
		// Test continue in for loop
		{
			`
            chal_bhai (bhai_sun i = 0; i < 3; i = i + 1) {
                agar (i == 1) {
                    aage_bhad;
                }
                bol_bhai(i);
            }
            `,
			"0\n2\n",
		},
		// Test infinite loop detection
		{
			`
            bhai_sun x = 0;
            jaha_tak (sach) {
                x = x + 1;
            }
            `,
			"Infinite loop detected! Check your loop condition.",
		},
	}

	for i, tt := range tests {
		env := object.NewEnvironment()
		l := lexer.New(tt.input)
		p := parser.New(l)
		program := p.ParseProgram()

		evaluated := Eval(program, env)

		if evaluated != nil && evaluated.Type() == object.ERROR_OBJ {
			if evaluated.(*object.Error).Message != tt.expected {
				t.Errorf("test %d - wrong error message. expected=%q, got=%q",
					i, tt.expected, evaluated.(*object.Error).Message)
			}
			continue
		}

		actual := env.OutputBuilder.String()
		if actual != tt.expected {
			t.Errorf("test %d - wrong output. expected=%q, got=%q",
				i, tt.expected, actual)
		}
	}
}
