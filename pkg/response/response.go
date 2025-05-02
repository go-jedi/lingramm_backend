package response

// Response represents a standardized structure for API responses.
// It includes a status flag, a message, and a data field that can hold any type of data.
type Response struct {
	Status  bool        `json:"status"`  // Indicates whether the request was successful.
	Message string      `json:"message"` // Provides a message describing the result of the request.
	Error   string      `json:"error"`   // Shows detailed error.
	Data    interface{} `json:"data"`    // Holds the actual data returned by the API.
}

// New creates and returns a new instance of Response with the provided status, message, and data.
func New(status bool, message string, err string, data interface{}) *Response {
	return &Response{
		Status:  status,
		Message: message,
		Error:   err,
		Data:    data,
	}
}
