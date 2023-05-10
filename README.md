# monkey
A Go interpreter and compiler for the Monkey esolang.

Monkey has a C-like syntax, supports variable bindings, prefix and infix operators, has first-class and higher-order functions, can handle closures with ease and has integers, booleans, arrays and hashes built-in. The language was created for Thorsten Ball's [*Writing An Interpreter In Go*](https://interpreterbook.com) and [*Writing A Compiler In Go*](https://compilerbook.com/) texts.

## Implementation Notes
### Source Files vs. REPL
Support is provided for executing source files as well as launching the REPL:
  * `go run main.go` will open the REPL
  * `go run main.go ./hello_world.mnk` will execute the specified source file

By default, execution is done using the VM created from the compiler book. The `--no-vm` flag can be used to instead execute using the interpreter created from the interpreter book.

### *Writing an Interpreter In Go*
* Macro support from [*The Lost Chapter: A Macro System For Monkey*](https://interpreterbook.com/lost) is implemented
