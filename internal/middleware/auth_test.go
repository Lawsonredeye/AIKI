package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"aiki/internal/pkg/jwt"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAuthMiddleware(t *testing.T) {
	e := echo.New()
	jwtManager := jwt.NewManager("test-secret", 15*time.Minute, 7*24*time.Hour)

	t.Run("valid token", func(t *testing.T) {
		userID := int32(123)
		email := "test@example.com"

		token, err := jwtManager.GenerateAccessToken(userID, email)
		require.NoError(t, err)

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		handler := func(c echo.Context) error {
			assert.Equal(t, userID, c.Get("user_id"))
			assert.Equal(t, email, c.Get("user_email"))
			return c.String(http.StatusOK, "OK")
		}

		middleware := Auth(jwtManager)
		err = middleware(handler)(c)

		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("missing authorization header", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		handler := func(c echo.Context) error {
			return c.String(http.StatusOK, "Should not reach here")
		}

		middleware := Auth(jwtManager)
		err := middleware(handler)(c)

		require.NoError(t, err)
		assert.Equal(t, http.StatusUnauthorized, rec.Code)
	})

	t.Run("invalid token format", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set("Authorization", "InvalidFormat token")
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		handler := func(c echo.Context) error {
			return c.String(http.StatusOK, "Should not reach here")
		}

		middleware := Auth(jwtManager)
		err := middleware(handler)(c)

		require.NoError(t, err)
		assert.Equal(t, http.StatusUnauthorized, rec.Code)
	})

	t.Run("invalid token", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set("Authorization", "Bearer invalid-token")
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		handler := func(c echo.Context) error {
			return c.String(http.StatusOK, "Should not reach here")
		}

		middleware := Auth(jwtManager)
		err := middleware(handler)(c)

		require.NoError(t, err)
		assert.Equal(t, http.StatusUnauthorized, rec.Code)
	})

	t.Run("expired token", func(t *testing.T) {
		expiredManager := jwt.NewManager("test-secret", -1*time.Hour, 7*24*time.Hour)
		userID := int32(123)
		email := "test@example.com"

		expiredToken, err := expiredManager.GenerateAccessToken(userID, email)
		require.NoError(t, err)

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set("Authorization", "Bearer "+expiredToken)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		handler := func(c echo.Context) error {
			return c.String(http.StatusOK, "Should not reach here")
		}

		middleware := Auth(jwtManager)
		err = middleware(handler)(c)

		require.NoError(t, err)
		assert.Equal(t, http.StatusUnauthorized, rec.Code)
	})
}
