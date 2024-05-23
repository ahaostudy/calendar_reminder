package jwt

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestJWT(t *testing.T) {
	id := uint(123456)

	token, err := GenerateToken(id)
	assert.NoError(t, err)

	parseId, ok := ParseToken(token)
	assert.True(t, ok)
	assert.Equal(t, id, parseId)
}
