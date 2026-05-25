package postgres

import (
	"database/sql"
	"errors"

	"tmossDev.github.com/eco-system/product-management/backend/domain/promotion/model"
	"tmossDev.github.com/eco-system/product-management/backend/domain/promotion/repository"
	"tmossDev.github.com/eco-system/shared-components/backend/package/datastore"
	"tmossDev.github.com/eco-system/shared-components/backend/package/datastore/flows"
	"tmossDev.github.com/eco-system/shared-components/backend/package/logger"
	"tmossDev.github.com/eco-system/shared-components/backend/package/types"
)

type PromotionRepository struct {
	store datastore.DataStore
}

func NewPostgresPromotionRepository(store datastore.DataStore) repository.PromotionRepository {
	return &PromotionRepository{
		store: store,
	}
}

func (repo *PromotionRepository) Shutdown() {
	if err := repo.store.Close(); err != nil {
		logger.Errorf("Unable to close promotion repo: %s", err.Error())
	}
}

func (repo *PromotionRepository) GetPromotionSettings() (*model.PromotionSettings, error) {
	stmt, err := flows.GetReaderStatement("GetPromotionSettings", GetPromotionSettings, repo.store)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	var settings model.PromotionSettings
	err = stmt.QueryRow().Scan(
		&settings.PromotionsEnabled,
		&settings.UpdatedUser,
		&settings.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, types.NewNoTFoundOrNoRecordError()
		}

		logger.Errorf("Unable to get promotion settings: %s", err.Error())
		return nil, types.NewInternalServerError()
	}

	return &settings, nil
}

func (repo *PromotionRepository) UpdatePromotionSettings(settings model.PromotionSettings) error {
	result, err := repo.store.GetConnection().GetWriter().Exec(
		UpdatePromotionSettings,
		settings.PromotionsEnabled,
		settings.UpdatedUser,
	)
	if err != nil {
		logger.Errorf("Unable to update promotion settings: %s", err.Error())
		return types.NewInternalServerError()
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return types.NewNoTFoundOrNoRecordError()
	}

	return nil
}
