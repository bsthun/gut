package gut

import "fmt"

// Must ensures that a function call succeeds by panicking if an error is returned.
// It works with a single return value and an error.
func Must[T any](val T, err error) T {
	if err != nil {
		panic(fmt.Sprintf("Must failed: %v", err))
	}
	return val
}

// Must2 is a generic version for functions returning two values and an error.
func Must2[T1, T2 any](val1 T1, val2 T2, err error) (T1, T2) {
	if err != nil {
		panic(fmt.Sprintf("Must2 failed: %v", err))
	}
	return val1, val2
}

// Must3 is a generic version for functions returning three values and an error.
func Must3[T1, T2, T3 any](val1 T1, val2 T2, val3 T3, err error) (T1, T2, T3) {
	if err != nil {
		panic(fmt.Sprintf("Must3 failed: %v", err))
	}
	return val1, val2, val3
}
