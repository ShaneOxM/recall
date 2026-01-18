package protocol

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"
)

// generateID creates a unique ID with timestamp prefix for sortability.
// Format: {unix_millis}-{random_hex}
func generateID() string {
	timestamp := time.Now().UnixMilli()
	randomBytes := make([]byte, 4)
	rand.Read(randomBytes)
	return fmt.Sprintf("%d-%s", timestamp, hex.EncodeToString(randomBytes))
}
