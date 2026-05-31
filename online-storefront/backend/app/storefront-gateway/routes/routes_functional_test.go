package routes_test

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/kataras/iris/v12"
	"tmossDev.github.com/eco-system/online-storefront/backend/app/storefront-gateway/client"
	"tmossDev.github.com/eco-system/online-storefront/backend/app/storefront-gateway/routes"
	"tmossDev.github.com/eco-system/product-management/backend/domain/product/model"
	sharedConstants "tmossDev.github.com/eco-system/shared-components/backend/package/constants"
	"tmossDev.github.com/eco-system/shared-components/backend/package/transport/http/middleware"
	userModel "tmossDev.github.com/eco-system/shared-components/backend/package/user/model"
)

type fakeUserClient struct {
	loginRequestID    string
	logoutRequestID   string
	registerRequestID string
	userRequestID     string
	token             string
}

func (client *fakeUserClient) Register(requestID string, _ []byte) (*userModel.LoginResponse, error) {
	client.registerRequestID = requestID
	return client.loginResponse("New Customer"), nil
}

func (client *fakeUserClient) Login(requestID string, _ []byte) (*userModel.LoginResponse, error) {
	client.loginRequestID = requestID
	return client.loginResponse("Storefront Customer"), nil
}

func (client *fakeUserClient) Logout(requestID string, jwt string) error {
	client.logoutRequestID = requestID
	client.token = jwt
	return nil
}

func (client *fakeUserClient) User(requestID string, jwt string) (*userModel.UserResponse, error) {
	client.userRequestID = requestID
	client.token = jwt
	return &userModel.UserResponse{ID: 2, FirstName: "Storefront", LastName: "Customer", Email: "customer@test.com"}, nil
}

func (client *fakeUserClient) Shutdown() {}

func (client *fakeUserClient) loginResponse(name string) *userModel.LoginResponse {
	return &userModel.LoginResponse{
		Jwt:         "storefront-token",
		AccessToken: "storefront-token",
		ExpireAt:    time.Now().Add(24 * time.Hour).Unix(),
		User: userModel.AuthUserResponse{
			ID:    "2",
			Name:  name,
			Email: "customer@test.com",
			Role:  "Customer",
		},
	}
}

type fakeProductClient struct {
	requestID string
	products  []model.ProductResponse
}

func (productClient *fakeProductClient) ListProducts(requestID string) ([]model.ProductResponse, error) {
	productClient.requestID = requestID
	return productClient.products, nil
}

func (productClient *fakeProductClient) GetProduct(requestID string, productID string) (*model.ProductResponse, error) {
	productClient.requestID = requestID
	for _, product := range productClient.products {
		if productID == "1" && product.ID == 1 || productID == "2" && product.ID == 2 {
			return &product, nil
		}
	}

	return &productClient.products[0], nil
}

func (productClient *fakeProductClient) GetProductMedia(requestID string, _ string) (*client.ProductMedia, error) {
	productClient.requestID = requestID
	body := "image bytes"
	return &client.ProductMedia{
		Body:          io.NopCloser(strings.NewReader(body)),
		ContentLength: int64(len(body)),
		ContentType:   "image/jpeg",
	}, nil
}

func (productClient *fakeProductClient) Shutdown() {}

type testServer struct {
	app           *iris.Application
	userClient    *fakeUserClient
	productClient *fakeProductClient
}

func newTestServer(t *testing.T) *testServer {
	t.Helper()

	app := iris.New()
	app.Use(middleware.RequestIDMiddleware)

	userClient := &fakeUserClient{}
	productClient := &fakeProductClient{
		products: []model.ProductResponse{
			{ID: 1, SKU: "ACTIVE-1", Name: "Active Product", Status: "Active"},
			{ID: 2, SKU: "DRAFT-1", Name: "Draft Product", Status: "Draft"},
		},
	}
	routes.Setup(app, userClient, productClient)

	if err := app.Build(); err != nil {
		t.Fatalf("build iris app: %v", err)
	}

	return &testServer{app: app, userClient: userClient, productClient: productClient}
}

func (server *testServer) doJSON(method string, path string, token string, body any, requestID string) *httptest.ResponseRecorder {
	var requestBody bytes.Buffer
	if body != nil {
		_ = json.NewEncoder(&requestBody).Encode(body)
	}

	request := httptest.NewRequest(method, path, &requestBody)
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set(sharedConstants.CTXRequestIdKey, requestID)
	if token != "" {
		request.Header.Set("Authorization", "Bearer "+token)
	}

	response := httptest.NewRecorder()
	server.app.ServeHTTP(response, request)
	return response
}

func TestRegister(t *testing.T) {
	server := newTestServer(t)
	response := server.doJSON(http.MethodPost, "/api/register", "", map[string]string{"email": "customer@test.com"}, "register-id")

	if response.Code != http.StatusOK || server.userClient.registerRequestID != "register-id" {
		t.Fatalf("expected delegated registration, got status %d and request id %q", response.Code, server.userClient.registerRequestID)
	}
	if response.Header().Get("Set-Cookie") == "" {
		t.Fatal("expected registration to set a session cookie")
	}
}

func TestLogin(t *testing.T) {
	server := newTestServer(t)
	response := server.doJSON(http.MethodPost, "/api/auth/login", "", map[string]string{"email": "customer@test.com"}, "login-id")

	if response.Code != http.StatusOK || server.userClient.loginRequestID != "login-id" {
		t.Fatalf("expected delegated login, got status %d and request id %q", response.Code, server.userClient.loginRequestID)
	}
}

func TestLogout(t *testing.T) {
	server := newTestServer(t)
	response := server.doJSON(http.MethodPost, "/api/auth/logout", "storefront-token", nil, "logout-id")

	if response.Code != http.StatusOK || server.userClient.token != "storefront-token" || server.userClient.logoutRequestID != "logout-id" {
		t.Fatalf("expected delegated logout, got status %d, token %q and request id %q", response.Code, server.userClient.token, server.userClient.logoutRequestID)
	}
}

func TestCurrentUser(t *testing.T) {
	server := newTestServer(t)
	response := server.doJSON(http.MethodGet, "/api/users/me", "storefront-token", nil, "user-id")

	if response.Code != http.StatusOK || server.userClient.userRequestID != "user-id" {
		t.Fatalf("expected delegated current user lookup, got status %d and request id %q", response.Code, server.userClient.userRequestID)
	}
}

func TestListProductsOnlyReturnsActiveCatalogEntries(t *testing.T) {
	server := newTestServer(t)
	response := server.doJSON(http.MethodGet, "/api/products", "", nil, "products-id")

	var products []model.ProductResponse
	if err := json.Unmarshal(response.Body.Bytes(), &products); err != nil {
		t.Fatalf("decode products response: %v", err)
	}
	if response.Code != http.StatusOK || len(products) != 1 || products[0].SKU != "ACTIVE-1" {
		t.Fatalf("expected only active products, got status %d and %#v", response.Code, products)
	}
}

func TestProductDetailsHidesDraftEntries(t *testing.T) {
	server := newTestServer(t)
	response := server.doJSON(http.MethodGet, "/api/products/2", "", nil, "draft-id")

	if response.Code != http.StatusNotFound {
		t.Fatalf("expected hidden draft product to return %d, got %d: %s", http.StatusNotFound, response.Code, response.Body.String())
	}
}

func TestGetProductMedia(t *testing.T) {
	server := newTestServer(t)
	response := server.doJSON(http.MethodGet, "/api/product-media/products/1/photo.jpg", "", nil, "media-id")

	if response.Code != http.StatusOK || response.Header().Get("Content-Type") != "image/jpeg" || response.Body.String() != "image bytes" {
		t.Fatalf("expected proxied media, got status %d, content type %q and body %q", response.Code, response.Header().Get("Content-Type"), response.Body.String())
	}
}
