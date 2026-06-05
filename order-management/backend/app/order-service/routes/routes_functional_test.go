package routes_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/kataras/iris/v12"
	"tmossDev.github.com/eco-system/order-management/backend/app/order-service/routes"
	"tmossDev.github.com/eco-system/order-management/backend/domain/order/model"
)

type fakeOrderService struct {
	action  string
	orderID uint64
	body    string
}

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
	order.Status = "Fulfilled"
	return &order, nil
}

func orderResponse(orderID uint64) model.OrderResponse {
	return model.OrderResponse{ID: orderID, UserID: 42, CartID: 7, Status: "Created", Items: []model.OrderItem{}}
}

func newTestServer(t *testing.T) (*iris.Application, *fakeOrderService) {
	t.Helper()
	app := iris.New()
	orderService := &fakeOrderService{}
	routes.Setup(app, orderService)
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
	response := httptest.NewRecorder()
	app.ServeHTTP(response, request)
	return response
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
	response := doJSON(app, http.MethodPut, "/api/orders/9/status", map[string]any{"status": "Fulfilled"})
	if response.Code != http.StatusOK || service.action != "status" || service.orderID != 9 {
		t.Fatalf("expected order status update, got status %d and service %#v", response.Code, service)
	}
}
