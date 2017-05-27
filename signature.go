package middlesign

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
)

// sign signs source with secret
func sign(source, secret string) []byte {
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(source))
	return mac.Sum(nil)
}

// SignString signs source with secret and return HEX representation of signature
func SignString(source, secret string) string {
	return hex.EncodeToString(sign(source, secret))
}

// IsSignatureValid validates signature using source and key
func IsSignatureValid(source, secret, signature string) (bool, error) {
	expected := sign(source, secret)
	sign, err := hex.DecodeString(signature)
	if err != nil {
		return false, err
	}

	return hmac.Equal(sign, expected), nil
}
