package gen

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"
)

func GenerateUniqueFilename(extension string) string {
	timestamp := time.Now().Unix()
	randomBytes := make([]byte, 6)
	rand.Read(randomBytes)
	randomStr := hex.EncodeToString(randomBytes)

	filename := fmt.Sprintf("%d-%s", timestamp, randomStr)
	return filename + "." + extension
}
