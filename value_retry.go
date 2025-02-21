package gut

import (
	"fmt"
	"time"
)

func Retry(attempts int, delay time.Duration, fn func() error) error {
	var err error
	for i := 0; i < attempts; i++ {
		if err = fn(); err == nil {
			return nil
		}
		time.Sleep(delay)
	}
	return fmt.Errorf("all %d attempts failed: %w", attempts, err)
}
