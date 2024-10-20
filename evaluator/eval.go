package evaluator

import (
	"fmt"
	"strconv"

	"github.com/ankush-web-eng/brolang/parser"
)

func Eval(node parser.ASTNode) int {
	switch exp := node.(type) {
	case *parser.NumberLiteral:
		// Convert the number literal to an integer
		val, err := strconv.Atoi(exp.Value)
		if err != nil {
			fmt.Println("Error converting value to integer:", err)
			return 0
		}
		return val

	case *parser.BinaryExpression:
		left := Eval(exp.Left)
		right := Eval(exp.Right)

		// Process the operator (currently only supports addition)
		switch exp.Operator.Literal {
		case "+":
			return left + right
		default:
			fmt.Println("Unsupported operator:", exp.Operator.Literal)
			return 0
		}

	default:
		fmt.Println("Unknown expression type")
		return 0
	}
}
