package gut

import (
	"encoding/json"
	"errors"
	"math"
)

const (
	idEncoderCharset = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	idEncoderLength  = 11
	idEncoderBase    = len(idEncoderCharset)
)

type Id uint64

func (id *Id) MarshalJSON() ([]byte, error) {
	return json.Marshal(IdEncode(uint64(*id)))
}

func (id *Id) UnmarshalJSON(data []byte) error {
	var encoded string
	if err := json.Unmarshal(data, &encoded); err != nil {
		return err
	}

	decoded, err := IdDecode(encoded)
	if err != nil {
		return err
	}

	*id = Id(decoded)
	return nil
}

var idEncoderKey []byte // 16-byte key
var idEncoderPaddingLength = 1

func SetIdEncoderKey(key []byte) error {
	if len(key) != 16 {
		return errors.New("key must be exactly 16 bytes")
	}

	idEncoderKey = key
	return nil
}

func SetIdEncoderPaddingLength(length int) {
	idEncoderPaddingLength = length
}

func IdEncode(id uint64) string {
	// * extract 8 uint16 values from our key
	k := make([]uint16, 8)
	for i := 0; i < 8; i++ {
		k[i] = uint16(idEncoderKey[i*2])<<8 | uint16(idEncoderKey[i*2+1])
	}

	// * apply the transformation
	encrypted := id

	// * xor with first key part
	encrypted ^= uint64(k[0])<<48 | uint64(k[1])<<32 | uint64(k[2])<<16 | uint64(k[3])

	// * rotate bits right by 17
	encrypted = (encrypted >> 17) | (encrypted << (64 - 17))

	// * xor with second key part
	encrypted ^= uint64(k[4])<<48 | uint64(k[5])<<32 | uint64(k[6])<<16 | uint64(k[7])

	// * rotate bits left by 13
	encrypted = (encrypted << 13) | (encrypted >> (64 - 13))

	// * add the first key part
	encrypted += uint64(k[0])<<48 | uint64(k[1])<<32 | uint64(k[2])<<16 | uint64(k[3])

	// * convert to string
	base := Base62(encrypted)

	// * generate padding
	padding := EncodePadding(encrypted, k)

	return base + padding
}

func IdDecode(encoded string) (uint64, error) {
	// * check if the string has the correct length
	if len(encoded) != idEncoderLength+idEncoderPaddingLength {
		return 0, errors.New("invalid encoded string length")
	}

	// * split into base part and padding
	basePart := encoded[:11]
	paddingPart := encoded[11:]

	// * convert from to uint64
	encrypted, err := Base62Parse(basePart)
	if err != nil {
		return 0, err
	}

	// * extract 8 uint16 values from our key
	k := make([]uint16, 8)
	for i := 0; i < 8; i++ {
		k[i] = uint16(idEncoderKey[i*2])<<8 | uint16(idEncoderKey[i*2+1])
	}

	// * verify padding
	expectedPadding := EncodePadding(encrypted, k)
	if paddingPart != expectedPadding {
		return 0, errors.New("invalid encoded string: padding verification failed")
	}

	// * subtract the first key part
	encrypted -= uint64(k[0])<<48 | uint64(k[1])<<32 | uint64(k[2])<<16 | uint64(k[3])

	// * rotate bits right by 13
	encrypted = (encrypted >> 13) | (encrypted << (64 - 13))

	// * xor with second key part
	encrypted ^= uint64(k[4])<<48 | uint64(k[5])<<32 | uint64(k[6])<<16 | uint64(k[7])

	// * rotate bits left by 17
	encrypted = (encrypted << 17) | (encrypted >> (64 - 17))

	// * xor with first key part
	encrypted ^= uint64(k[0])<<48 | uint64(k[1])<<32 | uint64(k[2])<<16 | uint64(k[3])

	return encrypted, nil
}

func EncodePadding(encrypted uint64, k []uint16) string {
	// * create deterministic padding with more entropy sources
	padValue := encrypted

	// * round 1: mix with rotating key pattern
	padValue ^= uint64(k[7])<<48 | uint64(k[6])<<32 | uint64(k[5])<<16 | uint64(k[4])
	padValue = (padValue << 11) | (padValue >> (64 - 11))

	// * round 2: add cross-key mixing
	padValue ^= uint64(k[3])<<40 | uint64(k[2])<<24 | uint64(k[1])<<8 | uint64(k[0])
	padValue = (padValue >> 19) | (padValue << (64 - 19))

	// * round 3: final entropy mixing with bit shifts
	padValue ^= (padValue >> 23) ^ (padValue >> 31) ^ (padValue >> 41)
	padValue += uint64(k[0])<<32 | uint64(k[7])<<16 | uint64(k[3])

	// * generate padding with additional variety
	result := make([]byte, idEncoderPaddingLength)

	// * generate each character
	for i := 0; i < idEncoderPaddingLength; i++ {
		result[i] = idEncoderCharset[padValue%uint64(idEncoderBase)]
		padValue = (padValue >> 7) ^ (padValue << 3) ^ uint64(k[i%8])
	}

	return string(result)
}

func Base62(num uint64) string {
	result := make([]byte, idEncoderLength)

	for i := 0; i < 11; i++ {
		result[i] = idEncoderCharset[0]
	}

	i := 10
	for num > 0 && i >= 0 {
		result[i] = idEncoderCharset[num%uint64(idEncoderBase)]
		num /= uint64(idEncoderBase)
		i--
	}

	return string(result)
}

func Base62Parse(str string) (uint64, error) {
	var result uint64 = 0

	for i := 0; i < len(str); i++ {
		var val int = -1
		for j := 0; j < idEncoderBase; j++ {
			if idEncoderCharset[j] == str[i] {
				val = j
				break
			}
		}

		if val == -1 {
			return 0, errors.New("invalid character in encoded string")
		}

		if result > math.MaxUint64/uint64(idEncoderBase) {
			return 0, errors.New("overflow occurred during decoding")
		}

		result = result*uint64(idEncoderBase) + uint64(val)
	}

	return result, nil
}
