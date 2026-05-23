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

	return &product, nil
}

func scanProductRows(rows *sql.Rows) ([]model.ProductResponse, error) {
	products := make([]model.ProductResponse, 0)
	for rows.Next() {
		var product model.ProductResponse
		var photosJSON []byte
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
			&product.CreatedUser,
			&product.CreatedAt,
			&product.UpdatedUser,
			&product.UpdatedAt,
		); err != nil {
			logger.Errorf("Unable to scan product row: %s", err.Error())
			return nil, types.NewInternalServerError()
		}

		product.Photos = decodeProductPhotos(photosJSON)
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
