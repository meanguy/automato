package token

type (
	TokenType int

	Token struct {
		Type TokenType
		Str  string
		Line int
	}
)

const (
	// Single character tokens.
	LeftParen TokenType = iota + 1
	RightParen
	LeftBrace
	RightBrace
	Comma
	Dot
	Minus
	Plus
	Semicolon
	Slash
	Star
	Percent

	// 1-2 character tokens.
	Bang
	BangEqual
	Equal
	EqualEqual
	Greater
	GreaterEqual
	Less
	LessEqual

	// Literals.
	Identifier
	String
	Number

	// Keywords.
	And
	Class
	Else
	False
	For
	Fun
	If
	Nil
	Or
	Print
	Return
	Super
	This
	True
	Var
	While

	// Sentinel tokens.
	Error
	EOF
)
