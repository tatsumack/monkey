package repl

import (
	"bufio"
	"fmt"
	"io"
	"monkey/evaluator"
	"monkey/lexer"
	"monkey/parser"
)

const PROMPT = ">> "

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)

	for {
		fmt.Printf(PROMPT)
		scanned := scanner.Scan()
		if !scanned {
			return
		}

		line := scanner.Text()
		if line == "q" || line == "quit" {
			fmt.Printf("quit repl\n")
			break
		}

		l := lexer.New(line)
		p := parser.New(l)

		program := p.ParseProgram()
		if len(p.Errors()) > 0 {
			printParserErrors(out, p.Errors())
			continue
		}

		evaluated := evaluator.Eval(program)
		if evaluated != nil {
			io.WriteString(out, evaluated.Inspect())
			io.WriteString(out, "\n")
		}
	}
}

func printParserErrors(out io.Writer, errors []string) {
	io.WriteString(out, "parse errors:\n")
	for _, msg := range errors {
		io.WriteString(out, fmt.Sprintf("\t%s\n", msg))
	}
}
