package postgres

import (
	"database/sql"
	"encoding/json"
	"errors"

	"tmossDev.github.com/eco-system/product-management/backend/domain/product/model"
	"tmossDev.github.com/eco-system/product-management/backend/domain/product/repository"
	"tmossDev.github.com/eco-system/shared-components/backend/package/datastore"
	"tmossDev.github.com/eco-system/shared-components/backend/package/datastore/flows"
	"tmossDev.github.com/eco-system/shared-components/backend/package/logger"
	"tmossDev.github.com/eco-system/shared-components/backend/package/types"
)

type ProductRepository struct {
	store datastore.DataStore
}

func NewPostgresProductRepository(store datastore.DataStore) repository.ProductRepository {
	return &ProductRepository{
		store: store,
	}
}

func (repo *ProductRepository) Shutdown() {
	if err := repo.store.Close(); err != nil {
		logger.Errorf("Unable to close product repo: %s", err.Error())
	}
}

func mapRowToProduct(row *sql.Row) (*model.ProductResponse, error) {
	var product model.ProductResponse
	var photosJSON []byte
	var discountsJSON []byte

	if row.Err() != nil {
		return nil, row.Err()
	}

	err := row.Scan(
		&product.ID,
		&product.SKU,
		&product.Name,
		&product.ShortDescription,
		&product.Description,
		&product.Category,
		&product.PriceCents,
		&product.Currency,
		&product.InventoryCount,
		&product.Status,
		&photosJSON,
		&discountsJSON,
		&product.CreatedUser,
		&product.CreatedAt,
		&product.UpdatedUser,
		&product.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, types.NewNoTFoundOrNoRecordError()
		}

		logger.Errorf("Unable to map product response: %s", err.Error())
		return nil, types.NewInternalServerError()
	}

	product.Photos = decodeProductPhotos(photosJSON)
	product.Discounts = decodeDiscounts(discountsJSON)

	return &product, nil
}

func scanProductRows(rows *sql.Rows) ([]model.ProductResponse, error) {
	products := make([]model.ProductResponse, 0)
	for rows.Next() {
		var product model.ProductResponse
		var photosJSON []byte
		var discountsJSON []byte
		if err := rows.Scan(
			&product.ID,
			&product.SKU,
			&product.Name,
			&product.ShortDescription,
			&product.Description,
			&product.Category,
			&product.PriceCents,
			&product.Currency,
			&product.InventoryCount,
			&product.Status,
			&photosJSON,
			&discountsJSON,
			&product.CreatedUser,
			&product.CreatedAt,
			&product.UpdatedUser,
			&product.UpdatedAt,
		); err != nil {
			logger.Errorf("Unable to scan product row: %s", err.Error())
			return nil, types.NewInternalServerError()
		}

		product.Photos = decodeProductPhotos(photosJSON)
		product.Discounts = decodeDiscounts(discountsJSON)
		products = append(products, product)
	}

	if err := rows.Err(); err != nil {
		logger.Errorf("Unable to read product rows: %s", err.Error())
		return nil, types.NewInternalServerError()
	}

	return products, nil
}

func (repo *ProductRepository) List() ([]model.ProductResponse, error) {
	stmt, err := flows.GetReaderStatement("ListProducts", ListProducts, repo.store)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query()
	if err != nil {
		logger.Errorf("Unable to list products: %s", err.Error())
		return nil, types.NewInternalServerError()
	}
	defer rows.Close()

	return scanProductRows(rows)
}

func (repo *ProductRepository) GetByID(productID uint64) (*model.ProductResponse, error) {
	stmt, err := flows.GetReaderStatement("GetProductByID", GetByID, repo.store)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	return mapRowToProduct(stmt.QueryRow(productID))
}

func (repo *ProductRepository) Create(product model.ProductRequest, creatingUserID uint64) (*model.ProductResponse, error) {
	var insertedID uint64
	err := repo.store.GetConnection().GetWriter().QueryRow(
		CreateProduct,
		product.SKU,
		product.Name,
		product.ShortDescription,
		product.Description,
		product.Category,
		product.PriceCents,
		product.Currency,
		product.InventoryCount,
		product.Status,
		encodeProductPhotos(product.Photos),
		creatingUserID,
	).Scan(&insertedID)
	if err != nil {
		logger.Errorf("Unable to create product: %s", err.Error())
		return nil, types.NewInternalServerError()
	}

	return repo.GetByID(insertedID)
}

func (repo *ProductRepository) Update(productID uint64, product model.ProductUpdateRequest, updatingUserID uint64) (*model.ProductResponse, error) {
	result, err := repo.store.GetConnection().GetWriter().Exec(
		UpdateProduct,
		product.SKU,
		product.Name,
		product.ShortDescription,
		product.Description,
		product.Category,
		product.PriceCents,
		product.Currency,
		product.InventoryCount,
		product.Status,
		encodeProductPhotos(product.Photos),
		updatingUserID,
		productID,
	)
	if err != nil {
		logger.Errorf("Unable to update product: %s", err.Error())
		return nil, types.NewInternalServerError()
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return nil, types.NewNoTFoundOrNoRecordError()
	}

	return repo.GetByID(productID)
}

func (repo *ProductRepository) Delete(productID uint64, deletingUserID uint64) error {
	result, err := repo.store.GetConnection().GetWriter().Exec(DeleteProduct, deletingUserID, productID)
	if err != nil {
		logger.Errorf("Unable to delete product: %s", err.Error())
		return types.NewInternalServerError()
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return types.NewNoTFoundOrNoRecordError()
	}

	return nil
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

func (repo *ProductRepository) ListDiscounts() ([]model.Discount, error) {
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

func (repo *ProductRepository) GetDiscountByID(discountID uint64) (*model.Discount, error) {
	stmt, err := flows.GetReaderStatement("GetDiscountByID", GetDiscountByID, repo.store)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	return mapRowToDiscount(stmt.QueryRow(discountID))
}

func (repo *ProductRepository) CreateDiscount(discount model.DiscountRequest, creatingUserID uint64) (*model.Discount, error) {
	tx, err := repo.store.GetConnection().GetWriter().Begin()
	if err != nil {
		logger.Errorf("Unable to start discount create transaction: %s", err.Error())
		return nil, types.NewInternalServerError()
	}
	defer tx.Rollback()

	var insertedID uint64
	err = tx.QueryRow(
		CreateDiscount,
		discount.Name,
		discount.Description,
		discount.DiscountType,
		discount.Scope,
		discount.PercentageBasisPoints,
		discount.AmountCents,
		discount.Currency,
		discount.MinProductCount,
		discount.StartsAt,
		discount.EndsAt,
		discount.Status,
		creatingUserID,
	).Scan(&insertedID)
	if err != nil {
		logger.Errorf("Unable to create discount: %s", err.Error())
		return nil, types.NewInternalServerError()
	}

	if err := replaceDiscountProducts(tx, insertedID, discount.ProductIDs); err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		logger.Errorf("Unable to commit discount create transaction: %s", err.Error())
		return nil, types.NewInternalServerError()
	}

	return repo.GetDiscountByID(insertedID)
}

func (repo *ProductRepository) UpdateDiscount(discountID uint64, discount model.DiscountUpdateRequest, updatingUserID uint64) (*model.Discount, error) {
	tx, err := repo.store.GetConnection().GetWriter().Begin()
	if err != nil {
		logger.Errorf("Unable to start discount update transaction: %s", err.Error())
		return nil, types.NewInternalServerError()
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
		discount.MinProductCount,
		discount.StartsAt,
		discount.EndsAt,
		discount.Status,
		updatingUserID,
		discountID,
	)
	if err != nil {
		logger.Errorf("Unable to update discount: %s", err.Error())
		return nil, types.NewInternalServerError()
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return nil, types.NewNoTFoundOrNoRecordError()
	}

	if err := replaceDiscountProducts(tx, discountID, discount.ProductIDs); err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		logger.Errorf("Unable to commit discount update transaction: %s", err.Error())
		return nil, types.NewInternalServerError()
	}

	return repo.GetDiscountByID(discountID)
}

func (repo *ProductRepository) DeleteDiscount(discountID uint64, deletingUserID uint64) error {
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

func encodeProductPhotos(photos []model.ProductPhotoRequest) string {
	if photos == nil {
		return "[]"
	}

	encoded, err := json.Marshal(photos)
	if err != nil {
		logger.Errorf("Unable to encode product photos: %s", err.Error())
		return "[]"
	}

	return string(encoded)
}

func decodeProductPhotos(photosJSON []byte) []model.ProductPhoto {
	if len(photosJSON) == 0 {
		return []model.ProductPhoto{}
	}

	var photos []model.ProductPhoto
	if err := json.Unmarshal(photosJSON, &photos); err != nil {
		logger.Errorf("Unable to decode product photos: %s", err.Error())
		return []model.ProductPhoto{}
	}

	return photos
}

func decodeDiscounts(discountsJSON []byte) []model.Discount {
	if len(discountsJSON) == 0 {
		return []model.Discount{}
	}

	var discounts []model.Discount
	if err := json.Unmarshal(discountsJSON, &discounts); err != nil {
		logger.Errorf("Unable to decode product discounts: %s", err.Error())
		return []model.Discount{}
	}

	return discounts
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
