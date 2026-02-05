package session

type Payload struct {
	Key      string `json:"key,omitempty"`
	Nonce    string `json:"nonce,omitempty"`
	Data     string `json:"data,omitempty"`
	TagLen   int    `json:"tagLen,omitempty"`
	Padding  string `json:"padding,omitempty"`
	NonceLen int    `json:"nonceLen,omitempty"`
}
