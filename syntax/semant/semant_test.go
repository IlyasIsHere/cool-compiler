package semant_test

import (
	"cool-compiler/lexer"
	"cool-compiler/parser"
	"cool-compiler/semant"
	"fmt"
	"strings"
	"testing"
)

func validMainClass() string {
	return `
        class Main {
            main(): Object {
                0
            };
        };
    `
}

func withMainClass(code string) string {
	return code + validMainClass()
}

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
        class Main {
            main(): Object { 0 };
        };
        class Main {};
    `,
			expected: []string{"class Main is already defined"},
		},
		{
			name: "Inherit from basic class",
			code: `
				class Bad inherits Int {};
				class Main {
					main(): Object { 0 };
				};
            `,
			expected: []string{"class Bad cannot inherit from built-in class Int"},
		},
		{
			name: "Inherit from undefined class",
			code: `
                class Main inherits UndefinedClass {
					main(): Object { 0 };
				};
            `,
			expected: []string{"undefined class UndefinedClass"},
		},
		{
			name: "Main inherits from A",
			code: `
                class A {};
                class Main inherits A {
					main(): Object { 0 };
				};
            `,
			expected: []string{},
		},
		{
			name: "Main inherits main method from parent",
			code: `
                class A {
                    main(): Object { 0 };
                };
                class Main inherits A {};
            `,
			expected: []string{"class Main must define method 'main' with 0 parameters"},
		},
		{
			name: "main method doesn't respect signature",
			code: `
                class Main {
                    main(a: A): Object { 1 };
                };
            `,
			expected: []string{"main method must have 0 parameters", "undefined type A in formal parameter of method main"},
		},
		{
			name: "undefined type in method parameters",
			code: `
                class Main {
                    main(): Object { 1 };
                    test(a: A): String { "hhh" }; 
                };
            `,
			expected: []string{"undefined type A in formal parameter of method test"},
		},
		{
			name: "Inheritance cycle detection",
			code: `
                class A inherits B {};
                class B inherits A {};
				class Main {
					main(): Object { 0 };
				};
            `,
			expected: []string{"inheritance cycle detected"},
		},
		{
			name: "Method return type mismatch",
			code: `
                class Main {
                    test(): String { 5 };
					main(): Object { 0 };
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
                    test(): String {
                        let x: String <- "test" in x
                    };
					main(): Object { 0 };
                };
            `,
			expected: []string{},
		},
		{
			name: "Attribute initialization type mismatch",
			code: `
                class Main {
                    x: Int <- "hello";
					main(): Object { 0 };
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
					main(): Object { 0 };
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
					main(): Object { 0 };
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
					main(): Object { 0 };
                };
            `,
			expected: []string(nil),
		},
		{
			name: "Valid SELF_TYPE usage",
			code: `
                class Main {
                    test(): SELF_TYPE {
                        new SELF_TYPE
                    };
					main(): Object { 0 };
                };
                class Other {
                    x: Main <- (new Main).test();
                };
            `,
			expected: []string(nil),
		},

		// Let Expression Tests
		{
			name: "Let binding type mismatch",
			code: `
                class Main {
                    test(): Int {
                        let x: Int <- "string" in x
                    };
					main(): Object { 0 };
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
					main(): Object { 0 };
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
					main(): Object { 0 };
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
					main(): Object { 0 };
                };
            `,
			expected: []string{"undefined type UndefinedType in case branch"},
		},

		// Object and IO method dispatch tests
		{
			name: "Valid Object method dispatch",
			code: `
                class Main {
                    test(): Object {
                        (new Main)@Object.abort()
                    };
					main(): Object { 0 };
                };
            `,
			expected: []string(nil),
		},
		{
			name: "Valid Object method dispatch 2",
			code: `
                class Main {
                    test(): Object {
                        (new Main).abort()
                    };
					main(): Object { 0 };
                };
            `,
			expected: []string(nil),
		},
		{
			name: "Invalid static dispatch method",
			code: `
                class Main {
                    test(): Object {
                        (new Main)@Object.invalid_method()
                    };
					main(): Object { 0 };
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
					main(): Object { 0 };
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
		{
			name: "undefined type in attribute",
			code: `
                class Main {
                    x: A;
					main(): Object { 0 };
                };
            `,
			expected: []string{"undefined type A for attribute x"},
		},
		{
			name: "Attribute redefinition in child class",
			code: withMainClass(`
                class A {
                    x: Int;
                };
                class B inherits A {
                    x: String;
                };
            `),
			expected: []string{"attribute x is already defined in a parent class of B"},
		},
		// {
		// 	name: "Redefine SELF_TYPE",
		// 	code: `
		//         class SELF_TYPE {
		//             main(): Object { 0 };
		//         };
		//         class Main {
		//             main(): Object { 0 };
		//         };
		//     `,
		// 	expected: []string{"cannot inherit from special type SELF_TYPE", "redefinition of basic class SELF_TYPE"},
		// },
		// {
		// 	name: "Inheritance from SELF_TYPE",
		// 	code: `
		//         class A inherits SELF_TYPE {};
		//         class Main {
		//             main(): Object { 0 };
		//         };
		//     `,
		// 	expected: []string{"cannot inherit from special type SELF_TYPE"},
		// },
		{
			name: "Let with no init expression, undefined type",
			code: `
                class Main {
                    main(): Object {
                        let x: UndefinedType in 0
                    };
                };
            `,
			expected: []string{"undefined type UndefinedType in let binding"},
		},
		{
			name: "Let with init expression, undefined type",
			code: `
                class Main {
                    main(): Object {
                        let x: UndefinedType <- 1 in 0
                    };
                };
            `,
			expected: []string{"undefined type UndefinedType in let binding"},
		},
		{
			name: "Assignment type mismatch",
			code: `
                class Main {
                    x: Int;
                    main(): Object {
                        x <- "string"
                    };
                };
            `,
			expected: []string{"type mismatch in assignment: variable 'x' has type Int but was assigned value of type String"},
		},
		{
			name: "Assignment to undefined variable",
			code: `
                class Main {
                    main(): Object {
                        y <- 1
                    };
                };
            `,
			expected: []string{"undefined identifier y in assignment"},
		},
		{
			name: "Method call on undefined object",
			code: `
                class Main {
                    main(): Object {
                        a.test()
                    };
                };
            `,
			expected: []string{"undefined identifier a", "undefined method test in Object"},
		},
		{
			name: "Static dispatch on undefined class",
			code: `
                class Main {
                    main(): Object {
                        (new Main)@Undefined.test()
                    };
                };
            `,
			expected: []string{"undefined type Undefined in static dispatch"},
		},
		{
			name: "if statement with non-boolean condition",
			code: `
                class Main {
                    main(): Object {
                        if 1 then 1 else 0 fi
                    };
                };
            `,
			expected: []string{"condition of if statement is of type Int, expected Bool"},
		},
		{
			name: "while loop with non-boolean condition",
			code: `
                class Main {
                    main(): Object {
                        while 1 loop 0 pool
                    };
                };
            `,
			expected: []string{"condition of if statement is of type Int, expected Bool"},
		},
		{
			name: "redefine object",
			code: `
				class Object {};
                class Main {
                    main(): Object {
						0
					 };
                };
            `,
			expected: []string{"class Object is already defined"},
		},
		{
			name: "redefine IO",
			code: `
				class IO {};
                class Main {
                    main(): Object {
						0
					 };
                };
            `,
			expected: []string{"class IO is already defined"},
		},
		{
			name: "Valid isvoid check",
			code: `
                class Main {
                    main(): Object {
                        isvoid (new Main)
                    };
                };
            `,
			expected: []string{},
		},
		{
			name: "Invalid void check",
			code: `
                class Main {
                    main(): Object {
                        isvoid 5
                    };
                };
            `,
			expected: []string{},
		},
		{
			name: "Call method with wrong number of arguments",
			code: `
				class A {
					foo(x: Int): Int { x };
				};
                class Main {
                    main(): Object {
                        (new A).foo()
                    };
                };
            `,
			expected: []string{"method foo expects 1 arguments, but got 0"},
		},
		{
			name: "Call method with wrong argument type",
			code: `
				class A {
					foo(x: Int): Int { x };
				};
                class Main {
                    main(): Object {
                        (new A).foo("string")
                    };
                };
            `,
			expected: []string{"argument 1 of method foo expects type Int, but got String"},
		},
		{
			name: "Override method with wrong number of arguments",
			code: withMainClass(`
				class A {
					foo(x: Int): Int { x };
				};
				class B inherits A {
					foo(): Int { 0 };
				};
			`),
			expected: []string{"method foo overrides parent method but has different number of parameters (0 vs 1)"},
		},
		{
			name: "Override method with wrong argument type",
			code: withMainClass(`
				class A {
					foo(x: Int): Int { x };
				};
				class B inherits A {
					foo(x: String): Int { 0 };
				};
			`),
			expected: []string{"method foo overrides parent method but parameter 1 has different type (String vs Int)"},
		},
	}

	// testNum := len(tests) - 1

	// for _, tt := range tests[testNum : testNum+1] {
	for _, tt := range tests {
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
