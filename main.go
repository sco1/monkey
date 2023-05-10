package main

import (
	"fmt"
	"io"
	"monkey/evaluator"
	"monkey/lexer"
	"monkey/object"
	"monkey/parser"
	"monkey/repl"
	"os"
	"os/user"
)

func main() {
	switch len(os.Args) {
	case 1:
		user, err := user.Current()
		if err != nil {
			panic(err)
		}
		fmt.Printf("Hello %s! This is the Monkey programming language!\n", user.Username)
		fmt.Printf("Feel free to type in commands\n")
		repl.Start(os.Stdin, os.Stdout)
	case 2:
		src, err := os.ReadFile(os.Args[1])
		if err != nil {
			msg := fmt.Sprintf("Could not parse file: %s", os.Args[1])
			panic(msg)
		}

		runSource(src, os.Stdout)

	default:
		fmt.Printf("Only one argument is supported. Got: %d", len(os.Args)-1)
	}
}

func runSource(src []byte, out io.Writer) {
	l := lexer.New(string(src))
	p := parser.New(l)

	program := p.ParseProgram()
	if len(p.Errors()) != 0 {
		printParserErrors(out, p.Errors())
		os.Exit(1)
	}

	env := object.NewEnvironment()
	macroEnv := object.NewEnvironment()

	evaluator.DefineMacros(program, macroEnv)
	expanded := evaluator.ExpandMacros(program, macroEnv)

	evaluated := evaluator.Eval(expanded, env)
	if evaluated != nil {
		io.WriteString(out, evaluated.Inspect())
	}
}

func printParserErrors(out io.Writer, errors []string) {
	io.WriteString(out, "Error(s) occured while parsing:\n")
	for _, msg := range errors {
		io.WriteString(out, "\t"+msg+"\n")
	}
}
