package vm_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/meanguy/automato/internal/value"
	"github.com/meanguy/automato/internal/vm"
)

func TestVMPushAndPopValue(t *testing.T) {
	vm := vm.NewVM()
	expected := value.Value(3.14159)

	vm.Push(expected)
	actual := vm.Pop()

	assert.Equal(t, expected, actual)
}

func TestVMInterpretSource(t *testing.T) {
	vm := vm.NewVM()
	expected := value.Value(5)

	assert.NoError(t, vm.Interpret("3+2"))
	assert.Equal(t, expected, vm.Pop())
}
