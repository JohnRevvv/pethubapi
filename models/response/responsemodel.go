package response

type ResponseModel struct {
	RetCode string      `json:"retCode"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

type GameResponseModel struct {
	
	RetCode string      `json:"retCode"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

type AdopterResponseModel struct {
	RetCode string      `json:"retCode"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

type ShelterResponseModel struct {
	RetCode string      `json:"retCode"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}