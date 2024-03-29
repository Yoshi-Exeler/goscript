package goscript

import "fmt"

func expectLength[T any](arr []T, length int, expectation string) {
	if len(arr) != length {
		panic(fmt.Sprintf("assertion failed: %v. expected length %v but got %v", expectation, length, len(arr)))
	}
}

func expectValue[T comparable](value T, expectedValue T) {
	if value != expectedValue {
		panic(fmt.Sprintf("assertion failed expected %v but got %v", value, expectedValue))
	}
}

func expectPanic(handler func()) {
	defer func() {
		err := recover()
		if err == nil {
			panic("expected panic while executing handler but the handler ran successfully")
		}
	}()
	handler()
}
