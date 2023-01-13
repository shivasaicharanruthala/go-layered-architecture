package product

import (
	"errors"
	"layeredArchitecture/model"
	"layeredArchitecture/service"
	"layeredArchitecture/store"
	"reflect"
)

type dataStore struct {
	productsStore  store.ProductsInterface
	brandStore    store.BrandInterface
}

func New(pdb store.ProductsInterface, bdb store.BrandInterface) service.ServiceInterface {
	return &dataStore{productsStore: pdb, brandStore: bdb}
}

func (s *dataStore) GetById(id int) (*model.Products, error) {
	productsList, err := s.productsStore.GetByIdProducts(id)
	if err != nil {
		return &model.Products{}, errors.New("product Id not found") // nil,err
	}
	p := &model.Products{Id: 0, Brand: model.Brand{Id: 0}}
	if reflect.DeepEqual(productsList, p){
		return &model.Products{}, errors.New("product Id not found")
	}
	brandList, err := s.brandStore.GetByIdBrand(productsList.Brand.Id)
	if err != nil {
		return productsList, errors.New("Brand Id not found") //nil, err
	}
	productsList.Brand.Name = brandList.Name
	return productsList, nil
}

func (s *dataStore) PostProduct(p model.Products) (*model.Products, error) {
	brandId, err := s.brandStore.CheckBrandExsists(p.Brand.Name)
	if err != nil {
		return &model.Products{}, err
	}
	if brandId == 0 { // brand doesnt exsists
		brandDetails := model.Brand{Name: p.Brand.Name}
		newBrandDetails, err := s.brandStore.PostBrand(brandDetails)
		if err != nil {
			return &model.Products{}, err
		}
		p.Brand = *newBrandDetails

		// create a new product
		newProductDetails, err := s.productsStore.PostProduct(p)
		if err != nil {
			return &model.Products{}, err
		}
		return newProductDetails, nil
	} else { // brand exsists
		p.Brand.Id = brandId
		newProductDetails, err := s.productsStore.PostProduct(p)
		if err != nil {
			return &model.Products{}, err
		}
		return newProductDetails, nil
	}
}

func (s *dataStore) DeleteProductById(id int) (int, error) {
	res, err := s.productsStore.DeleteById(id)
	if err != nil {
		return 0, err
	}
	return res, nil
}

func (s *dataStore) UpdateProductById(p model.Products) (*model.Products, error) {
	res, err := s.productsStore.UpdateProduct(p)
	if err != nil {
		return nil, err
	}
	if res !=0 {
		updatedDetails, err := s.GetById(p.Id)
		if err != nil {
			return nil, err
		}
		return updatedDetails, nil
	}
	return nil, errors.New("same data or id doesnt matched")
}
