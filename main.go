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
	"runtime"
	"strings"
	"time"
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

	// Apply LLVM optimization passes
	optimizedOutputFilename, sizeBeforeOpt, sizeAfterOpt, execTimeBeforeOpt, execTimeAfterOpt := applyOptimizationPasses(outputFilename)

	// Compile the optimized LLVM IR to an executable
	var cmd *exec.Cmd
	outputName := "output.exe"
	if runtime.GOOS == "windows" {
		cmd = exec.Command("clang", "-Wno-deprecated", optimizedOutputFilename, "-o", outputName, "-llegacy_stdio_definitions")
	} else {
		cmd = exec.Command("clang", "-Wno-deprecated", optimizedOutputFilename, "-o", outputName)
	}

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		fmt.Printf("Warning: Compilation with clang failed. You may need to install LLVM/Clang.\n")
	} else {
		fmt.Printf("Compilation successful. Executable: %s\n", outputName)
	}

	// Report optimization metrics
	printOptimizationReport(sizeBeforeOpt, sizeAfterOpt, execTimeBeforeOpt, execTimeAfterOpt)
}

// applyOptimizationPasses applies LLVM optimization passes to the IR code
// and returns the name of the optimized file along with metrics
func applyOptimizationPasses(inputFile string) (string, int64, int64, time.Duration, time.Duration) {
	// Create unoptimized executable for baseline comparison
	unoptimizedExe := "unoptimized.exe"
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("clang", "-Wno-deprecated", inputFile, "-o", unoptimizedExe, "-llegacy_stdio_definitions")
	} else {
		cmd = exec.Command("clang", "-Wno-deprecated", inputFile, "-o", unoptimizedExe)
	}
	err := cmd.Run()
	if err != nil {
		fmt.Printf("Warning: Failed to create unoptimized executable for comparison: %v\n", err)
	}

	// Get file size before optimization
	unoptimizedInfo, err := os.Stat(inputFile)
	var sizeBeforeOpt int64
	if err != nil {
		fmt.Printf("Warning: Could not get file size: %v\n", err)
	} else {
		sizeBeforeOpt = unoptimizedInfo.Size()
	}

	// Measure execution time of unoptimized code (if possible)
	var execTimeBeforeOpt time.Duration
	if _, err := os.Stat(unoptimizedExe); err == nil {
		startTime := time.Now()
		execCmd := exec.Command("./" + unoptimizedExe)
		execCmd.Run()
		execTimeBeforeOpt = time.Since(startTime)
	}

	// Output filename for the optimized IR
	optimizedOutputFilename := "output_optimized.ll"

	// Apply function inlining optimization
	fmt.Println("Applying function inlining optimization pass...")
	inliningCmd := exec.Command("opt", "-passes=inline", "-inline-threshold=1000", inputFile, "-S", "-o", "output_inline.ll")

	inliningCmd.Stdout = os.Stdout
	inliningCmd.Stderr = os.Stderr
	err = inliningCmd.Run()
	if err != nil {
		fmt.Printf("Warning: Function inlining optimization failed: %v\n", err)
		fmt.Println("Continuing with unoptimized code...")
		optimizedOutputFilename = inputFile // Fallback to unoptimized file
	} else {
		// Apply constant propagation optimization on top of inlining
		fmt.Println("Applying constant propagation optimization pass...")
		constPropCmd := exec.Command("opt", "-passes=ipsccp", "output_inline.ll", "-S", "-o", optimizedOutputFilename)
		constPropCmd.Stdout = os.Stdout
		constPropCmd.Stderr = os.Stderr
		err = constPropCmd.Run()
		if err != nil {
			fmt.Printf("Warning: Constant propagation optimization failed: %v\n", err)
			fmt.Println("Continuing with only function inlining optimization...")
			optimizedOutputFilename = "output_inline.ll" // Fallback to inlined file
		}
	}

	// Get file size after optimization
	optimizedInfo, err := os.Stat(optimizedOutputFilename)
	var sizeAfterOpt int64
	if err != nil {
		fmt.Printf("Warning: Could not get optimized file size: %v\n", err)
	} else {
		sizeAfterOpt = optimizedInfo.Size()
	}

	// Create optimized executable
	optimizedExe := "optimized.exe"
	if runtime.GOOS == "windows" {
		cmd = exec.Command("clang", "-Wno-deprecated", optimizedOutputFilename, "-o", optimizedExe, "-llegacy_stdio_definitions")
	} else {
		cmd = exec.Command("clang", "-Wno-deprecated", optimizedOutputFilename, "-o", optimizedExe)
	}
	err = cmd.Run()

	// Measure execution time of optimized code (if possible)
	var execTimeAfterOpt time.Duration
	if err == nil {
		startTime := time.Now()
		execCmd := exec.Command("./" + optimizedExe)
		execCmd.Run()
		execTimeAfterOpt = time.Since(startTime)
	}

	return optimizedOutputFilename, sizeBeforeOpt, sizeAfterOpt, execTimeBeforeOpt, execTimeAfterOpt
}

// printOptimizationReport prints a report of the optimization metrics
func printOptimizationReport(sizeBeforeOpt, sizeAfterOpt int64, execTimeBeforeOpt, execTimeAfterOpt time.Duration) {
	fmt.Println("\n=== LLVM Optimization Report ===")
	fmt.Println("Optimization Passes Applied:")
	fmt.Println("1. Function Inlining (-passes=inline)")
	fmt.Println("   - Inlines function calls to reduce call overhead")
	fmt.Println("   - Enables more opportunities for other optimizations")
	fmt.Println("2. Constant Propagation (-passes=ipsccp)")
	fmt.Println("   - Identifies and propagates constant values through the program")
	fmt.Println("   - Eliminates computations that can be performed at compile time")

	if sizeBeforeOpt > 0 && sizeAfterOpt > 0 {
		sizeDiff := sizeBeforeOpt - sizeAfterOpt
		sizePercent := float64(sizeDiff) / float64(sizeBeforeOpt) * 100
		fmt.Printf("\nIR Code Size Impact:\n")
		fmt.Printf("- Before optimization: %d bytes\n", sizeBeforeOpt)
		fmt.Printf("- After optimization: %d bytes\n", sizeAfterOpt)
		if sizeDiff > 0 {
			fmt.Printf("- Reduction: %d bytes (%.2f%%)\n", sizeDiff, sizePercent)
		} else {
			fmt.Printf("- Increase: %d bytes (%.2f%%)\n", -sizeDiff, -sizePercent)
		}
	}

	if execTimeBeforeOpt > 0 && execTimeAfterOpt > 0 {
		timeDiff := execTimeBeforeOpt - execTimeAfterOpt
		timePercent := float64(timeDiff) / float64(execTimeBeforeOpt) * 100
		fmt.Printf("\nExecution Time Impact:\n")
		fmt.Printf("- Before optimization: %v\n", execTimeBeforeOpt)
		fmt.Printf("- After optimization: %v\n", execTimeAfterOpt)
		if timeDiff > 0 {
			fmt.Printf("- Speedup: %v (%.2f%%)\n", timeDiff, timePercent)
		} else {
			fmt.Printf("- Slowdown: %v (%.2f%%)\n", -timeDiff, -timePercent)
		}
	}

	fmt.Println("\nNote: Performance impact may vary depending on the specific COOL program.")
	fmt.Println("For small test programs, the impact might be minimal.")
	fmt.Println("=================================")
}
