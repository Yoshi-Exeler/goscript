package goscript

// func TestResolveExpression(t *testing.T) {
// 	// this tree represents 5 + 5 * 2
// 	five := uint8(5)
// 	two := uint8(2)
// 	fiveExpr := &Expression{
// 		Operator: BO_CONSTANT,
// 		Value: &BinaryTypedValue{
// 			Type:  BT_UINT8,
// 			Value: &five,
// 		},
// 	}
// 	twoExpr := &Expression{
// 		Operator: BO_CONSTANT,
// 		Value: &BinaryTypedValue{
// 			Type:  BT_UINT8,
// 			Value: &two,
// 		},
// 	}
// 	multExpr := &Expression{
// 		LeftExpression:  fiveExpr,
// 		RightExpression: twoExpr,
// 		Operator:        BO_MULTIPLY,
// 		Value: &BinaryTypedValue{
// 			Type:  BT_UINT8,
// 			Value: 0,
// 		},
// 	}
// 	plusExpr := &Expression{
// 		LeftExpression:  fiveExpr,
// 		RightExpression: multExpr,
// 		Operator:        BO_PLUS,
// 		Value: &BinaryTypedValue{
// 			Type:  BT_UINT8,
// 			Value: 0,
// 		},
// 	}
// 	// now resolve it
// 	rt := Runtime{}
// 	result := rt.ResolveExpression(plusExpr)
// 	if result.Type != BT_UINT8 {
// 		t.Fatalf("expression resolution failed, expected type %v but got %v", BT_UINT8, result.Type)
// 	}
// 	resV := result.Value.(uint8)
// 	if resV != 15 {
// 		t.Fatalf("expression resolution failed, expected result %v but got %v", 15, resV)
// 	}
// }
