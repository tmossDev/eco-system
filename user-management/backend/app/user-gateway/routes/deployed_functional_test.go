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

func TestDeployedGatewayFunctional(t *testing.T) {
	baseURL := os.Getenv("GATEWAY_BASE_URL")
	if baseURL == "" {
		t.Skip("GATEWAY_BASE_URL is required for deployed gateway integration tests")
	}

	client := &http.Client{Timeout: 15 * time.Second}
	email := getEnv("FUNCTIONAL_TEST_EMAIL", "admin@test.com")
	password := getEnv("FUNCTIONAL_TEST_PASSWORD", "password")

	loginPayload := deployedRequest(t, client, http.MethodPost, baseURL+"/api/auth/login", http.StatusOK, "", map[string]any{
		"email":    email,
		"password": password,
	})
	token := requiredStringField(t, loginPayload, "accessToken")

	deployedRequest(t, client, http.MethodGet, baseURL+"/api/users", http.StatusOK, token, nil)
	deployedRequest(t, client, http.MethodGet, baseURL+"/api/users/2", http.StatusOK, token, nil)

	createdPayload := deployedRequest(t, client, http.MethodPost, baseURL+"/api/users", http.StatusCreated, token, map[string]any{
		"name":   "CI Functional User",
		"email":  "ci.functional@test.com",
		"role":   "User",
		"status": "Active",
	})
	createdUserID := requiredStringField(t, createdPayload, "id")

	deployedRequest(t, client, http.MethodPut, baseURL+"/api/users/"+createdUserID, http.StatusOK, token, map[string]any{
		"name":   "CI Functional User Updated",
		"email":  "ci.functional@test.com",
		"role":   "User",
		"status": "Active",
	})
	deployedRequest(t, client, http.MethodDelete, baseURL+"/api/users/"+createdUserID, http.StatusNoContent, token, nil)
	deployedRequest(t, client, http.MethodPost, baseURL+"/api/auth/logout", http.StatusOK, token, nil)
}

func deployedRequest(
	t *testing.T,
	client *http.Client,
	method string,
	url string,
	expectedStatus int,
	token string,
	body any,
) map[string]any {
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

	requestID := fmt.Sprintf("ci-%s-%s-%s", getEnv("GITHUB_RUN_ID", "local"), method, request.URL.Path)
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

	var responsePayload map[string]any
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

func requiredStringField(t *testing.T, payload map[string]any, field string) string {
	t.Helper()

	value, ok := payload[field].(string)
	if !ok || value == "" {
		t.Fatalf("expected response field %q, got %#v", field, payload)
	}

	return value
}

func getEnv(key string, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}

	return fallback
}
