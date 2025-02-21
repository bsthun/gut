package gut

import (
	crand "crypto/rand"
	"math/big"
	"math/rand"
	"strings"
	"time"
)

var Rand = rand.New(rand.NewSource(time.Now().UnixNano()))

func Random(characters string, number int) *string {
	var generated strings.Builder
	for range number {
		random := Rand.Intn(len(characters))
		randomChar := characters[random]
		generated.WriteString(string(randomChar))
	}

	var str = generated.String()
	return &str
}

func RandomSecure(characters string, number int) *string {
	var generated strings.Builder
	for range number {
		index, _ := crand.Int(crand.Reader, big.NewInt(int64(len(characters))))
		generated.WriteByte(characters[index.Int64()])
	}

	result := generated.String()
	return &result
}

var RandomSet = struct {
	Num           string
	MixedAlphaNum string
	UpperAlpha    string
	UpperAlphaNum string
}{
	"0123456789",
	"ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789",
	"ABCDEFGHIJKLMNOPQRSTUVWXYZ",
	"0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ",
}
