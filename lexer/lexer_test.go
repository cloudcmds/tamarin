package lexer

import (
	"errors"
	"fmt"
	"testing"

	"github.com/cloudcmds/tamarin/token"
	"github.com/stretchr/testify/require"
)

func TestNull(t *testing.T) {
	input := "a = null;"
	tests := []struct {
		expectedType    token.Type
		expectedLiteral string
	}{
		{token.IDENT, "a"},
		{token.ASSIGN, "="},
		{token.NULL, "null"},
		{token.SEMICOLON, ";"},
		{token.EOF, ""},
	}
	l := New(input)
	for i, tt := range tests {
		tok, err := l.NextToken()
		require.Nil(t, err)
		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong, expected=%q, got=%q", i, tt.expectedType, tok.Type)
		}
		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - Literal wrong, expected=%q, got=%q", i, tt.expectedLiteral, tok.Literal)
		}
	}
}

func TestNextToken1(t *testing.T) {
	input := "%=+(){},;?|| &&`/foo`++--***=.."

	tests := []struct {
		expectedType    token.Type
		expectedLiteral string
	}{
		{token.MOD, "%"},
		{token.ASSIGN, "="},
		{token.PLUS, "+"},
		{token.LPAREN, "("},
		{token.RPAREN, ")"},
		{token.LBRACE, "{"},
		{token.RBRACE, "}"},
		{token.COMMA, ","},
		{token.SEMICOLON, ";"},
		{token.QUESTION, "?"},
		{token.OR, "||"},
		{token.AND, "&&"},
		{token.BACKTICK, "/foo"},
		{token.PLUS_PLUS, "++"},
		{token.MINUS_MINUS, "--"},
		{token.POW, "**"},
		{token.ASTERISK_EQUALS, "*="},
		{token.PERIOD, "."},
		{token.PERIOD, "."},
		{token.EOF, ""},
	}
	l := New(input)
	for i, tt := range tests {
		tok, err := l.NextToken()
		require.Nil(t, err)
		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong, expected=%q, got=%q", i, tt.expectedType, tok.Type)
		}
		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - Literal wrong, expected=%q, got=%q", i, tt.expectedLiteral, tok.Literal)
		}
	}
}

func TestNextToken2(t *testing.T) {
	input := `let five=5;
let ten =10;
let add = func(x, y){
  x+y
};
let result = add(five, ten);
!- *5;
5<10>5;

if(5<10){
	return true;
}else{
	return false;
}
10 == 10;
10 != 9;
"foobar"
"foo bar"
[1,2];
{"foo":"bar"}
1.2
0.5
0.3
世界
for
2 >= 1
1 <= 3
break
`
	tests := []struct {
		expectedType    token.Type
		expectedLiteral string
	}{
		{token.LET, "let"},
		{token.IDENT, "five"},
		{token.ASSIGN, "="},
		{token.INT, "5"},
		{token.SEMICOLON, ";"},
		{token.NEWLINE, "\n"},
		{token.LET, "let"},
		{token.IDENT, "ten"},
		{token.ASSIGN, "="},
		{token.INT, "10"},
		{token.SEMICOLON, ";"},
		{token.NEWLINE, "\n"},
		{token.LET, "let"},
		{token.IDENT, "add"},
		{token.ASSIGN, "="},
		{token.FUNC, "func"},
		{token.LPAREN, "("},
		{token.IDENT, "x"},
		{token.COMMA, ","},
		{token.IDENT, "y"},
		{token.RPAREN, ")"},
		{token.LBRACE, "{"},
		{token.NEWLINE, "\n"},
		{token.IDENT, "x"},
		{token.PLUS, "+"},
		{token.IDENT, "y"},
		{token.NEWLINE, "\n"},
		{token.RBRACE, "}"},
		{token.SEMICOLON, ";"},
		{token.NEWLINE, "\n"},
		{token.LET, "let"},
		{token.IDENT, "result"},
		{token.ASSIGN, "="},
		{token.IDENT, "add"},
		{token.LPAREN, "("},
		{token.IDENT, "five"},
		{token.COMMA, ","},
		{token.IDENT, "ten"},
		{token.RPAREN, ")"},
		{token.SEMICOLON, ";"},
		{token.NEWLINE, "\n"},
		{token.BANG, "!"},
		{token.MINUS, "-"},
		{token.ASTERISK, "*"},
		{token.INT, "5"},
		{token.SEMICOLON, ";"},
		{token.NEWLINE, "\n"},
		{token.INT, "5"},
		{token.LT, "<"},
		{token.INT, "10"},
		{token.GT, ">"},
		{token.INT, "5"},
		{token.SEMICOLON, ";"},
		{token.NEWLINE, "\n"},
		{token.NEWLINE, "\n"},
		{token.IF, "if"},
		{token.LPAREN, "("},
		{token.INT, "5"},
		{token.LT, "<"},
		{token.INT, "10"},
		{token.RPAREN, ")"},
		{token.LBRACE, "{"},
		{token.NEWLINE, "\n"},
		{token.RETURN, "return"},
		{token.TRUE, "true"},
		{token.SEMICOLON, ";"},
		{token.NEWLINE, "\n"},
		{token.RBRACE, "}"},
		{token.ELSE, "else"},
		{token.LBRACE, "{"},
		{token.NEWLINE, "\n"},
		{token.RETURN, "return"},
		{token.FALSE, "false"},
		{token.SEMICOLON, ";"},
		{token.NEWLINE, "\n"},
		{token.RBRACE, "}"},
		{token.NEWLINE, "\n"},
		{token.INT, "10"},
		{token.EQ, "=="},
		{token.INT, "10"},
		{token.SEMICOLON, ";"},
		{token.NEWLINE, "\n"},
		{token.INT, "10"},
		{token.NOT_EQ, "!="},
		{token.INT, "9"},
		{token.SEMICOLON, ";"},
		{token.NEWLINE, "\n"},
		{token.STRING, "foobar"},
		{token.NEWLINE, "\n"},
		{token.STRING, "foo bar"},
		{token.NEWLINE, "\n"},
		{token.LBRACKET, "["},
		{token.INT, "1"},
		{token.COMMA, ","},
		{token.INT, "2"},
		{token.RBRACKET, "]"},
		{token.SEMICOLON, ";"},
		{token.NEWLINE, "\n"},
		{token.LBRACE, "{"},
		{token.STRING, "foo"},
		{token.COLON, ":"},
		{token.STRING, "bar"},
		{token.RBRACE, "}"},
		{token.NEWLINE, "\n"},
		{token.FLOAT, "1.2"},
		{token.NEWLINE, "\n"},
		{token.FLOAT, "0.5"},
		{token.NEWLINE, "\n"},
		{token.FLOAT, "0.3"},
		{token.NEWLINE, "\n"},
		{token.IDENT, "世界"},
		{token.NEWLINE, "\n"},
		{token.FOR, "for"},
		{token.NEWLINE, "\n"},
		{token.INT, "2"},
		{token.GT_EQUALS, ">="},
		{token.INT, "1"},
		{token.NEWLINE, "\n"},
		{token.INT, "1"},
		{token.LT_EQUALS, "<="},
		{token.INT, "3"},
		{token.NEWLINE, "\n"},
		{token.BREAK, "break"},
		{token.NEWLINE, "\n"},
		{token.EOF, ""},
	}
	l := New(input)
	for i, tt := range tests {
		tok, err := l.NextToken()
		require.Nil(t, err)
		if tok.Type != tt.expectedType {
			fmt.Println(tok.Literal)
			t.Fatalf("tests[%d] - tokentype wrong, expected=%q, got=%q", i, tt.expectedType, tok.Type)
		}
		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - Literal wrong, expected=%q, got=%q", i, tt.expectedLiteral, tok.Literal)
		}
	}
}

func TestUnicodeLexer(t *testing.T) {
	input := `世界`
	l := New(input)
	tok, err := l.NextToken()
	require.Nil(t, err)
	if tok.Type != token.IDENT {
		t.Fatalf("token type wrong, expected=%q, got=%q", token.IDENT, tok.Type)
	}
	if tok.Literal != "世界" {
		t.Fatalf("token literal wrong, expected=%q, got=%q", "世界", tok.Literal)
	}
}

func TestString(t *testing.T) {
	input := `"\n\r\t\\\""`

	tests := []struct {
		expectedType    token.Type
		expectedLiteral string
	}{
		{token.STRING, "\n\r\t\\\""},
		{token.EOF, ""},
	}
	l := New(input)
	for i, tt := range tests {
		tok, err := l.NextToken()
		require.Nil(t, err)
		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong, expected=%q, got=%q", i, tt.expectedType, tok.Type)
		}
		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - Literal wrong, expected=%q, got=%q", i, tt.expectedLiteral, tok.Literal)
		}
	}

}
func TestSimpleComment(t *testing.T) {
	input := `=+// This is a comment
// This is still a comment
# I like comments
let a = 1; # This is a comment too.
// This is a final
// comment on two-lines`

	tests := []struct {
		expectedType    token.Type
		expectedLiteral string
	}{
		{token.ASSIGN, "="},
		{token.PLUS, "+"},
		{token.NEWLINE, "\n"},
		{token.NEWLINE, "\n"},
		{token.NEWLINE, "\n"},
		{token.LET, "let"},
		{token.IDENT, "a"},
		{token.ASSIGN, "="},
		{token.INT, "1"},
		{token.SEMICOLON, ";"},
		{token.NEWLINE, "\n"},
		{token.NEWLINE, "\n"},
		{token.EOF, ""},
	}
	l := New(input)
	for i, tt := range tests {
		tok, err := l.NextToken()
		require.Nil(t, err)
		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong, expected=%q, got=%q", i, tt.expectedType, tok.Type)
		}
		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - Literal wrong, expected=%q, got=%q", i, tt.expectedLiteral, tok.Literal)
		}
	}
}

func TestMultiLineComment(t *testing.T) {
	input := `=+/* This is a comment

We're still in a comment
let c = 2; */
let a = 1;
// This isa comment
// This is still a comment.
/* Now a multi-line again
   Which is two-lines
 */`

	tests := []struct {
		expectedType    token.Type
		expectedLiteral string
	}{
		{token.ASSIGN, "="},
		{token.PLUS, "+"},
		{token.NEWLINE, "\n"},
		{token.LET, "let"},
		{token.IDENT, "a"},
		{token.ASSIGN, "="},
		{token.INT, "1"},
		{token.SEMICOLON, ";"},
		{token.NEWLINE, "\n"},
		{token.NEWLINE, "\n"},
		{token.NEWLINE, "\n"},
		{token.EOF, ""},
	}
	l := New(input)
	for i, tt := range tests {
		tok, err := l.NextToken()
		require.Nil(t, err)
		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong, expected=%q, got=%q", i, tt.expectedType, tok.Type)
		}
		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - Literal wrong, expected=%q, got=%q", i, tt.expectedLiteral, tok.Literal)
		}
	}
}

func TestIntegers(t *testing.T) {
	input := `10 0x10 0xF0 0xFE 0b0101 0xFF 0b101 0xFF;`

	tests := []struct {
		expectedType    token.Type
		expectedLiteral string
	}{
		{token.INT, "10"},
		{token.INT, "0x10"},
		{token.INT, "0xF0"},
		{token.INT, "0xFE"},
		{token.INT, "0b0101"},
		{token.INT, "0xFF"},
		{token.INT, "0b101"},
		{token.INT, "0xFF"},
		{token.SEMICOLON, ";"},
		{token.EOF, ""},
	}
	l := New(input)
	for i, tt := range tests {
		tok, err := l.NextToken()
		require.Nil(t, err)
		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong, expected=%q, got=%q", i, tt.expectedType, tok.Type)
		}
		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - Literal wrong, expected=%q, got=%q", i, tt.expectedLiteral, tok.Literal)
		}
	}
}

// Test that the shebang-line is handled specially.
func TestShebang(t *testing.T) {
	input := `#!/bin/monkey
10;`

	tests := []struct {
		expectedType    token.Type
		expectedLiteral string
	}{
		{token.NEWLINE, "\n"},
		{token.INT, "10"},
		{token.SEMICOLON, ";"},
		{token.EOF, ""},
	}
	l := New(input)
	for i, tt := range tests {
		tok, err := l.NextToken()
		require.Nil(t, err)
		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong, expected=%q, got=%q", i, tt.expectedType, tok.Type)
		}
		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - Literal wrong, expected=%q, got=%q", i, tt.expectedLiteral, tok.Literal)
		}
	}
}

// TestMoreHandling does nothing real, but it bumps our coverage!
func TestMoreHandling(t *testing.T) {
	input := `#!/bin/monkey
1 += 1;
2 -= 2;
3 /= 3;
x */ 3;

let t = true;
let f = false;

if ( t && f ) { puts( "What?" ); }
if ( t || f ) { puts( "What?" ); }

let a = 1;
a++;

let b = a % 1;
b--;
b -= 2;

if ( a<3 ) { puts( "Blah!"); }
if ( a>3 ) { puts( "Blah!"); }

let b = 3;
b**b;
b *= 3;
if ( b <= 3  ) { puts "blah\n" }
if ( b >= 3  ) { puts "blah\n" }

let a = "steve";
let a = "steve\n";
let a = "steve\t";
let a = "steve\r";
let a = "steve\\";
let a = "steve\"";
let c = 3.113;
.;`
	l := New(input)
	tok, _ := l.NextToken()
	for tok.Type != token.EOF {
		tok, _ = l.NextToken()
	}
}

// TestDotMethod ensures that identifiers are parsed correctly for the
// case where we need to split at periods.
func TestDotMethod(t *testing.T) {
	input := `
foo.bar();
baz.qux();
`

	tests := []struct {
		expectedType    token.Type
		expectedLiteral string
	}{
		{token.NEWLINE, "\n"},
		{token.IDENT, "foo"},
		{token.PERIOD, "."},
		{token.IDENT, "bar"},
		{token.LPAREN, "("},
		{token.RPAREN, ")"},
		{token.SEMICOLON, ";"},
		{token.NEWLINE, "\n"},
		{token.IDENT, "baz"},
		{token.PERIOD, "."},
		{token.IDENT, "qux"},
		{token.LPAREN, "("},
		{token.RPAREN, ")"},
		{token.SEMICOLON, ";"},
		{token.NEWLINE, "\n"},
		{token.EOF, ""},
	}
	l := New(input)
	for i, tt := range tests {
		tok, err := l.NextToken()
		require.Nil(t, err)
		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong, expected=%q, got=%q", i, tt.expectedType, tok.Type)
		}
		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - Literal wrong, expected=%q, got=%q", i, tt, tok)
		}
	}
}

// TestRegexp ensures a simple regexp can be parsed.
func TestRegexp(t *testing.T) {
	input := `if ( f ~= /steve/i )
if ( f ~= /steve/m )
if ( f ~= /steve/mi )
if ( f !~ /steve/mi )
if ( f ~= /steve/miiiiiiiiiiiiiiiiimmmmmmmmmmmmmiiiii )`

	tests := []struct {
		expectedType    token.Type
		expectedLiteral string
	}{
		{token.IF, "if"},
		{token.LPAREN, "("},
		{token.IDENT, "f"},
		{token.CONTAINS, "~="},
		{token.REGEXP, "(?i)steve"},
		{token.RPAREN, ")"},
		{token.NEWLINE, "\n"},
		{token.IF, "if"},
		{token.LPAREN, "("},
		{token.IDENT, "f"},
		{token.CONTAINS, "~="},
		{token.REGEXP, "(?m)steve"},
		{token.RPAREN, ")"},
		{token.NEWLINE, "\n"},
		{token.IF, "if"},
		{token.LPAREN, "("},
		{token.IDENT, "f"},
		{token.CONTAINS, "~="},
		{token.REGEXP, "(?mi)steve"},
		{token.RPAREN, ")"},
		{token.NEWLINE, "\n"},
		{token.IF, "if"},
		{token.LPAREN, "("},
		{token.IDENT, "f"},
		{token.NOT_CONTAINS, "!~"},
		{token.REGEXP, "(?mi)steve"},
		{token.RPAREN, ")"},
		{token.NEWLINE, "\n"},
		{token.IF, "if"},
		{token.LPAREN, "("},
		{token.IDENT, "f"},
		{token.CONTAINS, "~="},
		{token.REGEXP, "(?mi)steve"},
		{token.RPAREN, ")"},
		{token.EOF, ""},
	}
	l := New(input)
	for i, tt := range tests {
		tok, _ := l.NextToken()
		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong, expected=%q, got=%q", i, tt.expectedType, tok.Type)
		}
		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - Literal wrong, expected=%q, got=%q", i, tt.expectedLiteral, tok.Literal)
		}
	}
}

// TestIllegalRegexp is designed to look for an unterminated/illegal regexp
func TestIllegalRegexp(t *testing.T) {
	input := `if f ~= /steve`

	tests := []struct {
		expectedType    token.Type
		expectedLiteral string
		expectedError   error
	}{
		{token.IF, "if", nil},
		{token.IDENT, "f", nil},
		{token.CONTAINS, "~=", nil},
		{token.REGEXP, "steve", errors.New("unterminated regular expression")},
		{token.EOF, "", nil},
	}
	l := New(input)
	for _, tt := range tests {
		tok, err := l.NextToken()
		require.Equal(t, tt.expectedError, err)
		require.Equal(t, tt.expectedType, tok.Type)
		require.Equal(t, tt.expectedLiteral, tok.Literal)
	}
}

// TestDiv is designed to test that a division is recognized; that it is
// not confused with a regular-expression.
func TestDiv(t *testing.T) {
	input := `a = b / c;
a = 3/4;`

	tests := []struct {
		expectedType    token.Type
		expectedLiteral string
	}{
		{token.IDENT, "a"},
		{token.ASSIGN, "="},
		{token.IDENT, "b"},
		{token.SLASH, "/"},
		{token.IDENT, "c"},
		{token.SEMICOLON, ";"},
		{token.NEWLINE, "\n"},
		{token.IDENT, "a"},
		{token.ASSIGN, "="},
		{token.INT, "3"},
		{token.SLASH, "/"},
		{token.INT, "4"},
		{token.SEMICOLON, ";"},
		{token.EOF, ""},
	}
	l := New(input)
	for i, tt := range tests {
		tok, _ := l.NextToken()
		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong, expected=%q, got=%q", i, tt.expectedType, tok.Type)
		}
		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - Literal wrong, expected=%q, got=%q", i, tt.expectedLiteral, tok.Literal)
		}
	}
}

func TestLineNumbers(t *testing.T) {
	l := New("ab + cd\n foo+=111")
	tests := []struct {
		expectedType     token.Type
		expectedLiteral  string
		expectedLine     int
		expectedStartPos int
		expectedEndPos   int
	}{
		{token.IDENT, "ab", 0, 0, 1},
		{token.PLUS, "+", 0, 3, 3},
		{token.IDENT, "cd", 0, 5, 6},
		{token.NEWLINE, "\n", 0, 7, 7},
		{token.IDENT, "foo", 1, 1, 3},
		{token.PLUS_EQUALS, "+=", 1, 4, 5},
		{token.INT, "111", 1, 6, 8},
		{token.EOF, "", 1, 9, 9},
	}
	for i, tt := range tests {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			tok, err := l.NextToken()
			require.Nil(t, err)
			require.Equal(t, tt.expectedType, tok.Type)
			require.Equal(t, tt.expectedLiteral, tok.Literal)
			// require.Equal(t, tt.expectedLine, tok.Line) // FIXME
			require.Equal(t, tt.expectedStartPos, tok.StartPosition.Column)
			require.Equal(t, tt.expectedEndPos, tok.EndPosition.Column)
		})
	}
}

func TestTokenLengths(t *testing.T) {
	tests := []struct {
		input            string
		expectedType     token.Type
		expectedLiteral  string
		expectedLine     int
		expectedStartPos int
		expectedEndPos   int
	}{
		{"abc", token.IDENT, "abc", 0, 0, 2},
		{"111", token.INT, "111", 0, 0, 2},
		{"1.1", token.FLOAT, "1.1", 0, 0, 2},
		{`"b"`, token.STRING, "b", 0, 0, 2},
		{"for", token.FOR, "for", 0, 0, 2},
		{"let", token.LET, "let", 0, 0, 2},
		{"false", token.FALSE, "false", 0, 0, 4},
		{"import", token.IMPORT, "import", 0, 0, 5},
		{">=", token.GT_EQUALS, ">=", 0, 0, 1},
		{" \n", token.NEWLINE, "\n", 0, 1, 1},
		{" {", token.LBRACE, "{", 0, 1, 1},
		{" ++", token.PLUS_PLUS, "++", 0, 1, 2},
	}
	for i, tt := range tests {
		t.Run(fmt.Sprintf("%d-%s", i, tt.input), func(t *testing.T) {
			l := New(tt.input)
			tok, err := l.NextToken()
			require.Nil(t, err)
			require.Equal(t, tt.expectedType, tok.Type)
			require.Equal(t, tt.expectedLiteral, tok.Literal)
			// require.Equal(t, tt.expectedLine, tok.Line) // FIXME
			require.Equal(t, tt.expectedStartPos, tok.StartPosition.Column)
			require.Equal(t, tt.expectedEndPos, tok.EndPosition.Column)
		})
	}
}

func TestStringTypes(t *testing.T) {
	tests := []struct {
		input           string
		expectedType    token.Type
		expectedLiteral string
	}{
		{`"\"foo'"`, token.STRING, "\"foo'"},
		{`'"foo\''`, token.STRING, "\"foo'"},
		{"`foo`", token.BACKTICK, "foo"},
		{"\"\\nhey\"", token.STRING, "\nhey"},
	}
	for i, tt := range tests {
		t.Run(fmt.Sprintf("%d-%s", i, tt.input), func(t *testing.T) {
			l := New(tt.input)
			tok, err := l.NextToken()
			require.Nil(t, err)
			require.Equal(t, tt.expectedType, tok.Type)
			require.Equal(t, tt.expectedLiteral, tok.Literal)
		})
	}
}
