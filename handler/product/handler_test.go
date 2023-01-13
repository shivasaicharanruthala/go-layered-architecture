package product

import (
	"bytes"
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"io/ioutil"
	"layeredArchitecture/model"
	"layeredArchitecture/service"
	"net/http/httptest"
	"reflect"
	"testing"
)

type TestHandlers struct {
	id           int
	m *gomock.Call
	reqMethod    string
	path         string
	actualBody   *model.Products
	expectedBody string
	err          responseErr
	statusCode   int
}

func TestGetByIdHandler(t *testing.T) { // correct response // id not found  //brand id not found  //id type not found
	ctrl := gomock.NewController(t)
	ps := service.NewMockServiceInterface(ctrl)

	testCases := []TestHandlers{
		{1, ps.EXPECT().GetById(1).Return(&model.Products{Id: 1, Name: "biscuts", Brand: model.Brand{Id: 1, Name: "brand1"}}, nil),"GET", "/product?id=1", &model.Products{Id: 1, Name: "biscuts", Brand: model.Brand{Id: 1, Name: "brand1"}}, `{"productId":1,"productName":"biscuts","brand":{"brandId":1,"brandName":"brand1"}}`, responseErr{}, 200},
		{5, ps.EXPECT().GetById(5).Return(&model.Products{}, errors.New("Id not found")) ,"GET", "/product?id=5", &model.Products{}, `{"code":404,"message":"bad request"}`, responseErr{StatusCode: 404, Err: "Id not found"}, 404},
		{3, ps.EXPECT().GetById(3).Return(&model.Products{Id: 3, Name: "icecream", Brand: model.Brand{Id: 4, Name: ""}}, errors.New("brand id not found")),"GET", "/product?id=3", &model.Products{Id: 3, Name: "icecream", Brand: model.Brand{Id: 4, Name: ""}}, `{"code":404,"message":"bad request"}`, responseErr{StatusCode: 404, Err: "brand id not found"}, 404},
		{9, nil,"GET", "/product?id=abc", &model.Products{Id: 3, Name: "icecream", Brand: model.Brand{Id: 4, Name: ""}}, `{"code":400,"message":"invalid parameter"}`, responseErr{StatusCode: 400, Err: "brand id not found"}, 400},
	}

	for i, tc := range testCases {
		//ps.EXPECT().GetById(tc.id).Return(tc.actualBody, tc.err)
		serv := New(ps) //check

		w := httptest.NewRecorder()
		r := httptest.NewRequest(tc.reqMethod, tc.path, nil)

		serv.GetById(w, r)
		result := w.Result()
		res, err := ioutil.ReadAll(result.Body)

		if !reflect.DeepEqual(string(res), tc.expectedBody) {
			t.Errorf("Test failed %v got %v but wanted %v", i, string(res), tc.expectedBody)
		}

		if !reflect.DeepEqual(w.Code, tc.statusCode) {
			t.Errorf("Test failed %v got %v but wanted %v", i, w.Code, tc.statusCode)
		}

		if !reflect.DeepEqual(err, nil) {
			t.Errorf("test failed %v, this shoulnt throw an error but %v", i, err)
		}
	}
}

type TestCreate struct {
	// product service interface
	productDetailsInp model.Products
	productDetailsOut *model.Products
	prodError         error

	// handler
	path        string
	method      string
	body        []byte
	expextedOut string
	statusCode  int
}

func TestProductHandler_Create(t *testing.T) {
	testCases := []TestCreate{
		{model.Products{Name: "hide and seek", Brand: model.Brand{Name: "brand1"}}, &model.Products{Id: 3, Name: "hide and seek", Brand: model.Brand{Id: 1, Name: "brand1"}}, nil, "/product", "POST", []byte(`{"productName": "hide and seek", "brand":{"brandName":"brand1"}}`), `{"productId":3,"productName":"hide and seek","brand":{"brandId":1,"brandName":"brand1"}}`, 201},
		{model.Products{Name: "hide and seek"}, &model.Products{}, errors.New("brand is missing"), "/product", "POST", []byte(`{"productName": "hide and seek"}`), `{"code":400,"message":"failed to create new product"}`, 404},
		//{model.Products{Name: "dairymilk", Brand: model.Brand{Name: "brand2"}}, &model.Products{Id: 3, Name: "dairymilk", Brand: model.Brand{Id: 2, Name: "brand2"}}, nil, "/product", "POST", []byte(`{"productName:  "dairymilk"}`), "invalid body", 400},
	}
	ctrl := gomock.NewController(t)
	ps := service.NewMockServiceInterface(ctrl)

	for i, tc := range testCases {
		ps.EXPECT().PostProduct(tc.productDetailsInp).Return(tc.productDetailsOut, tc.prodError)
		serv := New(ps)

		w := httptest.NewRecorder()
		r := httptest.NewRequest(tc.method, tc.path, bytes.NewBuffer(tc.body))

		serv.Create(w, r)
		result := w.Result()
		res, err := ioutil.ReadAll(result.Body)

		if !reflect.DeepEqual(string(res), tc.expextedOut) {
			t.Errorf("Test failed %v got %v but wanted %v", i, string(res), tc.expextedOut)
		}

		if !reflect.DeepEqual(w.Code, tc.statusCode) {
			t.Errorf("Test failed %v got %v but wanted %v", i, w.Code, tc.statusCode)
		}

		if !reflect.DeepEqual(err, nil) {
			t.Errorf("test failed %v, this shoulnt throw an error but %v", i, err)
		}
	}
}

type TestDelete struct{
	prodID int
	rowsAffected int
	err error
	m *gomock.Call
	urlParams string
	expextedOut string
	statusCode int
}
func TestProductHandler_Delete(t *testing.T) {
	ctrl := gomock.NewController(t)
	ps := service.NewMockServiceInterface(ctrl)

	testCases := []TestDelete{
		{4,1,nil,ps.EXPECT().DeleteProductById(4).Return(1, nil),"4", `{}`, 201},
		{19,0,nil,ps.EXPECT().DeleteProductById(19).Return(0, nil),"19", `{"code":500,"message":"Error in deleting ot id doesnt exsists"}`, 500},
		{20,0,nil,nil,"abc", `{"code":400,"message":"not a proper ID"}`, 404},
	}

	for i,tc := range testCases{
		serv := New(ps)

		w := httptest.NewRecorder()
		r := httptest.NewRequest("DELETE", "/product/"+tc.urlParams, nil)

		vars := map[string]string{
			"id": tc.urlParams,
		}
		r = mux.SetURLVars(r, vars)

		serv.Delete(w,r)
		result := w.Result()
		res, err := ioutil.ReadAll(result.Body)

		if !reflect.DeepEqual(string(res), tc.expextedOut) {
			t.Errorf("Test failed %v got %v but wanted %v", i, string(res), tc.expextedOut)
		}

		if !reflect.DeepEqual(w.Code, tc.statusCode) {
			t.Errorf("Test failed %v got %v but wanted %v", i, w.Code, tc.statusCode)
		}

		if !reflect.DeepEqual(err, nil) {
			t.Errorf("test failed %v, this shoulnt throw an error but %v", i, err)
		}
	}
}

type TestUpdate struct {
	prodId int
	urlParams string
	prodDetails model.Products
	body []byte
	expextedOut string
	statusCode int
	output *model.Products
	err error
	err2 error

}
func TestProductHandler_Update(t *testing.T) {
	ctrl := gomock.NewController(t)
	ps := service.NewMockServiceInterface(ctrl)

	testCases := []TestUpdate{
		{1,"1", model.Products{Name: "goodday", Brand: model.Brand{Id: 3}}, []byte(`{"productName": "goodday", "brand": {"brandId": 3}}`), `{"productId":1,"productName":"goodday","brand":{"brandId":3,"brandName":""}}`, 201, &model.Products{Id: 1, Name: "goodday", Brand: model.Brand{Id: 3}},nil,nil},
		{1,"abc", model.Products{Name: "goodday", Brand: model.Brand{Id: 3}}, []byte(`{"productName": "goodday", "brand": {"brandId": 3}}`), `{"code":404,"message":"this type ID not accepted"}`, 404, &model.Products{},errors.New("err"), nil},
		//{1,"2", model.Products{Name: "goodday", Brand: model.Brand{Id: 3}}, []byte(`{"productName": "goodday", "brand": {"brandId": 3}}`), `{"code":500,"message":"product not updated"}`, 500, &model.Products{},nil, errors.New("err")},
	}

	for i,tc := range testCases {
		if tc.err ==nil {
			tc.prodDetails.Id = tc.prodId
			ps.EXPECT().UpdateProductById(tc.prodDetails).Return(tc.output, tc.err2)
		}
		serv := New(ps)
		w := httptest.NewRecorder()
		r := httptest.NewRequest("UPDATE", "/product/"+tc.urlParams, bytes.NewBuffer(tc.body))

		vars := map[string]string{
			"id": tc.urlParams,
		}
		r = mux.SetURLVars(r, vars)
		serv.Update(w,r)
		result := w.Result()
		res, err := ioutil.ReadAll(result.Body)
		//fmt.Println("--------",res)
		if !reflect.DeepEqual(string(res), tc.expextedOut) {
			t.Errorf("Test failed %v got %v but wanted %v", i, string(res), tc.expextedOut)
		}

		if !reflect.DeepEqual(w.Code, tc.statusCode) {
			t.Errorf("Test failed %v got %v but wanted %v", i, w.Code, tc.statusCode)
		}

		if !reflect.DeepEqual(err, nil) {
			t.Errorf("test failed %v, this shoulnt throw an error but %v", i, err)
		}

	}
}


