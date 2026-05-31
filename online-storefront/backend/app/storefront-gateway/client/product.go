package client

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"tmossDev.github.com/eco-system/product-management/backend/domain/product/model"
	sharedConstants "tmossDev.github.com/eco-system/shared-components/backend/package/constants"
	"tmossDev.github.com/eco-system/shared-components/backend/package/types"
)

const (
	defaultProductServiceURL = "http://product-service:8080"
	defaultProductGatewayURL = "http://product-gateway:8080"
)

type ProductMedia struct {
	Body          io.ReadCloser
	ContentLength int64
	ContentType   string
}

type ProductClient interface {
	ListProducts(requestID string) ([]model.ProductResponse, error)
	GetProduct(requestID string, productID string) (*model.ProductResponse, error)
	GetProductMedia(requestID string, objectKey string) (*ProductMedia, error)
	Shutdown()
}

type HTTPProductClient struct {
	serviceURL string
	gatewayURL string
	httpClient *http.Client
}

func NewHTTPProductClient(serviceURL string, gatewayURL string) *HTTPProductClient {
	if serviceURL == "" {
		serviceURL = defaultProductServiceURL
	}
	if gatewayURL == "" {
		gatewayURL = defaultProductGatewayURL
	}

	return &HTTPProductClient{
		serviceURL: strings.TrimRight(serviceURL, "/"),
		gatewayURL: strings.TrimRight(gatewayURL, "/"),
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (client *HTTPProductClient) Shutdown() {}

func (client *HTTPProductClient) ListProducts(requestID string) ([]model.ProductResponse, error) {
	response, err := client.get(client.serviceURL+"/api/products", requestID)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode < http.StatusOK || response.StatusCode >= http.StatusMultipleChoices {
		return nil, mapStatusError("product service", response.StatusCode)
	}

	var products []model.ProductResponse
	if err := json.NewDecoder(response.Body).Decode(&products); err != nil {
		return nil, types.NewInternalServerError()
	}

	return products, nil
}

func (client *HTTPProductClient) GetProduct(requestID string, productID string) (*model.ProductResponse, error) {
	response, err := client.get(client.serviceURL+"/api/products/"+url.PathEscape(productID), requestID)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode < http.StatusOK || response.StatusCode >= http.StatusMultipleChoices {
		return nil, mapStatusError("product service", response.StatusCode)
	}

	var product model.ProductResponse
	if err := json.NewDecoder(response.Body).Decode(&product); err != nil {
		return nil, types.NewInternalServerError()
	}

	return &product, nil
}

func (client *HTTPProductClient) GetProductMedia(requestID string, objectKey string) (*ProductMedia, error) {
	escapedObjectKey, err := escapeObjectKey(objectKey)
	if err != nil {
		return nil, err
	}

	response, err := client.get(client.gatewayURL+"/api/product-media/"+escapedObjectKey, requestID)
	if err != nil {
		return nil, err
	}

	if response.StatusCode < http.StatusOK || response.StatusCode >= http.StatusMultipleChoices {
		response.Body.Close()
		return nil, mapStatusError("product gateway", response.StatusCode)
	}

	return &ProductMedia{
		Body:          response.Body,
		ContentLength: response.ContentLength,
		ContentType:   response.Header.Get("Content-Type"),
	}, nil
}

func (client *HTTPProductClient) get(endpoint string, requestID string) (*http.Response, error) {
	request, err := http.NewRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, types.NewInternalServerError()
	}
	request.Header.Set(sharedConstants.CTXRequestIdKey, requestID)

	return client.httpClient.Do(request)
}

func escapeObjectKey(objectKey string) (string, error) {
	segments := strings.Split(strings.Trim(objectKey, "/"), "/")
	for index, segment := range segments {
		if segment == "" || segment == "." || segment == ".." {
			return "", types.NewInvalidInputError()
		}
		segments[index] = url.PathEscape(segment)
	}

	return strings.Join(segments, "/"), nil
}
