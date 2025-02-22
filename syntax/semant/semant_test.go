package semant_test

import (
	"cool-compiler/lexer"
	"cool-compiler/parser"
	"cool-compiler/semant"
	"fmt"
	"strings"
	"testing"
)

func TestSemanticAnalysis(t *testing.T) {
	tests := []struct {
		name     string
		code     string
		expected []string // Expected error substrings
	}{
		// Class Structure Tests
		{
			name: "Duplicate class definition",
			code: `
                class Main {};
                class Main {};
            `,
			expected: []string{"class Main is already defined"},
		},
		{
			name: "Inherit from basic class",
			code: `
                class Bad inherits Int {};
            `,
			expected: []string{"class Bad cannot inherit from built-in class Int"},
		},
		{
			name: "Inherit from undefined class",
			code: `
                class Main inherits UndefinedClass {};
            `,
			expected: []string{"undefined class UndefinedClass"},
		},
		{
			name: "hhh",
			code: `
                class A inherits B {};
                class B {};
            `,
			expected: []string{"inheritance cycle detected"},
		},
		{
			name: "Inheritance cycle detection",
			code: `
                class A inherits B {};
                class B inherits A {};
            `,
			expected: []string{"inheritance cycle detected"},
		},
		{
			name: "Method return type mismatch",
			code: `
                class Main {
                    test(): String { 5 };
                };
            `,
			expected: []string{"method test is expected to return String, found Int"},
		},

		// Attribute Validation Tests
		{
			name: "Attribute shadowing in nested scope",
			code: `
                class Main {
                    x: Int;
                    test(): Int {
                        let x: String <- "test" in x
                    };
                };
            `,
			expected: []string{},
		},
		{
			name: "Attribute initialization type mismatch",
			code: `
                class Main {
                    x: Int <- "hello";
                };
            `,
			expected: []string{"attribute x cannot be of type String, expected Int"},
		},

		// Expression Type Checking
		{
			name: "Arithmetic type mismatch",
			code: `
                class Main {
                    test(): Int {
                        "string" + 5
                    };
                };
            `,
			expected: []string{"arithmetic operation on non-Int types: String + Int"},
		},
		{
			name: "Comparison type mismatch",
			code: `
                class Main {
                    test(): Bool {
                        5 < "string"
                    };
                };
            `,
			expected: []string{"comparison between incompatible types: Int < String"},
		},

		// Special Cases
		{
			name: "SELF_TYPE validation",
			code: `
                class Main {
                    test(): SELF_TYPE {
                        new SELF_TYPE
                    };
                };
            `,
			expected: []string(nil),
		},
		{
			name: "Invalid SELF_TYPE usage",
			code: `
                class Main {
                    test(): SELF_TYPE {
                        new SELF_TYPE
                    };
                };
                class Other {
                    x: Main <- (new Main).test();
                };
            `,
			expected: []string(nil), // Should be valid
		},

		// Let Expression Tests
		{
			name: "Let binding type mismatch",
			code: `
                class Main {
                    test(): Int {
                        let x: Int <- "string" in x
                    };
                };
            `,
			expected: []string{"Let binding with wrong type"},
		},

		// Dispatch Expression Tests
		{
			name: "Dynamic dispatch on undefined method",
			code: `
                class Main {
                    test(): Object {
                        (new Main).undefined_method()
                    };
                };
            `,
			expected: []string{"undefined method undefined_method in Main"},
		},

		// Case Expression Tests
		{
			name: "Case with duplicate branch types",
			code: `
                class Main {
                    test(x: Object): Object {
                        case x of
                            y: Object => 1;
                            z: Object => 2;
                        esac
                    };
                };
            `,
			expected: []string{"duplicate branch type Object in case expression"},
		},
		{
			name: "Case with undefined branch type",
			code: `
                class Main {
                    test(x: Object): Object {
                        case x of
                            y: UndefinedType => 1;
                        esac
                    };
                };
            `,
			expected: []string{"undefined type UndefinedType in case expression"},
		},

		// Object and IO method dispatch tests
		{
			name: "Valid Object method dispatch",
			code: `
                class Main {
                    test(): Object {
                        (new Main)@Object.abort()
                    };
                };
            `,
			expected: []string(nil), // Should be valid
		},
		{
			name: "Invalid static dispatch method",
			code: `
                class Main {
                    test(): Object {
                        (new Main)@Object.invalid_method()
                    };
                };
            `,
			expected: []string{"undefined method invalid_method in Object"},
		},
		{
			name: "Valid IO method dispatch",
			code: `
                class Main {
                    test(): IO {
                        (new IO).out_string("Hello")
                    };
                };
            `,
			expected: []string(nil), // Should be valid
		},
		{
			name: "Valid IO method chaining",
			code: `
                class Main {
                    main(): IO {
                        (new IO).out_string("Hello").out_int(42)
                    };
                };
            `,
			expected: []string(nil), // Valid SELF_TYPE return chaining
		},
		{
			name: "Invalid main method",
			code: `
                class Main {
                    main(x: Int): Object { x };
                };
            `,
			expected: []string{"main method must have 0 parameters"},
		},
	}

	testNum := 3

	for _, tt := range tests[testNum : testNum+1] {
		// for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := parser.New(lexer.NewLexer(strings.NewReader(tt.code)))
			program := p.ParseProgram()
			if len(p.Errors()) > 0 {
				t.Fatalf("parser errors: %v", p.Errors())
			}

			analyser := semant.NewSemanticAnalyser()
			analyser.Analyze(program)

			// Printing the errors
			fmt.Printf("Semantic errors:\n")
			for _, err := range analyser.Errors() {
				fmt.Printf("  %s\n", err)
			}

			// Verify expected errors
			if len(tt.expected) != len(analyser.Errors()) {
				t.Errorf("expected %d errors, got %d", len(tt.expected), len(analyser.Errors()))
			}

			for _, expectedErr := range tt.expected {
				found := false
				for _, actualErr := range analyser.Errors() {
					if strings.Contains(actualErr, expectedErr) {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("missing expected error containing: %q", expectedErr)
				}
			}
		})
	}
}
