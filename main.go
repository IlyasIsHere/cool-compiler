package main

import (
	"cool-compiler/codegen"
	"cool-compiler/lexer"
	"cool-compiler/parser"
	"cool-compiler/semant"
	"fmt"
	"os"
	"strings"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s <cool_file.cl>\n", os.Args[0])
		os.Exit(1)
	}

	input, err := os.ReadFile(os.Args[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file: %v\n", err)
		os.Exit(1)
	}

	l := lexer.NewLexer(strings.NewReader(string(input)))
	p := parser.New(l)
	program := p.ParseProgram()

	if len(p.Errors()) != 0 {
		fmt.Fprintf(os.Stderr, "Parser Errors:\n")
		for _, msg := range p.Errors() {
			fmt.Fprintf(os.Stderr, "%s\n", msg)
		}
		os.Exit(1)
	}

	sa := semant.NewSemanticAnalyser()
	sa.Analyze(program)

	if len(sa.Errors()) != 0 {
		fmt.Fprintf(os.Stderr, "Semantic Errors:\n")
		for _, msg := range sa.Errors() {
			fmt.Fprintf(os.Stderr, "%s\n", msg)
		}
		os.Exit(1)
	}

	// Generate LLVM IR for the program
	module, err := codegen.Generate(program)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Code generation error: %v\n", err)
		os.Exit(1)
	}

	// Output LLVM IR to a file
	outputFilename := "output.ll"
	err = os.WriteFile(outputFilename, []byte(module.String()), 0644)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error writing output file: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("LLVM IR code generated successfully and written to %s\n", outputFilename)
	fmt.Println("To compile with clang, install LLVM and Clang, then run:")
	fmt.Printf("clang %s runtime.c -o output\n", outputFilename)
}
