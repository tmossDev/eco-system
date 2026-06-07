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

func TestDeployedOrderGatewayFunctional(t *testing.T) {
	baseURL := os.Getenv("ORDER_GATEWAY_BASE_URL")
	if baseURL == "" {
		t.Skip("ORDER_GATEWAY_BASE_URL is required for deployed order gateway integration tests")
	}

	client := &http.Client{Timeout: 15 * time.Second}
	email := orderGetEnv("ORDER_FUNCTIONAL_TEST_EMAIL", "admin@test.com")
	password := orderGetEnv("ORDER_FUNCTIONAL_TEST_PASSWORD", "password")

	loginPayload := orderDeployedRequest(t, client, http.MethodPost, baseURL+"/api/auth/login", http.StatusOK, "", map[string]any{
		"email":    email,
		"password": password,
	})
	token := orderRequiredStringField(t, loginPayload, "accessToken")

	ordersPayload := orderDeployedRequest(t, client, http.MethodGet, baseURL+"/api/orders", http.StatusOK, token, nil)
	orders := orderRequiredListResponse(t, ordersPayload)
	orderID := orderRequiredNumber(t, orders[0], "id")
	orderStatus := orderRequiredStringField(t, orders[0], "status")

	orderDeployedRequest(t, client, http.MethodGet, fmt.Sprintf("%s/api/orders/%d", baseURL, orderID), http.StatusOK, token, nil)
	orderDeployedRequest(t, client, http.MethodPut, fmt.Sprintf("%s/api/orders/%d/status", baseURL, orderID), http.StatusOK, token, map[string]any{
		"status": orderStatus,
	})
	orderDeployedRequest(t, client, http.MethodPost, baseURL+"/api/auth/logout", http.StatusOK, token, nil)
}

func orderDeployedRequest(
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

	requestID := fmt.Sprintf("ci-order-%s-%s-%s", orderGetEnv("GITHUB_RUN_ID", "local"), method, request.URL.Path)
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
	if response.StatusCode != http.StatusNoContent {
		if err := json.NewDecoder(response.Body).Decode(&responsePayload); err != nil {
			t.Fatalf("decode %s %s response: %v", method, request.URL.Path, err)
		}
	}

	if response.StatusCode != expectedStatus {
		t.Fatalf("%s %s returned %d, expected %d: %#v", method, request.URL.Path, response.StatusCode, expectedStatus, responsePayload)
	}

	t.Logf("Functional test passed: %s %s (%d)", method, request.URL.Path, response.StatusCode)

	return responsePayload
}

func orderRequiredStringField(t *testing.T, payload any, field string) string {
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

func orderRequiredNumber(t *testing.T, payload any, field string) uint64 {
	t.Helper()

	objectPayload, ok := payload.(map[string]any)
	if !ok {
		t.Fatalf("expected object response with field %q, got %#v", field, payload)
	}

	value, ok := objectPayload[field].(float64)
	if !ok || value == 0 {
		t.Fatalf("expected numeric response field %q, got %#v", field, payload)
	}

	return uint64(value)
}

func orderRequiredListResponse(t *testing.T, payload any) []any {
	t.Helper()

	listPayload, ok := payload.([]any)
	if !ok || len(listPayload) == 0 {
		t.Fatalf("expected non-empty list response, got %#v", payload)
	}

	return listPayload
}

func orderGetEnv(key string, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}

	return fallback
}
