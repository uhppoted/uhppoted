package test_vectors

type Test struct {
	TestName string     `json:"name"`
	Request  Request    `json:"request"`
	Response []Response `json:"responses"`
}

type Request struct {
	Values  map[string]any `json:"values"`
	Message []byte         `json:"message"`
}

type Response struct {
	Values  map[string]any `json:"values"`
	Message []byte         `json:"message"`
}

var Tests = []Test{
	GetAllControllers,
	GetController,
	SetIP,
	GetTime,
	SetTime,
}
