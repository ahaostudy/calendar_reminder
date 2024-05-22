package sha256

import (
	"crypto/sha256"
	"encoding/hex"
	"os"
)

const EncryptSaltKey = "ENCRYPT_SALT"

var salt = os.Getenv(EncryptSaltKey)

func Encrypt(password string) string {
	hash := sha256.New()
	hash.Write([]byte(password + salt))
	return hex.EncodeToString(hash.Sum(nil))
}
