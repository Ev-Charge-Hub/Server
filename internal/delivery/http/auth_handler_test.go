package http_test

import (
	"Ev-Charge-Hub/Server/utils"
	"net/http"
	"net/http/httptest"
	"testing"

	deliveryHttp "Ev-Charge-Hub/Server/internal/delivery/http"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func setupSecurityRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	r.GET("/security/validate-token", deliveryHttp.TokenValidationHandler)
	return r
}

func TestValidateToken_Success(t *testing.T) {
	router := setupSecurityRouter()

	// 🔐 Generate valid token
	token, err := utils.CreateToken("12345", "testuser", "ADMIN")
	assert.NoError(t, err)

	req := httptest.NewRequest("GET", "/security/validate-token", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
	assert.Contains(t, resp.Body.String(), `"valid":true`)
	assert.Contains(t, resp.Body.String(), `"role":"ADMIN"`)
}

func TestValidateToken_MissingHeader(t *testing.T) {
	router := setupSecurityRouter()

	req := httptest.NewRequest("GET", "/security/validate-token", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusUnauthorized, resp.Code)
	assert.Contains(t, resp.Body.String(), "Authorization header missing")
}

func TestValidateToken_InvalidToken(t *testing.T) {
	router := setupSecurityRouter()

	// ❌ Use an invalid token
	req := httptest.NewRequest("GET", "/security/validate-token", nil)
	req.Header.Set("Authorization", "Bearer invalid.token.here")

	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusUnauthorized, resp.Code)
	assert.Contains(t, resp.Body.String(), "Invalid token")
}
