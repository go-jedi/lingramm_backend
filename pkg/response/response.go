package response

// Response represents a standardized structure for API responses.
// It includes a status flag, a message, and a data field that can hold any type of data.
type Response[T any] struct {
	Status  bool   `json:"status"`  // Indicates whether the request was successful.
	Message string `json:"message"` // Provides a message describing the result of the request.
	Error   string `json:"error"`   // Shows detailed error.
	Data    T      `json:"data"`    // Holds the actual data returned by the API.
}

// New creates and returns a new instance of Response with the provided status, message, and data.
func New[T any](status bool, message string, err string, data T) *Response[T] {
	return &Response[T]{
		Status:  status,
		Message: message,
		Error:   err,
		Data:    data,
	}
}
