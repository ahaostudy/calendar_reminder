package controller

type Response struct {
	StatusCode int    `json:"status_code"`
	StatusMsg  string `json:"status_msg,omitempty"`
	Data       any    `json:"data,omitempty"`
}

const (
	StatusCodeSuccess = 0

	StatusCodeInvalidParams = 11
	StatusCodeNotLogin      = 12

	StatusCodeOperationFailed = 20
)

var statusMessages = map[int]string{
	StatusCodeSuccess: "success",

	StatusCodeInvalidParams: "invalid params",
	StatusCodeNotLogin:      "user does not login",

	StatusCodeOperationFailed: "operation failed",
}

// WithStatusCode constructs a response body with a status code that is automatically populated with status messages
func WithStatusCode(statusCode int) *Response {
	resp := &Response{StatusCode: statusCode}
	if msg, ok := statusMessages[statusCode]; ok {
		resp.StatusMsg = msg
	}
	return resp
}

// WithStatus manually specify response status code and status message
func WithStatus(statusCode int, statusMsg string) *Response {
	return &Response{StatusCode: statusCode, StatusMsg: statusMsg}
}

// Success normal response
func Success(data any) *Response {
	return &Response{StatusCode: StatusCodeSuccess, StatusMsg: statusMessages[StatusCodeSuccess], Data: data}
}
