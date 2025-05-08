package tests

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"testing"
)

func TestUser1_CreateValidUser(t *testing.T) {
	user, token := createTestUser(t)
	defer deleteTestUser(t, user.ID, token)
	assert.NotZero(t, user.ID)
	assert.NotEmpty(t, token)
}

// Создание пользователя с существующим email
func TestUser2_CreateDuplicateUser(t *testing.T) {
	user, token := createTestUser(t)
	defer deleteTestUser(t, user.ID, token)

	payload := map[string]interface{}{
		"name":     user.Name,
		"email":    user.Email,
		"age":      user.Age,
		"password": "testpassword",
	}
	resp := doRequest(t, "POST", baseURL+"/users", "", payload)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
}

func TestUser3_GetUserByID(t *testing.T) {
	user, token := createTestUser(t)
	defer deleteTestUser(t, user.ID, token)

	resp := doRequest(t, "GET", fmt.Sprintf("%s/users/%d", baseURL, user.ID), token, nil)
	defer resp.Body.Close()
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestUser4_UpdateUserName(t *testing.T) {
	user, token := createTestUser(t)
	defer deleteTestUser(t, user.ID, token)

	update := map[string]interface{}{
		"name":  "Updated Name",
		"email": user.Email,
		"age":   user.Age,
	}
	resp := doRequest(t, "PUT", fmt.Sprintf("%s/users/%d", baseURL, user.ID), token, update)
	defer resp.Body.Close()
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestUser5_DeleteUser(t *testing.T) {
	user, token := createTestUser(t)

	deleteTestUser(t, user.ID, token)

	resp := doRequest(t, "GET", fmt.Sprintf("%s/users/%d", baseURL, user.ID), token, nil)
	defer resp.Body.Close()
	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
}

func TestUser6_GetUsersList(t *testing.T) {
	user, token := createTestUser(t)
	defer deleteTestUser(t, user.ID, token)

	resp := doRequest(t, "GET", baseURL+"/users", token, nil)
	defer resp.Body.Close()
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestUser7_GetUsersWithPagination(t *testing.T) {
	user, token := createTestUser(t)
	defer deleteTestUser(t, user.ID, token)

	resp := doRequest(t, "GET", baseURL+"/users?page=1&limit=2", token, nil)
	defer resp.Body.Close()
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestUser8_UpdateUserInvalidEmail(t *testing.T) {
	user, token := createTestUser(t)
	defer deleteTestUser(t, user.ID, token)

	update := map[string]interface{}{
		"email": "not-an-email",
	}
	resp := doRequest(t, "PUT", fmt.Sprintf("%s/users/%d", baseURL, user.ID), token, update)
	defer resp.Body.Close()
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestUser9_CreateUserWithEmptyName(t *testing.T) {
	payload := map[string]interface{}{
		"name":     "",
		"email":    randomEmail(),
		"age":      25,
		"password": "123456",
	}
	resp := doRequest(t, "POST", baseURL+"/users", "", payload)
	defer resp.Body.Close()
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestUser10_CreateUserWithInvalidAge(t *testing.T) {
	payload := map[string]interface{}{
		"name":     "Invalid Age",
		"email":    randomEmail(),
		"age":      -5,
		"password": "123456",
	}
	resp := doRequest(t, "POST", baseURL+"/users", "", payload)
	defer resp.Body.Close()
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestUser11_CreateUserWithoutPassword(t *testing.T) {
	payload := map[string]interface{}{
		"name":  "No Password",
		"email": randomEmail(),
		"age":   22,
	}
	resp := doRequest(t, "POST", baseURL+"/users", "", payload)
	defer resp.Body.Close()
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestUser12_GetNonExistingUser(t *testing.T) {
	user, token := createTestUser(t)
	defer deleteTestUser(t, user.ID, token)

	nonExistentID := 999999
	resp := doRequest(t, "GET", fmt.Sprintf("%s/users/%d", baseURL, nonExistentID), token, nil)
	defer resp.Body.Close()
	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
}

func TestUser13_UpdateOtherUser(t *testing.T) {
	user1, token1 := createTestUser(t)
	defer deleteTestUser(t, user1.ID, token1)

	user2, token2 := createTestUser(t)
	defer deleteTestUser(t, user2.ID, token2)

	update := map[string]interface{}{
		"name": "Hacked Name",
	}
	resp := doRequest(t, "PUT", fmt.Sprintf("%s/users/%d", baseURL, user2.ID), token1, update)
	defer resp.Body.Close()
	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
}

func TestUser14_DeleteOtherUser(t *testing.T) {
	user1, token1 := createTestUser(t)
	defer deleteTestUser(t, user1.ID, token1)

	user2, token2 := createTestUser(t)
	defer deleteTestUser(t, user2.ID, token2)

	req, _ := http.NewRequest("DELETE", fmt.Sprintf("%s/users/%d", baseURL, user2.ID), nil)
	req.Header.Set("Authorization", "Bearer "+token1)
	client := &http.Client{}
	resp, err := client.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
}

func TestUser15_UpdateUserWithoutToken(t *testing.T) {
	user, token := createTestUser(t)
	defer deleteTestUser(t, user.ID, token)

	update := map[string]interface{}{
		"name": "No Token Update",
	}
	resp := doRequest(t, "PUT", fmt.Sprintf("%s/users/%d", baseURL, user.ID), "", update)
	defer resp.Body.Close()
	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
}

func TestUser16_DeleteUserWithoutToken(t *testing.T) {
	user, token := createTestUser(t)
	defer deleteTestUser(t, user.ID, token)

	req, _ := http.NewRequest("DELETE", fmt.Sprintf("%s/users/%d", baseURL, user.ID), nil)
	client := &http.Client{}
	resp, err := client.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
}

func TestUser17_UpdateUserWithEmptyPayload(t *testing.T) {
	user, token := createTestUser(t)
	defer deleteTestUser(t, user.ID, token)

	update := map[string]interface{}{}
	resp := doRequest(t, "PUT", fmt.Sprintf("%s/users/%d", baseURL, user.ID), token, update)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestUser18_ListUsersUnauthorized(t *testing.T) {
	resp := doRequest(t, "GET", baseURL+"/users", "", nil)
	defer resp.Body.Close()
	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
}

func TestUser19_CreateUserWithInvalidEmail(t *testing.T) {
	payload := map[string]interface{}{
		"name":     "Invalid Email",
		"email":    "invalid-email",
		"age":      25,
		"password": "pass123",
	}
	resp := doRequest(t, "POST", baseURL+"/users", "", payload)
	defer resp.Body.Close()
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestUser20_UpdateUserAgeToZero(t *testing.T) {
	user, token := createTestUser(t)
	defer deleteTestUser(t, user.ID, token)

	update := map[string]interface{}{
		"age": 0,
	}
	resp := doRequest(t, "PUT", fmt.Sprintf("%s/users/%d", baseURL, user.ID), token, update)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}
