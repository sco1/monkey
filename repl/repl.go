package repl

import (
	"bufio"
	"fmt"
	"io"
	"monkey/compiler"
	"monkey/evaluator"
	"monkey/lexer"
	"monkey/object"
	"monkey/parser"
	"monkey/vm"
)

const PROMPT = ">> "

func Start(in io.Reader, out io.Writer, noVM bool) {
	fmt.Printf("Welcome to the Monkey programming language! ")
	if noVM {
		fmt.Printf("(Interpreter Mode)\n")
	} else {
		fmt.Printf("(VM Mode)\n")
	}
	fmt.Printf("Feel free to type in commands\n")

	scanner := bufio.NewScanner(in)

	// Global states, not sure how to do this in a conditional
	env := object.NewEnvironment()
	macroEnv := object.NewEnvironment()

	constants := []object.Object{}
	globals := make([]object.Object, vm.GlobalsSize)
	symbolTable := compiler.NewSymbolTable()

	for {
		fmt.Fprint(out, PROMPT)
		scanned := scanner.Scan()
		if !scanned {
			return
		}

		line := scanner.Text()
		l := lexer.New(line)
		p := parser.New(l)

		program := p.ParseProgram()
		if len(p.Errors()) != 0 {
			printParserErrors(out, p.Errors())
			continue
		}

		if noVM {
			evaluator.DefineMacros(program, macroEnv)
			expanded := evaluator.ExpandMacros(program, macroEnv)

			evaluated := evaluator.Eval(expanded, env)
			if evaluated != nil {
				io.WriteString(out, evaluated.Inspect())
				io.WriteString(out, "\n")
			}
		} else {
			comp := compiler.NewWithState(symbolTable, constants)
			err := comp.Compile(program)
			if err != nil {
				fmt.Fprintf(out, "Compilation failed:\n %s\n", err)
				continue
			}

			code := comp.Bytecode()
			constants = code.Constants
			machine := vm.NewWithGlobalsStore(code, globals)
			err = machine.Run()
			if err != nil {
				fmt.Fprintf(out, "Executing bytecode failed:\n %s\n", err)
				continue
			}

			lastPopped := machine.LastPoppedStackElem()
			io.WriteString(out, lastPopped.Inspect())
			io.WriteString(out, "\n")
		}
	}
}

const MONKEY_FACE = `            __,__
   .--.  .-"     "-.  .--.
  / .. \/  .-. .-.  \/ .. \
 | |  '|  /   Y   \  |'  | |
 | \   \  \ 0 | 0 /  /   / |
  \ '- ,\.-"""""""-./, -' /
   ''-' /_   ^ ^   _\ '-''
       |  \._   _./  |
       \   \ '~' /   /
        '._ '-=-' _.'
           '-----'
`

func printParserErrors(out io.Writer, errors []string) {
	io.WriteString(out, MONKEY_FACE)
	io.WriteString(out, "Woops! We ran into some monkey business here!\n")
	io.WriteString(out, " parser errors:\n")
	for _, msg := range errors {
		io.WriteString(out, "\t"+msg+"\n")
	}
}
