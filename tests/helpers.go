package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/require"
	"io"
	"math/rand"
	"net/http"
	"testing"
	"time"
)

const baseURL = "http://localhost:8080"

type User struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Age      int    `json:"age"`
	Password string `json:"-"` // не сериализуется
}

type AuthResponse struct {
	Token string `json:"token"`
}

type Order struct {
	ID        int     `json:"id"`
	UserID    int     `json:"user_id"`
	Product   string  `json:"product"`
	Quantity  int     `json:"quantity"`
	Price     float64 `json:"price"`
	CreatedAt string  `json:"created_at"`
}

// --------------------------------- Utility Functions ---------------------------------

func userToRegisterPayload(user User) map[string]interface{} {
	return map[string]interface{}{
		"name":     user.Name,
		"email":    user.Email,
		"age":      user.Age,
		"password": user.Password,
	}
}

func userToLoginPayload(user User) map[string]string {
	return map[string]string{
		"email":    user.Email,
		"password": user.Password,
	}
}

func randomEmail() string {
	rand.Seed(time.Now().UnixNano())
	return fmt.Sprintf("test_user_%d@example.com", rand.Intn(1_000_000))
}

func readAndCloseBody(t *testing.T, body io.ReadCloser) string {
	defer func(body io.ReadCloser) {
		err := body.Close()
		if err != nil {

		}
	}(body)
	data, err := io.ReadAll(body)
	require.NoError(t, err)
	return string(data)
}

// --------------------------------- JWT Related ---------------------------------

func generateExpiredToken(t *testing.T, email string) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": email,
		"exp":   time.Now().Add(-1 * time.Hour).Unix(), // expired 1h ago
		"iat":   time.Now().Add(-2 * time.Hour).Unix(),
	})
	secret := "veryverystrongkeytojwthello"
	tokenString, err := token.SignedString([]byte(secret))
	require.NoError(t, err)
	return tokenString
}

// --------------------------------- HTTP Helpers ---------------------------------

func doRequest(t *testing.T, method, url, token string, body interface{}) *http.Response {
	var buf io.Reader
	if body != nil {
		jsonBody, err := json.Marshal(body)
		require.NoError(t, err)
		buf = bytes.NewBuffer(jsonBody)
	}

	req, err := http.NewRequest(method, url, buf)
	require.NoError(t, err)

	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	require.NoError(t, err)

	return resp
}

// --------------------------------- User/Test Helpers ---------------------------------

func createTestUser(t *testing.T) (User, string) {
	email := randomEmail()
	password := "testpassword"

	newUser := User{
		Name:     "Test User",
		Email:    email,
		Age:      30,
		Password: password,
	}

	// Register
	resp := doRequest(t, "POST", baseURL+"/users", "", userToRegisterPayload(newUser))
	defer resp.Body.Close()
	require.Equal(t, 201, resp.StatusCode)

	var createdUser User
	err := json.NewDecoder(resp.Body).Decode(&createdUser)
	require.NoError(t, err)

	// Login
	loginResp := doRequest(t, "POST", baseURL+"/auth/login", "", userToLoginPayload(newUser))
	defer loginResp.Body.Close()
	require.Equal(t, 200, loginResp.StatusCode)

	var auth AuthResponse
	err = json.NewDecoder(loginResp.Body).Decode(&auth)
	require.NoError(t, err)

	return createdUser, auth.Token
}

func deleteTestUser(t *testing.T, userID int, token string) {
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/users/%d", baseURL, userID), nil)
	require.NoError(t, err)
	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	require.Equal(t, 204, resp.StatusCode)
}

func createTestOrder(t *testing.T, userID int, token string, order map[string]interface{}) {
	resp := doRequest(t, "POST", fmt.Sprintf("%s/users/%d/orders", baseURL, userID), token, order)
	defer resp.Body.Close()
	require.Equal(t, 201, resp.StatusCode)
}
