package opcode

type OpCode uint8

const (
	OpReturn OpCode = iota + 1
	OpConstant
	OpConstantLong
	OpNegate
	OpAdd
	OpSubtract
	OpMultiply
	OpDivide
)
