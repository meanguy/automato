package scanner_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/meanguy/automato/internal/scanner"
	"github.com/meanguy/automato/internal/scanner/token"
)

func TestScanToken(t *testing.T) {
	testCases := []struct {
		text     string
		expected token.TokenType
	}{
		{"", token.EOF},
		{"// comment", token.EOF},
		{"\n", token.EOF},
		{"   \t", token.EOF},
		{"\nfoo", token.Identifier},
		{`"nonterminating`, token.Error},
		{`....`, token.Dot},
		{`"foo"`, token.String},
		{`"foo" + "bar"`, token.String},
		{"200", token.Number},
		{"1.23", token.Number},
		{"foo", token.Identifier},
		{"false", token.False},
		{"super", token.Super},
		{"return", token.Return},
		{"(", token.LeftParen},
		{";", token.Semicolon},
		{"/", token.Slash},
		{"-", token.Minus},
		{"!=", token.BangEqual},
		{"=", token.Equal},
		{"==", token.EqualEqual},
		{">=", token.GreaterEqual},
	}

	for index, tc := range testCases {
		t.Run(fmt.Sprintf("%d - %s", index, tc.text), func(t *testing.T) {
			scan := scanner.NewScanner(tc.text)
			assert.Equal(t, tc.expected, scan.ScanToken().Type)
		})
	}
}

func TestScanTokenMultipleTokens(t *testing.T) {
	testCases := []struct {
		text     string
		expected []token.Token
	}{
		{
			text: "5",
			expected: []token.Token{
				{Type: token.Number, Line: 1, Str: "5"},
			},
		},
		{
			text: "3+2",
			expected: []token.Token{
				{Type: token.Number, Line: 1, Str: "3"},
				{Type: token.Plus, Line: 1, Str: "+"},
				{Type: token.Number, Line: 1, Str: "2"},
			},
		},
		{
			text: "8 >= 5",
			expected: []token.Token{
				{Type: token.Number, Line: 1, Str: "8"},
				{Type: token.GreaterEqual, Line: 1, Str: ">="},
				{Type: token.Number, Line: 1, Str: "5"},
			},
		},
		{
			text: "if (foobar == 5) {\n\treturn \"hello\"\n}",
			expected: []token.Token{
				{Type: token.If, Line: 1, Str: "if"},
				{Type: token.LeftParen, Line: 1, Str: "("},
				{Type: token.Identifier, Line: 1, Str: "foobar"},
				{Type: token.EqualEqual, Line: 1, Str: "=="},
				{Type: token.Number, Line: 1, Str: "5"},
				{Type: token.RightParen, Line: 1, Str: ")"},
				{Type: token.LeftBrace, Line: 1, Str: "{"},
				{Type: token.Return, Line: 2, Str: "return"},
				{Type: token.String, Line: 2, Str: "\"hello\""},
				{Type: token.RightBrace, Line: 3, Str: "}"},
			},
		},
		{
			text: "print(\"hello\", \"world\")",
			expected: []token.Token{
				{Type: token.Print, Line: 1, Str: "print"},
				{Type: token.LeftParen, Line: 1, Str: "("},
				{Type: token.String, Line: 1, Str: "\"hello\""},
				{Type: token.Comma, Line: 1, Str: ","},
				{Type: token.String, Line: 1, Str: "\"world\""},
				{Type: token.RightParen, Line: 1, Str: ")"},
			},
		},
		{
			text: "fun fib(n) {\n\tif n <= 1 {\n\t\treturn n\n\t} else {\n\t\treturn fib(n-1) + fib(n-2)\n\t}\n}",
			expected: []token.Token{
				{Type: token.Fun, Line: 1, Str: "fun"},
				{Type: token.Identifier, Line: 1, Str: "fib"},
				{Type: token.LeftParen, Line: 1, Str: "("},
				{Type: token.Identifier, Line: 1, Str: "n"},
				{Type: token.RightParen, Line: 1, Str: ")"},
				{Type: token.LeftBrace, Line: 1, Str: "{"},
				{Type: token.If, Line: 2, Str: "if"},
				{Type: token.Identifier, Line: 2, Str: "n"},
				{Type: token.LessEqual, Line: 2, Str: "<="},
				{Type: token.Number, Line: 2, Str: "1"},
				{Type: token.LeftBrace, Line: 2, Str: "{"},
				{Type: token.Return, Line: 3, Str: "return"},
				{Type: token.Identifier, Line: 3, Str: "n"},
				{Type: token.RightBrace, Line: 4, Str: "}"},
				{Type: token.Else, Line: 4, Str: "else"},
				{Type: token.LeftBrace, Line: 4, Str: "{"},
				{Type: token.Return, Line: 5, Str: "return"},
				{Type: token.Identifier, Line: 5, Str: "fib"},
				{Type: token.LeftParen, Line: 5, Str: "("},
				{Type: token.Identifier, Line: 5, Str: "n"},
				{Type: token.Minus, Line: 5, Str: "-"},
				{Type: token.Number, Line: 5, Str: "1"},
				{Type: token.RightParen, Line: 5, Str: ")"},
				{Type: token.Plus, Line: 5, Str: "+"},
				{Type: token.Identifier, Line: 5, Str: "fib"},
				{Type: token.LeftParen, Line: 5, Str: "("},
				{Type: token.Identifier, Line: 5, Str: "n"},
				{Type: token.Minus, Line: 5, Str: "-"},
				{Type: token.Number, Line: 5, Str: "2"},
				{Type: token.RightParen, Line: 5, Str: ")"},
				{Type: token.RightBrace, Line: 6, Str: "}"},
				{Type: token.RightBrace, Line: 7, Str: "}"},
			},
		},
		{
			text: "while n < 10 {\n\tvar t = n % 3\n\tn = n + t\n}",
			expected: []token.Token{
				{Type: token.While, Line: 1, Str: "while"},
				{Type: token.Identifier, Line: 1, Str: "n"},
				{Type: token.Less, Line: 1, Str: "<"},
				{Type: token.Number, Line: 1, Str: "10"},
				{Type: token.LeftBrace, Line: 1, Str: "{"},
				{Type: token.Var, Line: 2, Str: "var"},
				{Type: token.Identifier, Line: 2, Str: "t"},
				{Type: token.Equal, Line: 2, Str: "="},
				{Type: token.Identifier, Line: 2, Str: "n"},
				{Type: token.Percent, Line: 2, Str: "%"},
				{Type: token.Number, Line: 2, Str: "3"},
				{Type: token.Identifier, Line: 3, Str: "n"},
				{Type: token.Equal, Line: 3, Str: "="},
				{Type: token.Identifier, Line: 3, Str: "n"},
				{Type: token.Plus, Line: 3, Str: "+"},
				{Type: token.Identifier, Line: 3, Str: "t"},
				{Type: token.RightBrace, Line: 4, Str: "}"},
			},
		},
		{
			text: "class Foo {}",
			expected: []token.Token{
				{Type: token.Class, Line: 1, Str: "class"},
				{Type: token.Identifier, Line: 1, Str: "Foo"},
				{Type: token.LeftBrace, Line: 1, Str: "{"},
				{Type: token.RightBrace, Line: 1, Str: "}"},
			},
		},
		{
			text: "service != nil and !(service.has_error() or service.waiting())",
			expected: []token.Token{
				{Type: token.Identifier, Line: 1, Str: "service"},
				{Type: token.BangEqual, Line: 1, Str: "!="},
				{Type: token.Nil, Line: 1, Str: "nil"},
				{Type: token.And, Line: 1, Str: "and"},
				{Type: token.Bang, Line: 1, Str: "!"},
				{Type: token.LeftParen, Line: 1, Str: "("},
				{Type: token.Identifier, Line: 1, Str: "service"},
				{Type: token.Dot, Line: 1, Str: "."},
				{Type: token.Identifier, Line: 1, Str: "has_error"},
				{Type: token.LeftParen, Line: 1, Str: "("},
				{Type: token.RightParen, Line: 1, Str: ")"},
				{Type: token.Or, Line: 1, Str: "or"},
				{Type: token.Identifier, Line: 1, Str: "service"},
				{Type: token.Dot, Line: 1, Str: "."},
				{Type: token.Identifier, Line: 1, Str: "waiting"},
				{Type: token.LeftParen, Line: 1, Str: "("},
				{Type: token.RightParen, Line: 1, Str: ")"},
			},
		},
		{
			text: "this.store = super.get_data_store()",
			expected: []token.Token{
				{Type: token.This, Line: 1, Str: "this"},
				{Type: token.Dot, Line: 1, Str: "."},
				{Type: token.Identifier, Line: 1, Str: "store"},
				{Type: token.Equal, Line: 1, Str: "="},
				{Type: token.Super, Line: 1, Str: "super"},
				{Type: token.Dot, Line: 1, Str: "."},
				{Type: token.Identifier, Line: 1, Str: "get_data_store"},
				{Type: token.LeftParen, Line: 1, Str: "("},
				{Type: token.RightParen, Line: 1, Str: ")"},
			},
		},
		{
			text: "if true { return 2 > 1 * 1 }",
			expected: []token.Token{
				{Type: token.If, Line: 1, Str: "if"},
				{Type: token.True, Line: 1, Str: "true"},
				{Type: token.LeftBrace, Line: 1, Str: "{"},
				{Type: token.Return, Line: 1, Str: "return"},
				{Type: token.Number, Line: 1, Str: "2"},
				{Type: token.Greater, Line: 1, Str: ">"},
				{Type: token.Number, Line: 1, Str: "1"},
				{Type: token.Star, Line: 1, Str: "*"},
				{Type: token.Number, Line: 1, Str: "1"},
				{Type: token.RightBrace, Line: 1, Str: "}"},
			},
		},
		{
			text: "$",
			expected: []token.Token{
				{Type: token.Error, Line: 1, Str: "unexpected character '$'"},
			},
		},
		{
			text: "\"multiline\nstring\"",
			expected: []token.Token{
				{Type: token.String, Line: 2, Str: "\"multiline\nstring\""},
			},
		},
	}

	for index, tc := range testCases {
		t.Run(fmt.Sprintf("%d - %s", index, tc.text), func(t *testing.T) {
			scan := scanner.NewScanner(tc.text)

			actual := make([]token.Token, 0, len(tc.expected))
			for i := 0; i < len(tc.expected); i++ {
				actual = append(actual, scan.ScanToken())
			}

			assert.Equal(t, len(tc.expected), len(actual))
			for i := 0; i < len(tc.expected); i++ {
				assert.Equal(t, tc.expected[i], actual[i])
			}
		})
	}
}
