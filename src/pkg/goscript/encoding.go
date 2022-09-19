package goscript

import (
	"fmt"

	"github.com/Yoshi-Exeler/goscript/src/pkg/encoding"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
)

func EncodeProgram(program *Program) ([]byte, error) {
	encProg := encoding.Program{
		SymbolTableSize: uint64(program.SymbolTableSize),
	}
	for _, op := range program.Operations {
		encProg.Operations = append(encProg.Operations, encodeOp(op))
	}
	return proto.Marshal(&encProg)
}

func encodeOp(op BinaryOperation) *encoding.BinaryOperation {
	encOp := encoding.BinaryOperation{
		Type: uint32(op.Type),
	}
	for _, arg := range op.Args {
		encOp.Args = append(encOp.Args, encodeAny(arg))
	}
	return &encOp
}

func encodeAny(arg any) *anypb.Any {
	switch val := arg.(type) {
	case *Expression:
		buff, err := proto.Marshal(encodeExpr(val))
		if err != nil {
			panic("failed to encode expressiont to proto buffer expression")
		}
		return &anypb.Any{
			Value: buff,
		}
	case BinaryType:
		return encodeU64Container(uint64(val))
	case *uint8:
		return encodeU64Container(uint64(*val))
	case *uint16:
		return encodeU64Container(uint64(*val))
	case *uint32:
		return encodeU64Container(uint64(*val))
	case *uint64:
		return encodeU64Container(uint64(*val))
	case *int8:
		return encodeU64Container(uint64(*val))
	case *int16:
		return encodeU64Container(uint64(*val))
	case *int32:
		return encodeU64Container(uint64(*val))
	case *int64:
		return encodeU64Container(uint64(*val))
	case *[]*BinaryTypedValue:
		arr := &encoding.ArrayContainer{}
		for _, elem := range *val {
			arr.Values = append(arr.Values, &encoding.BinaryTypedValue{
				Type:  uint32(elem.Type),
				Value: encodeAny(elem.Value),
			})
		}
		buff, err := proto.Marshal(arr)
		if err != nil {
			panic("failed to encode array to proto buffer array")
		}
		return &anypb.Any{
			Value: buff,
		}
	case *FunctionArgument:
		buff, err := proto.Marshal(&encoding.FunctionArgument{
			Expression: encodeExpr(val.Expression),
			SymbolRef:  uint64(val.SymbolRef),
		})
		if err != nil {
			panic("failed to encode function arg to proto buffer function arg")
		}
		return &anypb.Any{
			Value: buff,
		}
	case nil:
		return &anypb.Any{}
	case int:
		return encodeU64Container(uint64(val))
	case *string:
		buff, err := proto.Marshal(&encoding.StringContainer{
			Value: string(*val),
		})
		if err != nil {
			panic("failed to encode string to proto buffer string")
		}
		return &anypb.Any{
			Value: buff,
		}
	default:
		panic(fmt.Sprintf("unknown type %#v in encode arg", arg))
	}
}

func encodeU64Container(u uint64) *anypb.Any {
	buff, err := proto.Marshal(&encoding.U64Container{
		Value: u,
	})
	if err != nil {
		panic("failed to encode int to proto buffer uint32")
	}
	return &anypb.Any{
		Value: buff,
	}
}

func encodeExpr(expr *Expression) *encoding.Expression {
	encExpr := encoding.Expression{
		Ref:      uint64(expr.Ref),
		Operator: uint32(expr.Operator),
	}
	if expr.Value != nil {
		encExpr.Value = &encoding.BinaryTypedValue{
			Type:  uint32(expr.Value.Type),
			Value: encodeAny(expr.Value.Value),
		}
	}
	for _, arg := range expr.Args {
		encExpr.Args = append(encExpr.Args, encodeAny(arg))
	}
	if expr.LeftExpression != nil {
		encExpr.Left = encodeExpr(expr.LeftExpression)
	}
	if expr.RightExpression != nil {
		encExpr.Right = encodeExpr(expr.RightExpression)
	}
	return &encExpr
}
