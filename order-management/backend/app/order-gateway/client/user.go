package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	sharedConstants "tmossDev.github.com/eco-system/shared-components/backend/package/constants"
	"tmossDev.github.com/eco-system/shared-components/backend/package/types"
	userModel "tmossDev.github.com/eco-system/shared-components/backend/package/user/model"
)

const defaultUserServiceURL = "http://user-service:8080"

type AuthClient interface {
	Login(requestID string, body []byte) (*userModel.LoginResponse, error)
	Logout(requestID string, jwt string) error
	Shutdown()
}

type HTTPAuthClient struct {
	baseURL    string
	httpClient *http.Client
}

func NewHTTPAuthClient(baseURL string) *HTTPAuthClient {
	if baseURL == "" {
		baseURL = defaultUserServiceURL
	}

	return &HTTPAuthClient{
		baseURL: strings.TrimRight(baseURL, "/"),
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (client *HTTPAuthClient) Shutdown() {}

func (client *HTTPAuthClient) Login(requestID string, body []byte) (*userModel.LoginResponse, error) {
	request, err := http.NewRequest(http.MethodPost, client.baseURL+"/api/auth/login", bytes.NewReader(body))
	if err != nil {
		return nil, types.NewInternalServerError()
	}
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set(sharedConstants.CTXRequestIdKey, requestID)

	response, err := client.httpClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode < http.StatusOK || response.StatusCode >= http.StatusMultipleChoices {
		return nil, mapStatusError(response.StatusCode)
	}

	var loginResponse userModel.LoginResponse
	if err := json.NewDecoder(response.Body).Decode(&loginResponse); err != nil {
		return nil, types.NewInternalServerError()
	}

	return &loginResponse, nil
}

func (client *HTTPAuthClient) Logout(requestID string, jwt string) error {
	request, err := http.NewRequest(http.MethodPost, client.baseURL+"/api/auth/logout", nil)
	if err != nil {
		return types.NewInternalServerError()
	}
	request.Header.Set("Authorization", "Bearer "+jwt)
	request.Header.Set(sharedConstants.CTXRequestIdKey, requestID)

	response, err := client.httpClient.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode < http.StatusOK || response.StatusCode >= http.StatusMultipleChoices {
		return mapStatusError(response.StatusCode)
	}

	_, _ = io.Copy(io.Discard, response.Body)
	return nil
}

func mapStatusError(statusCode int) error {
	switch statusCode {
	case http.StatusBadRequest:
		return types.NewInvalidInputError()
	case http.StatusUnauthorized, http.StatusForbidden:
		return types.NewUnauthorizedError()
	default:
		if statusCode >= http.StatusInternalServerError {
			return types.NewInternalServerError()
		}
		return fmt.Errorf("user service returned status %d", statusCode)
	}
}
