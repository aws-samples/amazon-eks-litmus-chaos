package handler

import (
	"encoding/json"
	"io"
	"like-service/models"
)

type LikeRequest struct {
	Id int `json:"id"`
}

type LikesResponse struct {
	Likes    *[]models.Like `json:"likes"`
	Hostname string         `json:"hostname"`
}

func (p *LikesResponse) ToJSON(w io.Writer) error {
	encoder := json.NewEncoder(w)
	return encoder.Encode(p)
}

func (p *LikeRequest) FromJSON(r io.Reader) error {
	decoder := json.NewDecoder(r)
	return decoder.Decode(p)
}
