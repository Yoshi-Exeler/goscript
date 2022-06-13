package goscript

// applyOperator applies the specified operator to the specified values, assuming that the operation has been type checked before
func applyOperator(l *BinaryTypedValue, r *BinaryTypedValue, op BinaryOperator) *BinaryTypedValue {
	switch op {
	case BO_PLUS:
		switch l.Type {
		case BT_INT8:
			return &BinaryTypedValue{
				Type:  BT_INT8,
				Value: genericPlus[int8](l.Value, r.Value),
			}
		case BT_INT16:
			return &BinaryTypedValue{
				Type:  BT_INT16,
				Value: genericPlus[int16](l.Value, r.Value),
			}
		case BT_INT32:
			return &BinaryTypedValue{
				Type:  BT_INT32,
				Value: genericPlus[int32](l.Value, r.Value),
			}
		case BT_INT64:
			return &BinaryTypedValue{
				Type:  BT_INT64,
				Value: genericPlus[int64](l.Value, r.Value),
			}
		case BT_UINT8:
			return &BinaryTypedValue{
				Type:  BT_UINT8,
				Value: genericPlus[uint8](l.Value, r.Value),
			}
		case BT_UINT16:
			return &BinaryTypedValue{
				Type:  BT_UINT16,
				Value: genericPlus[uint16](l.Value, r.Value),
			}
		case BT_UINT32:
			return &BinaryTypedValue{
				Type:  BT_UINT32,
				Value: genericPlus[uint32](l.Value, r.Value),
			}
		case BT_UINT64:
			return &BinaryTypedValue{
				Type:  BT_UINT64,
				Value: genericPlus[uint64](l.Value, r.Value),
			}
		case BT_BYTE:
			return &BinaryTypedValue{
				Type:  BT_BYTE,
				Value: genericPlus[byte](l.Value, r.Value),
			}
		case BT_FLOAT32:
			return &BinaryTypedValue{
				Type:  BT_FLOAT32,
				Value: genericPlus[float32](l.Value, r.Value),
			}
		case BT_FLOAT64:
			return &BinaryTypedValue{
				Type:  BT_FLOAT64,
				Value: genericPlus[float64](l.Value, r.Value),
			}
		default:
			panic("invalid type for plus operator")
		}
	case BO_MINUS:
		switch l.Type {
		case BT_INT8:
			return &BinaryTypedValue{
				Type:  BT_INT8,
				Value: genericMinus[int8](l.Value, r.Value),
			}
		case BT_INT16:
			return &BinaryTypedValue{
				Type:  BT_INT16,
				Value: genericMinus[int16](l.Value, r.Value),
			}
		case BT_INT32:
			return &BinaryTypedValue{
				Type:  BT_INT32,
				Value: genericMinus[int32](l.Value, r.Value),
			}
		case BT_INT64:
			return &BinaryTypedValue{
				Type:  BT_INT64,
				Value: genericMinus[int64](l.Value, r.Value),
			}
		case BT_UINT8:
			return &BinaryTypedValue{
				Type:  BT_UINT8,
				Value: genericMinus[uint8](l.Value, r.Value),
			}
		case BT_UINT16:
			return &BinaryTypedValue{
				Type:  BT_UINT16,
				Value: genericMinus[uint16](l.Value, r.Value),
			}
		case BT_UINT32:
			return &BinaryTypedValue{
				Type:  BT_UINT32,
				Value: genericMinus[uint32](l.Value, r.Value),
			}
		case BT_UINT64:
			return &BinaryTypedValue{
				Type:  BT_UINT64,
				Value: genericMinus[uint64](l.Value, r.Value),
			}
		case BT_BYTE:
			return &BinaryTypedValue{
				Type:  BT_BYTE,
				Value: genericMinus[byte](l.Value, r.Value),
			}
		case BT_FLOAT32:
			return &BinaryTypedValue{
				Type:  BT_FLOAT32,
				Value: genericMinus[float32](l.Value, r.Value),
			}
		case BT_FLOAT64:
			return &BinaryTypedValue{
				Type:  BT_FLOAT64,
				Value: genericMinus[float64](l.Value, r.Value),
			}
		default:
			panic("invalid type for minus operator")
		}
	case BO_MULTIPLY:
		switch l.Type {
		case BT_INT8:
			return &BinaryTypedValue{
				Type:  BT_INT8,
				Value: genericMultiply[int8](l.Value, r.Value),
			}
		case BT_INT16:
			return &BinaryTypedValue{
				Type:  BT_INT16,
				Value: genericMultiply[int16](l.Value, r.Value),
			}
		case BT_INT32:
			return &BinaryTypedValue{
				Type:  BT_INT32,
				Value: genericMultiply[int32](l.Value, r.Value),
			}
		case BT_INT64:
			return &BinaryTypedValue{
				Type:  BT_INT64,
				Value: genericMultiply[int64](l.Value, r.Value),
			}
		case BT_UINT8:
			return &BinaryTypedValue{
				Type:  BT_UINT8,
				Value: genericMultiply[uint8](l.Value, r.Value),
			}
		case BT_UINT16:
			return &BinaryTypedValue{
				Type:  BT_UINT16,
				Value: genericMultiply[uint16](l.Value, r.Value),
			}
		case BT_UINT32:
			return &BinaryTypedValue{
				Type:  BT_UINT32,
				Value: genericMultiply[uint32](l.Value, r.Value),
			}
		case BT_UINT64:
			return &BinaryTypedValue{
				Type:  BT_UINT64,
				Value: genericMultiply[uint64](l.Value, r.Value),
			}
		case BT_BYTE:
			return &BinaryTypedValue{
				Type:  BT_BYTE,
				Value: genericMultiply[byte](l.Value, r.Value),
			}
		case BT_FLOAT32:
			return &BinaryTypedValue{
				Type:  BT_FLOAT32,
				Value: genericMultiply[float32](l.Value, r.Value),
			}
		case BT_FLOAT64:
			return &BinaryTypedValue{
				Type:  BT_FLOAT64,
				Value: genericMultiply[float64](l.Value, r.Value),
			}
		default:
			panic("invalid type for multiply operator")
		}
	case BO_DIVIDE:
		switch l.Type {
		case BT_INT8:
			return &BinaryTypedValue{
				Type:  BT_INT8,
				Value: genericDivide[int8](l.Value, r.Value),
			}
		case BT_INT16:
			return &BinaryTypedValue{
				Type:  BT_INT16,
				Value: genericDivide[int16](l.Value, r.Value),
			}
		case BT_INT32:
			return &BinaryTypedValue{
				Type:  BT_INT32,
				Value: genericDivide[int32](l.Value, r.Value),
			}
		case BT_INT64:
			return &BinaryTypedValue{
				Type:  BT_INT64,
				Value: genericDivide[int64](l.Value, r.Value),
			}
		case BT_UINT8:
			return &BinaryTypedValue{
				Type:  BT_UINT8,
				Value: genericDivide[uint8](l.Value, r.Value),
			}
		case BT_UINT16:
			return &BinaryTypedValue{
				Type:  BT_UINT16,
				Value: genericDivide[uint16](l.Value, r.Value),
			}
		case BT_UINT32:
			return &BinaryTypedValue{
				Type:  BT_UINT32,
				Value: genericDivide[uint32](l.Value, r.Value),
			}
		case BT_UINT64:
			return &BinaryTypedValue{
				Type:  BT_UINT64,
				Value: genericDivide[uint64](l.Value, r.Value),
			}
		case BT_BYTE:
			return &BinaryTypedValue{
				Type:  BT_BYTE,
				Value: genericDivide[byte](l.Value, r.Value),
			}
		case BT_FLOAT32:
			return &BinaryTypedValue{
				Type:  BT_FLOAT32,
				Value: genericDivide[float32](l.Value, r.Value),
			}
		case BT_FLOAT64:
			return &BinaryTypedValue{
				Type:  BT_FLOAT64,
				Value: genericDivide[float64](l.Value, r.Value),
			}
		default:
			panic("invalid type for divide operator")
		}
	case BO_EQUALS:
		switch l.Type {
		case BT_INT8:
			return &BinaryTypedValue{
				Type:  BT_BOOLEAN,
				Value: genericEquals[int8](l.Value, r.Value),
			}
		case BT_INT16:
			return &BinaryTypedValue{
				Type:  BT_BOOLEAN,
				Value: genericEquals[int16](l.Value, r.Value),
			}
		case BT_INT32:
			return &BinaryTypedValue{
				Type:  BT_BOOLEAN,
				Value: genericEquals[int32](l.Value, r.Value),
			}
		case BT_INT64:
			return &BinaryTypedValue{
				Type:  BT_BOOLEAN,
				Value: genericEquals[int64](l.Value, r.Value),
			}
		case BT_UINT8:
			return &BinaryTypedValue{
				Type:  BT_BOOLEAN,
				Value: genericEquals[uint8](l.Value, r.Value),
			}
		case BT_UINT16:
			return &BinaryTypedValue{
				Type:  BT_BOOLEAN,
				Value: genericEquals[uint16](l.Value, r.Value),
			}
		case BT_UINT32:
			return &BinaryTypedValue{
				Type:  BT_BOOLEAN,
				Value: genericEquals[uint32](l.Value, r.Value),
			}
		case BT_UINT64:
			return &BinaryTypedValue{
				Type:  BT_BOOLEAN,
				Value: genericEquals[uint64](l.Value, r.Value),
			}
		case BT_BYTE:
			return &BinaryTypedValue{
				Type:  BT_BOOLEAN,
				Value: genericEquals[byte](l.Value, r.Value),
			}
		case BT_FLOAT32:
			return &BinaryTypedValue{
				Type:  BT_BOOLEAN,
				Value: genericEquals[float32](l.Value, r.Value),
			}
		case BT_FLOAT64:
			return &BinaryTypedValue{
				Type:  BT_BOOLEAN,
				Value: genericEquals[float64](l.Value, r.Value),
			}
		case BT_CHAR:
			return &BinaryTypedValue{
				Type:  BT_BOOLEAN,
				Value: genericEquals[rune](l.Value, r.Value),
			}
		case BT_STRING:
			return &BinaryTypedValue{
				Type:  BT_BOOLEAN,
				Value: genericEquals[string](l.Value, r.Value),
			}
		case BT_BOOLEAN:
			return &BinaryTypedValue{
				Type:  BT_BOOLEAN,
				Value: genericEquals[bool](l.Value, r.Value),
			}
		default:
			panic("invalid type for equals operator")
		}
	case BO_GREATER:
		switch l.Type {
		case BT_INT8:
			return &BinaryTypedValue{
				Type:  BT_BOOLEAN,
				Value: genericGreater[int8](l.Value, r.Value),
			}
		case BT_INT16:
			return &BinaryTypedValue{
				Type:  BT_BOOLEAN,
				Value: genericGreater[int16](l.Value, r.Value),
			}
		case BT_INT32:
			return &BinaryTypedValue{
				Type:  BT_BOOLEAN,
				Value: genericGreater[int32](l.Value, r.Value),
			}
		case BT_INT64:
			return &BinaryTypedValue{
				Type:  BT_BOOLEAN,
				Value: genericGreater[int64](l.Value, r.Value),
			}
		case BT_UINT8:
			return &BinaryTypedValue{
				Type:  BT_BOOLEAN,
				Value: genericGreater[uint8](l.Value, r.Value),
			}
		case BT_UINT16:
			return &BinaryTypedValue{
				Type:  BT_BOOLEAN,
				Value: genericGreater[uint16](l.Value, r.Value),
			}
		case BT_UINT32:
			return &BinaryTypedValue{
				Type:  BT_BOOLEAN,
				Value: genericGreater[uint32](l.Value, r.Value),
			}
		case BT_UINT64:
			return &BinaryTypedValue{
				Type:  BT_BOOLEAN,
				Value: genericGreater[uint64](l.Value, r.Value),
			}
		case BT_BYTE:
			return &BinaryTypedValue{
				Type:  BT_BOOLEAN,
				Value: genericGreater[byte](l.Value, r.Value),
			}
		case BT_FLOAT32:
			return &BinaryTypedValue{
				Type:  BT_BOOLEAN,
				Value: genericGreater[float32](l.Value, r.Value),
			}
		case BT_FLOAT64:
			return &BinaryTypedValue{
				Type:  BT_BOOLEAN,
				Value: genericGreater[float64](l.Value, r.Value),
			}
		case BT_CHAR:
			return &BinaryTypedValue{
				Type:  BT_BOOLEAN,
				Value: genericGreater[rune](l.Value, r.Value),
			}
		default:
			panic("invalid type for equals operator")
		}
	case BO_LESSER:
		switch l.Type {
		case BT_INT8:
			return &BinaryTypedValue{
				Type:  BT_BOOLEAN,
				Value: genericLesser[int8](l.Value, r.Value),
			}
		case BT_INT16:
			return &BinaryTypedValue{
				Type:  BT_BOOLEAN,
				Value: genericLesser[int16](l.Value, r.Value),
			}
		case BT_INT32:
			return &BinaryTypedValue{
				Type:  BT_BOOLEAN,
				Value: genericLesser[int32](l.Value, r.Value),
			}
		case BT_INT64:
			return &BinaryTypedValue{
				Type:  BT_BOOLEAN,
				Value: genericLesser[int64](l.Value, r.Value),
			}
		case BT_UINT8:
			return &BinaryTypedValue{
				Type:  BT_BOOLEAN,
				Value: genericLesser[uint8](l.Value, r.Value),
			}
		case BT_UINT16:
			return &BinaryTypedValue{
				Type:  BT_BOOLEAN,
				Value: genericLesser[uint16](l.Value, r.Value),
			}
		case BT_UINT32:
			return &BinaryTypedValue{
				Type:  BT_BOOLEAN,
				Value: genericLesser[uint32](l.Value, r.Value),
			}
		case BT_UINT64:
			return &BinaryTypedValue{
				Type:  BT_BOOLEAN,
				Value: genericLesser[uint64](l.Value, r.Value),
			}
		case BT_BYTE:
			return &BinaryTypedValue{
				Type:  BT_BOOLEAN,
				Value: genericLesser[byte](l.Value, r.Value),
			}
		case BT_FLOAT32:
			return &BinaryTypedValue{
				Type:  BT_BOOLEAN,
				Value: genericLesser[float32](l.Value, r.Value),
			}
		case BT_FLOAT64:
			return &BinaryTypedValue{
				Type:  BT_BOOLEAN,
				Value: genericLesser[float64](l.Value, r.Value),
			}
		case BT_CHAR:
			return &BinaryTypedValue{
				Type:  BT_BOOLEAN,
				Value: genericLesser[rune](l.Value, r.Value),
			}
		default:
			panic("invalid type for equals operator")
		}
	case BO_GREATER_EQUALS:
		switch l.Type {
		case BT_INT8:
			return &BinaryTypedValue{
				Type:  BT_BOOLEAN,
				Value: genericGreaterEquals[int8](l.Value, r.Value),
			}
		case BT_INT16:
			return &BinaryTypedValue{
				Type:  BT_BOOLEAN,
				Value: genericGreaterEquals[int16](l.Value, r.Value),
			}
		case BT_INT32:
			return &BinaryTypedValue{
				Type:  BT_BOOLEAN,
				Value: genericGreaterEquals[int32](l.Value, r.Value),
			}
		case BT_INT64:
			return &BinaryTypedValue{
				Type:  BT_BOOLEAN,
				Value: genericGreaterEquals[int64](l.Value, r.Value),
			}
		case BT_UINT8:
			return &BinaryTypedValue{
				Type:  BT_BOOLEAN,
				Value: genericGreaterEquals[uint8](l.Value, r.Value),
			}
		case BT_UINT16:
			return &BinaryTypedValue{
				Type:  BT_BOOLEAN,
				Value: genericGreaterEquals[uint16](l.Value, r.Value),
			}
		case BT_UINT32:
			return &BinaryTypedValue{
				Type:  BT_BOOLEAN,
				Value: genericGreaterEquals[uint32](l.Value, r.Value),
			}
		case BT_UINT64:
			return &BinaryTypedValue{
				Type:  BT_BOOLEAN,
				Value: genericGreaterEquals[uint64](l.Value, r.Value),
			}
		case BT_BYTE:
			return &BinaryTypedValue{
				Type:  BT_BOOLEAN,
				Value: genericGreaterEquals[byte](l.Value, r.Value),
			}
		case BT_FLOAT32:
			return &BinaryTypedValue{
				Type:  BT_BOOLEAN,
				Value: genericGreaterEquals[float32](l.Value, r.Value),
			}
		case BT_FLOAT64:
			return &BinaryTypedValue{
				Type:  BT_BOOLEAN,
				Value: genericGreaterEquals[float64](l.Value, r.Value),
			}
		case BT_CHAR:
			return &BinaryTypedValue{
				Type:  BT_BOOLEAN,
				Value: genericGreaterEquals[rune](l.Value, r.Value),
			}
		default:
			panic("invalid type for equals operator")
		}
	case BO_LESSER_EQUALS:
		switch l.Type {
		case BT_INT8:
			return &BinaryTypedValue{
				Type:  BT_BOOLEAN,
				Value: genericLesserEquals[int8](l.Value, r.Value),
			}
		case BT_INT16:
			return &BinaryTypedValue{
				Type:  BT_BOOLEAN,
				Value: genericLesserEquals[int16](l.Value, r.Value),
			}
		case BT_INT32:
			return &BinaryTypedValue{
				Type:  BT_BOOLEAN,
				Value: genericLesserEquals[int32](l.Value, r.Value),
			}
		case BT_INT64:
			return &BinaryTypedValue{
				Type:  BT_BOOLEAN,
				Value: genericLesserEquals[int64](l.Value, r.Value),
			}
		case BT_UINT8:
			return &BinaryTypedValue{
				Type:  BT_BOOLEAN,
				Value: genericLesserEquals[uint8](l.Value, r.Value),
			}
		case BT_UINT16:
			return &BinaryTypedValue{
				Type:  BT_BOOLEAN,
				Value: genericLesserEquals[uint16](l.Value, r.Value),
			}
		case BT_UINT32:
			return &BinaryTypedValue{
				Type:  BT_BOOLEAN,
				Value: genericLesserEquals[uint32](l.Value, r.Value),
			}
		case BT_UINT64:
			return &BinaryTypedValue{
				Type:  BT_BOOLEAN,
				Value: genericLesserEquals[uint64](l.Value, r.Value),
			}
		case BT_BYTE:
			return &BinaryTypedValue{
				Type:  BT_BOOLEAN,
				Value: genericLesserEquals[byte](l.Value, r.Value),
			}
		case BT_FLOAT32:
			return &BinaryTypedValue{
				Type:  BT_BOOLEAN,
				Value: genericLesserEquals[float32](l.Value, r.Value),
			}
		case BT_FLOAT64:
			return &BinaryTypedValue{
				Type:  BT_BOOLEAN,
				Value: genericLesserEquals[float64](l.Value, r.Value),
			}
		case BT_CHAR:
			return &BinaryTypedValue{
				Type:  BT_BOOLEAN,
				Value: genericLesserEquals[rune](l.Value, r.Value),
			}
		default:
			panic("invalid type for equals operator")
		}
	default:
		panic("unrecognized operator")
	}
}

func genericEquals[T comparable](l any, r any) bool {
	resL := l.(T)
	resR := r.(T)
	return resL == resR
}

func genericGreater[T Numeric](l any, r any) bool {
	resL := l.(T)
	resR := r.(T)
	return resL > resR
}

func genericLesser[T Numeric](l any, r any) bool {
	resL := l.(T)
	resR := r.(T)
	return resL < resR
}

func genericGreaterEquals[T Numeric](l any, r any) bool {
	resL := l.(T)
	resR := r.(T)
	return resL >= resR
}

func genericLesserEquals[T Numeric](l any, r any) bool {
	resL := l.(T)
	resR := r.(T)
	return resL <= resR
}
