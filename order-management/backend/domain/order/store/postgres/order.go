package postgres

import (
	"database/sql"
	"errors"

	"tmossDev.github.com/eco-system/order-management/backend/domain/order/model"
	"tmossDev.github.com/eco-system/order-management/backend/domain/order/repository"
	"tmossDev.github.com/eco-system/shared-components/backend/package/datastore"
	"tmossDev.github.com/eco-system/shared-components/backend/package/logger"
	"tmossDev.github.com/eco-system/shared-components/backend/package/types"
)

type OrderRepository struct {
	store datastore.DataStore
}

func NewPostgresOrderRepository(store datastore.DataStore) repository.OrderRepository {
	return &OrderRepository{store: store}
}

func (repo *OrderRepository) ListOrders() ([]model.OrderResponse, error) {
	rows, err := repo.store.GetConnection().GetWriter().Query(listOrders)
	if err != nil {
		return nil, internalError("list orders", err)
	}
	defer rows.Close()

	orders := []model.OrderResponse{}
	for rows.Next() {
		var order model.OrderResponse
		if err := rows.Scan(&order.ID, &order.UserID, &order.CartID, &order.Status, &order.ItemCount, &order.SubtotalCents, &order.Currency, &order.CreatedAt, &order.UpdatedAt); err != nil {
			return nil, internalError("scan order", err)
		}
		order.Items = []model.OrderItem{}
		orders = append(orders, order)
	}
	if err := rows.Err(); err != nil {
		return nil, internalError("read orders", err)
	}

	for index := range orders {
		items, err := repo.loadOrderItems(orders[index].ID)
		if err != nil {
			return nil, err
		}
		orders[index].Items = items
	}
	return orders, nil
}

func (repo *OrderRepository) GetOrder(orderID uint64) (*model.OrderResponse, error) {
	return repo.loadOrder(orderID)
}

func (repo *OrderRepository) UpdateStatus(orderID uint64, status string) (*model.OrderResponse, error) {
	order, err := scanOrder(repo.store.GetConnection().GetWriter().QueryRow(updateOrderStatus, orderID, status))
	if err != nil {
		return nil, err
	}

	items, err := repo.loadOrderItems(order.ID)
	if err != nil {
		return nil, err
	}
	order.Items = items
	return order, nil
}

func scanOrder(row *sql.Row) (*model.OrderResponse, error) {
	var order model.OrderResponse
	if err := row.Scan(&order.ID, &order.UserID, &order.CartID, &order.Status, &order.ItemCount, &order.SubtotalCents, &order.Currency, &order.CreatedAt, &order.UpdatedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, types.NewNoTFoundOrNoRecordError()
		}
		return nil, internalError("load order", err)
	}
	order.Items = []model.OrderItem{}
	return &order, nil
}

func (repo *OrderRepository) loadOrder(orderID uint64) (*model.OrderResponse, error) {
	order, err := scanOrder(repo.store.GetConnection().GetWriter().QueryRow(getOrder, orderID))
	if err != nil {
		return nil, err
	}

	items, err := repo.loadOrderItems(order.ID)
	if err != nil {
		return nil, err
	}
	order.Items = items
	return order, nil
}

func (repo *OrderRepository) loadOrderItems(orderID uint64) ([]model.OrderItem, error) {
	rows, err := repo.store.GetConnection().GetWriter().Query(getOrderItems, orderID)
	if err != nil {
		return nil, internalError("load order items", err)
	}
	defer rows.Close()

	items := []model.OrderItem{}
	for rows.Next() {
		var item model.OrderItem
		if err := rows.Scan(&item.ProductID, &item.SKU, &item.Name, &item.Quantity, &item.PriceCents, &item.Currency, &item.LineTotal, &item.ThumbnailURL); err != nil {
			return nil, internalError("scan order item", err)
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, internalError("read order items", err)
	}
	return items, nil
}

func internalError(action string, err error) error {
	logger.Errorf("Unable to %s: %s", action, err.Error())
	return types.NewInternalServerError()
}
