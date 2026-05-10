package types

// Default is a struct used in helper functions
type Default struct {
	RequestID string `json:"requestID"`
	Status    string `json:"status"`
	Message   string `json:"message"`
}

type DefaultData struct {
	RequestID string      `json:"requestID"`
	Status    string      `json:"status"`
	Message   string      `json:"message"`
	Data      interface{} `json:"data"`
}

type DefaultList struct {
	RequestID string      `json:"requestID"`
	Status    string      `json:"status"`
	Message   string      `json:"message"`
	Page      int64       `json:"page"`
	PageSize  int64       `json:"pageSize"`
	Total     int64       `json:"total"`
	Data      interface{} `json:"data"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type DefaultErrorResponse struct {
	RequestId string `json:"requestId"`
	Code      string `json:"code"`
	Error     string `json:"error"`
}
