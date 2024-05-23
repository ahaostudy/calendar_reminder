package sha256

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSha256(t *testing.T) {
	k := Encrypt("123456")
	assert.Equal(t, k, "8d969eef6ecad3c29a3a629280e686cf0c3f5d5a86aff3ca12020c923adc6c92")
}
