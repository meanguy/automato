package debug

import (
	"fmt"
	"io"

	"github.com/meanguy/automato/internal/mem"
	"github.com/meanguy/automato/internal/opcode"
	"github.com/meanguy/automato/internal/value"
)

func DisassembleChunk(w io.Writer, chunk *mem.Chunk, header string) {
	fmt.Fprintf(w, "== %v ==\n", header)

	length := len(chunk.Code)
	for offset := 0; offset < length; {
		offset = DisassembleInstruction(w, chunk, offset)
	}
}

//nolint:cyclop // interpreting opcodes is necessarily complex
func DisassembleInstruction(w io.Writer, chunk *mem.Chunk, offset int) int {
	fmt.Fprintf(w, "%04d ", offset)

	if offset > 0 && chunk.Lines[offset] == chunk.Lines[offset-1] {
		fmt.Fprint(w, "   | ")
	} else {
		fmt.Fprintf(w, "%4d ", chunk.Lines[offset])
	}

	op := opcode.OpCode(chunk.Read(offset))
	switch op {
	case opcode.OpReturn:
		return simpleInstruction(w, "OpReturn", offset)
	case opcode.OpConstant:
		return constantInstruction(w, "OpConstant", chunk, offset)
	case opcode.OpConstantLong:
		return constantLongInstruction(w, "OpConstantLong", chunk, offset)
	case opcode.OpNegate:
		return simpleInstruction(w, "OpNegate", offset)
	case opcode.OpAdd:
		return simpleInstruction(w, "OpAdd", offset)
	case opcode.OpSubtract:
		return simpleInstruction(w, "OpSubtract", offset)
	case opcode.OpMultiply:
		return simpleInstruction(w, "OpMultiply", offset)
	case opcode.OpDivide:
		return simpleInstruction(w, "OpDivide", offset)
	default:
		fmt.Fprintf(w, "Unknown opcode: %d\n", op)

		return offset + 1
	}
}

func DisassembleStack(w io.Writer, stack []value.Value) {
	fmt.Fprintf(w, "+---| ")

	for _, val := range stack {
		fmt.Fprintf(w, "[ %v ]", val)
	}

	fmt.Fprintln(w)
}

func constantInstruction(w io.Writer, name string, chunk *mem.Chunk, offset int) int {
	constantID := chunk.Read(offset + 1)

	fmt.Fprintf(w, "%-16s %4d '%v'\n", name, constantID, chunk.Constants[constantID])

	return offset + 2
}

func constantLongInstruction(w io.Writer, name string, chunk *mem.Chunk, offset int) int {
	constantID := chunk.ReadWord(offset + 1)

	fmt.Fprintf(w, "%-16s %4d '%v'\n", name, constantID, chunk.Constants[constantID])

	return offset + 3
}

func simpleInstruction(w io.Writer, name string, offset int) int {
	fmt.Fprintf(w, "%s\n", name)

	return offset + 1
}
