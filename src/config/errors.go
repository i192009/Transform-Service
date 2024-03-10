package config

import "gitlab.zixel.cn/go/framework"

var ErrorCodes = map[int]string{
	1000: "Server internal error.",
	1001: "Database error.",

	30001: "Invalid request format.",
	30002: "Request validate failed.",
	30003: "Job not found.",
	30004: "Duplicate request",

	/// New Job
	30100: "Create job failed.",
	30101: "No Member Interest",
	30102: "targetFormat is not support",

	/// Get Job details
	30200: "Request header error.",
	30201: "Too many tasks, try again later. ",
	30202: "add Servers is  error. ",
	30203: "config is  error. ",
}

func NewErrorNo(code int, message string, err error) *ErrorNo {
	return &ErrorNo{Code: code, Message: message, Err: err}
}

type ErrorNo struct {
	Code    int
	Message string
	Err     error
}

func (err *ErrorNo) ErrCode() int {
	if err.Code == 0 {
		return 1001
	}
	return err.Code
}
func (err *ErrorNo) ErrMsg() string {
	if err.Message == "" {
		return ErrorCodes[err.Code]
	}
	return err.Message
}
func (err *ErrorNo) Error() error {
	if err.Err == nil {
		if err.Message == "" {
			err.Message = ErrorCodes[err.Code]
		}
		err.Err = framework.NewServiceError(err.Code, err.Message)
	}
	return err.Err
}
