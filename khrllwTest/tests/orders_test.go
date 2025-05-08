package tests

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func Test1_CreateOrderSuccess(t *testing.T) {
	user, token := createTestUser(t)
	defer deleteTestUser(t, user.ID, token)

	order := map[string]interface{}{
		"product":  "Monitor",
		"quantity": 2,
		"price":    299.99,
	}

	resp := doRequest(t, "POST", fmt.Sprintf("%s/users/%d/orders", baseURL, user.ID), token, order)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusCreated, resp.StatusCode)
}

func Test2_GetOrderList(t *testing.T) {
	user, token := createTestUser(t)
	defer deleteTestUser(t, user.ID, token)

	createTestOrder(t, user.ID, token, map[string]interface{}{"product": "Keyboard", "quantity": 1, "price": 49.99})

	resp := doRequest(t, "GET", fmt.Sprintf("%s/users/%d/orders", baseURL, user.ID), token, nil)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func Test3_GetOrderByID(t *testing.T) {
	user, token := createTestUser(t)
	defer deleteTestUser(t, user.ID, token)

	createTestOrder(t, user.ID, token, map[string]interface{}{"product": "Mouse", "quantity": 1, "price": 19.99})

	resp := doRequest(t, "GET", fmt.Sprintf("%s/users/%d/orders/", baseURL, user.ID), token, nil)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func Test4_CreateOrderWithoutAuth(t *testing.T) {
	user, token := createTestUser(t)
	defer deleteTestUser(t, user.ID, token)

	order := map[string]interface{}{
		"product":  "Tablet",
		"quantity": 1,
		"price":    500.00,
	}

	resp := doRequest(t, "POST", fmt.Sprintf("%s/users/%d/orders", baseURL, user.ID), "", order)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
}

func Test5_CreateOrderWithInvalidPayload(t *testing.T) {
	user, token := createTestUser(t)
	defer deleteTestUser(t, user.ID, token)

	invalidOrder := map[string]interface{}{
		"product":  "",
		"quantity": -2,
	}

	resp := doRequest(t, "POST", fmt.Sprintf("%s/users/%d/orders", baseURL, user.ID), token, invalidOrder)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func Test6_GetOrdersWithPagination(t *testing.T) {
	user, token := createTestUser(t)
	defer deleteTestUser(t, user.ID, token)

	for i := 0; i < 5; i++ {
		createTestOrder(t, user.ID, token, map[string]interface{}{
			"product":  fmt.Sprintf("Product%d", i),
			"quantity": 1,
			"price":    100.00,
		})
	}

	resp := doRequest(t, "GET", fmt.Sprintf("%s/users/%d/orders?page=1&limit=2", baseURL, user.ID), token, nil)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func Test7_GetOrderInvalidID(t *testing.T) {
	user, token := createTestUser(t)
	defer deleteTestUser(t, user.ID, token)

	resp := doRequest(t, "GET", fmt.Sprintf("%s/users/%d/orders/9999", baseURL, user.ID), token, nil)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
}

func Test8_CreateMultipleOrders(t *testing.T) {
	user, token := createTestUser(t)
	defer deleteTestUser(t, user.ID, token)

	for i := 0; i < 3; i++ {
		createTestOrder(t, user.ID, token, map[string]interface{}{
			"product":  fmt.Sprintf("Item %d", i),
			"quantity": i + 1,
			"price":    float64(10 * (i + 1)),
		})
	}

	resp := doRequest(t, "GET", fmt.Sprintf("%s/users/%d/orders", baseURL, user.ID), token, nil)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func Test9_OrderAccessByAnotherUser(t *testing.T) {
	user1, token1 := createTestUser(t)
	defer deleteTestUser(t, user1.ID, token1)

	user2, token2 := createTestUser(t)
	defer deleteTestUser(t, user2.ID, token2)

	createTestOrder(t, user1.ID, token1, map[string]interface{}{"product": "Speaker", "quantity": 1, "price": 80.00})

	resp := doRequest(t, "GET", fmt.Sprintf("%s/users/%d/orders/", baseURL, user1.ID), token2, nil)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
}

func Test10_OrderCreationNegativePrice(t *testing.T) {
	user, token := createTestUser(t)
	defer deleteTestUser(t, user.ID, token)

	order := map[string]interface{}{
		"product":  "InvalidItem",
		"quantity": 1,
		"price":    -50.00,
	}

	resp := doRequest(t, "POST", fmt.Sprintf("%s/users/%d/orders", baseURL, user.ID), token, order)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func Test11_OrderCreationZeroQuantity(t *testing.T) {
	user, token := createTestUser(t)
	defer deleteTestUser(t, user.ID, token)

	order := map[string]interface{}{
		"product":  "ZeroQty",
		"quantity": 0,
		"price":    20.00,
	}

	resp := doRequest(t, "POST", fmt.Sprintf("%s/users/%d/orders", baseURL, user.ID), token, order)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func Test12_CreateOrderInvalidUserID(t *testing.T) {
	user, token := createTestUser(t)
	defer deleteTestUser(t, user.ID, token)

	order := map[string]interface{}{
		"product":  "Book",
		"quantity": 1,
		"price":    30.00,
	}

	resp := doRequest(t, "POST", baseURL+"/users/99999/orders", token, order)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
}

func Test13_GetOrdersNoOrders(t *testing.T) {
	user, token := createTestUser(t)
	defer deleteTestUser(t, user.ID, token)

	resp := doRequest(t, "GET", fmt.Sprintf("%s/users/%d/orders", baseURL, user.ID), token, nil)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func Test14_UpdateNonExistentOrder(t *testing.T) {
	user, token := createTestUser(t)
	defer deleteTestUser(t, user.ID, token)

	update := map[string]interface{}{
		"product":  "NewProduct",
		"quantity": 1,
		"price":    999.99,
	}

	resp := doRequest(t, "PUT", fmt.Sprintf("%s/users/%d/orders/9999", baseURL, user.ID), token, update)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
}
