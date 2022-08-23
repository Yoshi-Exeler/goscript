package goscript

// applyOperator applies the specified operator to the specified values, assuming that the operation has been type checked before
func applyOperator(l *BinaryTypedValue, r *BinaryTypedValue, op BinaryOperator, v *BinaryTypedValue) *BinaryTypedValue {
	switch op {
	case BO_PLUS:
		switch l.Type {
		case BT_INT8:
			genericPlus[int8](l.Value, r.Value, v)
			return v
		case BT_INT16:
			genericPlus[int16](l.Value, r.Value, v)
			return v
		case BT_INT32:
			genericPlus[int32](l.Value, r.Value, v)
			return v
		case BT_INT64:
			genericPlus[int64](l.Value, r.Value, v)
			return v
		case BT_UINT8:
			genericPlus[uint8](l.Value, r.Value, v)
			return v
		case BT_UINT16:
			genericPlus[uint16](l.Value, r.Value, v)
			return v
		case BT_UINT32:
			genericPlus[uint32](l.Value, r.Value, v)
			return v
		case BT_UINT64:
			genericPlus[uint64](l.Value, r.Value, v)
			return v
		case BT_BYTE:
			genericPlus[byte](l.Value, r.Value, v)
			return v
		case BT_FLOAT32:
			genericPlus[float32](l.Value, r.Value, v)
			return v
		case BT_FLOAT64:
			genericPlus[float64](l.Value, r.Value, v)
			return v
		default:
			panic("invalid type for plus operator")
		}
	case BO_MINUS:
		switch l.Type {
		case BT_INT8:
			genericMinus[int8](l.Value, r.Value, v)
			return v
		case BT_INT16:
			genericMinus[int16](l.Value, r.Value, v)
			return v
		case BT_INT32:
			genericMinus[int32](l.Value, r.Value, v)
			return v
		case BT_INT64:
			genericMinus[int64](l.Value, r.Value, v)
			return v
		case BT_UINT8:
			genericMinus[uint8](l.Value, r.Value, v)
			return v
		case BT_UINT16:
			genericMinus[uint16](l.Value, r.Value, v)
			return v
		case BT_UINT32:
			genericMinus[uint32](l.Value, r.Value, v)
			return v
		case BT_UINT64:
			genericMinus[uint64](l.Value, r.Value, v)
			return v
		case BT_BYTE:
			genericMinus[byte](l.Value, r.Value, v)
			return v
		case BT_FLOAT32:
			genericMinus[float32](l.Value, r.Value, v)
			return v
		case BT_FLOAT64:
			genericMinus[float64](l.Value, r.Value, v)
			return v
		default:
			panic("invalid type for minus operator")
		}
	case BO_MULTIPLY:
		switch l.Type {
		case BT_INT8:
			genericMultiply[int8](l.Value, r.Value, v.Value)
			return v
		case BT_INT16:
			genericMultiply[int16](l.Value, r.Value, v.Value)
			return v
		case BT_INT32:
			genericMultiply[int32](l.Value, r.Value, v.Value)
			return v
		case BT_INT64:
			genericMultiply[int64](l.Value, r.Value, v.Value)
			return v
		case BT_UINT8:
			genericMultiply[uint8](l.Value, r.Value, v.Value)
			return v
		case BT_UINT16:
			genericMultiply[uint16](l.Value, r.Value, v.Value)
			return v
		case BT_UINT32:
			genericMultiply[uint32](l.Value, r.Value, v.Value)
			return v
		case BT_UINT64:
			genericMultiply[uint64](l.Value, r.Value, v.Value)
			return v
		case BT_BYTE:
			genericMultiply[byte](l.Value, r.Value, v.Value)
			return v
		case BT_FLOAT32:
			genericMultiply[float32](l.Value, r.Value, v.Value)
			return v
		case BT_FLOAT64:
			genericMultiply[float64](l.Value, r.Value, v.Value)
			return v
		default:
			panic("invalid type for multiply operator")
		}
	case BO_DIVIDE:
		switch l.Type {
		case BT_INT8:
			genericDivide[int8](l.Value, r.Value, v)
			return v
		case BT_INT16:
			genericDivide[int16](l.Value, r.Value, v)
			return v
		case BT_INT32:
			genericDivide[int32](l.Value, r.Value, v)
			return v
		case BT_INT64:
			genericDivide[int64](l.Value, r.Value, v)
			return v
		case BT_UINT8:
			genericDivide[uint8](l.Value, r.Value, v)
			return v
		case BT_UINT16:
			genericDivide[uint16](l.Value, r.Value, v)
			return v
		case BT_UINT32:
			genericDivide[uint32](l.Value, r.Value, v)
			return v
		case BT_UINT64:
			genericDivide[uint64](l.Value, r.Value, v)
			return v
		case BT_BYTE:
			genericDivide[byte](l.Value, r.Value, v)
			return v
		case BT_FLOAT32:
			genericDivide[float32](l.Value, r.Value, v)
			return v
		case BT_FLOAT64:
			genericDivide[float64](l.Value, r.Value, v)
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
	return *l.(*T) == *r.(*T)
}

func genericGreater[T Numeric](l any, r any) bool {
	return *l.(*T) > *r.(*T)
}

func genericLesser[T Numeric](l any, r any) bool {
	lptr := l.(*T)
	rptr := r.(*T)
	return *lptr < *rptr

}

func genericGreaterEquals[T Numeric](l any, r any) bool {
	return *l.(*T) >= *r.(*T)
}

func genericLesserEquals[T Numeric](l any, r any) bool {
	return *l.(*T) <= *r.(*T)
}

func genericPlus[T Numeric](l any, r any, v *BinaryTypedValue) {
	*v.Value.(*T) = *l.(*T) + *r.(*T)
}

func genericMinus[T Numeric](l any, r any, v *BinaryTypedValue) {
	*v.Value.(*T) = *l.(*T) - *r.(*T)
}

func genericMultiply[T Numeric](l any, r any, v any) {
	vptr := v.(*T)
	lptr := l.(*T)
	rptr := r.(*T)
	*vptr = *lptr * *rptr
}

func genericDivide[T Numeric](l any, r any, v *BinaryTypedValue) {
	result := float64(*l.(*T)) / float64(*r.(*T))
	v.Value = &result
}
