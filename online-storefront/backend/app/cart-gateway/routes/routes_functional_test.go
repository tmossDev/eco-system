package routes_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/kataras/iris/v12"
	"tmossDev.github.com/eco-system/online-storefront/backend/app/cart-gateway/routes"
	"tmossDev.github.com/eco-system/online-storefront/backend/domain/cart/model"
	orderModel "tmossDev.github.com/eco-system/online-storefront/backend/domain/order/model"
	userConstants "tmossDev.github.com/eco-system/shared-components/backend/package/user/constants"
	"tmossDev.github.com/eco-system/shared-components/backend/package/utils"
)

type fakeCartService struct {
	action    string
	userID    uint64
	productID uint64
	body      string
}

type fakeOrderService struct {
	action string
	userID uint64
}

func (service *fakeCartService) GetCurrent(userID uint64) (*model.CartResponse, error) {
	service.action, service.userID = "get", userID
	return cartResponse(userID), nil
}

func (service *fakeCartService) AddItem(userID uint64, body string) (*model.CartResponse, error) {
	service.action, service.userID, service.body = "add", userID, body
	return cartResponse(userID), nil
}

func (service *fakeCartService) UpdateItem(userID uint64, productID uint64, body string) (*model.CartResponse, error) {
	service.action, service.userID, service.productID, service.body = "update", userID, productID, body
	return cartResponse(userID), nil
}

func (service *fakeCartService) RemoveItem(userID uint64, productID uint64) (*model.CartResponse, error) {
	service.action, service.userID, service.productID = "remove", userID, productID
	return cartResponse(userID), nil
}

func (service *fakeCartService) Clear(userID uint64) (*model.CartResponse, error) {
	service.action, service.userID = "clear", userID
	return cartResponse(userID), nil
}

func (service *fakeOrderService) Checkout(userID uint64) (*orderModel.OrderResponse, error) {
	service.action, service.userID = "checkout", userID
	return &orderModel.OrderResponse{ID: 99, UserID: userID, CartID: 1, Status: "Order Submitted", Items: []orderModel.OrderItem{}}, nil
}

func (service *fakeOrderService) ListOrders(userID uint64) ([]orderModel.OrderResponse, error) {
	service.action, service.userID = "orders", userID
	return []orderModel.OrderResponse{{ID: 99, UserID: userID, CartID: 1, Status: "Order Submitted", Items: []orderModel.OrderItem{}}}, nil
}

func cartResponse(userID uint64) *model.CartResponse {
	return &model.CartResponse{ID: 1, UserID: userID, Items: []model.CartItem{}}
}

func newTestServer(t *testing.T) (*iris.Application, *fakeCartService, *fakeOrderService, string) {
	t.Helper()
	app := iris.New()
	cartService := &fakeCartService{}
	orderService := &fakeOrderService{}
	routes.Setup(app, cartService, orderService)
	if err := app.Build(); err != nil {
		t.Fatalf("build iris app: %v", err)
	}
	token, _, err := utils.GenerateJwt("42", userConstants.PASSWORD_SECRET_HASHING_KEY)
	if err != nil {
		t.Fatalf("generate jwt: %v", err)
	}
	return app, cartService, orderService, token
}

func doJSON(app *iris.Application, method string, path string, token string, body any) *httptest.ResponseRecorder {
	var requestBody bytes.Buffer
	if body != nil {
		_ = json.NewEncoder(&requestBody).Encode(body)
	}
	request := httptest.NewRequest(method, path, &requestBody)
	if token != "" {
		request.Header.Set("Authorization", "Bearer "+token)
	}
	response := httptest.NewRecorder()
	app.ServeHTTP(response, request)
	return response
}

func TestCartRoutesRequireAuthentication(t *testing.T) {
	app, _, _, _ := newTestServer(t)
	response := doJSON(app, http.MethodGet, "/api/cart", "", nil)
	if response.Code != http.StatusUnauthorized {
		t.Fatalf("expected unauthorized response, got %d: %s", response.Code, response.Body.String())
	}
}

func TestGetCurrentCart(t *testing.T) {
	app, service, _, token := newTestServer(t)
	response := doJSON(app, http.MethodGet, "/api/cart", token, nil)
	if response.Code != http.StatusOK || service.action != "get" || service.userID != 42 {
		t.Fatalf("expected current cart lookup, got status %d and service %#v", response.Code, service)
	}
}

func TestAddCartItem(t *testing.T) {
	app, service, _, token := newTestServer(t)
	response := doJSON(app, http.MethodPost, "/api/cart/items", token, map[string]any{"product_id": 7, "quantity": 2})
	if response.Code != http.StatusCreated || service.action != "add" || service.userID != 42 {
		t.Fatalf("expected cart item add, got status %d and service %#v", response.Code, service)
	}
}

func TestUpdateCartItem(t *testing.T) {
	app, service, _, token := newTestServer(t)
	response := doJSON(app, http.MethodPut, "/api/cart/items/7", token, map[string]any{"quantity": 3})
	if response.Code != http.StatusOK || service.action != "update" || service.productID != 7 {
		t.Fatalf("expected cart item update, got status %d and service %#v", response.Code, service)
	}
}

func TestRemoveCartItem(t *testing.T) {
	app, service, _, token := newTestServer(t)
	response := doJSON(app, http.MethodDelete, "/api/cart/items/7", token, nil)
	if response.Code != http.StatusOK || service.action != "remove" || service.productID != 7 {
		t.Fatalf("expected cart item removal, got status %d and service %#v", response.Code, service)
	}
}

func TestClearCart(t *testing.T) {
	app, service, _, token := newTestServer(t)
	response := doJSON(app, http.MethodDelete, "/api/cart", token, nil)
	if response.Code != http.StatusOK || service.action != "clear" {
		t.Fatalf("expected cart clear, got status %d and service %#v", response.Code, service)
	}
}

func TestCheckoutCart(t *testing.T) {
	app, _, service, token := newTestServer(t)
	response := doJSON(app, http.MethodPost, "/api/cart/checkout", token, nil)
	if response.Code != http.StatusCreated || service.action != "checkout" || service.userID != 42 {
		t.Fatalf("expected cart checkout, got status %d and service %#v", response.Code, service)
	}
}

func TestListOrders(t *testing.T) {
	app, _, service, token := newTestServer(t)
	response := doJSON(app, http.MethodGet, "/api/orders", token, nil)
	if response.Code != http.StatusOK || service.action != "orders" || service.userID != 42 {
		t.Fatalf("expected order list, got status %d and service %#v", response.Code, service)
	}
}
