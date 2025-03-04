# COOL Compiler

A compiler for the COOL (Classroom Object-Oriented Language) programming language that generates LLVM IR code and executables.

## Overview

This compiler translates COOL programs into LLVM IR, which is then compiled to native executables using Clang. The project implements a complete compilation pipeline:

1. **Lexical Analysis** - Converts source code into tokens
2. **Parsing** - Builds an Abstract Syntax Tree (AST) from tokens, using Pratt parsing algorithm to manage precedences
3. **Semantic Analysis** - Performs type checking and validates program semantics
4. **Code Generation** - Translates the AST to LLVM IR code
5. **LLVM Optimization** - Applies optimization passes to improve code quality

## About COOL

COOL (Classroom Object-Oriented Language) is a small language designed for teaching compiler construction. Key features include:

- Strong static typing
- Class-based object orientation with inheritance
- Built-in error handling
- Simple syntax with minimal keywords

Example COOL program (Hello World):
```
class Main {
   main(): Object {
      let message : String <- "Hello, COOL World!\n" in
         (new IO).out_string(message)
   };
};
```

## Prerequisites

- Go programming language
- LLVM/Clang for compiling the generated IR code
- LLVM Optimizer (`opt`) for applying optimization passes

## Installation

1. **Install Go**:
   - Download from [golang.org](https://golang.org/dl/)
   - Verify installation: `go version`

2. **Install LLVM/Clang**:
   - **Windows**: Download and install from the [LLVM website](https://releases.llvm.org/download.html)
   - **macOS**: `brew install llvm`
   - **Linux**: `sudo apt-get install llvm clang`

3. **Clone the repository**:
   ```
   git clone <repository-url>
   cd cool-compiler
   ```


## Usage

```
go run main.go <cool_file.cl>
```

This will:
1. Parse the COOL source file
2. Perform semantic analysis
3. Generate LLVM IR code (saved to `output.ll`)
4. Apply LLVM optimization passes (saved to `output_optimized.ll`)
5. Compile the optimized IR code to an executable (`optimized.exe`)

Then you can run the executable.

## Project Structure

- `lexer/` - Tokenizes COOL source code
- `parser/` - Parses tokens into an AST
- `ast/` - Defines the Abstract Syntax Tree structures
- `semant/` - Performs semantic analysis and type checking
- `codegen/` - Generates LLVM IR from the AST
- `cool_examples/` - Example COOL programs for testing
  - Includes samples for basic features, inheritance, recursion, etc.
- `OPTIMIZATION.md` - Details about the optimization strategies
- `cool-manual.pdf` - Reference manual for the COOL language

## Features

- Full support for COOL language features:
  - Classes and inheritance
  - Static and dynamic dispatch
  - Control structures (if, while, case)
  - Basic operations (+, -, *, /, <, <=, =)
  - Standard library functions
- Optimization:
  - Function inlining optimization
  - Constant propagation optimization
  - Performance and size metrics reporting

## Optimization Passes

The compiler applies the following LLVM optimization passes:

1. **Function Inlining** - Replaces function calls with the function body, eliminating call overhead
2. **Constant Propagation** - Identifies constant values and propagates them through the code

These optimizations can significantly improve code performance and reduce code size, particularly for complex programs. See `OPTIMIZATION.md` for more details about the optimization implementation and its impact.

## Performance Reporting

After compilation, the compiler automatically generates a performance report that includes:

- Size comparison between optimized and unoptimized IR code
- Execution time comparison between optimized and unoptimized executables
- Percentage improvements in both metrics

This helps in understanding the impact of optimization passes on your COOL programs.

## Examples

See the `cool_examples/` directory for sample COOL programs that demonstrate language features, including:

- `hello.cl` - Basic Hello World program
- `factorial.cl` - Recursive factorial calculation
- `fibonacci.cl` - Fibonacci sequence generation
- `inheritance.cl` - Class inheritance demonstration
- `case_*.cl` - Case expression usage examples
- `optimization_test.cl` - Program demonstrating optimization benefits

### Optimization Test

Run the optimization test program to see the impact of the optimization passes:

```
go run main.go cool_examples/optimization_test.cl
```

This program includes functions specifically designed to benefit from the implemented optimization passes and reports metrics on the optimization impact.

## Dependencies

The project relies on the following Go packages:
- `github.com/llir/llvm` - For LLVM IR generation

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

## Troubleshooting

- **Missing LLVM/Clang**: If you encounter errors related to LLVM tools not being found, ensure they are properly installed and available in your system PATH.
- **Go dependency issues**: If you encounter missing package errors, run `go mod tidy` to resolve dependencies.

## Future Improvements

- Additional optimization passes
- Support for COOL language extensions
- Improved error reporting with suggestions