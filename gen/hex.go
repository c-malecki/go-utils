package gen

import (
	"crypto/rand"
	"encoding/hex"
	"strings"
)

func GenHexCode() (string, error) {
	b := make([]byte, 25)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return strings.ToUpper(hex.EncodeToString(b)), nil
}
