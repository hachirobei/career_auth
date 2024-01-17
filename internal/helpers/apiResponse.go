package helpers

// ApiResponse structure for standardizing API responses
type ApiResponse struct {
    Code    int         `json:"code"`
    Success bool        `json:"success"`
    Record  interface{} `json:"record,omitempty"`
    Message string      `json:"message,omitempty"`
}

// NewApiResponse creates a new ApiResponse
func NewApiResponse(code int, success bool, record interface{}, message string) ApiResponse {
    return ApiResponse{
        Code:    code,
        Success: success,
        Record:  record,
        Message: message,
    }
}