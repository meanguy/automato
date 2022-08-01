package vm

import "github.com/meanguy/automato/internal/scanner/token"

type (
	parseFn func(*parser)

	parseRule struct {
		prefix     parseFn
		infix      parseFn
		precedence precedence
	}
)

//nolint:gochecknoglobals // read-only lookup table
var parseRulesTable map[token.TokenType]parseRule

//nolint:gochecknoinits // init required to sidestep circulate dependency with the parser.
// While this is ugly, having a lookup table is a major optimization. prefix/infix functions
// act as receivers for the parser methods.
func init() {
	parseRulesTable = map[token.TokenType]parseRule{
		token.LeftParen: {
			prefix:     func(p *parser) { p.grouping() },
			precedence: noPrecedence,
		},
		token.RightParen: {precedence: noPrecedence},
		token.LeftBrace:  {precedence: noPrecedence},
		token.RightBrace: {precedence: noPrecedence},
		token.Comma:      {precedence: noPrecedence},
		token.Dot:        {precedence: noPrecedence},
		token.Minus: {
			prefix:     func(p *parser) { p.unary() },
			infix:      func(p *parser) { p.binary() },
			precedence: termPrecedence,
		},
		token.Plus: {
			infix:      func(p *parser) { p.binary() },
			precedence: termPrecedence,
		},
		token.Semicolon: {precedence: noPrecedence},
		token.Slash: {
			infix:      func(p *parser) { p.binary() },
			precedence: farctorPrecedence,
		},
		token.Star: {
			infix:      func(p *parser) { p.binary() },
			precedence: farctorPrecedence,
		},
		token.Bang:         {precedence: noPrecedence},
		token.BangEqual:    {precedence: noPrecedence},
		token.Equal:        {precedence: noPrecedence},
		token.EqualEqual:   {precedence: noPrecedence},
		token.Greater:      {precedence: noPrecedence},
		token.GreaterEqual: {precedence: noPrecedence},
		token.Less:         {precedence: noPrecedence},
		token.LessEqual:    {precedence: noPrecedence},
		token.Identifier:   {precedence: noPrecedence},
		token.String:       {precedence: noPrecedence},
		token.Number: {
			prefix:     func(p *parser) { p.number() },
			precedence: noPrecedence,
		},
		token.And:    {precedence: noPrecedence},
		token.Class:  {precedence: noPrecedence},
		token.Else:   {precedence: noPrecedence},
		token.False:  {precedence: noPrecedence},
		token.For:    {precedence: noPrecedence},
		token.Fun:    {precedence: noPrecedence},
		token.If:     {precedence: noPrecedence},
		token.Nil:    {precedence: noPrecedence},
		token.Or:     {precedence: noPrecedence},
		token.Print:  {precedence: noPrecedence},
		token.Return: {precedence: noPrecedence},
		token.Super:  {precedence: noPrecedence},
		token.This:   {precedence: noPrecedence},
		token.True:   {precedence: noPrecedence},
		token.Var:    {precedence: noPrecedence},
		token.While:  {precedence: noPrecedence},
		token.Error:  {precedence: noPrecedence},
		token.EOF:    {precedence: noPrecedence},
	}
}

func getParseRule(tokenType token.TokenType) *parseRule {
	rule := parseRulesTable[tokenType]

	return &rule
}
