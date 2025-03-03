package gut

import (
	"encoding/json"
	"errors"
	"math"
)

const (
	idEncoderCharset       = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	idEncoderEncodedLength = 11 // ceil(log_62(2^64))
	idEncoderBase          = len(idEncoderCharset)
)

type Id uint64

func (id *Id) MarshalJSON() ([]byte, error) {
	println(*id)
	return json.Marshal(EncodeId(uint64(*id)))
}

func (id *Id) UnmarshalJSON(data []byte) error {
	var encoded string
	if err := json.Unmarshal(data, &encoded); err != nil {
		return err
	}

	decoded, err := Decode(encoded)
	if err != nil {
		return err
	}

	*id = Id(decoded)
	return nil
}

var idEncoderKey []byte // 16-byte key

func SetIdEncoderKey(key []byte) error {
	if len(key) != 16 {
		return errors.New("key must be exactly 16 bytes")
	}

	idEncoderKey = key
	return nil
}

func EncodeId(id uint64) string {
	// Create a simple reversible transformation
	// Using a combination of XOR, bit shifting, and addition

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

	// * convert to base62 string
	return toBase62(encrypted)
}

func Decode(encoded string) (uint64, error) {
	// * check if the string has the correct length
	if len(encoded) != idEncoderEncodedLength {
		return 0, errors.New("invalid encoded string length")
	}

	// * convert from base62 string to uint64
	encrypted, err := fromBase62(encoded)
	if err != nil {
		return 0, err
	}

	// * extract 8 uint16 values from our key
	k := make([]uint16, 8)
	for i := 0; i < 8; i++ {
		k[i] = uint16(idEncoderKey[i*2])<<8 | uint16(idEncoderKey[i*2+1])
	}

	// # reverse the transformation

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

func toBase62(num uint64) string {
	result := make([]byte, idEncoderEncodedLength)

	for i := 0; i < idEncoderEncodedLength; i++ {
		result[i] = idEncoderCharset[0]
	}

	i := idEncoderEncodedLength - 1
	for num > 0 && i >= 0 {
		result[i] = idEncoderCharset[num%uint64(idEncoderBase)]
		num /= uint64(idEncoderBase)
		i--
	}

	return string(result)
}

func fromBase62(str string) (uint64, error) {
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
