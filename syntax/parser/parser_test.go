package parser

import (
	"cool-compiler/ast"
	"cool-compiler/lexer"
	"strings"
	"testing"
)

func newParserFromInput(input string) *Parser {
	l := lexer.NewLexer(strings.NewReader(input))
	return New(l)
}

func checkParserErrors(t *testing.T, p *Parser, i int) {
	errors := p.Errors()
	if len(errors) > 0 {
		t.Errorf("parser has %d errors for test case %d", len(errors), i)
		for _, msg := range errors {
			t.Errorf("parser error: %q", msg)
		}
		t.FailNow()
	}
}

func TestClassParser(t *testing.T) {
	tests := []struct {
		input          string
		expectedName   string
		expectedParent string
		shouldFail     bool
		errorContains  string
	}{
		{
			input:          "class Main {};",
			expectedName:   "Main",
			expectedParent: "",
		},
		{
			input:          "class A {age:Int<-30;};",
			expectedName:   "A",
			expectedParent: "",
		},
		{
			input:          "class B {func(): Int { 5 };};",
			expectedName:   "B",
			expectedParent: "",
		},
		{
			input:          "class B inherits A {func(): Object { 10 };};",
			expectedName:   "B",
			expectedParent: "A",
		},
		{
			input: `class Complex {
				x: Int <- 5;
				y: Int <- 10;
				init(): Complex { self };
				add(other: Complex): Complex { self };
			};`,
			expectedName:   "Complex",
			expectedParent: "",
		},
		{
			input:          `class D inherits C inherits B {};`,
			shouldFail:     true,
			errorContains:  "Expected next token to be",
			expectedName:   "",
			expectedParent: "",
		},
		{
			input: `class List inherits Collection {
				head: Int;
				tail: List;
				cons(x: Int): List { self };
			};`,
			expectedName:   "List",
			expectedParent: "Collection",
		},
		{
			input:          `class 123Invalid {};`,
			shouldFail:     true,
			errorContains:  "Expected next token to be TYPEID",
			expectedName:   "",
			expectedParent: "",
		},
	}

	for i, tt := range tests {
		parser := newParserFromInput(tt.input)
		class := parser.ParseClass()

		if tt.shouldFail {
			errors := parser.Errors()
			if len(errors) == 0 {
				t.Errorf("test [%d] expected parsing errors but got none", i)
				continue
			}

			foundExpectedError := false
			for _, err := range errors {
				if strings.Contains(err, tt.errorContains) {
					foundExpectedError = true
					break
				}
			}

			if !foundExpectedError {
				t.Errorf("test [%d] expected error containing '%s', got: %v",
					i, tt.errorContains, errors)
			}
			continue
		}

		checkParserErrors(t, parser, i)

		if class.Name.Value != tt.expectedName {
			t.Fatalf("[%q]: expected class name to be %q got %q",
				tt.input, tt.expectedName, class.Name.Value)
		}

		if class.Parent != nil {
			if class.Parent.Value != tt.expectedParent {
				t.Fatalf("[%q]: expected class parent to be %q got %q",
					tt.input, tt.expectedParent, class.Parent.Value)
			}
		} else if tt.expectedParent != "" {
			t.Fatalf("[%q]: expected class parent to be %q got nil",
				tt.input, tt.expectedParent)
		}
	}
}

func TestFormalParsing(t *testing.T) {
	tests := []struct {
		input         string
		expectedNames []string
		expectedTypes []string
	}{
		{
			input:         "(var1:Integer)",
			expectedNames: []string{"var1"},
			expectedTypes: []string{"Integer"},
		},
		{
			input:         "(var1:Integer,var2:Boolean,var3:String)",
			expectedNames: []string{"var1", "var2", "var3"},
			expectedTypes: []string{"Integer", "Boolean", "String"},
		},
	}

	for _, tt := range tests {
		parser := newParserFromInput(tt.input)
		formals := parser.parseFormals()

		if len(parser.errors) > 0 {
			for _, err := range parser.errors {
				t.Errorf("Parsing Error %s\n", err)
			}
			t.Fatalf("[%q]: Found errors while parsing", tt.input)
		}

		if len(formals) != len(tt.expectedNames) {
			t.Fatalf("[%q]: expected %d formals got %d: %v", tt.input, len(tt.expectedNames), len(formals), formals)
		}

		for i, formal := range formals {
			if formal.Name.Value != tt.expectedNames[i] {
				t.Fatalf("[%q]: expected formal name to be %q got %q", tt.input, tt.expectedNames[i], formal.Name.Value)
			}
			if formal.TypeDecl.Value != tt.expectedTypes[i] {
				t.Fatalf("[%q]: expected formal type to be %q got %q", tt.input, tt.expectedNames[i], formal.Name.Value)
			}
		}
	}
}

func TestMethodParsing(t *testing.T) {
	tests := []struct {
		input               string
		expectedMethodName  string
		expectedFormalNames []string
		expectedFormalTypes []string
		expectedMethodType  string
		expectedExpression  string
	}{
		{
			input:               "main(): Void { true };",
			expectedMethodName:  "main",
			expectedFormalNames: []string{},
			expectedFormalTypes: []string{},
			expectedMethodType:  "Void",
			expectedExpression:  "true",
		},
		{
			input:               "sum(a:Integer,b:Integer): Integer { a + b };",
			expectedMethodName:  "sum",
			expectedFormalNames: []string{"a", "b"},
			expectedFormalTypes: []string{"Integer", "Integer"},
			expectedMethodType:  "Integer",
			expectedExpression:  "(a + b)",
		},
		{
			input:               "factorial(n:Integer): Integer { if n = 0 then 1 else n * factorial(n-1) fi };",
			expectedMethodName:  "factorial",
			expectedFormalNames: []string{"n"},
			expectedFormalTypes: []string{"Integer"},
			expectedMethodType:  "Integer",
			expectedExpression:  "if (n = 0) then 1 else (n * factorial((n - 1))) fi",
		},
		{
			input:               "print(msg:String): Object { { out_string(msg); true; } };",
			expectedMethodName:  "print",
			expectedFormalNames: []string{"msg"},
			expectedFormalTypes: []string{"String"},
			expectedMethodType:  "Object",
			expectedExpression:  "{ out_string(msg); true; }",
		},
		{
			input:               "init(x:Int): Object { let temp:Int <- x in { value <- temp; self; } };",
			expectedMethodName:  "init",
			expectedFormalNames: []string{"x"},
			expectedFormalTypes: []string{"Int"},
			expectedMethodType:  "Object",
			expectedExpression:  "let temp : Int <- x in { (value <- temp); self; }",
		},
		{
			input: `complex(): Object { 
				{
					let x: Int <- 5 in x + 1;
					if x < 10 then x else 0 fi;
					while x < 0 loop x <- x - 1 pool;
				}
			};`,
			expectedMethodName:  "complex",
			expectedFormalNames: []string{},
			expectedFormalTypes: []string{},
			expectedMethodType:  "Object",
			expectedExpression:  "{ let x : Int <- 5 in (x + 1); if (x < 10) then x else 0 fi; while (x < 0) loop (x <- (x - 1)) pool; }",
		},
		{
			input: `nested(): Int {
				{
					if true then
						{ let x: Int <- 1 in x + 2; }
					else
						{ let y: Int <- 3 in y + 4; }
					fi;
				}
			};`,
			expectedMethodName:  "nested",
			expectedFormalNames: []string{},
			expectedFormalTypes: []string{},
			expectedMethodType:  "Int",
			expectedExpression:  "{ if true then { let x : Int <- 1 in (x + 2); } else { let y : Int <- 3 in (y + 4); } fi; }",
		},
	}

	for i, tt := range tests {
		parser := newParserFromInput(tt.input)
		method := parser.parseMethod()
		checkParserErrors(t, parser, i)

		if method.Name.Value != tt.expectedMethodName {
			t.Fatalf("[%q]: Expected method name to be %q found %q", tt.input, tt.expectedMethodName, method.Name.Value)
		}

		if len(method.Formals) != len(tt.expectedFormalNames) {
			t.Fatalf("[%q]: Expected %d formals, found %d", tt.input, len(tt.expectedFormalNames), len(method.Formals))
		}

		for i, formal := range method.Formals {
			if formal.Name.Value != tt.expectedFormalNames[i] {
				t.Fatalf("[%q]: Expected formal name to be %q found %q", tt.input, tt.expectedFormalNames[i], formal.Name.Value)
			}
			if formal.TypeDecl.Value != tt.expectedFormalTypes[i] {
				t.Fatalf("[%q]: Expected formal type to be %q found %q", tt.input, tt.expectedFormalTypes[i], formal.TypeDecl.Value)
			}
		}

		if method.TypeDecl.Value != tt.expectedMethodType {
			t.Fatalf("[%q]: Expected method type to be %q found %q", tt.input, tt.expectedMethodType, method.TypeDecl.Value)
		}

		if method.Expression == nil {
			t.Fatalf("[%q]: Method body expression cannot be nil", tt.input)
		}

		actualExpression := SerializeExpression(method.Expression)
		if actualExpression != tt.expectedExpression {
			t.Fatalf("[%q]: Expected method body to be %q found %q", tt.input, tt.expectedExpression, actualExpression)
		}
	}
}

func TestAttributeParsing(t *testing.T) {
	tests := []struct {
		input              string
		expectedName       string
		expectedType       string
		shouldFail         bool
		errorContains      string
		expectedExpression ast.Expression
	}{
		{
			input:        "firstName:String",
			expectedName: "firstName",
			expectedType: "String",
		},
		{
			input:        "age:Integer<-0",
			expectedName: "age",
			expectedType: "Integer",
		},
		{
			input:         "invalid: <- 5",
			shouldFail:    true,
			errorContains: "Expected next token to be TYPEID",
		},
	}

	for i, tt := range tests {
		parser := newParserFromInput(tt.input)
		attribute := parser.parseAttribute()

		if tt.shouldFail {
			errors := parser.Errors()
			if len(errors) == 0 {
				t.Errorf("test [%d] expected parsing errors but got none", i)
				continue
			}

			foundExpectedError := false
			for _, err := range errors {
				if strings.Contains(err, tt.errorContains) {
					foundExpectedError = true
					break
				}
			}

			if !foundExpectedError {
				t.Errorf("test [%d] expected error containing '%s', got: %v",
					i, tt.errorContains, errors)
			}
			continue
		}

		checkParserErrors(t, parser, i)
		if attribute.Name.Value != tt.expectedName {
			t.Fatalf("[%q]: Expected attribute name to be %q got %q", tt.input, tt.expectedName, attribute.Name.Value)
		}
		if attribute.TypeDecl.Value != tt.expectedType {
			t.Fatalf("[%q]: Expected attribute type to be %q got %q", tt.input, tt.expectedType, attribute.TypeDecl.Value)
		}
	}
}

func TestExpressionParsing(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"5", "5"},
		{`"hello world"`, `"hello world"`},
		{"true", "true"},
		{"false", "false"},
		{"x", "x"},
		{"not true", "(not true)"},
		{"1 + 2", "(1 + 2)"},
		{"1 < 2", "(1 < 2)"},
		{"1 <= 2", "(1 <= 2)"},
		{"~1", "(~ 1)"},
		{"1 = 2", "(1 = 2)"},
		{"1 * 2", "(1 * 2)"},
		{"isvoid 1", "isvoid 1"},
		{"1 / 2", "(1 / 2)"},
		{"(1 + 2)", "(1 + 2)"},
		{"new Object", "new Object"},
		{"x <- 5", "(x <- 5)"},
		{"if true then 1 else 2 fi", "if true then 1 else 2 fi"},
		{"while true loop 1 pool", "while true loop 1 pool"},
		{"1 + 2 * 3", "(1 + (2 * 3))"},
		{"(1 + 2) * 3", "((1 + 2) * 3)"},
		{"1 * 2 + 3", "((1 * 2) + 3)"},
		{"1 * 2 * 3", "((1 * 2) * 3)"},
		{"1 + 2 + 3", "((1 + 2) + 3)"},
		{"1 / 2 * 3 + 4", "(((1 / 2) * 3) + 4)"},
		{"not (1 < 2)", "(not (1 < 2))"},
		{"not true = false", "((not true) = false)"},
		{"1 < 2 = true", "((1 < 2) = true)"},
		{"if 1 + 2 < 3 * 4 then 5 else 6 fi", "if ((1 + 2) < (3 * 4)) then 5 else 6 fi"},
		{"while 1 + 1 = 2 loop x + 1 pool", "while ((1 + 1) = 2) loop (x + 1) pool"},
		{"obj.method(1 + 2, 3 * 4)", "obj.method((1 + 2), (3 * 4))"},
		{"obj@Type.method(1 + 2 * 3)", "obj@Type.method((1 + (2 * 3)))"},
		{"x <- 1 + 2 * 3", "(x <- (1 + (2 * 3)))"},
		{"x <- y <- 1 + 2", "(x <- (y <- (1 + 2)))"},
	}

	for i, tt := range tests {
		p := newParserFromInput(tt.input)
		checkParserErrors(t, p, i)

		expression := p.parseExpression(LOWEST)
		actual := SerializeExpression(expression)
		if actual != tt.expected {
			t.Errorf("test [%d] expected expression to be '%s', got '%s'", i, tt.expected, actual)
		}
	}

}

func TestComplexExpressionParsing(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			input:    `obj.method1().method2().method3(1, 2, 3)`,
			expected: "obj.method1().method2().method3(1, 2, 3)",
		},
		{
			input: `case x of
				n: Int => n + 1;
				s: String => s.length();
				o: Object => 0;
			esac`,
			expected: "case x of n : Int => (n + 1); s : String => s.length(); o : Object => 0; esac",
		},
		{
			input:    `let x: Int <- 1, y: Int <- 2 in x + y * (3 + 4)`,
			expected: "let x : Int <- 1, y : Int <- 2 in (x + (y * (3 + 4)))",
		},
		{
			input:    `if (not isvoid x) then x@Type.foo(1) else y.bar(2) fi`,
			expected: "if (not isvoid x) then x@Type.foo(1) else y.bar(2) fi",
		},
	}

	for i, tt := range tests {
		p := newParserFromInput(tt.input)
		checkParserErrors(t, p, i)

		expression := p.parseExpression(LOWEST)
		actual := SerializeExpression(expression)
		if actual != tt.expected {
			t.Errorf("test [%d] expected expression to be '%s', got '%s'", i, tt.expected, actual)
		}
	}
}
