package types

type HttpResponseType struct {
	Title  string      `json:"title" jsonschema_description:"title of the response"`
	Status int         `json:"status"`
	Data   interface{} `json:"data"`
	Error  string      `json:"error"`
}

func NewHttpResponse(title, error string, stat int, data interface{}) HttpResponseType {
	return HttpResponseType{
		Title:  title,
		Status: stat,
		Data:   data,
		Error:  error,
	}
}
