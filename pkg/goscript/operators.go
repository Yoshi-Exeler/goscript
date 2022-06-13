package goscript

// applyOperator applies the specified operator to the specified values, assuming that the operation has been type checked before
func applyOperator(l *BinaryTypedValue, r *BinaryTypedValue, op BinaryOperator, v *BinaryTypedValue) *BinaryTypedValue {
	switch op {
	case BO_PLUS:
		switch l.Type {
		case BT_INT8:
			v.Value = genericPlus[int8](l.Value, r.Value)
			return v
		case BT_INT16:
			v.Value = genericPlus[int16](l.Value, r.Value)
			return v
		case BT_INT32:
			v.Value = genericPlus[int32](l.Value, r.Value)
			return v
		case BT_INT64:
			v.Value = genericPlus[int64](l.Value, r.Value)
			return v
		case BT_UINT8:
			v.Value = genericPlus[uint8](l.Value, r.Value)
			return v
		case BT_UINT16:
			v.Value = genericPlus[uint16](l.Value, r.Value)
			return v
		case BT_UINT32:
			v.Value = genericPlus[uint32](l.Value, r.Value)
			return v
		case BT_UINT64:
			v.Value = genericPlus[uint64](l.Value, r.Value)
			return v
		case BT_BYTE:
			v.Value = genericPlus[byte](l.Value, r.Value)
			return v
		case BT_FLOAT32:
			v.Value = genericPlus[float32](l.Value, r.Value)
			return v
		case BT_FLOAT64:
			v.Value = genericPlus[float64](l.Value, r.Value)
			return v
		default:
			panic("invalid type for plus operator")
		}
	case BO_MINUS:
		switch l.Type {
		case BT_INT8:
			v.Value = genericMinus[int8](l.Value, r.Value)
			return v
		case BT_INT16:
			v.Value = genericMinus[int16](l.Value, r.Value)
			return v
		case BT_INT32:
			v.Value = genericMinus[int32](l.Value, r.Value)
			return v
		case BT_INT64:
			v.Value = genericMinus[int64](l.Value, r.Value)
			return v
		case BT_UINT8:
			v.Value = genericMinus[uint8](l.Value, r.Value)
			return v
		case BT_UINT16:
			v.Value = genericMinus[uint16](l.Value, r.Value)
			return v
		case BT_UINT32:
			v.Value = genericMinus[uint32](l.Value, r.Value)
			return v
		case BT_UINT64:
			v.Value = genericMinus[uint64](l.Value, r.Value)
			return v
		case BT_BYTE:
			v.Value = genericMinus[byte](l.Value, r.Value)
			return v
		case BT_FLOAT32:
			v.Value = genericMinus[float32](l.Value, r.Value)
			return v
		case BT_FLOAT64:
			v.Value = genericMinus[float64](l.Value, r.Value)
			return v
		default:
			panic("invalid type for minus operator")
		}
	case BO_MULTIPLY:
		switch l.Type {
		case BT_INT8:
			v.Value = genericMultiply[int8](l.Value, r.Value)
			return v
		case BT_INT16:
			v.Value = genericMultiply[int16](l.Value, r.Value)
			return v
		case BT_INT32:
			v.Value = genericMultiply[int32](l.Value, r.Value)
			return v
		case BT_INT64:
			v.Value = genericMultiply[int64](l.Value, r.Value)
			return v
		case BT_UINT8:
			v.Value = genericMultiply[uint8](l.Value, r.Value)
			return v
		case BT_UINT16:
			v.Value = genericMultiply[uint16](l.Value, r.Value)
			return v
		case BT_UINT32:
			v.Value = genericMultiply[uint32](l.Value, r.Value)
			return v
		case BT_UINT64:
			v.Value = genericMultiply[uint64](l.Value, r.Value)
			return v
		case BT_BYTE:
			v.Value = genericMultiply[byte](l.Value, r.Value)
			return v
		case BT_FLOAT32:
			v.Value = genericMultiply[float32](l.Value, r.Value)
			return v
		case BT_FLOAT64:
			v.Value = genericMultiply[float64](l.Value, r.Value)
			return v
		default:
			panic("invalid type for multiply operator")
		}
	case BO_DIVIDE:
		switch l.Type {
		case BT_INT8:
			v.Value = genericDivide[int8](l.Value, r.Value)
			return v
		case BT_INT16:
			v.Value = genericDivide[int16](l.Value, r.Value)
			return v
		case BT_INT32:
			v.Value = genericDivide[int32](l.Value, r.Value)
			return v
		case BT_INT64:
			v.Value = genericDivide[int64](l.Value, r.Value)
			return v
		case BT_UINT8:
			v.Value = genericDivide[uint8](l.Value, r.Value)
			return v
		case BT_UINT16:
			v.Value = genericDivide[uint16](l.Value, r.Value)
			return v
		case BT_UINT32:
			v.Value = genericDivide[uint32](l.Value, r.Value)
			return v
		case BT_UINT64:
			v.Value = genericDivide[uint64](l.Value, r.Value)
			return v
		case BT_BYTE:
			v.Value = genericDivide[byte](l.Value, r.Value)
			return v
		case BT_FLOAT32:
			v.Value = genericDivide[float32](l.Value, r.Value)
			return v
		case BT_FLOAT64:
			v.Value = genericDivide[float64](l.Value, r.Value)
			return v
		default:
			panic("invalid type for divide operator")
		}
	case BO_EQUALS:
		switch l.Type {
		case BT_INT8:
			v.Value = genericEquals[int8](l.Value, r.Value)
			return v
		case BT_INT16:
			v.Value = genericEquals[int16](l.Value, r.Value)
			return v
		case BT_INT32:
			v.Value = genericEquals[int32](l.Value, r.Value)
			return v
		case BT_INT64:
			v.Value = genericEquals[int64](l.Value, r.Value)
			return v
		case BT_UINT8:
			v.Value = genericEquals[uint8](l.Value, r.Value)
			return v
		case BT_UINT16:
			v.Value = genericEquals[uint16](l.Value, r.Value)
			return v
		case BT_UINT32:
			v.Value = genericEquals[uint32](l.Value, r.Value)
			return v
		case BT_UINT64:
			v.Value = genericEquals[uint64](l.Value, r.Value)
			return v
		case BT_BYTE:
			v.Value = genericEquals[byte](l.Value, r.Value)
			return v
		case BT_FLOAT32:
			v.Value = genericEquals[float32](l.Value, r.Value)
			return v
		case BT_FLOAT64:
			v.Value = genericEquals[float64](l.Value, r.Value)
			return v
		case BT_CHAR:
			v.Value = genericEquals[rune](l.Value, r.Value)
			return v
		case BT_STRING:
			v.Value = genericEquals[string](l.Value, r.Value)
			return v
		case BT_BOOLEAN:
			v.Value = genericEquals[int8](l.Value, r.Value)
			return v
		default:
			panic("invalid type for equals operator")
		}
	case BO_GREATER:
		switch l.Type {
		case BT_INT8:
			v.Value = genericGreater[int8](l.Value, r.Value)
			return v
		case BT_INT16:
			v.Value = genericGreater[int16](l.Value, r.Value)
			return v
		case BT_INT32:
			v.Value = genericGreater[int32](l.Value, r.Value)
			return v
		case BT_INT64:
			v.Value = genericGreater[int64](l.Value, r.Value)
			return v
		case BT_UINT8:
			v.Value = genericGreater[uint8](l.Value, r.Value)
			return v
		case BT_UINT16:
			v.Value = genericGreater[uint16](l.Value, r.Value)
			return v
		case BT_UINT32:
			v.Value = genericGreater[uint32](l.Value, r.Value)
			return v
		case BT_UINT64:
			v.Value = genericGreater[uint64](l.Value, r.Value)
			return v
		case BT_BYTE:
			v.Value = genericGreater[byte](l.Value, r.Value)
			return v
		case BT_FLOAT32:
			v.Value = genericGreater[float32](l.Value, r.Value)
			return v
		case BT_FLOAT64:
			v.Value = genericGreater[float64](l.Value, r.Value)
			return v
		case BT_CHAR:
			v.Value = genericGreater[rune](l.Value, r.Value)
			return v
		default:
			panic("invalid type for equals operator")
		}
	case BO_LESSER:
		switch l.Type {
		case BT_INT8:
			v.Value = genericLesser[int8](l.Value, r.Value)
			return v
		case BT_INT16:
			v.Value = genericLesser[int16](l.Value, r.Value)
			return v
		case BT_INT32:
			v.Value = genericLesser[int32](l.Value, r.Value)
			return v
		case BT_INT64:
			v.Value = genericLesser[int64](l.Value, r.Value)
			return v
		case BT_UINT8:
			v.Value = genericLesser[uint8](l.Value, r.Value)
			return v
		case BT_UINT16:
			v.Value = genericLesser[uint16](l.Value, r.Value)
			return v
		case BT_UINT32:
			v.Value = genericLesser[uint32](l.Value, r.Value)
			return v
		case BT_UINT64:
			v.Value = genericLesser[uint64](l.Value, r.Value)
			return v
		case BT_BYTE:
			v.Value = genericLesser[byte](l.Value, r.Value)
			return v
		case BT_FLOAT32:
			v.Value = genericLesser[float32](l.Value, r.Value)
			return v
		case BT_FLOAT64:
			v.Value = genericLesser[float64](l.Value, r.Value)
			return v
		case BT_CHAR:
			v.Value = genericLesser[rune](l.Value, r.Value)
			return v
		default:
			panic("invalid type for equals operator")
		}
	case BO_GREATER_EQUALS:
		switch l.Type {
		case BT_INT8:
			v.Value = genericGreaterEquals[int8](l.Value, r.Value)
			return v
		case BT_INT16:
			v.Value = genericGreaterEquals[int16](l.Value, r.Value)
			return v
		case BT_INT32:
			v.Value = genericGreaterEquals[int32](l.Value, r.Value)
			return v
		case BT_INT64:
			v.Value = genericGreaterEquals[int64](l.Value, r.Value)
			return v
		case BT_UINT8:
			v.Value = genericGreaterEquals[uint8](l.Value, r.Value)
			return v
		case BT_UINT16:
			v.Value = genericGreaterEquals[uint16](l.Value, r.Value)
			return v
		case BT_UINT32:
			v.Value = genericGreaterEquals[uint32](l.Value, r.Value)
			return v
		case BT_UINT64:
			v.Value = genericGreaterEquals[uint64](l.Value, r.Value)
			return v
		case BT_BYTE:
			v.Value = genericGreaterEquals[byte](l.Value, r.Value)
			return v
		case BT_FLOAT32:
			v.Value = genericGreaterEquals[float32](l.Value, r.Value)
			return v
		case BT_FLOAT64:
			v.Value = genericGreaterEquals[float64](l.Value, r.Value)
			return v
		case BT_CHAR:
			v.Value = genericGreaterEquals[rune](l.Value, r.Value)
			return v
		default:
			panic("invalid type for equals operator")
		}
	case BO_LESSER_EQUALS:
		switch l.Type {
		case BT_INT8:
			v.Value = genericLesserEquals[int8](l.Value, r.Value)
			return v
		case BT_INT16:
			v.Value = genericLesserEquals[int16](l.Value, r.Value)
			return v
		case BT_INT32:
			v.Value = genericLesserEquals[int32](l.Value, r.Value)
			return v
		case BT_INT64:
			v.Value = genericLesserEquals[int64](l.Value, r.Value)
			return v
		case BT_UINT8:
			v.Value = genericLesserEquals[uint8](l.Value, r.Value)
			return v
		case BT_UINT16:
			v.Value = genericLesserEquals[uint16](l.Value, r.Value)
			return v
		case BT_UINT32:
			v.Value = genericLesserEquals[uint32](l.Value, r.Value)
			return v
		case BT_UINT64:
			v.Value = genericLesserEquals[uint64](l.Value, r.Value)
			return v
		case BT_BYTE:
			v.Value = genericLesserEquals[byte](l.Value, r.Value)
			return v
		case BT_FLOAT32:
			v.Value = genericLesserEquals[float32](l.Value, r.Value)
			return v
		case BT_FLOAT64:
			v.Value = genericLesserEquals[float64](l.Value, r.Value)
			return v
		case BT_CHAR:
			v.Value = genericLesserEquals[rune](l.Value, r.Value)
			return v
		default:
			panic("invalid type for equals operator")
		}
	default:
		panic("unrecognized operator")
	}
}

func genericEquals[T comparable](l any, r any) bool {
	return l.(T) == r.(T)
}

func genericGreater[T Numeric](l any, r any) bool {
	return l.(T) > r.(T)
}

func genericLesser[T Numeric](l any, r any) bool {
	return l.(T) < r.(T)
}

func genericGreaterEquals[T Numeric](l any, r any) bool {
	return l.(T) >= r.(T)
}

func genericLesserEquals[T Numeric](l any, r any) bool {
	return l.(T) <= r.(T)
}

func genericPlus[T Numeric](l any, r any) T {
	return l.(T) + r.(T)
}

func genericMinus[T Numeric](l any, r any) T {
	return l.(T) - r.(T)
}

func genericMultiply[T Numeric](l any, r any) T {
	return l.(T) * r.(T)
}

func genericDivide[T Numeric](l any, r any) float64 {
	return float64(l.(T)) / float64(r.(T))
}
