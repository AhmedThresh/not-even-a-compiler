package repl

import (
	"bufio"
	"fmt"
	"io"

	"github.com/AhmedThresh/not-even-a-compiler/pkg/eval"
	"github.com/AhmedThresh/not-even-a-compiler/pkg/lexer"
	"github.com/AhmedThresh/not-even-a-compiler/pkg/object"
	"github.com/AhmedThresh/not-even-a-compiler/pkg/parser"
)

const Prompt = ">>"

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)
	env := object.NewEnvironment()

	for {
		fmt.Print(Prompt)
		scanned := scanner.Scan()
		if !scanned {
			return
		}

		line := scanner.Text()
		lexer := lexer.NewLexer(line)

		parser := parser.NewParser(lexer)
		program := parser.ParseProgram()

		if len(parser.Errors()) != 0 {
			printParserErrors(out, parser.Errors())
			continue
		}

		evaluated := eval.Eval(program, env)
		if evaluated != nil {
			io.WriteString(out, evaluated.Inspect())
			io.WriteString(out, "\n")
		}
	}
}

func printParserErrors(out io.Writer, errors []string) {
	io.WriteString(out, "Woops! We ran into some monkey business here!\n")
	io.WriteString(out, " parser errors:\n")
	for _, msg := range errors {
		io.WriteString(out, "\t"+msg+"\n")
	}
}
