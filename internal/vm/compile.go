package vm

import (
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"

	"github.com/meanguy/automato/internal/debug"
	"github.com/meanguy/automato/internal/mem"
	"github.com/meanguy/automato/internal/opcode"
	"github.com/meanguy/automato/internal/scanner"
	"github.com/meanguy/automato/internal/scanner/token"
	"github.com/meanguy/automato/internal/value"
)

type (
	parser struct {
		current  token.Token
		previous token.Token
		chunk    mem.Chunk
		scan     *scanner.Scanner
		debug    bool
		fatal    bool
		err      error
	}

	precedence int
)

const (
	noPrecedence precedence = iota + 1
	assignmentPrecedence
	// orPrecedence
	// andPrecedence
	// equalityPrecedence
	// comparisonPrecedence.
	termPrecedence
	farctorPrecedence
	unaryPrecedence
	// callPrecedence
	// primaryPrecedence.
)

func newParser(scan *scanner.Scanner, debug bool) *parser {
	return &parser{
		current:  token.Token{},
		previous: token.Token{},
		chunk:    mem.Chunk{},
		scan:     scan,
		debug:    debug,
		fatal:    false,
		err:      nil,
	}
}

func (p *parser) compile() (*mem.Chunk, error) {
	p.advance()
	p.expression()
	p.consume(token.EOF, "expect end of expression")
	p.endCompiler()

	if err := p.Err(); err != nil {
		return nil, err
	}

	return &p.chunk, nil
}

func (p *parser) emitByte(b byte, line int) {
	p.chunk.Write(b, line)
}

func (p *parser) emitWord(word uint16, line int) {
	p.chunk.WriteWord(word, line)
}

func (p *parser) emitConstant(val value.Value, line int) {
	constantID := p.chunk.AddConstant(val)
	if constantID > math.MaxUint8 {
		p.emitOpCode(opcode.OpConstantLong, line)
		p.emitWord(uint16(constantID), line)
	} else {
		p.emitOpCode(opcode.OpConstant, line)
		p.emitByte(byte(constantID), line)
	}
}

func (p *parser) emitOpCode(op opcode.OpCode, line int) {
	p.chunk.WriteOp(op, line)
}

func (p *parser) endCompiler() {
	if p.debug && p.err == nil {
		debug.DisassembleChunk(os.Stderr, &p.chunk, "code")
	}

	p.emitOpCode(opcode.OpReturn, p.current.Line)
}

func (p *parser) parsePrecedence(prec precedence) {
	p.advance()

	prefix := getParseRule(p.previous.Type).prefix
	if prefix == nil {
		p.errorAtPrevious("expect expression")

		return
	}

	for prec <= getParseRule(p.current.Type).precedence {
		p.advance()

		infix := getParseRule(p.previous.Type).infix
		infix(p)
	}
}

func (p *parser) advance() {
	p.previous = p.current

	for {
		p.current = p.scan.ScanToken()

		if p.current.Type != token.Error {
			break
		}

		p.errorAtCurrent(p.current.Str)
	}
}

func (p *parser) consume(tokenType token.TokenType, msg string) {
	if p.current.Type == tokenType {
		p.advance()

		return
	}

	p.errorAtCurrent(msg)
}

func (p *parser) binary() {
	operatorType := p.previous.Type
	rule := getParseRule(operatorType)
	p.parsePrecedence(rule.precedence)

	//nolint:exhaustive // we only care about a couple of token types
	switch operatorType {
	case token.Plus:
		p.emitOpCode(opcode.OpAdd, p.previous.Line)
	case token.Minus:
		p.emitOpCode(opcode.OpSubtract, p.previous.Line)
	case token.Star:
		p.emitOpCode(opcode.OpMultiply, p.previous.Line)
	case token.Slash:
		p.emitOpCode(opcode.OpDivide, p.previous.Line)
	default:
		return
	}
}

func (p *parser) expression() {
	p.parsePrecedence(assignmentPrecedence)
}

func (p *parser) grouping() {
	p.expression()

	p.consume(token.RightParen, "expected ')' after expression")
}

func (p *parser) number() {
	val, err := strconv.ParseFloat(p.previous.Str, 64)
	if err != nil {
		p.errorAtCurrent("failed to parse number '%s': %v", p.previous.Str, err)
	}

	p.emitConstant(value.Value(val), p.previous.Line)
}

func (p *parser) unary() {
	operatorType := p.previous.Type

	p.parsePrecedence(unaryPrecedence)

	//nolint:exhaustive // we only care about a couple of token types
	switch operatorType {
	case token.Minus:
		p.emitOpCode(opcode.OpNegate, p.previous.Line)
	default:
	}
}

func (p *parser) errorAtCurrent(msg string, args ...any) {
	p.errorAtToken(p.current, msg, args...)
}

func (p *parser) errorAtPrevious(msg string, args ...any) {
	p.errorAtToken(p.previous, msg, args...)
}

func (p *parser) errorAtToken(t token.Token, msg string, args ...any) {
	if p.fatal { // skip reporting errors if we're already handling an error
		return
	}

	p.fatal = true

	var buf strings.Builder

	fmt.Fprintf(&buf, "[line %d] error", t.Line)

	//nolint:exhaustive // we only care about a couple of token types
	switch t.Type {
	case token.EOF:
		fmt.Fprintf(os.Stderr, " at end")
	case token.Error:
		break
	default:
		fmt.Fprintf(os.Stderr, " at '%v'", fmt.Sprintf(msg, args...))
	}

	fmt.Fprintf(&buf, ": ")
	fmt.Fprintf(&buf, msg, args...)

	p.err = fmt.Errorf(buf.String())
	fmt.Fprint(os.Stderr, p.err.Error())
}

func (p *parser) Err() error {
	return p.err
}
