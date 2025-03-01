# COOL Compiler

A compiler for the COOL (Classroom Object-Oriented Language) programming language that generates LLVM IR code and executables.

## Overview

This compiler translates COOL programs into LLVM IR, which is then compiled to native executables using Clang. The project implements a complete compilation pipeline:

1. **Lexical Analysis** - Converts source code into tokens
2. **Parsing** - Builds an Abstract Syntax Tree (AST) from tokens
3. **Semantic Analysis** - Performs type checking and validates program semantics
4. **Code Generation** - Translates the AST to LLVM IR code

## Prerequisites

- Go programming language
- LLVM/Clang for compiling the generated IR code

## Usage

```
go run main.go <cool_file.cl>
```

This will:
1. Parse the COOL source file
2. Perform semantic analysis
3. Generate LLVM IR code (saved to `output.ll`)
4. Compile the IR code to an executable (`output.exe`)
5. Automatically run the executable

## Project Structure

- `lexer/` - Tokenizes COOL source code
- `parser/` - Parses tokens into an AST
- `semant/` - Performs semantic analysis and type checking
- `codegen/` - Generates LLVM IR from the AST
- `cool_examples/` - Example COOL programs for testing

## Features

- Full support for COOL language features:
  - Classes and inheritance
  - Static and dynamic dispatch
  - Control structures (if, while, case)
  - Basic operations (+, -, *, /, <, <=, =)
  - Standard library functions

## Examples

See the `cool_examples/` directory for sample COOL programs that demonstrate language features.

# Running Tests

### Lexer tests

To run lexer tests, move to the `lexer` directory, and run:
- `go test`

### Parser tests

To run parser tests, move to the `parser` directory, and run:
- `go test`

### Semantic Analyzer tests

To run semantic analyzer tests, move to the `semant` directory, and run:
- `go test`

### Codegen tests

To run semantic analyzer tests, move to the `codegen` directory, and run:
- `go test`

