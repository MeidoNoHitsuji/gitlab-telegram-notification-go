package toggl

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"strings"
)

func HmacIsValid(message, signature, secret string) bool {
	messageMAC, _ := hex.DecodeString(strings.TrimPrefix(signature, "sha256="))

	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(message))
	expectedMAC := mac.Sum(nil)

	return hmac.Equal([]byte(messageMAC), expectedMAC)
}
