package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/ankush-web-eng/brolang/lexer"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print(">> ")
		input, _ := reader.ReadString('\n')
		l := lexer.New(input)
		for tok := l.NextToken(); tok.Type != lexer.EOF; tok = l.NextToken() {
			fmt.Printf("%+v\n", tok)
		}
	}
}
