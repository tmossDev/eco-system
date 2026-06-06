package service

import (
	"testing"

	"tmossDev.github.com/eco-system/order-management/backend/domain/order/model"
	"tmossDev.github.com/eco-system/shared-components/backend/package/validator"
)

type fakeOrderRepository struct {
	order         model.OrderResponse
	updatedID     uint64
	updatedStatus string
}

func (repo *fakeOrderRepository) ListOrders() ([]model.OrderResponse, error) {
	return []model.OrderResponse{repo.order}, nil
}

func (repo *fakeOrderRepository) GetOrder(orderID uint64) (*model.OrderResponse, error) {
	order := repo.order
	order.ID = orderID
	return &order, nil
}

func (repo *fakeOrderRepository) UpdateStatus(orderID uint64, status string) (*model.OrderResponse, error) {
	repo.updatedID = orderID
	repo.updatedStatus = status
	order := repo.order
	order.ID = orderID
	order.Status = status
	return &order, nil
}

func TestUpdateStatusAllowsValidTransition(t *testing.T) {
	repo := &fakeOrderRepository{order: model.OrderResponse{Status: model.OrderStatusSubmitted}}
	service := NewOrderService(validator.NewValidator(), repo)

	order, err := service.UpdateStatus(7, `{"status":"Order Confirmed"}`)
	if err != nil {
		t.Fatalf("expected update to succeed: %v", err)
	}
	if repo.updatedID != 7 || repo.updatedStatus != model.OrderStatusConfirmed || order.Status != model.OrderStatusConfirmed {
		t.Fatalf("expected confirmed update, got repo=%#v order=%#v", repo, order)
	}
}

func TestUpdateStatusRejectsInvalidTransition(t *testing.T) {
	repo := &fakeOrderRepository{order: model.OrderResponse{Status: model.OrderStatusSubmitted}}
	service := NewOrderService(validator.NewValidator(), repo)

	if _, err := service.UpdateStatus(7, `{"status":"Order Returned"}`); err == nil {
		t.Fatal("expected invalid transition to fail")
	}
	if repo.updatedStatus != "" {
		t.Fatalf("expected invalid transition not to persist, got %q", repo.updatedStatus)
	}
}
