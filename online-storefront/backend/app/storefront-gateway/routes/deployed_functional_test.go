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

func TestDeployedStorefrontGatewayFunctional(t *testing.T) {
	baseURL := os.Getenv("STOREFRONT_GATEWAY_BASE_URL")
	if baseURL == "" {
		t.Skip("STOREFRONT_GATEWAY_BASE_URL is required for deployed storefront gateway integration tests")
	}

	client := &http.Client{Timeout: 15 * time.Second}
	email := storefrontGetEnv("STOREFRONT_FUNCTIONAL_TEST_EMAIL", "admin@test.com")
	password := storefrontGetEnv("STOREFRONT_FUNCTIONAL_TEST_PASSWORD", "password")

	storefrontDeployedRequest(t, client, http.MethodGet, baseURL+"/health", http.StatusOK, "", nil)

	productsPayload := storefrontDeployedRequest(t, client, http.MethodGet, baseURL+"/api/products", http.StatusOK, "", nil)
	products := storefrontRequiredListResponse(t, productsPayload)
	productID := storefrontRequiredNumberID(t, products[0], "id")
	storefrontDeployedRequest(t, client, http.MethodGet, fmt.Sprintf("%s/api/products/%d", baseURL, productID), http.StatusOK, "", nil)

	loginPayload := storefrontDeployedRequest(t, client, http.MethodPost, baseURL+"/api/auth/login", http.StatusOK, "", map[string]any{
		"email":    email,
		"password": password,
	})
	token := storefrontRequiredStringField(t, loginPayload, "accessToken")

	storefrontDeployedRequest(t, client, http.MethodGet, baseURL+"/api/users/me", http.StatusOK, token, nil)
	storefrontDeployedRequest(t, client, http.MethodPost, baseURL+"/api/auth/logout", http.StatusOK, token, nil)
}

func storefrontDeployedRequest(
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

	requestID := fmt.Sprintf("ci-storefront-%s-%s-%s", storefrontGetEnv("GITHUB_RUN_ID", "local"), method, request.URL.Path)
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

func storefrontRequiredListResponse(t *testing.T, payload any) []any {
	t.Helper()
	listPayload, ok := payload.([]any)
	if !ok || len(listPayload) == 0 {
		t.Fatalf("expected non-empty list response, got %#v", payload)
	}
	return listPayload
}

func storefrontRequiredStringField(t *testing.T, payload any, field string) string {
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

func storefrontRequiredNumberID(t *testing.T, payload any, field string) uint64 {
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

func storefrontGetEnv(key string, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
