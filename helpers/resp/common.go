package resp

import "net/http"

func Response(data interface{}, code int, message interface{}) Resp {
	return Resp{
		Code:    code,
		Message: message,
		Data:    data,
	}
}

func Success(data interface{}, code int) Resp {
	return Resp{
		Code:    code,
		Message: http.StatusText(http.StatusOK),
		Data:    data,
	}
}

func BadRequest(message interface{}) Resp {
	return Resp{
		Code:    http.StatusBadRequest,
		Message: message,
		Data:    nil,
	}
}

func MissingHeader() Resp {
	return Resp{
		Code:    http.StatusUnauthorized,
		Message: http.StatusText(http.StatusUnauthorized),
		Data:    nil,
	}
}

func InternalServerError() Resp {
	return Resp{
		Code:    http.StatusInternalServerError,
		Message: http.StatusText(http.StatusInternalServerError),
		Data:    nil,
	}
}
