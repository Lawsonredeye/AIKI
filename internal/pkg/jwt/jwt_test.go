package jwt

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenerateAccessToken(t *testing.T) {
	secret := "test-secret"
	manager := NewManager(secret, 15*time.Minute, 7*24*time.Hour)

	userID := int32(123)
	email := "test@example.com"

	token, err := manager.GenerateAccessToken(userID, email)
	require.NoError(t, err)
	assert.NotEmpty(t, token)
}

func TestGenerateRefreshToken(t *testing.T) {
	manager := NewManager("test-secret", 15*time.Minute, 7*24*time.Hour)

	token1 := manager.GenerateRefreshToken()
	token2 := manager.GenerateRefreshToken()

	assert.NotEmpty(t, token1)
	assert.NotEmpty(t, token2)
	assert.NotEqual(t, token1, token2) // Tokens should be unique
}

func TestValidateToken(t *testing.T) {
	secret := "test-secret"
	manager := NewManager(secret, 15*time.Minute, 7*24*time.Hour)

	userID := int32(123)
	email := "test@example.com"

	token, err := manager.GenerateAccessToken(userID, email)
	require.NoError(t, err)

	t.Run("valid token", func(t *testing.T) {
		claims, err := manager.ValidateToken(token)
		require.NoError(t, err)
		assert.Equal(t, userID, claims.UserID)
		assert.Equal(t, email, claims.Email)
	})

	t.Run("invalid token", func(t *testing.T) {
		_, err := manager.ValidateToken("invalid-token")
		assert.Error(t, err)
	})

	t.Run("expired token", func(t *testing.T) {
		expiredManager := NewManager(secret, -1*time.Hour, 7*24*time.Hour)
		expiredToken, err := expiredManager.GenerateAccessToken(userID, email)
		require.NoError(t, err)

		_, err = manager.ValidateToken(expiredToken)
		assert.Error(t, err)
	})

	t.Run("wrong secret", func(t *testing.T) {
		wrongManager := NewManager("wrong-secret", 15*time.Minute, 7*24*time.Hour)
		_, err := wrongManager.ValidateToken(token)
		assert.Error(t, err)
	})
}

func TestGetRefreshTokenExpiry(t *testing.T) {
	refreshExpiry := 7 * 24 * time.Hour
	manager := NewManager("test-secret", 15*time.Minute, refreshExpiry)

	assert.Equal(t, refreshExpiry, manager.GetRefreshTokenExpiry())
}
