package client

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"time"

	sharedConstants "tmossDev.github.com/eco-system/shared-components/backend/package/constants"
	"tmossDev.github.com/eco-system/shared-components/backend/package/types"
	userModel "tmossDev.github.com/eco-system/shared-components/backend/package/user/model"
)

const defaultUserServiceURL = "http://user-service:8080"

type UserClient interface {
	Register(requestID string, body []byte) (*userModel.LoginResponse, error)
	Login(requestID string, body []byte) (*userModel.LoginResponse, error)
	Logout(requestID string, jwt string) error
	User(requestID string, jwt string) (*userModel.UserResponse, error)
	Shutdown()
}

type HTTPUserClient struct {
	baseURL    string
	httpClient *http.Client
}

func NewHTTPUserClient(baseURL string) *HTTPUserClient {
	if baseURL == "" {
		baseURL = defaultUserServiceURL
	}

	return &HTTPUserClient{
		baseURL: strings.TrimRight(baseURL, "/"),
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (client *HTTPUserClient) Shutdown() {}

func (client *HTTPUserClient) Register(requestID string, body []byte) (*userModel.LoginResponse, error) {
	return client.authenticate(requestID, "/api/register", body)
}

func (client *HTTPUserClient) Login(requestID string, body []byte) (*userModel.LoginResponse, error) {
	return client.authenticate(requestID, "/api/auth/login", body)
}

func (client *HTTPUserClient) authenticate(requestID string, path string, body []byte) (*userModel.LoginResponse, error) {
	request, err := http.NewRequest(http.MethodPost, client.baseURL+path, bytes.NewReader(body))
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
		return nil, mapStatusError("user service", response.StatusCode)
	}

	var loginResponse userModel.LoginResponse
	if err := json.NewDecoder(response.Body).Decode(&loginResponse); err != nil {
		return nil, types.NewInternalServerError()
	}

	return &loginResponse, nil
}

func (client *HTTPUserClient) Logout(requestID string, jwt string) error {
	request, err := client.authorizedRequest(http.MethodPost, "/api/auth/logout", requestID, jwt)
	if err != nil {
		return err
	}

	response, err := client.httpClient.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode < http.StatusOK || response.StatusCode >= http.StatusMultipleChoices {
		return mapStatusError("user service", response.StatusCode)
	}

	_, _ = io.Copy(io.Discard, response.Body)
	return nil
}

func (client *HTTPUserClient) User(requestID string, jwt string) (*userModel.UserResponse, error) {
	request, err := client.authorizedRequest(http.MethodGet, "/api/users/me", requestID, jwt)
	if err != nil {
		return nil, err
	}

	response, err := client.httpClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode < http.StatusOK || response.StatusCode >= http.StatusMultipleChoices {
		return nil, mapStatusError("user service", response.StatusCode)
	}

	var userResponse userModel.UserResponse
	if err := json.NewDecoder(response.Body).Decode(&userResponse); err != nil {
		return nil, types.NewInternalServerError()
	}

	return &userResponse, nil
}

func (client *HTTPUserClient) authorizedRequest(method string, path string, requestID string, jwt string) (*http.Request, error) {
	request, err := http.NewRequest(method, client.baseURL+path, nil)
	if err != nil {
		return nil, types.NewInternalServerError()
	}
	request.Header.Set("Authorization", "Bearer "+jwt)
	request.Header.Set(sharedConstants.CTXRequestIdKey, requestID)

	return request, nil
}
