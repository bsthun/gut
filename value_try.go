package gut

import "fmt"

// Try executes a function and recovers from panic, returning the error instead of crashing.
func Try(fn func()) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("panic recovered: %v", r)
		}
	}()
	fn()
	return nil
}
