//go:build integration

package routes_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"testing"
	"time"
)

func TestDeployedCartGatewayFunctional(t *testing.T) {
	cartBaseURL := os.Getenv("CART_GATEWAY_BASE_URL")
	storefrontBaseURL := os.Getenv("STOREFRONT_GATEWAY_BASE_URL")
	if cartBaseURL == "" || storefrontBaseURL == "" {
		t.Skip("CART_GATEWAY_BASE_URL and STOREFRONT_GATEWAY_BASE_URL are required for deployed cart gateway integration tests")
	}

	client := &http.Client{Timeout: 15 * time.Second}
	email := cartGetEnv("CART_FUNCTIONAL_TEST_EMAIL", "admin@test.com")
	password := cartGetEnv("CART_FUNCTIONAL_TEST_PASSWORD", "password")

	cartDeployedRequest(t, client, http.MethodGet, cartBaseURL+"/health", http.StatusOK, "", nil)
	cartDeployedRequest(t, client, http.MethodGet, cartBaseURL+"/api/cart", http.StatusUnauthorized, "", nil)

	loginPayload := cartDeployedRequest(t, client, http.MethodPost, storefrontBaseURL+"/api/auth/login", http.StatusOK, "", map[string]any{
		"email":    email,
		"password": password,
	})
	token := cartRequiredStringField(t, loginPayload, "accessToken")

	productsPayload := cartDeployedRequest(t, client, http.MethodGet, storefrontBaseURL+"/api/products", http.StatusOK, "", nil)
	product := cartRequiredListResponse(t, productsPayload)[0]
	productID := cartRequiredNumberID(t, product, "id")
	priceCents := cartRequiredNumber(t, product, "price_cents")

	cartDeployedRequest(t, client, http.MethodDelete, cartBaseURL+"/api/cart", http.StatusOK, token, nil)
	cartAssertTotals(t, cartDeployedRequest(t, client, http.MethodGet, cartBaseURL+"/api/cart", http.StatusOK, token, nil), 0, 0)

	cartAssertTotals(t, cartDeployedRequest(t, client, http.MethodPost, cartBaseURL+"/api/cart/items", http.StatusCreated, token, map[string]any{
		"product_id": productID,
		"quantity":   2,
	}), 2, priceCents*2)

	cartAssertTotals(t, cartDeployedRequest(t, client, http.MethodPut, fmt.Sprintf("%s/api/cart/items/%d", cartBaseURL, productID), http.StatusOK, token, map[string]any{
		"quantity": 3,
	}), 3, priceCents*3)

	cartAssertTotals(t, cartDeployedRequest(t, client, http.MethodDelete, fmt.Sprintf("%s/api/cart/items/%d", cartBaseURL, productID), http.StatusOK, token, nil), 0, 0)

	cartDeployedRequest(t, client, http.MethodPost, cartBaseURL+"/api/cart/items", http.StatusCreated, token, map[string]any{
		"product_id": productID,
		"quantity":   1,
	})
	orderPayload := cartDeployedRequest(t, client, http.MethodPost, cartBaseURL+"/api/cart/checkout", http.StatusCreated, token, nil)
	cartAssertOrder(t, orderPayload, 1, priceCents)
	orderHistoryPayload := cartDeployedRequest(t, client, http.MethodGet, cartBaseURL+"/api/orders", http.StatusOK, token, nil)
	cartAssertOrderHistory(t, orderHistoryPayload)
	cartAssertTotals(t, cartDeployedRequest(t, client, http.MethodGet, cartBaseURL+"/api/cart", http.StatusOK, token, nil), 0, 0)
	cartDeployedRequest(t, client, http.MethodPost, storefrontBaseURL+"/api/auth/logout", http.StatusOK, token, nil)
}

func cartDeployedRequest(
	t *testing.T,
	client *http.Client,
	method string,
	url string,
	expectedStatus int,
	token string,
	body any,
) any {
	t.Helper()

	var requestBody bytes.Buffer
	if body != nil {
		if err := json.NewEncoder(&requestBody).Encode(body); err != nil {
			t.Fatalf("encode request body: %v", err)
		}
	}

	request, err := http.NewRequest(method, url, &requestBody)
	if err != nil {
		t.Fatalf("create request: %v", err)
	}
	requestID := fmt.Sprintf("ci-cart-%s-%s-%s", cartGetEnv("GITHUB_RUN_ID", "local"), method, request.URL.Path)
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("request-id", requestID)
	if token != "" {
		request.Header.Set("Authorization", "Bearer "+token)
	}

	t.Logf("Functional test: %s %s expects %d request-id=%s", method, request.URL.Path, expectedStatus, requestID)

	response, err := client.Do(request)
	if err != nil {
		t.Fatalf("%s %s failed: %v", method, request.URL.Path, err)
	}
	defer response.Body.Close()

	var responsePayload any
	if err := json.NewDecoder(response.Body).Decode(&responsePayload); err != nil {
		t.Fatalf("decode %s %s response: %v", method, request.URL.Path, err)
	}
	if response.StatusCode != expectedStatus {
		t.Fatalf("%s %s returned %d, expected %d: %#v", method, request.URL.Path, response.StatusCode, expectedStatus, responsePayload)
	}

	t.Logf("Functional test passed: %s %s (%d)", method, request.URL.Path, response.StatusCode)
	return responsePayload
}

func cartRequiredListResponse(t *testing.T, payload any) []any {
	t.Helper()
	listPayload, ok := payload.([]any)
	if !ok || len(listPayload) == 0 {
		t.Fatalf("expected non-empty list response, got %#v", payload)
	}
	return listPayload
}

func cartRequiredStringField(t *testing.T, payload any, field string) string {
	t.Helper()
	objectPayload, ok := payload.(map[string]any)
	if !ok {
		t.Fatalf("expected object response with field %q, got %#v", field, payload)
	}
	value, ok := objectPayload[field].(string)
	if !ok || value == "" {
		t.Fatalf("expected response field %q, got %#v", field, payload)
	}
	return value
}

func cartRequiredNumberID(t *testing.T, payload any, field string) uint64 {
	t.Helper()
	return uint64(cartRequiredNumber(t, payload, field))
}

func cartRequiredNumber(t *testing.T, payload any, field string) int64 {
	t.Helper()
	objectPayload, ok := payload.(map[string]any)
	if !ok {
		t.Fatalf("expected object response with field %q, got %#v", field, payload)
	}
	value, ok := objectPayload[field].(float64)
	if !ok {
		t.Fatalf("expected numeric response field %q, got %#v", field, payload)
	}
	return int64(value)
}

func cartAssertTotals(t *testing.T, payload any, expectedItemCount int64, expectedSubtotalCents int64) {
	t.Helper()
	if actual := cartRequiredNumber(t, payload, "item_count"); actual != expectedItemCount {
		t.Fatalf("expected cart item count %d, got %d: %#v", expectedItemCount, actual, payload)
	}
	if actual := cartRequiredNumber(t, payload, "subtotal_cents"); actual != expectedSubtotalCents {
		t.Fatalf("expected cart subtotal %d, got %d: %#v", expectedSubtotalCents, actual, payload)
	}
}

func cartAssertOrder(t *testing.T, payload any, expectedItemCount int64, expectedSubtotalCents int64) {
	t.Helper()
	if actual := cartRequiredNumber(t, payload, "id"); actual == 0 {
		t.Fatalf("expected order id, got %#v", payload)
	}
	if status := cartRequiredStringField(t, payload, "status"); status != "Order Submitted" {
		t.Fatalf("expected Order Submitted order status, got %q: %#v", status, payload)
	}
	if actual := cartRequiredNumber(t, payload, "item_count"); actual != expectedItemCount {
		t.Fatalf("expected order item count %d, got %d: %#v", expectedItemCount, actual, payload)
	}
	if actual := cartRequiredNumber(t, payload, "subtotal_cents"); actual != expectedSubtotalCents {
		t.Fatalf("expected order subtotal %d, got %d: %#v", expectedSubtotalCents, actual, payload)
	}
}

func cartAssertOrderHistory(t *testing.T, payload any) {
	t.Helper()
	orders := cartRequiredListResponse(t, payload)
	if cartRequiredNumber(t, orders[0], "id") == 0 {
		t.Fatalf("expected first order id, got %#v", payload)
	}
}

func cartGetEnv(key string, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
