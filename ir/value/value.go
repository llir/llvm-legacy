// Package values provides a definition of LLVM IR values.
package value

import (
	"fmt"

	"github.com/llir/llvm/ir/types"
)

// TODO: Complete the list of value implementations.
//
//    *ir.BasicBlock

// A Value represents a computed value that may be used as an operand of other
// values. Some values can have a name and they belong to a function or a
// module.
//
// Value is one of the following types:
//
//    constant.Constant
//    *instruction.Local
type Value interface {
	fmt.Stringer
	// Type returns the type of the value.
	Type() types.Type
}
