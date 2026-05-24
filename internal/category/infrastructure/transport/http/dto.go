package http

type createCategoryRequest struct {
	Name string `json:"name"`
}

type updateCategoryRequest struct {
	Name string `json:"name"`
}

type errorResponse struct {
	Error struct {
		Code    string `json:"code"`
		Message string `json:"message"`
	} `json:"error"`
}

func newErrorResponse(code, message string) errorResponse {
	resp := errorResponse{}
	resp.Error.Code = code
	resp.Error.Message = message
	return resp
}
