package lexer

import (
	"strings"
	"testing"
)

func TestNextToken(t *testing.T) {
	tests := []struct {
		input             string
		expectedTokenType []TokenType
		expectedLiteral   []string
	}{
		{
			"class Main {};",
			[]TokenType{CLASS, TYPEID, LBRACE, RBRACE, SEMI, EOF},
			[]string{"class", "Main", "{", "}", ";", ""},
		},
		{
			"x <- true;-- One line comment\nx <- false;",
			[]TokenType{OBJECTID, ASSIGN, BOOL_CONST, SEMI, OBJECTID, ASSIGN, BOOL_CONST, SEMI, EOF},
			[]string{"x", "<-", "true", ";", "x", "<-", "false", ";", ""},
		},
		{
			"a <- 0; b   <- a <= \"1\\n\";",
			[]TokenType{OBJECTID, ASSIGN, INT_CONST, SEMI, OBJECTID, ASSIGN, OBJECTID, LE, STR_CONST, SEMI, EOF},
			[]string{"a", "<-", "0", ";", "b", "<-", "a", "<=", "1\n", ";", ""},
		},
		{
			"{true\n1\n\"some string\"\n}",
			[]TokenType{LBRACE, BOOL_CONST, INT_CONST, STR_CONST, RBRACE, EOF},
			[]string{"{", "true", "1", "some string", "}", ""},
		},
		{
			"{true\n1\n\"some string\"}",
			[]TokenType{LBRACE, BOOL_CONST, INT_CONST, STR_CONST, RBRACE, EOF},
			[]string{"{", "true", "1", "some string", "}", ""},
		},
		{
			"let a:A in true",
			[]TokenType{LET, OBJECTID, COLON, TYPEID, IN, BOOL_CONST, EOF},
			[]string{"let", "a", ":", "A", "in", "true", ""},
		},
		{
			"case a of b:B => false esac",
			[]TokenType{CASE, OBJECTID, OF, OBJECTID, COLON, TYPEID, DARROW, BOOL_CONST, ESAC, EOF},
			[]string{"case", "a", "of", "b", ":", "B", "=>", "false", "esac", ""},
		},
		{
			"class List inherits Object { };",
			[]TokenType{CLASS, TYPEID, INHERITS, TYPEID, LBRACE, RBRACE, SEMI, EOF},
			[]string{"class", "List", "inherits", "Object", "{", "}", ";", ""},
		},
		{
			"method(x : Int, y : String) : Bool { true };",
			[]TokenType{OBJECTID, LPAREN, OBJECTID, COLON, TYPEID, COMMA, OBJECTID, COLON, TYPEID, RPAREN, COLON, TYPEID, LBRACE, BOOL_CONST, RBRACE, SEMI, EOF},
			[]string{"method", "(", "x", ":", "Int", ",", "y", ":", "String", ")", ":", "Bool", "{", "true", "}", ";", ""},
		},
		{
			"while (*(* nested comment *)*) not isvoid x loop\n pool",
			[]TokenType{WHILE, NOT, ISVOID, OBJECTID, LOOP, POOL, EOF},
			[]string{"while", "not", "isvoid", "x", "loop", "pool", ""},
		},
		{
			"\"Hello\\tWorld\\n\"",
			[]TokenType{STR_CONST, EOF},
			[]string{"Hello\tWorld\n", ""},
		},
		{
			"if x < y then x <- y else y <- x fi",
			[]TokenType{IF, OBJECTID, LT, OBJECTID, THEN, OBJECTID, ASSIGN, OBJECTID, ELSE, OBJECTID, ASSIGN, OBJECTID, FI, EOF},
			[]string{"if", "x", "<", "y", "then", "x", "<-", "y", "else", "y", "<-", "x", "fi", ""},
		},
		{
			"1 + 2 * 3 / 4 - ~5",
			[]TokenType{INT_CONST, PLUS, INT_CONST, TIMES, INT_CONST, DIVIDE, INT_CONST, MINUS, NEG, INT_CONST, EOF},
			[]string{"1", "+", "2", "*", "3", "/", "4", "-", "~", "5", ""},
		},
		{
			"new Object.method()",
			[]TokenType{NEW, TYPEID, DOT, OBJECTID, LPAREN, RPAREN, EOF},
			[]string{"new", "Object", ".", "method", "(", ")", ""},
		},
		{
			"case x of\n y:Int => 1;\n z:String => \"hello\";\n esac",
			[]TokenType{CASE, OBJECTID, OF, OBJECTID, COLON, TYPEID, DARROW, INT_CONST, SEMI, OBJECTID, COLON, TYPEID, DARROW, STR_CONST, SEMI, ESAC, EOF},
			[]string{"case", "x", "of", "y", ":", "Int", "=>", "1", ";", "z", ":", "String", "=>", "hello", ";", "esac", ""},
		},
		{
			"x -- comment 1\n<- -- comment 2\ny -- comment 3\n",
			[]TokenType{OBJECTID, ASSIGN, OBJECTID, EOF},
			[]string{"x", "<-", "y", ""},
		},
		{
			"(* outer (* inner *) comment *) class",
			[]TokenType{CLASS, EOF},
			[]string{"class", ""},
		},
		{
			"x <= y = z < w",
			[]TokenType{OBJECTID, LE, OBJECTID, EQ, OBJECTID, LT, OBJECTID, EOF},
			[]string{"x", "<=", "y", "=", "z", "<", "w", ""},
		},
		{
			"(*(*(*(**)*)*)*)class",
			[]TokenType{CLASS, EOF},
			[]string{"class", ""},
		},
		{
			"(*(*(*(*(*(*(*(*hello*)*)*)*)*)*)*)*) world",
			[]TokenType{OBJECTID, EOF},
			[]string{"world", ""},
		},
		{
			"(* (* (* crazy (* nested (* comments *) here *) *) *) *) test",
			[]TokenType{OBJECTID, EOF},
			[]string{"test", ""},
		},
		{
			"\"This is an unterminated string...",
			[]TokenType{ERROR, EOF},
			[]string{"EOF in string constant", ""},
		},

		// Test EOF Inside Multi-Line Comment
		{
			"(* This comment never ends...",
			[]TokenType{ERROR, EOF},
			[]string{"Unterminated comment", ""},
		},
		{
			"\"Valid escape: \\x\"",
			[]TokenType{STR_CONST, EOF},
			[]string{"Valid escape: x", ""},
		},
		{
			"    \n\t   ",
			[]TokenType{EOF},
			[]string{""},
		},
		{
			"foo_bar123 <- 42;",
			[]TokenType{OBJECTID, ASSIGN, INT_CONST, SEMI, EOF},
			[]string{"foo_bar123", "<-", "42", ";", ""},
		},
		{
			"If x Then y Else z Fi",
			[]TokenType{IF, OBJECTID, THEN, OBJECTID, ELSE, OBJECTID, FI, EOF},
			[]string{"If", "x", "Then", "y", "Else", "z", "Fi", ""},
		},
		{
			"{ ( [ ] ) }",
			[]TokenType{LBRACE, LPAREN, ERROR, ERROR, RPAREN, RBRACE, EOF},
			[]string{"{", "(", "Unexpected character: [", "Unexpected character: ]", ")", "}", ""},
		},
		{
			"12345678901234567890",
			[]TokenType{ERROR, EOF},
			[]string{"Number out of range", ""},
		},

		{
			"a <- b;\n-- Comment line\nc <- d;",
			[]TokenType{OBJECTID, ASSIGN, OBJECTID, SEMI, OBJECTID, ASSIGN, OBJECTID, SEMI, EOF},
			[]string{"a", "<-", "b", ";", "c", "<-", "d", ";", ""},
		},
	}

	for _, tt := range tests {
		l := NewLexer(strings.NewReader(tt.input))
		for i, expTType := range tt.expectedTokenType {
			tok := l.NextToken()

			if tok.Type != expTType {
				t.Fatalf("[%q]: Wrong token type %d-th Token. expected=%s, got %s", tt.input, i, expTType, tok.Type)
			}

			if tok.Literal != tt.expectedLiteral[i] {
				t.Fatalf("[%q]: Wrong literal at test %d-it Token. expected=%q, got %q", tt.input, i, tt.expectedLiteral[i], tok.Literal)
			}
		}
	}
}
