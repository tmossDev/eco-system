package postgres

import (
	"database/sql"
	"encoding/json"
	"errors"

	"tmossDev.github.com/eco-system/product-management/backend/domain/promotion/model"
	"tmossDev.github.com/eco-system/product-management/backend/domain/promotion/repository"
	"tmossDev.github.com/eco-system/shared-components/backend/package/datastore"
	"tmossDev.github.com/eco-system/shared-components/backend/package/datastore/flows"
	"tmossDev.github.com/eco-system/shared-components/backend/package/logger"
	"tmossDev.github.com/eco-system/shared-components/backend/package/types"
)

type DiscountRepository struct {
	store datastore.DataStore
}

func NewPostgresDiscountRepository(store datastore.DataStore) repository.DiscountRepository {
	return &DiscountRepository{
		store: store,
	}
}

func (repo *DiscountRepository) Shutdown() {
	if err := repo.store.Close(); err != nil {
		logger.Errorf("Unable to close discount repo: %s", err.Error())
	}
}

func mapRowToDiscount(row *sql.Row) (*model.Discount, error) {
	var discount model.Discount
	var productIDsJSON []byte
	var percentageBasisPoints sql.NullInt64
	var amountCents sql.NullInt64

	if row.Err() != nil {
		return nil, row.Err()
	}

	err := row.Scan(
		&discount.ID,
		&discount.Name,
		&discount.Description,
		&discount.DiscountType,
		&discount.Scope,
		&percentageBasisPoints,
		&amountCents,
		&discount.Currency,
		&discount.BuyQuantity,
		&discount.FreeQuantity,
		&discount.MinProductCount,
		&discount.StartsAt,
		&discount.EndsAt,
		&discount.Status,
		&productIDsJSON,
		&discount.CreatedUser,
		&discount.CreatedAt,
		&discount.UpdatedUser,
		&discount.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, types.NewNoTFoundOrNoRecordError()
		}

		logger.Errorf("Unable to map discount response: %s", err.Error())
		return nil, types.NewInternalServerError()
	}

	if percentageBasisPoints.Valid {
		discount.PercentageBasisPoints = &percentageBasisPoints.Int64
	}
	if amountCents.Valid {
		discount.AmountCents = &amountCents.Int64
	}
	discount.ProductIDs = decodeProductIDs(productIDsJSON)

	return &discount, nil
}

func scanDiscountRows(rows *sql.Rows) ([]model.Discount, error) {
	discounts := make([]model.Discount, 0)
	for rows.Next() {
		var discount model.Discount
		var productIDsJSON []byte
		var percentageBasisPoints sql.NullInt64
		var amountCents sql.NullInt64
		if err := rows.Scan(
			&discount.ID,
			&discount.Name,
			&discount.Description,
			&discount.DiscountType,
			&discount.Scope,
			&percentageBasisPoints,
			&amountCents,
			&discount.Currency,
			&discount.BuyQuantity,
			&discount.FreeQuantity,
			&discount.MinProductCount,
			&discount.StartsAt,
			&discount.EndsAt,
			&discount.Status,
			&productIDsJSON,
			&discount.CreatedUser,
			&discount.CreatedAt,
			&discount.UpdatedUser,
			&discount.UpdatedAt,
		); err != nil {
			logger.Errorf("Unable to scan discount row: %s", err.Error())
			return nil, types.NewInternalServerError()
		}

		if percentageBasisPoints.Valid {
			discount.PercentageBasisPoints = &percentageBasisPoints.Int64
		}
		if amountCents.Valid {
			discount.AmountCents = &amountCents.Int64
		}
		discount.ProductIDs = decodeProductIDs(productIDsJSON)
		discounts = append(discounts, discount)
	}

	if err := rows.Err(); err != nil {
		logger.Errorf("Unable to read discount rows: %s", err.Error())
		return nil, types.NewInternalServerError()
	}

	return discounts, nil
}

func (repo *DiscountRepository) ListDiscounts() ([]model.Discount, error) {
	stmt, err := flows.GetReaderStatement("ListDiscounts", ListDiscounts, repo.store)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query()
	if err != nil {
		logger.Errorf("Unable to list discounts: %s", err.Error())
		return nil, types.NewInternalServerError()
	}
	defer rows.Close()

	return scanDiscountRows(rows)
}

func (repo *DiscountRepository) GetDiscountByID(discountID uint64) (*model.Discount, error) {
	stmt, err := flows.GetReaderStatement("GetDiscountByID", GetDiscountByID, repo.store)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	return mapRowToDiscount(stmt.QueryRow(discountID))
}

func (repo *DiscountRepository) CreateDiscount(discount *model.Discount) error {
	tx, err := repo.store.GetConnection().GetWriter().Begin()
	if err != nil {
		logger.Errorf("Unable to start discount create transaction: %s", err.Error())
		return types.NewInternalServerError()
	}
	defer tx.Rollback()

	err = tx.QueryRow(
		CreateDiscount,
		discount.Name,
		discount.Description,
		discount.DiscountType,
		discount.Scope,
		discount.PercentageBasisPoints,
		discount.AmountCents,
		discount.Currency,
		discount.BuyQuantity,
		discount.FreeQuantity,
		discount.MinProductCount,
		discount.StartsAt,
		discount.EndsAt,
		discount.Status,
		discount.CreatedUser,
	).Scan(&discount.ID)
	if err != nil {
		logger.Errorf("Unable to create discount: %s", err.Error())
		return types.NewInternalServerError()
	}

	if err := replaceDiscountProducts(tx, discount.ID, discount.ProductIDs); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		logger.Errorf("Unable to commit discount create transaction: %s", err.Error())
		return types.NewInternalServerError()
	}

	return nil
}

func (repo *DiscountRepository) UpdateDiscount(discount model.Discount) error {
	tx, err := repo.store.GetConnection().GetWriter().Begin()
	if err != nil {
		logger.Errorf("Unable to start discount update transaction: %s", err.Error())
		return types.NewInternalServerError()
	}
	defer tx.Rollback()

	result, err := tx.Exec(
		UpdateDiscount,
		discount.Name,
		discount.Description,
		discount.DiscountType,
		discount.Scope,
		discount.PercentageBasisPoints,
		discount.AmountCents,
		discount.Currency,
		discount.BuyQuantity,
		discount.FreeQuantity,
		discount.MinProductCount,
		discount.StartsAt,
		discount.EndsAt,
		discount.Status,
		discount.UpdatedUser,
		discount.ID,
	)
	if err != nil {
		logger.Errorf("Unable to update discount: %s", err.Error())
		return types.NewInternalServerError()
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return types.NewNoTFoundOrNoRecordError()
	}

	if err := replaceDiscountProducts(tx, discount.ID, discount.ProductIDs); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		logger.Errorf("Unable to commit discount update transaction: %s", err.Error())
		return types.NewInternalServerError()
	}

	return nil
}

func (repo *DiscountRepository) DeleteDiscount(discountID uint64, deletingUserID uint64) error {
	result, err := repo.store.GetConnection().GetWriter().Exec(DeleteDiscount, deletingUserID, discountID)
	if err != nil {
		logger.Errorf("Unable to delete discount: %s", err.Error())
		return types.NewInternalServerError()
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return types.NewNoTFoundOrNoRecordError()
	}

	return nil
}

func replaceDiscountProducts(tx *sql.Tx, discountID uint64, productIDs []uint64) error {
	if _, err := tx.Exec(ReplaceDiscountProducts, discountID); err != nil {
		logger.Errorf("Unable to clear discount products: %s", err.Error())
		return types.NewInternalServerError()
	}

	for _, productID := range productIDs {
		if _, err := tx.Exec(InsertDiscountProduct, discountID, productID); err != nil {
			logger.Errorf("Unable to assign discount product: %s", err.Error())
			return types.NewInternalServerError()
		}
	}

	return nil
}

func decodeProductIDs(productIDsJSON []byte) []uint64 {
	if len(productIDsJSON) == 0 {
		return []uint64{}
	}

	var productIDs []uint64
	if err := json.Unmarshal(productIDsJSON, &productIDs); err != nil {
		logger.Errorf("Unable to decode discount product ids: %s", err.Error())
		return []uint64{}
	}

	return productIDs
}
