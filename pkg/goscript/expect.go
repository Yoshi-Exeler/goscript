package goscript

import "fmt"

func expectLength[T any](arr []T, length int, expectation string) {
	if len(arr) != length {
		panic(fmt.Sprintf("assertion failed: %v. expected length %v but got %v", expectation, length, len(arr)))
	}
}
