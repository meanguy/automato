package scanner

import (
	"fmt"

	"github.com/meanguy/automato/internal/scanner/token"
)

type Scanner struct {
	Source []byte
	Cursor int
	Start  int
	Line   int
}

func NewScanner(source string) *Scanner {
	return &Scanner{
		Source: []byte(source),
		Cursor: 0,
		Start:  0,
		Line:   1,
	}
}

//nolint:cyclop // token parsing and scanning has a high degree of branching
// by design -- actual complexity should be hidden in helper methods.
func (s *Scanner) ScanToken() token.Token {
	s.skipWhitespace()
	s.skipComments()

	s.Start = s.Cursor

	if s.isAtEnd() {
		return s.makeToken(token.EOF)
	}

	char := s.advance()

	if s.isAlpha(char) {
		return s.makeIdentifierToken()
	}

	if s.isDigit(char) {
		return s.makeNumberLiteralToken()
	}

	switch char {
	case '(':
		return s.makeToken(token.LeftParen)
	case ')':
		return s.makeToken(token.RightParen)
	case '{':
		return s.makeToken(token.LeftBrace)
	case '}':
		return s.makeToken(token.RightBrace)
	case ';':
		return s.makeToken(token.Semicolon)
	case ',':
		return s.makeToken(token.Comma)
	case '.':
		return s.makeToken(token.Dot)
	case '-':
		return s.makeToken(token.Minus)
	case '+':
		return s.makeToken(token.Plus)
	case '/':
		return s.makeToken(token.Slash)
	case '*':
		return s.makeToken(token.Star)
	case '%':
		return s.makeToken(token.Percent)
	case '!':
		if s.match('=') {
			return s.makeToken(token.BangEqual)
		} else {
			return s.makeToken(token.Bang)
		}
	case '=':
		if s.match('=') {
			return s.makeToken(token.EqualEqual)
		} else {
			return s.makeToken(token.Equal)
		}
	case '<':
		if s.match('=') {
			return s.makeToken(token.LessEqual)
		} else {
			return s.makeToken(token.Less)
		}
	case '>':
		if s.match('=') {
			return s.makeToken(token.GreaterEqual)
		} else {
			return s.makeToken(token.Greater)
		}
	case '"':
		return s.makeStringLiteralToken()
	}

	return s.makeErrorTokenf("unexpected character '%v'", string(char))
}

func (s *Scanner) advance() byte {
	s.Cursor++

	return s.Source[s.Cursor-1]
}

func (s *Scanner) isAlpha(char byte) bool {
	return ('a' <= char && char <= 'z') ||
		('A' <= char && char <= 'Z') ||
		(char == '_')
}

func (s *Scanner) isAtEnd() bool {
	return len(s.Source) <= s.Cursor
}

func (s *Scanner) isDigit(char byte) bool {
	return '0' <= char && char <= '9'
}

func (s *Scanner) match(char byte) bool {
	if s.isAtEnd() {
		return false
	}

	if s.Source[s.Cursor] != char {
		return false
	}

	s.Cursor++

	return true
}

func (s *Scanner) checkKeyword(start int, prefix string, tokenType token.TokenType) token.TokenType {
	offset := s.Start + start
	slice := string(s.Source[offset : offset+len(prefix)])
	if slice == prefix {
		return tokenType
	}

	return token.Identifier
}

//nolint:cyclop // trie data structure is inherently highly-branching
func (s *Scanner) identifierType() token.TokenType {
	switch s.Source[s.Start] {
	case 'a':
		return s.checkKeyword(1, "nd", token.And)
	case 'c':
		return s.checkKeyword(1, "lass", token.Class)
	case 'e':
		return s.checkKeyword(1, "lse", token.Else)
	case 'i':
		return s.checkKeyword(1, "f", token.If)
	case 'n':
		return s.checkKeyword(1, "il", token.Nil)
	case 'o':
		return s.checkKeyword(1, "r", token.Or)
	case 'p':
		return s.checkKeyword(1, "rint", token.Print)
	case 'r':
		return s.checkKeyword(1, "eturn", token.Return)
	case 's':
		return s.checkKeyword(1, "uper", token.Super)
	case 'v':
		return s.checkKeyword(1, "ar", token.Var)
	case 'w':
		return s.checkKeyword(1, "hile", token.While)
	case 'f':
		if len(s.Source) > 1 {
			switch s.Source[s.Start+1] {
			case 'a':
				return s.checkKeyword(2, "lse", token.False)
			case 'o':
				return s.checkKeyword(2, "r", token.For)
			case 'u':
				return s.checkKeyword(2, "n", token.Fun)
			}
		}
	case 't':
		if len(s.Source) > 1 {
			switch s.Source[s.Start+1] {
			case 'h':
				return s.checkKeyword(2, "is", token.This)
			case 'r':
				return s.checkKeyword(2, "ue", token.True)
			}
		}
	}

	return token.Identifier
}

func (s *Scanner) makeErrorTokenf(msg string, args ...any) token.Token {
	return token.Token{
		Type: token.Error,
		Str:  fmt.Sprintf(msg, args...),
		Line: s.Line,
	}
}

func (s *Scanner) makeIdentifierToken() token.Token {
	for s.isAlpha(s.peek()) || s.isDigit(s.peek()) {
		_ = s.advance()
	}

	return s.makeToken(s.identifierType())
}

func (s *Scanner) makeNumberLiteralToken() token.Token {
	for s.isDigit(s.peek()) {
		_ = s.advance()
	}

	if s.peek() == '.' && s.isDigit(s.peekNext()) {
		_ = s.advance() // consume the .

		for s.isDigit(s.peek()) {
			_ = s.advance()
		}
	}

	return s.makeToken(token.Number)
}

func (s *Scanner) makeStringLiteralToken() token.Token {
	for s.peek() != '"' && !s.isAtEnd() {
		if s.peek() == '\n' {
			s.Line++
		}

		_ = s.advance()
	}

	if s.isAtEnd() {
		return s.makeErrorTokenf("unterminated string")
	}

	_ = s.advance()

	return s.makeToken(token.String)
}

func (s *Scanner) makeToken(tokenType token.TokenType) token.Token {
	return token.Token{
		Type: tokenType,
		Str:  string(s.Source[s.Start:s.Cursor]),
		Line: s.Line,
	}
}

func (s *Scanner) peek() byte {
	if s.isAtEnd() {
		return 0
	}

	return s.Source[s.Cursor]
}

func (s *Scanner) peekNext() byte {
	if s.Cursor >= len(s.Source)-1 {
		return 0
	}

	return s.Source[s.Cursor+1]
}

func (s *Scanner) skipComments() {
	for {
		char := s.peek()
		switch char {
		case '/':
			if s.peekNext() == '/' {
				for s.peek() != '\n' && !s.isAtEnd() {
					_ = s.advance()
				}
			} else {
				return
			}
		default:
			return
		}
	}
}

func (s *Scanner) skipWhitespace() {
	for {
		char := s.peek()
		switch char {
		case ' ':
			fallthrough
		case '\r':
			fallthrough
		case '\t':
			_ = s.advance()
		case '\n':
			s.Line++
			_ = s.advance()
		default:
			return
		}
	}
}
