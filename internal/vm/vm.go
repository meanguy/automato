package vm

import (
	"fmt"
	"os"

	"github.com/meanguy/automato/internal/debug"
	"github.com/meanguy/automato/internal/mem"
	"github.com/meanguy/automato/internal/opcode"
	"github.com/meanguy/automato/internal/scanner"
	"github.com/meanguy/automato/internal/value"
)

type (
	VM struct {
		Debug bool
		IP    int
		Chunk *mem.Chunk
		Stack []value.Value
	}

	VMOption func(*VM)
)

func NewVM(opts ...VMOption) *VM {
	vm := &VM{
		Debug: false,
		Chunk: nil,
		IP:    0,
	}

	for _, fn := range opts {
		fn(vm)
	}

	return vm
}

func EnableDebug() VMOption {
	return func(v *VM) {
		v.Debug = true
	}
}

func (v *VM) Interpret(source string) error {
	parser := newParser(scanner.NewScanner(source), v.Debug)
	chunk, err := parser.compile()
	if err != nil {
		return err
	}

	return v.InterpretChunk(chunk)
}

func (v *VM) InterpretChunk(chunk *mem.Chunk) error {
	v.Chunk = chunk
	v.IP = 0

	return v.run()
}

func (v *VM) Push(val value.Value) {
	v.Stack = append(v.Stack, val)
}

func (v *VM) Pop() value.Value {
	val := v.Stack[len(v.Stack)-1]
	v.Stack = v.Stack[:len(v.Stack)-1]

	return val
}

//nolint:cyclop // interpreting opcodes is necessarily complex
func (v *VM) run() error {
	for {
		instruction := opcode.OpCode(v.readByte())

		if v.Debug {
			debug.DisassembleStack(os.Stderr, v.Stack)
			debug.DisassembleInstruction(os.Stderr, v.Chunk, v.IP-1)
		}

		switch instruction {
		case opcode.OpReturn:
			fmt.Fprintf(os.Stderr, "%v\n", v.Pop())

			return nil
		case opcode.OpConstant:
			v.Push(v.Chunk.GetConstant(int(v.readByte())))
		case opcode.OpConstantLong:
			v.Push(v.Chunk.GetConstant(int(v.readWord())))
		case opcode.OpNegate:
			v.Push(-v.Pop())
		case opcode.OpAdd:
			v.binaryOp(func(lhs value.Value, rhs value.Value) value.Value { return lhs + rhs })
		case opcode.OpSubtract:
			v.binaryOp(func(lhs value.Value, rhs value.Value) value.Value { return lhs - rhs })
		case opcode.OpMultiply:
			v.binaryOp(func(lhs value.Value, rhs value.Value) value.Value { return lhs * rhs })
		case opcode.OpDivide:
			v.binaryOp(func(lhs value.Value, rhs value.Value) value.Value { return lhs / rhs })
		}
	}
}

func (v *VM) binaryOp(fn func(value.Value, value.Value) value.Value) {
	rhs := v.Pop()
	lhs := v.Pop()
	v.Push(fn(lhs, rhs))
}

func (v *VM) readByte() uint8 {
	b := v.Chunk.Read(v.IP)
	v.IP++

	return b
}

func (v *VM) readWord() uint16 {
	w := v.Chunk.ReadWord(v.IP)
	v.IP += 2

	return w
}
