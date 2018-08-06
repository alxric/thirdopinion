package config

// WSRequest is the request we send to the websocket
type WSRequest struct {
	Method   string    `json:"method"`
	Argument *Argument `json:"argument"`
	Vote     *Vote     `json:"vote"`
	Register *Register `json:"register"`
	User     *User     `json:"user"`
}

// WSResponse is the response from the websocket
type WSResponse struct {
	Error     string      `json:"error"`
	Msg       string      `json:"msg"`
	User      *User       `json:"user"`
	Arguments []*Argument `json:"arguments"`
}
