package tests

import (
	"encoding/base64"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"strconv"
	"testing"
	"time"
)

func Test1_ValidJWTLogin(t *testing.T) {
	user, token := createTestUser(t)
	defer deleteTestUser(t, user.ID, token)

	require.NotEmpty(t, token)
}

func Test2_LoginWithWrongPassword(t *testing.T) {
	user, token := createTestUser(t)
	defer deleteTestUser(t, user.ID, token)

	loginPayload := map[string]string{
		"email":    user.Email,
		"password": "wrongpassword",
	}
	resp := doRequest(t, "POST", baseURL+"/auth/login", "", loginPayload)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
}

func Test3_LoginWithUnknownEmail(t *testing.T) {
	loginPayload := map[string]string{
		"email":    "nonexistent@example.com",
		"password": "any",
	}
	resp := doRequest(t, "POST", baseURL+"/auth/login", "", loginPayload)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func Test4_RequestWithoutToken(t *testing.T) {
	resp := doRequest(t, "GET", baseURL+"/users", "", nil)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
}

func Test5_RequestWithInvalidTokenFormat(t *testing.T) {
	req, err := http.NewRequest("GET", baseURL+"/users", nil)
	require.NoError(t, err)
	req.Header.Set("Authorization", "InvalidTokenFormat")

	client := &http.Client{}
	resp, err := client.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
}

func Test6_RequestWithTamperedToken(t *testing.T) {
	user, token := createTestUser(t)
	defer deleteTestUser(t, user.ID, token)

	tampered := token[:len(token)-1] + "x"
	resp := doRequest(t, "GET", baseURL+"/users", tampered, nil)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
}

func Test7_ExpiredTokenSimulation(t *testing.T) {
	user, token := createTestUser(t)
	defer deleteTestUser(t, user.ID, token)

	expiredToken := generateExpiredToken(t, user.Email)

	resp := doRequest(t, "GET", baseURL+"/users", expiredToken, nil)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
}

func Test8_UseTokenOfAnotherUser(t *testing.T) {
	user1, token1 := createTestUser(t)
	defer deleteTestUser(t, user1.ID, token1)

	user2, token2 := createTestUser(t)
	defer deleteTestUser(t, user2.ID, token2)

	resp := doRequest(t, "GET", fmt.Sprintf("%s/users/%d", baseURL, user2.ID), token1, nil)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
}

func Test9_MissingBearerPrefix(t *testing.T) {
	user, token := createTestUser(t)
	defer deleteTestUser(t, user.ID, token)

	req, err := http.NewRequest("GET", baseURL+"/users", nil)
	require.NoError(t, err)
	req.Header.Set("Authorization", token) // нет Bearer

	client := &http.Client{}
	resp, err := client.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
}

func Test10_TokenReuseAfterUserDeletion(t *testing.T) {
	user, token := createTestUser(t)
	userID := strconv.Itoa(user.ID)

	deleteTestUser(t, user.ID, token)

	resp := doRequest(t, "GET", baseURL+"/users/"+userID, token, nil)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
}

func Test11_MalformedJWTStructure(t *testing.T) {
	malformedToken := "abc.def"

	resp := doRequest(t, "GET", baseURL+"/users", malformedToken, nil)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
}

func Test12_ValidTokenCanAccessProtectedRoute(t *testing.T) {
	user, token := createTestUser(t)
	defer deleteTestUser(t, user.ID, token)

	resp := doRequest(t, "GET", fmt.Sprintf("%s/users/%d", baseURL, user.ID), token, nil)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func Test13_TokenWithFutureIAT(t *testing.T) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": "test@example.com",
		"iat":   time.Now().Add(1 * time.Hour).Unix(),
		"exp":   time.Now().Add(2 * time.Hour).Unix(),
	})
	tokenString, err := token.SignedString([]byte("veryverystrongkeytojwthello"))
	require.NoError(t, err)

	resp := doRequest(t, "GET", baseURL+"/users", tokenString, nil)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
}

func Test14_TokenWithAlgNone(t *testing.T) {
	// Имитация попытки подделать токен с "alg": "none"
	header := base64.StdEncoding.EncodeToString([]byte(`{"alg":"none","typ":"JWT"}`))
	payload := base64.StdEncoding.EncodeToString([]byte(`{"email":"test@example.com","exp":9999999999}`))
	fakeToken := fmt.Sprintf("%s.%s.", header, payload)

	resp := doRequest(t, "GET", baseURL+"/users", fakeToken, nil)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
}

func Test15_TokenWithInvalidSignature(t *testing.T) {
	user, _ := createTestUser(t)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": user.Email,
		"exp":   time.Now().Add(1 * time.Hour).Unix(),
	})
	invalidToken, err := token.SignedString([]byte("wrongsecret"))
	require.NoError(t, err)

	resp := doRequest(t, "GET", baseURL+"/users", invalidToken, nil)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
}

func Test16_TokenInBodyInsteadOfHeader(t *testing.T) {
	user, token := createTestUser(t)
	defer deleteTestUser(t, user.ID, token)

	body := map[string]string{
		"token": token,
	}
	resp := doRequest(t, "POST", baseURL+"/users", "", body) // token не в header
	defer resp.Body.Close()

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func Test17_TokenValidForAnotherEndpoint(t *testing.T) {
	user, token := createTestUser(t)
	defer deleteTestUser(t, user.ID, token)

	order := map[string]interface{}{
		"product":  "Phone",
		"quantity": 1,
		"price":    999.9,
	}
	resp := doRequest(t, "POST", fmt.Sprintf("%s/users/%d/orders", baseURL, user.ID), token, order)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusCreated, resp.StatusCode)
}

func Test18_TokenWithExtraClaimsStillValid(t *testing.T) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": "test@example.com",
		"exp":   time.Now().Add(1 * time.Hour).Unix(),
		"role":  "admin",
	})
	tokenString, err := token.SignedString([]byte("veryverystrongkeytojwthello"))
	require.NoError(t, err)

	resp := doRequest(t, "GET", baseURL+"/users", tokenString, nil)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode) // Depends on backend behavior
}

func Test19_TokenReuseMultipleTimes(t *testing.T) {
	user, token := createTestUser(t)
	defer deleteTestUser(t, user.ID, token)

	for i := 0; i < 3; i++ {
		resp := doRequest(t, "GET", baseURL+"/users", token, nil)
		defer resp.Body.Close()
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	}
}
