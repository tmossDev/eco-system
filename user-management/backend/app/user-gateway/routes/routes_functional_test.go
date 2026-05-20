package routes_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/kataras/iris/v12"
	transportHTTP "tmossDev.github.com/eco-system/shared-components/backend/package/transport/http"
	"tmossDev.github.com/eco-system/shared-components/backend/package/transport/http/middleware"
	httpTypes "tmossDev.github.com/eco-system/shared-components/backend/package/transport/http/types"
	"tmossDev.github.com/eco-system/shared-components/backend/package/utils"
	"tmossDev.github.com/eco-system/user-management/backend/app/user-gateway/routes"
	userConstants "tmossDev.github.com/eco-system/user-management/backend/package/user/constants"
	"tmossDev.github.com/eco-system/user-management/backend/package/user/model"
)

type fakePublicUserService struct {
	loginResponse *model.LoginResponse
	logoutToken   string
}

func (service *fakePublicUserService) IsAuthenticated(string) error {
	return nil
}

func (service *fakePublicUserService) IsAuthorized(string, string) error {
	return nil
}

func (service *fakePublicUserService) User(string) (*model.UserResponse, error) {
	return &model.UserResponse{
		ID:        2,
		FirstName: "Jane",
		LastName:  "Doe",
		Email:     "admin@test.com",
		RoleID:    1,
	}, nil
}

func (service *fakePublicUserService) Logout(jwt string) error {
	service.logoutToken = jwt
	return nil
}

func (service *fakePublicUserService) Register(string) (*model.LoginResponse, error) {
	return service.loginResponse, nil
}

func (service *fakePublicUserService) Login(string, string) (*model.LoginResponse, error) {
	return service.loginResponse, nil
}

func (service *fakePublicUserService) Shutdown() {}

type fakePrivateUserService struct{}

func (service *fakePrivateUserService) UpdateUserInfo(uint64, string, uint64) (*model.UserResponse, error) {
	return nil, nil
}

func (service *fakePrivateUserService) UpdateUserPassword(uint64, string, uint64) (*model.UserResponse, error) {
	return nil, nil
}

func (service *fakePrivateUserService) Shutdown() {}

type gatewayTestServer struct {
	app           *iris.Application
	publicService *fakePublicUserService
	token         string
}

func newGatewayTestServer(t *testing.T) *gatewayTestServer {
	t.Helper()

	token, expireAt, err := utils.GenerateJwt("2", userConstants.PASSWORD_SECRET_HASHING_KEY)
	if err != nil {
		t.Fatalf("generate test jwt: %v", err)
	}

	publicService := &fakePublicUserService{
		loginResponse: &model.LoginResponse{
			Jwt:         token,
			AccessToken: token,
			ExpireAt:    expireAt,
			User: model.AuthUserResponse{
				ID:    "2",
				Name:  "Jane Doe",
				Email: "admin@test.com",
				Role:  "Admin",
			},
		},
	}

	app := iris.New()
	jwtMiddleware := transportHTTP.NewJWTMiddleware(httpTypes.JWTConfig{
		SecretKey:     []byte(userConstants.PASSWORD_SECRET_HASHING_KEY),
		TokenExpiry:   72 * time.Hour,
		SigningMethod: jwt.SigningMethodHS256,
		TokenPrefix:   "Bearer ",
	})

	app.Use(
		middleware.RequestIDMiddleware,
		jwtMiddleware([]string{"/auth/login", "/login", "/auth/logout", "/logout", "/refresh", "/health"}),
	)
	routes.Setup(app, publicService, &fakePrivateUserService{})

	if err := app.Build(); err != nil {
		t.Fatalf("build iris app: %v", err)
	}

	return &gatewayTestServer{
		app:           app,
		publicService: publicService,
		token:         token,
	}
}

func (server *gatewayTestServer) doJSON(method string, path string, token string, body any) *httptest.ResponseRecorder {
	var requestBody bytes.Buffer
	if body != nil {
		_ = json.NewEncoder(&requestBody).Encode(body)
	}

	request := httptest.NewRequest(method, path, &requestBody)
	request.Header.Set("Content-Type", "application/json")
	if token != "" {
		request.Header.Set("Authorization", "Bearer "+token)
	}

	response := httptest.NewRecorder()
	server.app.ServeHTTP(response, request)

	return response
}

func decodeJSONResponse(t *testing.T, response *httptest.ResponseRecorder) map[string]any {
	t.Helper()

	var payload map[string]any
	if err := json.Unmarshal(response.Body.Bytes(), &payload); err != nil {
		t.Fatalf("decode response body %q: %v", response.Body.String(), err)
	}

	return payload
}

func TestGatewayFunctionalLogin(t *testing.T) {
	server := newGatewayTestServer(t)

	response := server.doJSON(http.MethodPost, "/api/auth/login", "", map[string]any{
		"email":    "admin@test.com",
		"password": "password",
	})

	if response.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d: %s", http.StatusOK, response.Code, response.Body.String())
	}

	payload := decodeJSONResponse(t, response)
	if payload["accessToken"] == "" {
		t.Fatalf("expected login response to include accessToken, got %#v", payload)
	}
	if payload["jwt"] == "" {
		t.Fatalf("expected login response to include jwt, got %#v", payload)
	}
}

func TestGatewayFunctionalLogout(t *testing.T) {
	server := newGatewayTestServer(t)

	response := server.doJSON(http.MethodPost, "/api/auth/logout", server.token, nil)

	if response.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d: %s", http.StatusOK, response.Code, response.Body.String())
	}
	if server.publicService.logoutToken != server.token {
		t.Fatalf("expected logout token to be passed to service")
	}
}

func TestGatewayFunctionalListUsers(t *testing.T) {
	server := newGatewayTestServer(t)

	response := server.doJSON(http.MethodGet, "/api/users", server.token, nil)

	if response.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d: %s", http.StatusOK, response.Code, response.Body.String())
	}

	var users []map[string]any
	if err := json.Unmarshal(response.Body.Bytes(), &users); err != nil {
		t.Fatalf("decode users response: %v", err)
	}
	if len(users) == 0 {
		t.Fatal("expected seeded users")
	}
}

func TestGatewayFunctionalGetUser(t *testing.T) {
	server := newGatewayTestServer(t)

	response := server.doJSON(http.MethodGet, "/api/users/2", server.token, nil)

	if response.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d: %s", http.StatusOK, response.Code, response.Body.String())
	}

	payload := decodeJSONResponse(t, response)
	if payload["email"] != "admin@test.com" {
		t.Fatalf("expected admin user response, got %#v", payload)
	}
}

func TestGatewayFunctionalCreateUser(t *testing.T) {
	server := newGatewayTestServer(t)

	response := server.doJSON(http.MethodPost, "/api/users", server.token, map[string]any{
		"name":   "New User",
		"email":  "new.user@test.com",
		"role":   "User",
		"status": "Active",
	})

	if response.Code != http.StatusCreated {
		t.Fatalf("expected status %d, got %d: %s", http.StatusCreated, response.Code, response.Body.String())
	}

	payload := decodeJSONResponse(t, response)
	if payload["id"] == "" {
		t.Fatalf("expected created user id, got %#v", payload)
	}
}

func TestGatewayFunctionalEditUser(t *testing.T) {
	server := newGatewayTestServer(t)

	response := server.doJSON(http.MethodPut, "/api/users/2", server.token, map[string]any{
		"name":   "Jane Updated",
		"email":  "jane.updated@test.com",
		"role":   "Admin",
		"status": "Active",
	})

	if response.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d: %s", http.StatusOK, response.Code, response.Body.String())
	}

	payload := decodeJSONResponse(t, response)
	if payload["id"] != "2" {
		t.Fatalf("expected updated user id 2, got %#v", payload)
	}
	if payload["name"] != "Jane Updated" {
		t.Fatalf("expected updated user name, got %#v", payload)
	}
}

func TestGatewayFunctionalDeleteUser(t *testing.T) {
	server := newGatewayTestServer(t)

	response := server.doJSON(http.MethodDelete, "/api/users/2", server.token, nil)

	if response.Code != http.StatusNoContent {
		t.Fatalf("expected status %d, got %d: %s", http.StatusNoContent, response.Code, response.Body.String())
	}
}
