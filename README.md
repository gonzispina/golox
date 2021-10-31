# The Lox interpreter written in Go.

This is an implementation of the Lox Language [Crafting Interpreters](https://craftinginterpreters.com).

## Key differences

Most of the syntax is the same as proposed in the book except for:

- **No mandatory parenthesized expressions:** Just like Go, there's no need for parenthesized expressions in 'if' and 'for' statements
- **Go like 'for' statments:** The syntax of 'for' statements is the same as Go.
- **No 'while' statements:** I just don't like them.

## Extra features
* Multiline comments
* Enhanced error reporting
* `continue` statement and its corresponding error handling
* `break` statement and its corresponding error handling
* Uninitialized variable access is a runtime error
* Unused local variables and functions raises an error
* Lambda expressions
* Some other that I probably don't remember at the time of writing
