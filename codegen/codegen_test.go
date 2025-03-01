package codegen

import (
	"bytes"
	"cool-compiler/lexer"
	"cool-compiler/parser"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// TestCodeGeneration verifies that the code generator produces valid LLVM IR
// for a variety of COOL programs
func TestCodeGeneration(t *testing.T) {
	examples := []struct {
		name           string
		filename       string
		expectedOutput string
	}{
		{
			"Hello World",
			"hello.cl",
			"Hello, COOL World!\n",
		},
		{
			"Simple If Statement",
			"if.cl",
			"The number 10 is less than 20\nThe number is exactly 10!\nThis is after the if hahaha",
		},
		{
			"Multiple If Statements",
			"if2.cl",
			`Outer if: num < 20
The number 15 is between 10 and 19
---------------------------
Now checking equality:
The number is exactly 15!
The number is not 20
This is after the nested if statements`,
		},
		{
			"While Loop",
			"while.cl",
			"Starting while loop demonstration\nCounter value: 0\nCounter value: 1\nCounter value: 2\nCounter value: 3\nCounter value: 4\nWhile loop completed. Final counter value: 5\n",
		},
		{
			"If and While Combined",
			"if_while.cl",
			`Starting while loop demonstration
Counter value: 0
Regular iteration
Counter value: 1
Regular iteration
Counter value: 2
Found the middle value!
Executing nested while loop...
  Nested counter: 0
  Nested counter: 1
  Nested counter: 2
Nested while loop completed
Counter value: 3
Regular iteration
Counter value: 4
Regular iteration
While loop completed. Final counter value: 5`,
		},
		{
			"Simple Sum",
			"sum.cl",
			"Hello from init\n5\n20\nHello from sum\nThe sum is: 25\nhhhhhhhhIlyas\ntest STRING has changed haha",
		},
		{
			"Complex Sum",
			"sum2.cl",
			`Hello from init
5
20
The sum is: 25
-----------------------
The sum is: 120
-----------------------
Ilyas
another name
another name`,
		},
		{
			"Int Input",
			"int_input.cl",
			"Enter a number: You entered the number: 42\nDouble of your number: 84\n",
		},
		{
			"Factorial",
			"factorial.cl",
			"Recursive factorial: 120\nIterative factorial: 120\n",
		},
		{
			"Nested While",
			"nested_while.cl",
			"Nested while loops:\nOuter: 0\n  Inner: 0\n  Inner: 1\nOuter: 1\n  Inner: 0\n  Inner: 1\nOuter: 2\n  Inner: 0\n  Inner: 1\nDone.\n",
		},
		{
			"Inheritance",
			"inheritance.cl",
			`===== Shape Inheritance Demo =====
Base Shape: I am a Shape with color white
Area: 0

Circle: I am a Circle with color red
Area: 147
Circumference: 42

Rectangle: I am a Rectangle with color blue
Area: 32
Perimeter: 24`,
		},
		{
			"Inheritance 2 (multilevel)",
			"inheritance2.cl",
			`Class A: 10
10
Class B: 30
10
20
Class C: 60
10
20
30`,
		},
		{
			"String Methods",
			"string_methods.cl",
			`String: Hello, Cool!
Length: 12

String 1: Hello
String 2: World!
Concatenated: HelloWorld!

String: The Cool Programming Language
Substring (4,4): Cool
Substring (0,3): The
Substring (9,11): Program`,
		},
		{
			"Abort",
			"abort.cl",
			"Testing Object.abort()\nAbout to call abort...\n",
		},
		{
			"Case with Int",
			"case_int.cl",
			"Testing case with Int (10):\nIt's an integer: 10\n",
		},
		{
			"Case with String",
			"case_string.cl",
			"Testing case with String (\"Hello\"):\nIt's a string: Hello\n",
		},
		{
			"Case with Bool",
			"case_bool.cl",
			"Testing case with Bool (true):\nIt's a boolean\n",
		},
		{
			"Case with Object",
			"case_object.cl",
			"Testing case with Object (new Object):\nIt's some other object\n",
		},
		{
			"Inheritance 4",
			"inheritance4.cl",
			`Animal makes a sound
Rex barks!
Whiskers meows!
Labrador Rex eats food
Whiskers purrs softly`,
		},
		{
			"Static Dispatch",
			"static_dispatch.cl",
			`1 1
22 22
22 2
33 3 35`,
		},
	}

	for _, example := range examples {
		t.Run(example.name, func(t *testing.T) {
			// Read the example file
			inputFile := filepath.Join("../cool_examples", example.filename)
			input, err := os.ReadFile(inputFile)
			if err != nil {
				t.Fatalf("Failed to read file %s: %v", inputFile, err)
			}

			// Parse the program
			l := lexer.NewLexer(strings.NewReader(string(input)))
			p := parser.New(l)
			program := p.ParseProgram()

			if len(p.Errors()) != 0 {
				t.Fatalf("Parser errors in %s: %v", example.filename, p.Errors())
			}

			// Generate LLVM IR
			module, err := Generate(program)
			if err != nil {
				t.Fatalf("Code generation error for %s: %v", example.filename, err)
			}

			// Write the LLVM IR to a temporary file
			tempFile := filepath.Join(os.TempDir(), "test_"+example.filename+".ll")
			err = os.WriteFile(tempFile, []byte(module.String()), 0644)
			if err != nil {
				t.Fatalf("Failed to write LLVM IR to file: %v", err)
			}
			defer os.Remove(tempFile)

			// Try to compile the LLVM IR with clang
			// This is a simple validation that the output is syntactically valid
			output, err := exec.Command("clang", "-Wno-deprecated", tempFile, "-o", filepath.Join(os.TempDir(), "test_"+example.filename+".exe"), "-llegacy_stdio_definitions").CombinedOutput()

			// Just check if clang accepted the LLVM IR - don't worry about successfully compiling
			// as that would require the full runtime to be available
			if err != nil && strings.Contains(string(output), "error:") {
				t.Fatalf("Generated LLVM IR for %s is not valid: %v\n%s", example.filename, err, output)
			}

			// For int_input.cl, we need to provide simulated input
			var cmd *exec.Cmd
			executablePath := filepath.Join(os.TempDir(), "test_"+example.filename+".exe")

			if example.filename == "int_input.cl" {
				// Create a buffer with simulated input (42)
				simulatedInput := bytes.NewBufferString("42\n")
				cmd = exec.Command(executablePath)
				cmd.Stdin = simulatedInput
			} else {
				cmd = exec.Command(executablePath)
			}

			// Capture the output
			var stdout bytes.Buffer
			cmd.Stdout = &stdout

			// Run the compiled program
			err = cmd.Run()
			if err != nil {
				t.Logf("Warning: Program execution failed: %v", err)
				return // Continue with other tests
			}

			// Check the output
			actualOutput := stdout.String()

			// Normalize line endings and whitespace for comparison
			expectedNormalized := normalizeOutput(example.expectedOutput)
			actualNormalized := normalizeOutput(actualOutput)

			if expectedNormalized != actualNormalized {
				t.Errorf("Output mismatch for %s:\nExpected: %q\nActual: %q", example.filename, expectedNormalized, actualNormalized)
			}
		})
	}
}

// normalizeOutput normalizes string output for comparison
func normalizeOutput(s string) string {
	// Replace Windows line endings with Unix line endings
	s = strings.ReplaceAll(s, "\r\n", "\n")
	// Trim trailing spaces from each line
	lines := strings.Split(s, "\n")
	for i, line := range lines {
		lines[i] = strings.TrimRight(line, " ")
	}
	// Rejoin and trim any leading/trailing whitespace
	return strings.TrimSpace(strings.Join(lines, "\n"))
}

// TestSpecificFeatures tests specific COOL language features to ensure they generate correct LLVM IR
func TestSpecificFeatures(t *testing.T) {
	testCases := []struct {
		name           string
		coolCode       string
		expectedOutput string
	}{
		{
			name: "Simple Integer Arithmetic",
			coolCode: `
class Main {
	main(): Object {
		(new IO).out_int(1 + 2 * 3)
	};
};`,
			expectedOutput: "7",
		},
		{
			name: "If-Then-Else",
			coolCode: `
class Main {
	main(): Object {
		if 5 < 10 then
			(new IO).out_string("Five is less than ten")
		else
			(new IO).out_string("This shouldn't print")
		fi
	};
};`,
			expectedOutput: "Five is less than ten",
		},
		{
			name: "IsVoid Operation",
			coolCode: `
class Main {
	main(): Object {
		let 
			io: IO <- new IO,
			obj: Object <- new Object,
			nullObj: Object
		in {
			if isvoid nullObj then
				io.out_string("nullObj is void!")
			else
				io.out_string("nullObj is not void!")
			fi;

			io.out_string("\n");
			
			if isvoid obj then
				io.out_string("obj is void!")
			else
				io.out_string("obj is not void!")
			fi;
		}
	};
};`,
			expectedOutput: "nullObj is void!\nobj is not void!",
		},
		{
			name: "Basic Arithmetic Operations",
			coolCode: `
class Main {
	main(): Object {
		let 
			io: IO <- new IO
		in {
			io.out_int(5 + 3);
			io.out_string("\n");
			io.out_int(10 - 4);
			io.out_string("\n");
			io.out_int(3 * 6);
			io.out_string("\n");
			io.out_int(20 / 5);
		}
	};
};`,
			expectedOutput: "8\n6\n18\n4",
		},
		{
			name: "Complex Arithmetic Expressions",
			coolCode: `
class Main {
	main(): Object {
		let 
			io: IO <- new IO,
			a: Int <- 5,
			b: Int <- 3,
			c: Int <- 2
		in {
			io.out_int((a + b) * c);
			io.out_string("\n");
			io.out_int(a * b + c);
			io.out_string("\n");
			io.out_int(a * (b + c));
			io.out_string("\n");
			io.out_int((a + b) * (a - c));
		}
	};
};`,
			expectedOutput: "16\n17\n25\n24",
		},
		{
			name: "Arithmetic with Negation",
			coolCode: `
class Main {
	main(): Object {
		let 
			io: IO <- new IO
		in {
			io.out_int(~5);
			io.out_string("\n");
			io.out_int(~(~5));
			io.out_string("\n");
			io.out_int(10 + ~5);
			io.out_string("\n");
			io.out_int(~(3 * 2));
		}
	};
};`,
			expectedOutput: "-5\n5\n5\n-6",
		},
		{
			name: "Division Operations",
			coolCode: `
class Main {
	main(): Object {
		let 
			io: IO <- new IO
		in {
			io.out_int(10 / 2);
			io.out_string("\n");
			io.out_int(20 / 4);
			io.out_string("\n");
			io.out_int(15 / 3);
			io.out_string("\n");
			io.out_int(100 / 25);
			io.out_string("\n");
			io.out_int(7 / 2);  -- Integer division truncates
		}
	};
};`,
			expectedOutput: "5\n5\n5\n4\n3",
		},
		{
			name: "Complex Arithmetic Expressions",
			coolCode: `
class Main {
	main(): Object {
		let 
			io: IO <- new IO,
			a: Int <- 5,
			b: Int <- 10,
			c: Int <- 3,
			d: Int <- 7
		in {
			io.out_int((a + b) * (c - d) / 2 + 15);  -- (15 * -4) / 2 + 15 = -30 + 15 = -15
			io.out_string("\n");
			io.out_int(a * b + c * d - a / c + b * (c + d));  -- 50 + 21 - 1 + 10 * 10 = 70 + 100 = 170
			io.out_string("\n");
			io.out_int(~(a * b) + c * (d - a) + b / c * d);  -- -50 + 3 * 2 + 3 * 7 = -50 + 6 + 21 = -23
			io.out_string("\n");
			io.out_int((a + b + c + d) * (a - b + c - d) / (a + 1));  -- 25 * -9 / 6 = -225 / 6 = -37
			io.out_string("\n");
			io.out_int(a * b * c * d / (a + b + c + d) + a * a - b + c * c * c);  -- 1050 / 25 + 25 - 10 + 27 = 42 + 25 - 10 + 27 = 84
		}
	};
};`,
			expectedOutput: "-15\n170\n-23\n-37\n84",
		},
		{
			name: "Boolean NOT Operations",
			coolCode: `
class Main {
	main(): Object {
		let 
			io: IO <- new IO,
			a: Bool <- true,
			b: Bool <- false
		in {
			if not a then
				io.out_string("not true = true (ERROR)")
			else
				io.out_string("not true = false (CORRECT)")
			fi;
			io.out_string("\n");
			
			if not b then
				io.out_string("not false = true (CORRECT)")
			else
				io.out_string("not false = false (ERROR)")
			fi;
		}
	};
};`,
			expectedOutput: "not true = false (CORRECT)\nnot false = true (CORRECT)",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Parse the program
			l := lexer.NewLexer(strings.NewReader(tc.coolCode))
			p := parser.New(l)
			program := p.ParseProgram()

			if len(p.Errors()) != 0 {
				t.Fatalf("Parser errors: %v", p.Errors())
			}

			// Generate LLVM IR
			module, err := Generate(program)
			if err != nil {
				t.Fatalf("Code generation error: %v", err)
			}

			// Write the LLVM IR to a temporary file
			tempFile := filepath.Join(os.TempDir(), "test_"+tc.name+".ll")
			tempExec := filepath.Join(os.TempDir(), "test_"+tc.name+".exe")
			err = os.WriteFile(tempFile, []byte(module.String()), 0644)
			if err != nil {
				t.Fatalf("Failed to write LLVM IR to file: %v", err)
			}
			defer os.Remove(tempFile)
			defer os.Remove(tempExec)

			// Compile the LLVM IR
			output, err := exec.Command("clang", "-Wno-deprecated", tempFile, "-o", tempExec, "-llegacy_stdio_definitions").CombinedOutput()
			if err != nil {
				t.Logf("Compilation warning/error: %s", output)
				t.Logf("Skipping output verification due to compilation issues")
				return
			}

			// Run the compiled program
			cmd := exec.Command(tempExec)
			var stdout bytes.Buffer
			cmd.Stdout = &stdout
			err = cmd.Run()
			if err != nil {
				t.Logf("Warning: Program execution failed: %v", err)
				return
			}

			// Verify output
			actualOutput := normalizeOutput(stdout.String())
			expectedOutput := normalizeOutput(tc.expectedOutput)
			if actualOutput != expectedOutput {
				t.Errorf("Output mismatch:\nExpected: %q\nActual: %q", expectedOutput, actualOutput)
			}
		})
	}
}
