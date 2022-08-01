package mem_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/meanguy/automato/internal/mem"
	"github.com/meanguy/automato/internal/opcode"
	"github.com/meanguy/automato/internal/value"
)

func TestChunkAddAndGetConstant(t *testing.T) {
	chunk := mem.Chunk{}

	expected := value.Value(1.23)
	actual := chunk.GetConstant(chunk.AddConstant(expected))

	assert.Equal(t, expected, actual)
}

func TestChunkWriteAndRead(t *testing.T) {
	testCases := []struct {
		values []byte
	}{
		{values: []byte{0, 1, 2, 3}},
		{values: []byte{12}},
		{values: []byte{12}},
	}

	for _, tc := range testCases {
		chunk := mem.Chunk{}

		t.Run(fmt.Sprintf("len=%d", len(tc.values)), func(t *testing.T) {
			for _, val := range tc.values {
				chunk.Write(val, 0)
			}

			for offset, expected := range tc.values {
				actual := chunk.Read(offset)

				assert.Equal(t, expected, actual)
			}
		})
	}
}

func TestChunkWriteAndReadWord(t *testing.T) {
	testCases := []struct {
		values []uint16
	}{
		{values: []uint16{0, 1, 2, 3, 1024}},
		{values: []uint16{12, 1 << 13}},
		{values: []uint16{12}},
	}

	for _, tc := range testCases {
		chunk := mem.Chunk{}

		t.Run(fmt.Sprintf("len=%d", len(tc.values)), func(t *testing.T) {
			for _, val := range tc.values {
				chunk.WriteWord(val, 0)
			}

			for offset, expected := range tc.values {
				actual := chunk.ReadWord(offset * 2)

				assert.Equal(t, expected, actual)
			}
		})
	}
}

func TestChunkWriteAndReadOp(t *testing.T) {
	testCases := []struct {
		values []opcode.OpCode
	}{
		{values: []opcode.OpCode{opcode.OpConstant, opcode.OpConstant, opcode.OpMultiply, opcode.OpConstant, opcode.OpAdd}},
		{values: []opcode.OpCode{opcode.OpConstant, opcode.OpReturn}},
		{values: []opcode.OpCode{opcode.OpReturn}},
	}

	for _, tc := range testCases {
		chunk := mem.Chunk{}

		t.Run(fmt.Sprintf("len=%d", len(tc.values)), func(t *testing.T) {
			for _, val := range tc.values {
				chunk.WriteOp(val, 0)
			}

			for offset, expected := range tc.values {
				actual := chunk.ReadOp(offset)

				assert.Equal(t, expected, actual)
			}
		})
	}
}
