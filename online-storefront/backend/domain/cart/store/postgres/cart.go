package postgres

import (
	"database/sql"
	"errors"

	"tmossDev.github.com/eco-system/online-storefront/backend/domain/cart/model"
	"tmossDev.github.com/eco-system/online-storefront/backend/domain/cart/repository"
	"tmossDev.github.com/eco-system/shared-components/backend/package/datastore"
	"tmossDev.github.com/eco-system/shared-components/backend/package/logger"
	"tmossDev.github.com/eco-system/shared-components/backend/package/types"
)

type CartRepository struct {
	store datastore.DataStore
}

func NewPostgresCartRepository(store datastore.DataStore) repository.CartRepository {
	return &CartRepository{store: store}
}

func (repo *CartRepository) GetCurrent(userID uint64) (*model.CartResponse, error) {
	tx, err := repo.store.GetConnection().GetWriter().Begin()
	if err != nil {
		return nil, internalError("start cart transaction", err)
	}
	defer rollback(tx)

	cart, err := ensureCart(tx, userID)
	if err != nil {
		return nil, err
	}
	if err := tx.Commit(); err != nil {
		return nil, internalError("commit cart transaction", err)
	}

	return repo.loadCart(cart)
}

func (repo *CartRepository) AddItem(userID uint64, productID uint64, quantity int64) (*model.CartResponse, error) {
	return repo.changeCart(userID, func(tx *sql.Tx, cartID uint64) error {
		var savedProductID uint64
		if err := tx.QueryRow(addItem, cartID, productID, quantity).Scan(&savedProductID); err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return types.NewInvalidInputError()
			}
			return internalError("add cart item", err)
		}
		return nil
	})
}

func (repo *CartRepository) UpdateItem(userID uint64, productID uint64, quantity int64) (*model.CartResponse, error) {
	return repo.changeCart(userID, func(tx *sql.Tx, cartID uint64) error {
		var savedProductID uint64
		if err := tx.QueryRow(updateItem, cartID, productID, quantity).Scan(&savedProductID); err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return types.NewInvalidInputError()
			}
			return internalError("update cart item", err)
		}
		return nil
	})
}

func (repo *CartRepository) RemoveItem(userID uint64, productID uint64) (*model.CartResponse, error) {
	return repo.changeCart(userID, func(tx *sql.Tx, cartID uint64) error {
		result, err := tx.Exec(removeItem, cartID, productID)
		if err != nil {
			return internalError("remove cart item", err)
		}
		rows, _ := result.RowsAffected()
		if rows == 0 {
			return types.NewNoTFoundOrNoRecordError()
		}
		return nil
	})
}

func (repo *CartRepository) Clear(userID uint64) (*model.CartResponse, error) {
	return repo.changeCart(userID, func(tx *sql.Tx, cartID uint64) error {
		if _, err := tx.Exec(clearCart, cartID); err != nil {
			return internalError("clear cart", err)
		}
		return nil
	})
}

func (repo *CartRepository) changeCart(userID uint64, action func(*sql.Tx, uint64) error) (*model.CartResponse, error) {
	tx, err := repo.store.GetConnection().GetWriter().Begin()
	if err != nil {
		return nil, internalError("start cart transaction", err)
	}
	defer rollback(tx)

	cart, err := ensureCart(tx, userID)
	if err != nil {
		return nil, err
	}
	if err := action(tx, cart.ID); err != nil {
		return nil, err
	}
	if _, err := tx.Exec(touchCart, cart.ID); err != nil {
		return nil, internalError("touch cart", err)
	}
	if err := tx.Commit(); err != nil {
		return nil, internalError("commit cart transaction", err)
	}

	return repo.loadCartByUser(userID)
}

func ensureCart(tx *sql.Tx, userID uint64) (*model.CartResponse, error) {
	return scanCart(tx.QueryRow(getOrCreateCart, userID))
}

func scanCart(row *sql.Row) (*model.CartResponse, error) {
	var cart model.CartResponse
	if err := row.Scan(&cart.ID, &cart.UserID, &cart.CreatedAt, &cart.UpdatedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, err
		}
		return nil, internalError("load cart", err)
	}
	cart.Items = []model.CartItem{}
	return &cart, nil
}

func (repo *CartRepository) loadCartByUser(userID uint64) (*model.CartResponse, error) {
	cart, err := scanCart(repo.store.GetConnection().GetWriter().QueryRow(getCart, userID))
	if err != nil {
		return nil, err
	}
	return repo.loadCart(cart)
}

func (repo *CartRepository) loadCart(cart *model.CartResponse) (*model.CartResponse, error) {
	rows, err := repo.store.GetConnection().GetWriter().Query(getCartItems, cart.ID)
	if err != nil {
		return nil, internalError("load cart items", err)
	}
	defer rows.Close()

	currency := ""
	for rows.Next() {
		var item model.CartItem
		if err := rows.Scan(&item.ProductID, &item.SKU, &item.Name, &item.Quantity, &item.PriceCents, &item.Currency, &item.ThumbnailURL); err != nil {
			return nil, internalError("scan cart item", err)
		}
		item.LineTotal = item.PriceCents * item.Quantity
		cart.Items = append(cart.Items, item)
		cart.ItemCount += item.Quantity
		cart.SubtotalCents += item.LineTotal
		if currency == "" {
			currency = item.Currency
		}
	}
	if err := rows.Err(); err != nil {
		return nil, internalError("read cart items", err)
	}
	cart.Currency = currency
	return cart, nil
}

func internalError(action string, err error) error {
	logger.Errorf("Unable to %s: %s", action, err.Error())
	return types.NewInternalServerError()
}

func rollback(tx *sql.Tx) {
	if err := tx.Rollback(); err != nil && !errors.Is(err, sql.ErrTxDone) {
		logger.Errorf("Unable to rollback cart transaction: %s", err.Error())
	}
}
