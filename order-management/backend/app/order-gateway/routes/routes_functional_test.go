package routes_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/kataras/iris/v12"
	"tmossDev.github.com/eco-system/order-management/backend/app/order-gateway/routes"
	"tmossDev.github.com/eco-system/order-management/backend/domain/order/model"
	userModel "tmossDev.github.com/eco-system/shared-components/backend/package/user/model"
)

type fakeOrderService struct {
	action  string
	orderID uint64
	body    string
}

type fakeAuthClient struct{}

func (service *fakeOrderService) ListOrders() ([]model.OrderResponse, error) {
	service.action = "list"
	return []model.OrderResponse{orderResponse(1)}, nil
}

func (service *fakeOrderService) GetOrder(orderID uint64) (*model.OrderResponse, error) {
	service.action, service.orderID = "get", orderID
	order := orderResponse(orderID)
	return &order, nil
}

func (service *fakeOrderService) UpdateStatus(orderID uint64, body string) (*model.OrderResponse, error) {
	service.action, service.orderID, service.body = "status", orderID, body
	order := orderResponse(orderID)
	order.Status = "Order Fulfillment"
	return &order, nil
}

func (client *fakeAuthClient) Login(requestID string, body []byte) (*userModel.LoginResponse, error) {
	return &userModel.LoginResponse{Jwt: "jwt", AccessToken: "jwt"}, nil
}

func (client *fakeAuthClient) Logout(requestID string, jwt string) error {
	return nil
}

func (client *fakeAuthClient) Shutdown() {}

func orderResponse(orderID uint64) model.OrderResponse {
	return model.OrderResponse{ID: orderID, UserID: 42, CartID: 7, Status: "Order Submitted", Items: []model.OrderItem{}}
}

func newTestServer(t *testing.T) (*iris.Application, *fakeOrderService) {
	t.Helper()
	app := iris.New()
	orderService := &fakeOrderService{}
	routes.Setup(app, orderService, &fakeAuthClient{})
	if err := app.Build(); err != nil {
		t.Fatalf("build iris app: %v", err)
	}
	return app, orderService
}

func doJSON(app *iris.Application, method string, path string, body any) *httptest.ResponseRecorder {
	var requestBody bytes.Buffer
	if body != nil {
		_ = json.NewEncoder(&requestBody).Encode(body)
	}
	request := httptest.NewRequest(method, path, &requestBody)
	if path == "/api/auth/logout" {
		request.Header.Set("Authorization", "Bearer jwt")
	}
	response := httptest.NewRecorder()
	app.ServeHTTP(response, request)
	return response
}

func TestLogin(t *testing.T) {
	app, _ := newTestServer(t)
	response := doJSON(app, http.MethodPost, "/api/auth/login", map[string]any{"email": "admin@test.com", "password": "password"})
	if response.Code != http.StatusOK {
		t.Fatalf("expected login, got status %d: %s", response.Code, response.Body.String())
	}
}

func TestListOrders(t *testing.T) {
	app, service := newTestServer(t)
	response := doJSON(app, http.MethodGet, "/api/orders", nil)
	if response.Code != http.StatusOK || service.action != "list" {
		t.Fatalf("expected list orders, got status %d and service %#v", response.Code, service)
	}
}

func TestOrderDetails(t *testing.T) {
	app, service := newTestServer(t)
	response := doJSON(app, http.MethodGet, "/api/orders/9", nil)
	if response.Code != http.StatusOK || service.action != "get" || service.orderID != 9 {
		t.Fatalf("expected order details, got status %d and service %#v", response.Code, service)
	}
}

func TestUpdateOrderStatus(t *testing.T) {
	app, service := newTestServer(t)
	response := doJSON(app, http.MethodPut, "/api/orders/9/status", map[string]any{"status": "Order Fulfillment"})
	if response.Code != http.StatusOK || service.action != "status" || service.orderID != 9 {
		t.Fatalf("expected order status update, got status %d and service %#v", response.Code, service)
	}
}
