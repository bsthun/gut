package gut

import (
	"crypto/rand"
	"math"
	"testing"
)

func TestIdEncodingAllValues(t *testing.T) {
	// * set up a test key
	key := make([]byte, 16)
	_, err := rand.Read(key)
	if err != nil {
		t.Fatalf("failed to generate random key: %v", err)
	}

	err = SetIdEncoderKey(key)
	if err != nil {
		t.Fatalf("failed to set encoder key: %v", err)
	}

	SetIdEncoderPaddingLength(3)

	// * track encoded strings
	seen := make(map[string]uint64)

	// * test all uint64 values
	for id := uint64(0); id < math.MaxUint64; id++ {
		encoded := IdEncode(id)

		// * check for collision
		if originalId, exists := seen[encoded]; exists {
			t.Errorf("collision detected: id %d and id %d both encode to %s",
				originalId, id, encoded)
			return
		}
		seen[encoded] = id

		// * verify round trip
		decoded, err := IdDecode(encoded)
		if err != nil {
			t.Errorf("failed to decode %s: %v", encoded, err)
			return
		}
		if decoded != id {
			t.Errorf("round trip failed: %d -> %s -> %d", id, encoded, decoded)
			return
		}

		// * progress logging
		if id > 0 && id%1000000 == 0 {
			t.Logf("tested %d values, no collisions so far", id)
		}
	}

	// * test the final value (math.MaxUint64)
	id := uint64(math.MaxUint64)
	encoded := IdEncode(id)

	if originalId, exists := seen[encoded]; exists {
		t.Errorf("collision detected: id %d and id %d both encode to %s",
			originalId, id, encoded)
		return
	}

	decoded, err := IdDecode(encoded)
	if err != nil {
		t.Errorf("failed to decode %s: %v", encoded, err)
		return
	}
	if decoded != id {
		t.Errorf("round trip failed: %d -> %s -> %d", id, encoded, decoded)
		return
	}

	t.Logf("successfully tested all %d possible uint64 values with no collisions",
		uint64(math.MaxUint64))
}
