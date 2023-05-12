package error

type Error struct {
	HTTPCode int
	Message  string
}

type ErrorCollection struct {
	Errors []Error
}

func (ec *ErrorCollection) AddHTTPError(httpCode int, err error) {
	newError := Error{
		HTTPCode: httpCode,
		Message:  err.Error(),
	}
	ec.Errors = append(ec.Errors, newError)
}

func (ec *ErrorCollection) HasErrors() bool {
	if len(ec.Errors) > 0 {
		return true
	}
	return false
}
