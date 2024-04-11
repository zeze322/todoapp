package response

type Response struct {
	Status string `json:"status"`
	Error  string `json:"error"`
}

const (
	StatusError = "Error"
)

func Error(msg string) Response {
	return Response{Status: StatusError, Error: msg}
}
