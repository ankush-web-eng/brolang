package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/ankush-web-eng/brolang/lexer"
	"github.com/ankush-web-eng/brolang/parser"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print(">> ")
		input, _ := reader.ReadString('\n')
		l := lexer.New(input)
		p := parser.New(l)

		ast := p.ParseExpression()
		fmt.Printf("AST: %+v\n", ast)
	}
}
