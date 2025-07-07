package utils

type CommonResponse struct {
	Message string `json:"message"`
	Status  string `json:"status"`
}

type CommonResponseWithData struct {
	Message string      `json:"message"`
	Status  string      `json:"status"`
	Data    interface{} `json:"data"`
}
