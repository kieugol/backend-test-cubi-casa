package resp

type Resp struct {
	Code    int         `json:"code,omitempty"`
	Message interface{} `json:"message"`
	Data    interface{} `json:"data"`
}
