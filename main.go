package main

import (
	"flag"
	"fmt"
	"io"
	"monkey/compiler"
	"monkey/evaluator"
	"monkey/lexer"
	"monkey/object"
	"monkey/parser"
	"monkey/repl"
	"monkey/vm"
	"os"
)

func main() {
	var noVM bool
	flag.BoolVar(&noVM, "no-vm", false, "Execute using the interpreter path only")
	flag.Parse()

	switch len(flag.Args()) {
	case 0:
		repl.Start(os.Stdin, os.Stdout, noVM)
	case 1:
		filePath := flag.Args()[0]
		src, err := os.ReadFile(filePath)
		if err != nil {
			msg := fmt.Sprintf("Could not parse file: %s", filePath)
			panic(msg)
		}
		runSource(src, os.Stdout, noVM)
	default:
		fmt.Printf("Only one filepath is supported. Got: %d", len(flag.Args()))
	}
}

func runSource(src []byte, out io.Writer, noVM bool) {
	l := lexer.New(string(src))
	p := parser.New(l)

	program := p.ParseProgram()
	if len(p.Errors()) != 0 {
		printParserErrors(out, p.Errors())
		os.Exit(1)
	}

	if noVM {
		env := object.NewEnvironment()
		macroEnv := object.NewEnvironment()

		evaluator.DefineMacros(program, macroEnv)
		expanded := evaluator.ExpandMacros(program, macroEnv)

		evaluated := evaluator.Eval(expanded, env)
		if evaluated != nil {
			io.WriteString(out, evaluated.Inspect())
		}
	} else {
		comp := compiler.New()
		err := comp.Compile(program)
		if err != nil {
			fmt.Fprintf(out, "Compilation failed:\n %s\n", err)
		}

		machine := vm.New(comp.Bytecode())
		err = machine.Run()
		if err != nil {
			fmt.Fprintf(out, "Executing bytecode failed:\n %s\n", err)
		}

		stackTop := machine.StackTop()
		io.WriteString(out, stackTop.Inspect())
		io.WriteString(out, "\n")
	}
}

func printParserErrors(out io.Writer, errors []string) {
	io.WriteString(out, "Error(s) occured while parsing:\n")
	for _, msg := range errors {
		io.WriteString(out, "\t"+msg+"\n")
	}
}
