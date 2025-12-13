package outorouter

type HTTPError interface {
	error
	StatusCode() int
	Code() string
	Message() string
}

type httpError struct {
	statusCode uint16
	code       string
	message    string
}

func (e *httpError) Error() string   { return e.message }
func (e *httpError) StatusCode() int { return int(e.statusCode) }
func (e *httpError) Code() string    { return e.code }
func (e *httpError) Message() string { return e.message }

func NewHTTPError(statusCode uint16, code, message string) HTTPError {
	return &httpError{
		statusCode: statusCode,
		code:       code,
		message:    message,
	}
}

func BadRequestError(code, message string) HTTPError {
	return NewHTTPError(400, code, message)
}

func UnauthorizedError(code, message string) HTTPError {
	return NewHTTPError(401, code, message)
}

func ForbiddenError(code, message string) HTTPError {
	return NewHTTPError(403, code, message)
}

func NotFoundError(code, message string) HTTPError {
	return NewHTTPError(404, code, message)
}

func InternalServerError(code, message string) HTTPError {
	return NewHTTPError(500, code, message)
}

func ServiceUnavailableError(code, message string) HTTPError {
	return NewHTTPError(503, code, message)
}
