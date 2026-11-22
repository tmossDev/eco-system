package aws

import (
	"context"
	"encoding/json"

	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"tmossDev.github.com/eco-system/shared-components/backend/package/config"
	"tmossDev.github.com/eco-system/shared-components/backend/package/config/types"

	log "tmossDev.github.com/eco-system/shared-components/backend/package/logger"
)

type secretsClient interface {
	GetSecretValue(context.Context, *secretsmanager.GetSecretValueInput, ...func(*secretsmanager.Options)) (*secretsmanager.GetSecretValueOutput, error)
}

// Config loads application configuration from AWS Secrets Manager.
type Config struct {
	client secretsClient
}

func NewSecretConfigManager() config.Config {
	cfg, err := awsconfig.LoadDefaultConfig(context.Background())
	if err != nil {
		panic(err)
	}

	return &Config{
		client: secretsmanager.NewFromConfig(cfg),
	}
}

func (sh *Config) getSecretValue(secretName string) (*string, error) {
	input := &secretsmanager.GetSecretValueInput{
		SecretId: &secretName,
	}

	result, err := sh.client.GetSecretValue(context.Background(), input)
	if err != nil {
		return nil, err
	}

	var secretValue string
	if result == nil {
		return nil, err
	}

	if result.SecretString != nil {
		secretValue = *result.SecretString
	} else if result.SecretBinary != nil {
		secretValue = string(result.SecretBinary)
	}

	return &secretValue, err
}

func (sh *Config) GetConfig(configName string) (*types.ConfigModel, error) {
	log.Debugf("SYSTEM", "Loading secret '%s'...", configName)
	secretValue, err := sh.getSecretValue(configName)
	if err != nil {
		return nil, err
	}

	var secret map[string]interface{}
	err = json.Unmarshal([]byte(*secretValue), &secret)
	if err != nil {
		return nil, err
	}

	var mergedSecret interface{} = secret
	if secret["secrets"] != nil {
		// multi secrets list
		nestedSecrets := secret["secrets"].([]interface{})
		for _, v := range nestedSecrets {
			log.Debugf("SYSTEM", "Loading inner secret '%s'...", v)
			nestedSecretValue, err := sh.getSecretValue(v.(string))
			if err != nil {
				return nil, err
			}

			var nestedSecret map[string]interface{}
			err = json.Unmarshal([]byte(*nestedSecretValue), &nestedSecret)
			if err != nil {
				return nil, err
			}

			mergedSecret, err = merge(mergedSecret, nestedSecret)
			if err != nil {
				return nil, err
			}
		}
	}

	mergedSecretJson, err := json.Marshal(mergedSecret)
	if err != nil {
		return nil, err
	}

	var config types.ConfigModel
	err = json.Unmarshal(mergedSecretJson, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}

func merge(root, nested interface{}) (interface{}, error) {
	rootJson, err := json.Marshal(root)
	if err != nil {
		return nil, err
	}
	nestedJson, err := json.Marshal(nested)
	if err != nil {
		return nil, err
	}

	var rootData interface{}
	err = json.Unmarshal(rootJson, &rootData)
	if err != nil {
		return nil, err
	}
	var nestedData interface{}
	err = json.Unmarshal(nestedJson, &nestedData)
	if err != nil {
		return nil, err
	}

	return deepMerge(rootData, nestedData), nil
}

func deepMerge(rootData, nestedData interface{}) interface{} {
	switch rootData := rootData.(type) {
	case map[string]interface{}:
		nestedData, ok := nestedData.(map[string]interface{})
		if !ok {
			return rootData
		}
		for key, nestedValue := range nestedData {
			if rootValue, ok := rootData[key]; ok {
				rootData[key] = deepMerge(rootValue, nestedValue)
			} else {
				rootData[key] = nestedValue
			}
		}
	case nil:
		nestedData, ok := nestedData.(map[string]interface{})
		if ok {
			return nestedData
		}
	}

	return rootData
}
