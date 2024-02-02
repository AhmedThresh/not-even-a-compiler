package repl

import (
	"bufio"
	"fmt"
	"io"

	"github.com/AhmedThresh/not-even-a-compiler/pkg/lexer"
	"github.com/AhmedThresh/not-even-a-compiler/pkg/token"
)

const Prompt = ">>"

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)

	for {
		fmt.Print(Prompt)
		scanned := scanner.Scan()
		if !scanned {
			return
		}

		line := scanner.Text()
		lexer := lexer.NewLexer(line)
		for t := lexer.NextToken(); t.Type != token.EOF; t = lexer.NextToken() {
			fmt.Printf("%+v\n", t)
		}
	}
}
