package response

import "time"

func ErrorResponse(err error, message string, errors ...map[string]string) map[string]interface{} {
	errorMessage := ""
	if err != nil {
		errorMessage = err.Error()
	}
	_errors := map[string]string{}
	if len(errors) > 0 {
		_errors = errors[0]
	}
	response := map[string]interface{}{
		"detail":    message,
		"error":     errorMessage,
		"timestamp": time.Now().Format(time.RFC3339),
		"errors":    _errors,
	}
	return response
}
