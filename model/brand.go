package model

type Brand struct {
	Id   int    `json:"brandId"` // `json:"-"`
	Name string `json:"brandName"`
}
