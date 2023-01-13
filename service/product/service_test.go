package product

import (
	"errors"
	"github.com/golang/mock/gomock"
	"layeredArchitecture/model"
	"layeredArchitecture/store"
	"reflect"
	"testing"
)

type TestProducts struct {
	id                int
	productDataOutput *model.Products
	brandDataOutput   *model.Brand
	expectedOutput    *model.Products
	productErr        error
	brandError        error
	expectedError     error
}

func TestStore_GetById(t *testing.T) {
	ctrl := gomock.NewController(t)

	ps := store.NewMockProductsInterface(ctrl)
	bs := store.NewMockBrandInterface(ctrl)

	testcases := []TestProducts{
		{1, &model.Products{1, "biscuts", model.Brand{1, ""}}, &model.Brand{1, "brand1"}, &model.Products{1, "biscuts", model.Brand{1, "brand1"}}, nil, nil, error(nil)},
		{3, &model.Products{3, "icecream", model.Brand{5, ""}}, &model.Brand{}, &model.Products{3, "icecream", model.Brand{5, ""}}, nil, errors.New("Brand Id not found"), errors.New("Brand Id not found")},
		{5, &model.Products{}, &model.Brand{}, &model.Products{}, errors.New("Id not found"), errors.New("Id not found"), errors.New("product Id not found")},
	}

	for i, tc := range testcases {
		ps.EXPECT().GetByIdProducts(tc.id).Return(tc.productDataOutput, tc.productErr)
		if tc.productErr == nil {
			bs.EXPECT().GetByIdBrand(tc.productDataOutput.Brand.Id).Return(tc.brandDataOutput, tc.brandError)
		}
		prodService := New(ps, bs)
		data, err := prodService.GetById(tc.id)

		if !reflect.DeepEqual(err, tc.expectedError) {
			t.Errorf("Test Failed %v, Got %v wanted %v", i, err, tc.expectedError)
		}

		if !reflect.DeepEqual(data, tc.expectedOutput) {
			t.Errorf("Test Failed %v got %v but wanted %v", i, data, tc.expectedOutput)
		}
	}

}

type TestPostProduct struct {
	productInp      model.Products
	brandExsistsIdOut int
	brandExsistsErr error
	newBrandOut     *model.Brand
	postBrandOut    *model.Brand
	postBrandErr    error
	postProductOut  *model.Products
	postProductErr  error
	expectedOut     *model.Products
	expectedErr     error
}

func TestStore_PostProduct(t *testing.T) {
	ctrl := gomock.NewController(t)

	ps := store.NewMockProductsInterface(ctrl)
	bs := store.NewMockBrandInterface(ctrl)

	testCases := []TestPostProduct{
		{model.Products{Name: "5star", Brand: model.Brand{Name: "brand1"}}, 1, nil, nil, &model.Brand{}, nil, &model.Products{Id: 1, Name: "5star", Brand: model.Brand{Id: 1, Name: "brand1"}}, nil, &model.Products{Id: 1, Name: "5star", Brand: model.Brand{Id: 1,Name: "brand1"}}, nil},
		{model.Products{Name: "choclate", Brand: model.Brand{Name: "brand6"}}, 0, nil, nil, &model.Brand{Id: 12, Name: "brand6"}, nil, &model.Products{Id: 16, Name: "choclate", Brand: model.Brand{Id: 12, Name: "brand6"}}, nil, &model.Products{Id: 16, Name: "choclate", Brand: model.Brand{Id: 12,Name: "brand6"}}, nil},
		{model.Products{Name: "choclate", Brand: model.Brand{Name: "brand6"}}, 0, nil, nil, &model.Brand{}, errors.New("sql syntax error"), &model.Products{Id: 16, Name: "choclate", Brand: model.Brand{Id: 12, Name: "brand6"}}, nil, &model.Products{}, errors.New("sql syntax error")},
		{model.Products{Name: "choclate", Brand: model.Brand{Name: "brand6"}}, 0, nil, nil, &model.Brand{Id: 12, Name: "brand6"}, nil, &model.Products{}, errors.New("data store error in creating product"), &model.Products{}, errors.New("data store error in creating product")},
		{model.Products{Name: "choclate", Brand: model.Brand{Name: "brand6"}}, 1, nil, nil, &model.Brand{}, nil, &model.Products{}, errors.New("data store error in creating product"), &model.Products{}, errors.New("data store error in creating product")},
		{model.Products{Name: "choclate", Brand: model.Brand{Name: "brand6"}}, 0, errors.New("data store layer error in checking brand exists"), nil, &model.Brand{}, nil, &model.Products{}, nil, &model.Products{}, errors.New("data store layer error in checking brand exists")},
	}

	for i, tc := range testCases {
		bs.EXPECT().CheckBrandExsists(tc.productInp.Brand.Name).Return(tc.brandExsistsIdOut, tc.brandExsistsErr)
		if reflect.DeepEqual(tc.brandExsistsErr, nil){
			if tc.brandExsistsIdOut == 0 {
				// create brand
				newBrandDetails := model.Brand{Name: tc.productInp.Brand.Name}
				bs.EXPECT().PostBrand(newBrandDetails).Return(tc.postBrandOut, tc.postBrandErr)
				if reflect.DeepEqual(tc.postBrandErr, nil){
					// create product
					tc.productInp.Brand = *tc.postBrandOut
					ps.EXPECT().PostProduct(tc.productInp).Return(tc.postProductOut, tc.postProductErr)
				}

			} else {
				tc.productInp.Brand.Id = tc.brandExsistsIdOut
				ps.EXPECT().PostProduct(tc.productInp).Return(tc.postProductOut, tc.postProductErr)
			}
		}

		prodService := New(ps, bs)
		output, err := prodService.PostProduct(tc.productInp)
		if !reflect.DeepEqual(output, tc.expectedOut) {
			t.Errorf("Test Failed %v got %v but wanted %v", i, output, tc.expectedOut)
		}

		if !reflect.DeepEqual(err, tc.expectedErr) {
			t.Errorf("Test Failed %v, Got %v wanted %v", i, err, tc.expectedErr)
		}
	}
}

type TestDelete struct{
	productID int
	rowsAffected int
	err error
	expectedOut int
	expectedErr error
}
func TestStore_DeleteProductById(t *testing.T) {
	ctrl := gomock.NewController(t)
	ps := store.NewMockProductsInterface(ctrl)
	bs := store.NewMockBrandInterface(ctrl)

	testCases := []TestDelete{
		{1, 1, nil, 1, nil},
		{110, 0, nil, 0, nil},
		{2, 0, errors.New("data store layer error"), 0, errors.New("data store layer error")},
	}

	for i,tc := range testCases {
		ps.EXPECT().DeleteById(tc.productID).Return(tc.rowsAffected, tc.err)
		prodService := New(ps, bs)
		output, err := prodService.DeleteProductById(tc.productID)

		if !reflect.DeepEqual(output, tc.expectedOut) {
			t.Errorf("Test Failed %v got %v but wanted %v", i, output, tc.expectedOut)
		}

		if !reflect.DeepEqual(err, tc.expectedErr) {
			t.Errorf("Test Failed %v, Got %v wanted %v", i, err, tc.expectedErr)
		}
	}
}

type TestUpdateProduct struct {
	prodDetails model.Products
	rowsAffected int
	err error

}
//func TestStore_UpdateProductById(t *testing.T) {
//	testCases := []TestUpdateProduct{
//
//	}
//	ctrl := gomock.NewController(t)
//	ps := products.NewMockProductsInterface(ctrl)
//	bs := brand.NewMockBrandInterface(ctrl)
//
//	for i,tc := range testCases {
//		ps.EXPECT().UpdateProduct(tc.prodDetails).Return(tc.rowsAffected,tc.err)
//		prodService := New(ps, bs)
//		prodService.GetById(tc.prodDetails.Id)
//	}
//}