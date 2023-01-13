package model

type Products struct {
	Id    int    `json:"productId"`
	Name  string `json:"productName"`
	Brand Brand  `json:"brand"`
}
