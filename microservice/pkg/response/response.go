package response

type IResponse interface {
	BasicError(interface{}, int) ErrorResponse
	Data(int, interface{}) DataResponse
}

type Response struct{}

func New() *Response {
	return &Response{}
}

// ErrorResponse is a struct that contains error message and status
type ErrorResponse struct {
	Error  string `json:"error"`
	Status int    `json:"status"`
}

// DataResponse is a struct that contains data and status
type DataResponse struct {
	Data   interface{} `json:"data"`
	Status int         `json:"status"`
}

// BasicError is a function that returns a ErrorResponse
func (r *Response) BasicError(d interface{}, status int) ErrorResponse {
	var err ErrorResponse
	switch d.(type) {
	case error:
		err.Error = d.(error).Error()
	case string:
		err.Error = d.(string)
	case nil:
		err.Error = "unknown error"
	default:
		err.Error = "unknown error"
	}
	err.Status = status
	return err
}

// Data is a function that returns a DataResponse
func (r *Response) Data(status int, data interface{}) DataResponse {
	return DataResponse{
		Data:   data,
		Status: status,
	}
}
