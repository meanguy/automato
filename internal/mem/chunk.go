package mem

import (
	"github.com/meanguy/automato/internal/opcode"
	"github.com/meanguy/automato/internal/value"
)

type Chunk struct {
	Code      []uint8
	Constants []value.Value
	Lines     []int
}

func (c *Chunk) AddConstant(v value.Value) int {
	c.Constants = append(c.Constants, v)

	return len(c.Constants) - 1
}

func (c *Chunk) GetConstant(constantID int) value.Value {
	return c.Constants[constantID]
}

func (c *Chunk) Read(offset int) uint8 {
	return c.Code[offset]
}

func (c *Chunk) ReadWord(offset int) uint16 {
	major := uint16(c.Read(offset)) << 8
	minor := uint16(c.Read(offset + 1))

	return major | minor
}

func (c *Chunk) ReadOp(offset int) opcode.OpCode {
	return opcode.OpCode(c.Read(offset))
}

func (c *Chunk) Write(b uint8, line int) {
	c.Code = append(c.Code, b)
	c.Lines = append(c.Lines, line)
}

func (c *Chunk) WriteWord(w uint16, line int) {
	c.Write(uint8((w&0xff00)>>8), line)
	c.Write(uint8(w&0xff), line)
}

func (c *Chunk) WriteOp(op opcode.OpCode, line int) {
	c.Write(uint8(op), line)
}
