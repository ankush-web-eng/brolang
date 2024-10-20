package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/ankush-web-eng/brolang/evaluator"
	"github.com/ankush-web-eng/brolang/lexer"
	"github.com/ankush-web-eng/brolang/parser"
)

func main() {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print(">> ")
		input, _ := reader.ReadString('\n')

		// Tokenize the input
		l := lexer.New(input)

		// Parse the tokens into an AST
		p := parser.New(l)
		ast := p.ParseExpression()

		// Evaluate the AST and print the result
		result := evaluator.Eval(ast)
		fmt.Println(result)
	}
}
