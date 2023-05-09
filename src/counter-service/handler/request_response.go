package handler

import (
	"encoding/json"
	"io"
)

type CounterResponse struct {
	Count    int64  `json:"count"`
	Hostname string `json:"hostname"`
}

func (p *CounterResponse) ToJSON(w io.Writer) error {
	encoder := json.NewEncoder(w)
	return encoder.Encode(p)
}
