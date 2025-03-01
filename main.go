package main

import (
	"cool-compiler/codegen"
	"cool-compiler/lexer"
	"cool-compiler/parser"
	"cool-compiler/semant"
	"fmt"
	"log"
	"os"
	"os/exec"
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
		log.Fatal(err)
	}

	fmt.Printf("LLVM IR code generated successfully and written to %s\n", outputFilename)

	var cmd *exec.Cmd
	cmd = exec.Command("clang", "-Wno-deprecated", outputFilename, "-o", "output.exe", "-llegacy_stdio_definitions")

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		fmt.Printf("Warning: Compilation with clang failed. You may need to install LLVM/Clang.\n")
	} else {
		fmt.Println("Compilation successful. Executable: output.exe")

		// Run the compiled executable
		fmt.Println("Running the compiled program...")
		execCmd := exec.Command("./output.exe")
		execCmd.Stdout = os.Stdout
		execCmd.Stderr = os.Stderr

		if err := execCmd.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "Error running executable: %v\n", err)
			os.Exit(1)
		}
	}
}
