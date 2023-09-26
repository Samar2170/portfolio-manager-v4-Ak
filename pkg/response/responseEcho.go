package response

import "net/http"

func JSONResponseEcho(data interface{}) (int, interface{}) {
	return http.StatusOK, data
}

func DataResponseEcho() {}

func SuccessResponseEcho(message string) (int, map[string]string) {
	return http.StatusOK, map[string]string{
		"message": message,
	}
}

func BadRequestResponseEcho(err string) (int, map[string]string) {
	return http.StatusBadRequest, map[string]string{
		"message": "Bad Request",
		"error":   err,
	}
}

func UnauthorizedResponseEcho(err string) (int, map[string]string) {
	return http.StatusForbidden, map[string]string{
		"message": "Unauthorized",
		"error":   err,
	}
}

func NotFoundResponseEcho(err string) (int, map[string]string) {
	return http.StatusNotFound, map[string]string{
		"message": "Not Found",
		"error":   err,
	}
}

func InternalServerErrorResponseEcho(err string) (int, map[string]string) {
	return http.StatusInternalServerError, map[string]string{
		"message": "Internal Server Error",
		"error":   err,
	}
}

func MethodNotAllowedResponseEcho(err string) (int, map[string]string) {
	return http.StatusMethodNotAllowed, map[string]string{
		"message": "Method Not Allowed",
		"error":   err,
	}
}
