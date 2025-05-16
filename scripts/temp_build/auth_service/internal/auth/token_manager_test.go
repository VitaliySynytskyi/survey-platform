package auth

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewTokenManager(t *testing.T) {
	secret := "test_secret"
	tokenTTL := 15 * time.Minute
	manager := NewTokenManager(secret, tokenTTL)

	assert.NotNil(t, manager)
	assert.Equal(t, secret, manager.signingKey)
	assert.Equal(t, tokenTTL, manager.tokenTTL)
}

func TestTokenManager_GenerateToken(t *testing.T) {
	manager := NewTokenManager("test_secret", 15*time.Minute)
	userID := "user123"
	role := "user"

	token, err := manager.GenerateToken(userID, role)

	assert.NoError(t, err)
	assert.NotEmpty(t, token)

	// Verify token contents
	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return []byte("test_secret"), nil
	})

	assert.NoError(t, err)
	assert.True(t, parsedToken.Valid)

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	assert.True(t, ok)
	assert.Equal(t, userID, claims["sub"])
	assert.Equal(t, role, claims["role"])
}

func TestTokenManager_ValidateToken(t *testing.T) {
	manager := NewTokenManager("test_secret", 15*time.Minute)
	userID := "user123"
	role := "user"

	// Generate a valid token
	token, err := manager.GenerateToken(userID, role)
	require.NoError(t, err)

	// Validate the token
	claims, err := manager.ValidateToken(token)
	assert.NoError(t, err)
	assert.Equal(t, userID, claims["sub"])
	assert.Equal(t, role, claims["role"])

	// Test with invalid token
	invalidToken := "invalid.token.string"
	_, err = manager.ValidateToken(invalidToken)
	assert.Error(t, err)

	// Test with expired token
	expiredManager := NewTokenManager("test_secret", -10*time.Minute) // Token already expired
	expiredToken, err := expiredManager.GenerateToken(userID, role)
	require.NoError(t, err)
	_, err = manager.ValidateToken(expiredToken)
	assert.Error(t, err)

	// Test with token signed with wrong key
	wrongManager := NewTokenManager("wrong_secret", 15*time.Minute)
	wrongToken, err := wrongManager.GenerateToken(userID, role)
	require.NoError(t, err)
	_, err = manager.ValidateToken(wrongToken)
	assert.Error(t, err)
}

func TestTokenManager_RefreshToken(t *testing.T) {
	manager := NewTokenManager("test_secret", 15*time.Minute)
	userID := "user123"
	role := "user"

	// Тест 1: Успішне оновлення валідного токена
	t.Run("Valid Token Refresh", func(t *testing.T) {
		// Створюємо токен з коротким терміном дії
		token, err := manager.GenerateToken(userID, role)
		require.NoError(t, err)

		// Оновлюємо токен
		newToken, err := manager.RefreshToken(token)
		assert.NoError(t, err)
		assert.NotEqual(t, token, newToken) // Новий токен повинен відрізнятись від старого

		// Перевіряємо, що новий токен валідний
		claims, err := manager.ValidateToken(newToken)
		assert.NoError(t, err)
		assert.Equal(t, userID, claims["sub"])
		assert.Equal(t, role, claims["role"])

		// Перевіряємо, що термін дії нового токена встановлено правильно
		exp, ok := claims["exp"].(float64)
		assert.True(t, ok)
		assert.Greater(t, exp, float64(time.Now().Unix()))
	})

	// Тест 2: Спроба оновлення токена з неправильним форматом
	t.Run("Invalid Token Format", func(t *testing.T) {
		// Токен неправильного формату
		invalidToken := "invalid.token"
		_, err := manager.RefreshToken(invalidToken)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "token contains an invalid number of segments")
	})

	// Тест 3: Спроба оновлення простроченого токена
	t.Run("Expired Token", func(t *testing.T) {
		// Створюємо токен, який вже прострочений
		expiredManager := NewTokenManager("test_secret", -10*time.Minute) // Від'ємний TTL для простроченого токена
		expiredToken, err := expiredManager.GenerateToken(userID, role)
		require.NoError(t, err)

		// Спроба оновити прострочений токен
		// У реальному додатку це може бути дозволено або заборонено в залежності від політики безпеки
		// В нашому випадку ми тестуємо, що рішення обробляється консистентно
		_, err = manager.RefreshToken(expiredToken)

		// Якщо ваш RefreshToken метод повинен дозволяти оновлення простроченого токена:
		// assert.NoError(t, err)

		// Якщо ваш RefreshToken метод НЕ повинен дозволяти оновлення простроченого токена:
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "token is expired")
	})

	// Тест 4: Спроба оновлення токена з неправильним підписом
	t.Run("Invalid Signature", func(t *testing.T) {
		// Створюємо токен з іншим ключем
		wrongManager := NewTokenManager("wrong_secret", 15*time.Minute)
		wrongToken, err := wrongManager.GenerateToken(userID, role)
		require.NoError(t, err)

		// Пробуємо оновити токен з неправильним підписом
		_, err = manager.RefreshToken(wrongToken)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "signature is invalid")
	})
}
