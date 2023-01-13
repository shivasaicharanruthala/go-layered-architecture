package service

import (
	"layeredArchitecture/model"
)

type ServiceInterface interface {
	GetById(int) (*model.Products, error)
	PostProduct(model.Products) (*model.Products, error)
	DeleteProductById(int) (int, error)
	UpdateProductById(products model.Products) (*model.Products, error)
}
