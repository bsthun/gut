package gut

import (
	"log"

	"github.com/google/uuid"
)

func Ptr[T any](v T) *T {
	return &v
}

// Uint8Ptr converts a UUID string or *string to *[]uint8, with an optional default value.
func Uint8Ptr[T string | *string](id T, defaultValue ...string) (*[]uint8, error) {
	var strID string

	// Handle input types (string or *string)
	switch v := any(id).(type) {
	case string:
		strID = v
	case *string:
		if v == nil {
			return nil, nil // Safe handling for nil pointer
		}
		strID = *v
	default:
		return nil, nil // Invalid type
	}

	// Try parsing the UUID
	uuidBytes, err := uuid.Parse(strID)
	if err != nil {
		log.Printf("Invalid UUID: %v", err) // Log error

		// If a default value is provided, use it
		if len(defaultValue) > 0 {
			uuidBytes, err = uuid.Parse(defaultValue[0])
			if err != nil {
				return nil, err // Even the default is invalid, return error
			}
		} else {
			return nil, err // No default, return error
		}
	}

	bytes := uuidBytes[:]
	return &bytes, nil
}
