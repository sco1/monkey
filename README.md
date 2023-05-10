# monkey
A Go interpreter and compiler for the Monkey esolang.

Monkey has a C-like syntax, supports variable bindings, prefix and infix operators, has first-class and higher-order functions, can handle closures with ease and has integers, booleans, arrays and hashes built-in. The language was created for Thorsten Ball's [*Writing An Interpreter In Go*](https://interpreterbook.com) and [*Writing A Compiler In Go*](https://compilerbook.com/) texts.

## Implementation Notes
### *Writing an Interpreter In Go*
* Macro support from [*The Lost Chapter: A Macro System For Monkey*](https://interpreterbook.com/lost) is implemented
* Source file execution is supported in the form e.g. `go run main.go ./hello_world.mnk`; if no arguments are passed then the REPL is invoked.
