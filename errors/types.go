package errors

func BadRequest(message, detail string) *Error {
	return New(400, message, detail)
}

func Unauthorized(message, detail string) *Error {
	return New(401, message, detail)
}

func PaymentRequired(message, detail string) *Error {
	return New(402, message, detail)
}

func Forbidden(message, detail string) *Error {
	return New(403, message, detail)
}

func NotFound(message, detail string) *Error {
	return New(404, message, detail)
}

func Conflict(message, detail string) *Error {
	return New(409, message, detail)
}

func UnprocessableEntity(message, detail string) *Error {
	return New(422, message, detail)
}

func InternalServer(message, detail string) *Error {
	return New(500, message, detail)
}

func NotImplemented(message, detail string) *Error {
	return New(501, message, detail)
}

func ServiceUnavailable(message, detail string) *Error {
	return New(503, message, detail)
}
