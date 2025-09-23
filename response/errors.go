package response

type ErrorResp struct {
	Error string `json:"error"`
}

func NewErrorResp(msg string) ErrorResp {
	return ErrorResp{
		Error: msg,
	}
}