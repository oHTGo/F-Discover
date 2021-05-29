package interfaces

type ISuccess struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type ISuccessNoData struct {
	Message string `json:"message"`
}

type IFail struct {
	Message string `json:"message"`
}

type IFailWithErrors struct {
	Message string      `json:"message"`
	Errors  interface{} `json:"errors"`
}
