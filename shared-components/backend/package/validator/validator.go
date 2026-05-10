package validator

import (
	"encoding/json"

	"github.com/go-playground/validator/v10"
	"tmossDev.github.com/eco-system/shared-components/backend/package/logger"
	"tmossDev.github.com/eco-system/shared-components/backend/package/types"
	"tmossDev.github.com/eco-system/shared-components/backend/package/utils"
)

type Validator interface {
	MarshalAndValidateREQ(body string, request any) error
	GetValidator() validator.Validate
}

type SimpleValidator struct {
	Validate validator.Validate
}

func NewValidator() Validator {
	return &SimpleValidator{
		Validate: *validator.New(),
	}
}

func (validator *SimpleValidator) GetValidator() validator.Validate {
	return validator.Validate
}

func (validator *SimpleValidator) MarshalAndValidateREQ(body string, request any) error {
	body = utils.FormatJSONString(body)
	err := json.Unmarshal([]byte(body), &request)
	if err != nil {
		logger.Infof("Error marshaling request: %s", err.Error())
		return types.NewInvalidInputError()
	}

	err = validator.Validate.Struct(request)
	if err != nil {
		logger.Infof("Error validating: %s", err.Error())
		return types.NewInvalidInputError()
	}

	return nil
}
