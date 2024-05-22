package controller

type Response struct {
	StatusCode int    `json:"status_code"`
	StatusMsg  string `json:"status_msg,omitempty"`
	Data       any    `json:"data,omitempty"`
}

const (
	StatusCodeSuccess = 0

	StatusCodeInvalidParams = 10
	StatusCodeNotLogin      = 12

	StatusCodeOperationFailed = 20
)

var statusMessages = map[int]string{
	StatusCodeSuccess: "success",

	StatusCodeInvalidParams: "invalid params",
	StatusCodeNotLogin:      "user does not login",

	StatusCodeOperationFailed: "operation failed",
}

func WithStatusCode(statusCode int) *Response {
	resp := &Response{StatusCode: statusCode}
	if msg, ok := statusMessages[statusCode]; !ok {
		resp.StatusMsg = msg
	}
	return resp
}

func WithStatus(statusCode int, statusMsg string) *Response {
	return &Response{StatusCode: statusCode, StatusMsg: statusMsg}
}

func Success(data any) *Response {
	return &Response{StatusCode: StatusCodeSuccess, StatusMsg: statusMessages[StatusCodeSuccess], Data: data}
}
