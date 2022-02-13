package bank

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestPassword(t *testing.T) {
	var err error
	var a, b string
	a, err = SetPassword("secret")
	require.NoError(t, err)
	// Positive case
	assert.NoError(t, VerifyPassword(a, "secret"))
	// Negative cases
	assert.Equal(t, ErrPasswordMismatch, VerifyPassword(a, "alternative"))
	assert.Equal(t, ErrPasswordMismatch, VerifyPassword(a, "Secret"))
	assert.Equal(t, ErrPasswordMismatch, VerifyPassword(a, ""))
	// Variants
	b, err = SetPassword("alternative")
	require.NoError(t, err)
	assert.NotEqual(t, a, b)
	assert.NoError(t, VerifyPassword(b, "alternative"))
	assert.Equal(t, ErrPasswordMismatch, VerifyPassword(b, "secret"))
}
