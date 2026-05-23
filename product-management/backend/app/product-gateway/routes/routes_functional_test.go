package routes_test

import (
	"bytes"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/kataras/iris/v12"
	"tmossDev.github.com/eco-system/product-management/backend/app/product-gateway/routes"
	"tmossDev.github.com/eco-system/product-management/backend/domain/product/model"
	productService "tmossDev.github.com/eco-system/product-management/backend/domain/product/service"
	sharedConstants "tmossDev.github.com/eco-system/shared-components/backend/package/constants"
	transportHTTP "tmossDev.github.com/eco-system/shared-components/backend/package/transport/http"
	"tmossDev.github.com/eco-system/shared-components/backend/package/transport/http/middleware"
	httpTypes "tmossDev.github.com/eco-system/shared-components/backend/package/transport/http/types"
	userConstants "tmossDev.github.com/eco-system/shared-components/backend/package/user/constants"
	userModel "tmossDev.github.com/eco-system/shared-components/backend/package/user/model"
	"tmossDev.github.com/eco-system/shared-components/backend/package/utils"
)

type fakeProductService struct {
	products  []model.ProductResponse
	discounts []model.Discount
}

func newFakeProductService() *fakeProductService {
	return &fakeProductService{
		discounts: []model.Discount{
			{
				ID:                    1,
				Name:                  "Store-wide launch offer",
				DiscountType:          "Percentage",
				Scope:                 "Global",
				PercentageBasisPoints: int64Ptr(1000),
				MinProductCount:       1,
				Status:                "Active",
				CreatedUser:           1,
				CreatedAt:             time.Now().Format(time.RFC3339),
			},
		},
		products: []model.ProductResponse{
			{
				ID:               1,
				SKU:              "GEN-MUG-001",
				Name:             "Everyday Ceramic Mug",
				ShortDescription: "Durable mug for daily coffee.",
				Description:      "A durable mug for daily coffee.",
				Category:         "Home",
				PriceCents:       1299,
				Currency:         "USD",
				InventoryCount:   48,
				Status:           "Active",
				Photos: []model.ProductPhoto{
					{
						URL:          "https://example.com/mug.jpg",
						ThumbnailURL: "https://example.com/mug-thumb.jpg",
						AltText:      "Ceramic mug",
						IsPrimary:    true,
					},
				},
				CreatedUser: 1,
				CreatedAt:   time.Now().Format(time.RFC3339),
			},
		},
	}
}

func int64Ptr(value int64) *int64 {
	return &value
}

func (service *fakeProductService) ListProducts() ([]model.ProductResponse, error) {
	return service.products, nil
}

func (service *fakeProductService) GetProduct(productID uint64) (*model.ProductResponse, error) {
	for _, product := range service.products {
		if product.ID == productID {
			return &product, nil
		}
	}

	return &service.products[0], nil
}

func (service *fakeProductService) CreateProduct(string, uint64) (*model.ProductResponse, error) {
	product := service.products[0]
	product.ID = 2
	product.SKU = "NEW-PRODUCT-002"
	service.products = append([]model.ProductResponse{product}, service.products...)

	return &product, nil
}

func (service *fakeProductService) UpdateProduct(productID uint64, _ string, _ uint64) (*model.ProductResponse, error) {
	product := service.products[0]
	product.ID = productID
	product.Name = "Updated Product"

	return &product, nil
}

func (service *fakeProductService) DeleteProduct(uint64, uint64) error {
	return nil
}

func (service *fakeProductService) ListDiscounts() ([]model.Discount, error) {
	return service.discounts, nil
}

func (service *fakeProductService) GetDiscount(discountID uint64) (*model.Discount, error) {
	for _, discount := range service.discounts {
		if discount.ID == discountID {
			return &discount, nil
		}
	}

	return &service.discounts[0], nil
}

func (service *fakeProductService) CreateDiscount(string, uint64) (*model.Discount, error) {
	discount := service.discounts[0]
	discount.ID = 2
	service.discounts = append([]model.Discount{discount}, service.discounts...)

	return &discount, nil
}

func (service *fakeProductService) UpdateDiscount(discountID uint64, _ string, _ uint64) (*model.Discount, error) {
	discount := service.discounts[0]
	discount.ID = discountID
	discount.Name = "Updated Discount"

	return &discount, nil
}

func (service *fakeProductService) DeleteDiscount(uint64, uint64) error {
	return nil
}

func (service *fakeProductService) UploadProductPhoto(uint64, string, io.Reader, uint64) (*model.ProductResponse, error) {
	product := service.products[0]
	product.Photos = append(product.Photos, model.ProductPhoto{
		URL:          "/api/product-media/products/1/test-detail.jpg",
		ThumbnailURL: "/api/product-media/products/1/test-thumb.jpg",
		AltText:      product.Name,
		IsPrimary:    len(product.Photos) == 0,
	})

	return &product, nil
}

func (service *fakeProductService) GetProductMedia(string) (*productService.ProductMediaObject, error) {
	return nil, nil
}

func (service *fakeProductService) Shutdown() {}

type fakeAuthClient struct {
	loginResponse   *userModel.LoginResponse
	loginRequestID  string
	logoutRequestID string
	logoutToken     string
}

func (client *fakeAuthClient) Login(requestID string, _ []byte) (*userModel.LoginResponse, error) {
	client.loginRequestID = requestID
	return client.loginResponse, nil
}

func (client *fakeAuthClient) Logout(requestID string, jwt string) error {
	client.logoutRequestID = requestID
	client.logoutToken = jwt
	return nil
}

func (client *fakeAuthClient) Shutdown() {}

type productGatewayTestServer struct {
	app        *iris.Application
	authClient *fakeAuthClient
	token      string
}

func newProductGatewayTestServer(t *testing.T) *productGatewayTestServer {
	t.Helper()

	token, _, err := utils.GenerateJwt("1", userConstants.PASSWORD_SECRET_HASHING_KEY)
	if err != nil {
		t.Fatalf("generate test jwt: %v", err)
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
	authClient := &fakeAuthClient{
		loginResponse: &userModel.LoginResponse{
			Jwt:         token,
			AccessToken: token,
			ExpireAt:    time.Now().Add(24 * time.Hour).Unix(),
			User: userModel.AuthUserResponse{
				ID:    "1",
				Name:  "Product Admin",
				Email: "admin@test.com",
				Role:  "Admin",
			},
		},
	}
	routes.Setup(app, newFakeProductService(), authClient)

	if err := app.Build(); err != nil {
		t.Fatalf("build iris app: %v", err)
	}

	return &productGatewayTestServer{
		app:        app,
		authClient: authClient,
		token:      token,
	}
}

func (server *productGatewayTestServer) doJSON(method string, path string, token string, body any) *httptest.ResponseRecorder {
	return server.doJSONWithRequestID(method, path, token, body, "test-request-id")
}

func (server *productGatewayTestServer) doJSONWithRequestID(method string, path string, token string, body any, requestID string) *httptest.ResponseRecorder {
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

func decodeJSONResponse(t *testing.T, response *httptest.ResponseRecorder) map[string]any {
	t.Helper()

	var payload map[string]any
	if err := json.Unmarshal(response.Body.Bytes(), &payload); err != nil {
		t.Fatalf("decode response body %q: %v", response.Body.String(), err)
	}

	return payload
}

func TestProductGatewayFunctionalLogin(t *testing.T) {
	server := newProductGatewayTestServer(t)
	const requestID = "login-request-id"

	response := server.doJSONWithRequestID(http.MethodPost, "/api/auth/login", "", map[string]any{
		"email":    "admin@test.com",
		"password": "password",
	}, requestID)

	if response.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d: %s", http.StatusOK, response.Code, response.Body.String())
	}

	payload := decodeJSONResponse(t, response)
	if payload["accessToken"] == "" {
		t.Fatalf("expected login response to include accessToken, got %#v", payload)
	}
	if server.authClient.loginRequestID != requestID {
		t.Fatalf("expected login request id %q to be passed to user service, got %q", requestID, server.authClient.loginRequestID)
	}
}

func TestProductGatewayFunctionalLogout(t *testing.T) {
	server := newProductGatewayTestServer(t)
	const requestID = "logout-request-id"

	response := server.doJSONWithRequestID(http.MethodPost, "/api/auth/logout", server.token, nil, requestID)

	if response.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d: %s", http.StatusOK, response.Code, response.Body.String())
	}
	if server.authClient.logoutToken != server.token {
		t.Fatalf("expected logout token to be passed to user service")
	}
	if server.authClient.logoutRequestID != requestID {
		t.Fatalf("expected logout request id %q to be passed to user service, got %q", requestID, server.authClient.logoutRequestID)
	}
}

func TestProductGatewayFunctionalListProducts(t *testing.T) {
	server := newProductGatewayTestServer(t)

	response := server.doJSON(http.MethodGet, "/api/products", server.token, nil)

	if response.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d: %s", http.StatusOK, response.Code, response.Body.String())
	}

	var products []map[string]any
	if err := json.Unmarshal(response.Body.Bytes(), &products); err != nil {
		t.Fatalf("decode products response: %v", err)
	}
	if len(products) == 0 {
		t.Fatal("expected seeded products")
	}
}

func TestProductGatewayFunctionalGetProduct(t *testing.T) {
	server := newProductGatewayTestServer(t)

	response := server.doJSON(http.MethodGet, "/api/products/1", server.token, nil)

	if response.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d: %s", http.StatusOK, response.Code, response.Body.String())
	}

	payload := decodeJSONResponse(t, response)
	if payload["sku"] != "GEN-MUG-001" {
		t.Fatalf("expected product response, got %#v", payload)
	}
}

func TestProductGatewayFunctionalCreateProduct(t *testing.T) {
	server := newProductGatewayTestServer(t)

	response := server.doJSON(http.MethodPost, "/api/products", server.token, productPayload("NEW-PRODUCT-002"))

	if response.Code != http.StatusCreated {
		t.Fatalf("expected status %d, got %d: %s", http.StatusCreated, response.Code, response.Body.String())
	}

	payload := decodeJSONResponse(t, response)
	if payload["id"] == "" {
		t.Fatalf("expected created product id, got %#v", payload)
	}
}

func TestProductGatewayFunctionalEditProduct(t *testing.T) {
	server := newProductGatewayTestServer(t)

	response := server.doJSON(http.MethodPut, "/api/products/1", server.token, productPayload("GEN-MUG-001"))

	if response.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d: %s", http.StatusOK, response.Code, response.Body.String())
	}

	payload := decodeJSONResponse(t, response)
	if payload["name"] != "Updated Product" {
		t.Fatalf("expected updated product response, got %#v", payload)
	}
}

func TestProductGatewayFunctionalUploadProductPhoto(t *testing.T) {
	server := newProductGatewayTestServer(t)

	var body bytes.Buffer
	writer := multipart.NewWriter(&body)
	fileWriter, err := writer.CreateFormFile("file", "test.jpg")
	if err != nil {
		t.Fatalf("create multipart file: %v", err)
	}
	_, _ = fileWriter.Write([]byte("fake image bytes"))
	if err := writer.Close(); err != nil {
		t.Fatalf("close multipart writer: %v", err)
	}

	request := httptest.NewRequest(http.MethodPost, "/api/products/1/photos", &body)
	request.Header.Set("Content-Type", writer.FormDataContentType())
	request.Header.Set(sharedConstants.CTXRequestIdKey, "test-request-id")
	request.Header.Set("Authorization", "Bearer "+server.token)

	response := httptest.NewRecorder()
	server.app.ServeHTTP(response, request)

	if response.Code != http.StatusCreated {
		t.Fatalf("expected status %d, got %d: %s", http.StatusCreated, response.Code, response.Body.String())
	}
}

func TestProductGatewayFunctionalDeleteProduct(t *testing.T) {
	server := newProductGatewayTestServer(t)

	response := server.doJSON(http.MethodDelete, "/api/products/1", server.token, nil)

	if response.Code != http.StatusNoContent {
		t.Fatalf("expected status %d, got %d: %s", http.StatusNoContent, response.Code, response.Body.String())
	}
}

func productPayload(sku string) map[string]any {
	return map[string]any{
		"sku":               sku,
		"name":              "Functional Test Product",
		"short_description": "Short functional test product.",
		"description":       "Created by a functional test.",
		"category":          "Testing",
		"price_cents":       1499,
		"currency":          "USD",
		"inventory_count":   12,
		"status":            "Active",
		"photos": []map[string]any{
			{
				"url":           "https://example.com/product.jpg",
				"thumbnail_url": "https://example.com/product-thumb.jpg",
				"alt_text":      "Functional test product",
				"is_primary":    true,
			},
		},
	}
}
