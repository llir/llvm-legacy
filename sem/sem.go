// Package sem implements a static semantic analysis checker of LLVM IR modules.
package sem

import (
	"fmt"
	"strings"

	"github.com/llir/llvm/internal/enc"
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/irutil"
	"github.com/llir/llvm/ir/types"
	"github.com/pkg/errors"
)

// ErrorList represents a list of errors.
type ErrorList []error

// Error returns a string representation of the list of errors.
func (es ErrorList) Error() string {
	var errs []string
	for _, e := range es {
		errs = append(errs, e.Error())
	}
	return strings.Join(errs, "; ")
}

// Check performs static semantic analysis on the given LLVM IR module.
func Check(m *ir.Module) error {
	sem := &sem{}
	// check performs static semantic analysis on the given LLVM IR node.
	check := func(n interface{}) {
		switch n := n.(type) {
		case *ir.Global:
			sem.checkGlobal(n)
		case *ir.Function:
			sem.checkFunc(n)
		case *ir.Param:
			sem.checkParam(n)
		case *ir.BasicBlock:
			sem.checkBlock(n)
		case types.Type:
			sem.checkType(n)
		case constant.Constant:
			sem.checkConst(n)
		case ir.Instruction:
			sem.checkInst(n)
		case ir.Terminator:
			sem.checkTerm(n)
		}
	}
	irutil.Walk(m, check)
	if len(sem.errs) > 0 {
		return sem.errs
	}
	return nil
}

// sem represents a static semantic analysis checker for LLVM IR.
type sem struct {
	// List of identified errors.
	errs ErrorList
}

// Errorf formats according to a format specifier and appends the error to the
// list of identified semantic errors.
func (sem *sem) Errorf(format string, args ...interface{}) {
	err := errors.Errorf(format, args...)
	sem.errs = append(sem.errs, err)
}

// checkGlobal validates the semantics of the given global.
func (sem *sem) checkGlobal(global *ir.Global) {
	// Validate global name.
	if len(global.Name) == 0 {
		sem.Errorf("global name missing")
	} else if !isValidIdent(global.Name) {
		sem.Errorf("invalid global name `%v`", enc.Global(global.Name))
	}
	// Validate global variable type.
	content, elem := global.Content, global.Typ.Elem
	if !content.Equal(elem) {
		sem.Errorf("global variable content type `%v` and element type `%v` mismatch", content, elem)
	}
	// Validate global variable content type.
	if !isSingleValueType(content) && !isAggregateType(content) {
		sem.Errorf("invalid global content type; expected single value or aggregate type, got %T", content)
	}
	// Validate global variable initial value
	if init := global.Init; init != nil {
		if !content.Equal(init.Type()) {
			sem.Errorf("global variable content type `%v` and initial value type `%v` mismatch", content, init.Type())
		}
	}
}

// checkFunc validates the semantics of the given function.
func (sem *sem) checkFunc(f *ir.Function) {
	// TODO: Implement
	//panic("not yet implemented")
}

// checkParam validates the semantics of the given function parameter.
func (sem *sem) checkParam(param *ir.Param) {
	panic("not yet implemented")
}

// checkBlock validates the semantics of the given basic block.
func (sem *sem) checkBlock(param *ir.BasicBlock) {
	panic("not yet implemented")
}

// checkType validates the semantics of the given type.
func (sem *sem) checkType(t types.Type) {
	switch t := t.(type) {
	case *types.VoidType:
		// nothing to do.
	case *types.FuncType:
		// The return type of a function type is a void type or first class type -
		// except for label and metadata types.
		//
		// References:
		//    http://llvm.org/docs/LangRef.html#function-type
		if !types.IsVoid(t.Ret) && !isSingleValueType(t.Ret) && !isAggregateType(t.Ret) {
			sem.Errorf("invalid function return type; expected void, single value or aggregate type, got %T", t.Ret)
		}
		for _, param := range t.Params {
			if len(param.Name) > 0 && !isValidIdent(param.Name) {
				sem.Errorf("invalid function parameter name `%v`", enc.Local(param.Name))
			}
			if !isFirstClassType(param.Typ) {
				sem.Errorf("invalid function parameter; expected first class type, got %T", param.Typ)
			}
		}
	case *types.IntType:
		// Any bit width from 1 bit to 2²³-1 can be specified.
		//
		// References:
		//    http://llvm.org/docs/LangRef.html#integer-type
		const maxSize = 1<<23 - 1
		if t.Size < 1 {
			sem.Errorf("invalid integer type bit width; expected > 0, got %d", t.Size)
		} else if t.Size > maxSize {
			sem.Errorf("invalid integer type bit width; expected < 2^24, got %d", t.Size)
		}
	case *types.FloatType:
		switch t.Kind {
		case types.FloatKindIEEE_16:
		case types.FloatKindIEEE_32:
		case types.FloatKindIEEE_64:
		case types.FloatKindIEEE_128:
		case types.FloatKindDoubleExtended_80:
		case types.FloatKindDoubleDouble_128:
		default:
			sem.Errorf("invalid float type kind; expected half, float, double, fp128, x86_fp80 or ppc_fp128, got %v", t.Kind)
		}
	case *types.PointerType:
		if !types.IsFunc(t.Elem) && !isSingleValueType(t.Elem) && !isAggregateType(t.Elem) {
			sem.Errorf("invalid pointer element type; expected function, single value or aggregate type, got %T", t.Elem)
		}
	case *types.VectorType:
		// The number of elements is a constant integer value larger than 0; the
		// element type may be any integer, floating point or pointer type.
		// Vectors of size zero are not allowed.
		//
		// References:
		//    http://llvm.org/docs/LangRef.html#vector-type
		if !types.IsInt(t.Elem) && !types.IsFloat(t.Elem) && !types.IsPointer(t.Elem) {
			sem.Errorf("invalid vector element type; expected integer, floating-point or pointer type, got %T", t.Elem)
		}
	case *types.LabelType:
		// nothing to do.
	case *types.MetadataType:
		// nothing to do.
	case *types.ArrayType:
		if !isSingleValueType(t.Elem) && !isAggregateType(t.Elem) {
			sem.Errorf("invalid array element type; expected single value or aggregate type, got %T", t.Elem)
		}
	case *types.StructType:
		for _, field := range t.Fields {
			if !isSingleValueType(field) && !isAggregateType(field) {
				sem.Errorf("invalid struct field type; expected single value or aggregate type, got %T", field)
			}
		}
	case *types.NamedType:
		if len(t.Name) == 0 {
			sem.Errorf("type name missing")
		} else if !isValidIdent(t.Name) {
			sem.Errorf("invalid type name `%v`", enc.Local(t.Name))
		}
		// t.Def is validated when later traversed.
	default:
		panic(fmt.Errorf("support for type %T not yet implemented", t))
	}
}

// checkConst validates the semantics of the given constant.
func (sem *sem) checkConst(c constant.Constant) {
	switch c := c.(type) {
	// Simple constants.
	case *constant.Int:
		panic("not yet implemented")
	case *constant.Float:
		panic("not yet implemented")
	case *constant.Null:
		panic("not yet implemented")
	// Complex constants.
	case *constant.Vector:
		panic("not yet implemented")
	case *constant.Array:
		panic("not yet implemented")
	case *constant.Struct:
		panic("not yet implemented")
	case *constant.ZeroInitializer:
		panic("not yet implemented")
	// Binary expressions.
	case *constant.ExprAdd:
		panic("not yet implemented")
	case *constant.ExprFAdd:
		panic("not yet implemented")
	case *constant.ExprSub:
		panic("not yet implemented")
	case *constant.ExprFSub:
		panic("not yet implemented")
	case *constant.ExprMul:
		panic("not yet implemented")
	case *constant.ExprFMul:
		panic("not yet implemented")
	case *constant.ExprUDiv:
		panic("not yet implemented")
	case *constant.ExprSDiv:
		panic("not yet implemented")
	case *constant.ExprFDiv:
		panic("not yet implemented")
	case *constant.ExprURem:
		panic("not yet implemented")
	case *constant.ExprSRem:
		panic("not yet implemented")
	case *constant.ExprFRem:
		panic("not yet implemented")
	// Bitwise expressions.
	case *constant.ExprShl:
		panic("not yet implemented")
	case *constant.ExprLShr:
		panic("not yet implemented")
	case *constant.ExprAShr:
		panic("not yet implemented")
	case *constant.ExprAnd:
		panic("not yet implemented")
	case *constant.ExprOr:
		panic("not yet implemented")
	case *constant.ExprXor:
		panic("not yet implemented")
	// Memory expressions.
	case *constant.ExprGetElementPtr:
		panic("not yet implemented")
	// Conversion expressions.
	case *constant.ExprTrunc:
		panic("not yet implemented")
	case *constant.ExprZExt:
		panic("not yet implemented")
	case *constant.ExprSExt:
		panic("not yet implemented")
	case *constant.ExprFPTrunc:
		panic("not yet implemented")
	case *constant.ExprFPExt:
		panic("not yet implemented")
	case *constant.ExprFPToUI:
		panic("not yet implemented")
	case *constant.ExprFPToSI:
		panic("not yet implemented")
	case *constant.ExprUIToFP:
		panic("not yet implemented")
	case *constant.ExprSIToFP:
		panic("not yet implemented")
	case *constant.ExprPtrToInt:
		panic("not yet implemented")
	case *constant.ExprIntToPtr:
		panic("not yet implemented")
	case *constant.ExprBitCast:
		panic("not yet implemented")
	case *constant.ExprAddrSpaceCast:
		panic("not yet implemented")
	// Other expressions.
	case *constant.ExprICmp:
		panic("not yet implemented")
	case *constant.ExprFCmp:
		panic("not yet implemented")
	case *constant.ExprSelect:
		panic("not yet implemented")
	default:
		panic(fmt.Errorf("support for constant %T not yet implemented", c))
	}
}

// checkInst validates the semantics of the given instruction.
func (sem *sem) checkInst(inst ir.Instruction) {
	switch inst := inst.(type) {
	// Binary instructions.
	case *ir.InstAdd:
		panic("not yet implemented")
	case *ir.InstFAdd:
		panic("not yet implemented")
	case *ir.InstSub:
		panic("not yet implemented")
	case *ir.InstFSub:
		panic("not yet implemented")
	case *ir.InstMul:
		panic("not yet implemented")
	case *ir.InstFMul:
		panic("not yet implemented")
	case *ir.InstUDiv:
		panic("not yet implemented")
	case *ir.InstSDiv:
		panic("not yet implemented")
	case *ir.InstFDiv:
		panic("not yet implemented")
	case *ir.InstURem:
		panic("not yet implemented")
	case *ir.InstSRem:
		panic("not yet implemented")
	case *ir.InstFRem:
		panic("not yet implemented")
	// Bitwise instructions.
	case *ir.InstShl:
		panic("not yet implemented")
	case *ir.InstLShr:
		panic("not yet implemented")
	case *ir.InstAShr:
		panic("not yet implemented")
	case *ir.InstAnd:
		panic("not yet implemented")
	case *ir.InstOr:
		panic("not yet implemented")
	case *ir.InstXor:
		panic("not yet implemented")
	// Memory instructions.
	case *ir.InstAlloca:
		panic("not yet implemented")
	case *ir.InstLoad:
		panic("not yet implemented")
	case *ir.InstStore:
		panic("not yet implemented")
	case *ir.InstGetElementPtr:
		panic("not yet implemented")
	// Conversion instructions.
	case *ir.InstTrunc:
		panic("not yet implemented")
	case *ir.InstZExt:
		panic("not yet implemented")
	case *ir.InstSExt:
		panic("not yet implemented")
	case *ir.InstFPTrunc:
		panic("not yet implemented")
	case *ir.InstFPExt:
		panic("not yet implemented")
	case *ir.InstFPToUI:
		panic("not yet implemented")
	case *ir.InstFPToSI:
		panic("not yet implemented")
	case *ir.InstUIToFP:
		panic("not yet implemented")
	case *ir.InstSIToFP:
		panic("not yet implemented")
	case *ir.InstPtrToInt:
		panic("not yet implemented")
	case *ir.InstIntToPtr:
		panic("not yet implemented")
	case *ir.InstBitCast:
		panic("not yet implemented")
	case *ir.InstAddrSpaceCast:
		panic("not yet implemented")
	// Other instructions.
	case *ir.InstICmp:
		panic("not yet implemented")
	case *ir.InstFCmp:
		panic("not yet implemented")
	case *ir.InstPhi:
		panic("not yet implemented")
	case *ir.InstSelect:
		panic("not yet implemented")
	case *ir.InstCall:
		panic("not yet implemented")
	default:
		panic(fmt.Errorf("support for instruction %T not yet implemented", inst))
	}
}

// checkTerm validates the semantics of the given terminator.
func (sem *sem) checkTerm(term ir.Terminator) {
	switch term := term.(type) {
	case *ir.TermRet:
		panic("not yet implemented")
	case *ir.TermBr:
		panic("not yet implemented")
	case *ir.TermCondBr:
		panic("not yet implemented")
	case *ir.TermSwitch:
		panic("not yet implemented")
	case *ir.TermUnreachable:
		panic("not yet implemented")
	default:
		panic(fmt.Errorf("support for instruction %T not yet implemented", term))
	}
}

const (
	asciiLetter  = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	letter       = asciiLetter + "$-._"
	decimalDigit = "0123456789"
)

// isValidIdent reports whether the given identifier is valid.
func isValidIdent(ident string) bool {
	// TODO: Add support for quoted string identifiers.
	return isValidID(ident) || isValidName(ident)
}

// isValidID reports whether the given ID is valid.
func isValidID(id string) bool {
	// _decimals
	//    : _decimal_digit { _decimal_digit }
	// ;
	//
	// _id
	//    : _decimals
	// ;
	if len(id) < 1 {
		return false
	}
	for _, r := range id {
		const charset = decimalDigit
		if !strings.ContainsRune(charset, r) {
			return false
		}
	}
	return true
}

// isValidName reports whether the given name is valid.
func isValidName(name string) bool {
	// _ascii_letter
	//    : 'A' - 'Z'
	//    | 'a' - 'z'
	// ;
	//
	// _letter
	//    : _ascii_letter
	//    | '$'
	//    | '-'
	//    | '.'
	//    | '_'
	// ;
	//
	// _decimal_digit
	//    : '0' - '9'
	// ;
	//
	// _name
	//    : _letter { _letter | _decimal_digit }
	// ;
	if len(name) < 1 {
		return false
	}
	for i, r := range name {
		charset := letter
		if i > 0 {
			charset = letter + decimalDigit
		}
		if !strings.ContainsRune(charset, r) {
			return false
		}
	}
	return true
}

// isFirstClassType reports whether the given type is a first class type.
func isFirstClassType(t types.Type) bool {
	switch t := t.(type) {
	case *types.VoidType:
		return false
	case *types.FuncType:
		return false
	case *types.IntType:
		return true
	case *types.FloatType:
		return true
	case *types.PointerType:
		return true
	case *types.VectorType:
		return true
	case *types.LabelType:
		return true
	case *types.MetadataType:
		return true
	case *types.ArrayType:
		return true
	case *types.StructType:
		return true
	case *types.NamedType:
		return isFirstClassType(t.Def)
	default:
		panic(fmt.Errorf("support for type %T not yet implemented", t))
	}
}

// isSingleValueType reports whether the given type is a single value type.
func isSingleValueType(t types.Type) bool {
	switch t := t.(type) {
	case *types.VoidType:
		return false
	case *types.FuncType:
		return false
	case *types.IntType:
		return true
	case *types.FloatType:
		return true
	case *types.PointerType:
		return true
	case *types.VectorType:
		return true
	case *types.LabelType:
		return false
	case *types.MetadataType:
		return false
	case *types.ArrayType:
		return false
	case *types.StructType:
		return false
	case *types.NamedType:
		return isSingleValueType(t.Def)
	default:
		panic(fmt.Errorf("support for type %T not yet implemented", t))
	}
}

// isAggregateType reports whether the given type is an aggregate type.
func isAggregateType(t types.Type) bool {
	switch t := t.(type) {
	case *types.VoidType:
		return false
	case *types.FuncType:
		return false
	case *types.IntType:
		return true
	case *types.FloatType:
		return true
	case *types.PointerType:
		return true
	case *types.VectorType:
		return true
	case *types.LabelType:
		return false
	case *types.MetadataType:
		return false
	case *types.ArrayType:
		return true
	case *types.StructType:
		return true
	case *types.NamedType:
		return isAggregateType(t.Def)
	default:
		panic(fmt.Errorf("support for type %T not yet implemented", t))
	}
}
