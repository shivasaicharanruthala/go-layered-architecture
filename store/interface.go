package store

import "layeredArchitecture/model"

type BrandInterface interface {
	GetByIdBrand(int) (*model.Brand, error)
	PostBrand(brand model.Brand) (*model.Brand, error)
	CheckBrandExsists(string) (int, error)
}


type ProductsInterface interface {
	GetByIdProducts(int) (*model.Products, error)
	PostProduct(model.Products) (*model.Products, error)
	DeleteById(int) (int, error)
	UpdateProduct(products model.Products) (int, error)
}
