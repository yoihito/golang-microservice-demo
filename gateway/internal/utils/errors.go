package utils

type JSONError struct {
	Code    int    `json:"-"`
	Message string `json:"error"`
}

func (e JSONError) Error() string {
	return e.Message
}
