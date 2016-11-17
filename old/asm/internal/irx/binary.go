// generated by gen.go using 'go generate'; DO NOT EDIT.

// === [ Binary instructions ] =================================================
//
// References:
//    http://llvm.org/docs/LangRef.html#binary-operations

package irx

import (
	"github.com/llir/llvm/ir/instruction"
	"github.com/mewkiz/pkg/errutil"
)

// --- [ add ] -----------------------------------------------------------------

// NewAdd returns a new add instruction based on the given operand type and
// operands.
func NewAddInst(typ, xVal, yVal interface{}) (*instruction.Add, error) {
	x, err := NewValue(typ, xVal)
	if err != nil {
		return nil, errutil.Err(err)
	}
	y, err := NewValue(typ, yVal)
	if err != nil {
		return nil, errutil.Err(err)
	}
	return instruction.NewAdd(x, y)
}

// --- [ fadd ] ----------------------------------------------------------------

// NewFAdd returns a new fadd instruction based on the given operand type and
// operands.
func NewFAddInst(typ, xVal, yVal interface{}) (*instruction.FAdd, error) {
	x, err := NewValue(typ, xVal)
	if err != nil {
		return nil, errutil.Err(err)
	}
	y, err := NewValue(typ, yVal)
	if err != nil {
		return nil, errutil.Err(err)
	}
	return instruction.NewFAdd(x, y)
}

// --- [ sub ] -----------------------------------------------------------------

// NewSub returns a new sub instruction based on the given operand type and
// operands.
func NewSubInst(typ, xVal, yVal interface{}) (*instruction.Sub, error) {
	x, err := NewValue(typ, xVal)
	if err != nil {
		return nil, errutil.Err(err)
	}
	y, err := NewValue(typ, yVal)
	if err != nil {
		return nil, errutil.Err(err)
	}
	return instruction.NewSub(x, y)
}

// --- [ fsub ] ----------------------------------------------------------------

// NewFSub returns a new fsub instruction based on the given operand type and
// operands.
func NewFSubInst(typ, xVal, yVal interface{}) (*instruction.FSub, error) {
	x, err := NewValue(typ, xVal)
	if err != nil {
		return nil, errutil.Err(err)
	}
	y, err := NewValue(typ, yVal)
	if err != nil {
		return nil, errutil.Err(err)
	}
	return instruction.NewFSub(x, y)
}

// --- [ mul ] -----------------------------------------------------------------

// NewMul returns a new mul instruction based on the given operand type and
// operands.
func NewMulInst(typ, xVal, yVal interface{}) (*instruction.Mul, error) {
	x, err := NewValue(typ, xVal)
	if err != nil {
		return nil, errutil.Err(err)
	}
	y, err := NewValue(typ, yVal)
	if err != nil {
		return nil, errutil.Err(err)
	}
	return instruction.NewMul(x, y)
}

// --- [ fmul ] ----------------------------------------------------------------

// NewFMul returns a new fmul instruction based on the given operand type and
// operands.
func NewFMulInst(typ, xVal, yVal interface{}) (*instruction.FMul, error) {
	x, err := NewValue(typ, xVal)
	if err != nil {
		return nil, errutil.Err(err)
	}
	y, err := NewValue(typ, yVal)
	if err != nil {
		return nil, errutil.Err(err)
	}
	return instruction.NewFMul(x, y)
}

// --- [ udiv ] ----------------------------------------------------------------

// NewUDiv returns a new udiv instruction based on the given operand type and
// operands.
func NewUDivInst(typ, xVal, yVal interface{}) (*instruction.UDiv, error) {
	x, err := NewValue(typ, xVal)
	if err != nil {
		return nil, errutil.Err(err)
	}
	y, err := NewValue(typ, yVal)
	if err != nil {
		return nil, errutil.Err(err)
	}
	return instruction.NewUDiv(x, y)
}

// --- [ sdiv ] ----------------------------------------------------------------

// NewSDiv returns a new sdiv instruction based on the given operand type and
// operands.
func NewSDivInst(typ, xVal, yVal interface{}) (*instruction.SDiv, error) {
	x, err := NewValue(typ, xVal)
	if err != nil {
		return nil, errutil.Err(err)
	}
	y, err := NewValue(typ, yVal)
	if err != nil {
		return nil, errutil.Err(err)
	}
	return instruction.NewSDiv(x, y)
}

// --- [ fdiv ] ----------------------------------------------------------------

// NewFDiv returns a new fdiv instruction based on the given operand type and
// operands.
func NewFDivInst(typ, xVal, yVal interface{}) (*instruction.FDiv, error) {
	x, err := NewValue(typ, xVal)
	if err != nil {
		return nil, errutil.Err(err)
	}
	y, err := NewValue(typ, yVal)
	if err != nil {
		return nil, errutil.Err(err)
	}
	return instruction.NewFDiv(x, y)
}

// --- [ urem ] ----------------------------------------------------------------

// NewURem returns a new urem instruction based on the given operand type and
// operands.
func NewURemInst(typ, xVal, yVal interface{}) (*instruction.URem, error) {
	x, err := NewValue(typ, xVal)
	if err != nil {
		return nil, errutil.Err(err)
	}
	y, err := NewValue(typ, yVal)
	if err != nil {
		return nil, errutil.Err(err)
	}
	return instruction.NewURem(x, y)
}

// --- [ srem ] ----------------------------------------------------------------

// NewSRem returns a new srem instruction based on the given operand type and
// operands.
func NewSRemInst(typ, xVal, yVal interface{}) (*instruction.SRem, error) {
	x, err := NewValue(typ, xVal)
	if err != nil {
		return nil, errutil.Err(err)
	}
	y, err := NewValue(typ, yVal)
	if err != nil {
		return nil, errutil.Err(err)
	}
	return instruction.NewSRem(x, y)
}

// --- [ frem ] ----------------------------------------------------------------

// NewFRem returns a new frem instruction based on the given operand type and
// operands.
func NewFRemInst(typ, xVal, yVal interface{}) (*instruction.FRem, error) {
	x, err := NewValue(typ, xVal)
	if err != nil {
		return nil, errutil.Err(err)
	}
	y, err := NewValue(typ, yVal)
	if err != nil {
		return nil, errutil.Err(err)
	}
	return instruction.NewFRem(x, y)
}