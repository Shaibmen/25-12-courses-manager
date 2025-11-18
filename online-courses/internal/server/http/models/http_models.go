package models

type HttpResponse struct {
	Message string `json:"message"`
}

type HttpResponseWithData struct {
	Data any `json:"data"`
}
