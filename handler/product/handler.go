package product

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"layeredArchitecture/model"
	"layeredArchitecture/service"
	"net/http"
	"strconv"
)

type productHandler struct {
	productService service.ServiceInterface
}

func New(pHandler service.ServiceInterface) *productHandler {
	return &productHandler{productService: pHandler}
}


type responseErr struct {
	StatusCode int `json:"code"`
	Err string  `json:"message"`
}

type result struct {
	respBody string `json:"respBody"`
}

func (pHandle *productHandler) GetById(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	i, err := strconv.Atoi(id)
	if err != nil {

		resp := responseErr{StatusCode: 400, Err: "invalid parameter"}
		finalResp,_ := json.Marshal(resp)

		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write(finalResp)
		return
	}
	productDetails, err := pHandle.productService.GetById(i)
	fmt.Println("handler:  ",productDetails, err)
	if err != nil {
		resp := responseErr{StatusCode: 404, Err: "bad request"}
		finalResp,_ := json.Marshal(resp)

		w.WriteHeader(404)
		_, _ = w.Write(finalResp)
		return
	}
	detailsJson, _ := json.Marshal(productDetails)
	w.WriteHeader(200)
	_, _ = w.Write(detailsJson)
}

func (pHandle *productHandler) Create(w http.ResponseWriter, r *http.Request) {
	var product model.Products
	body, _ := ioutil.ReadAll(r.Body)

	err := json.Unmarshal(body, &product)
	if err != nil {
		resp := responseErr{StatusCode: 400, Err: "wrong body"}
		finalResp,_ := json.Marshal(resp)
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write(finalResp)
		return
	}
	resp, err := pHandle.productService.PostProduct(product)
	if err != nil {
		resp := responseErr{StatusCode: 400, Err: "failed to create new product"}
		finalResp,_ := json.Marshal(resp)
		w.WriteHeader(404)
		_, _ = w.Write(finalResp)
		return
	}

	respBody, _ := json.Marshal(resp)
	w.WriteHeader(201)
	_, _ = w.Write(respBody)
}

func (pHandle *productHandler) Delete(w http.ResponseWriter, r *http.Request) {
	prodId := mux.Vars(r)["id"]
	id, err := strconv.Atoi(prodId)
	if err != nil {
		resp := responseErr{StatusCode: 400, Err: "not a proper ID"}
		finalResp,_ := json.Marshal(resp)
		w.WriteHeader(404)
		_, _ = w.Write(finalResp)
		return
	}

	res, err := pHandle.productService.DeleteProductById(id)
	if err != nil || res == 0{
		resp := responseErr{StatusCode: 500, Err: "Error in deleting ot id doesnt exsists"}
		finalResp,_ := json.Marshal(resp)
		w.WriteHeader(500)
		_, _ = w.Write(finalResp)
		return
	}
	 resp := result{respBody: "product deleted"}
	 finalResp,_ := json.Marshal(resp)

	w.WriteHeader(201)
	_, _ = w.Write(finalResp)
}

func (pHandle *productHandler) Update(w http.ResponseWriter, r *http.Request) {
	// validate url params
	prodId := mux.Vars(r)["id"]
	id, err := strconv.Atoi(prodId)
	if err != nil {
		resp := responseErr{StatusCode: 404, Err: "this type ID not accepted"}
		finalResp,_ := json.Marshal(resp)

		w.WriteHeader(404)
		_, _ = w.Write(finalResp)
		return
	}

	//validate body
	var productDetails model.Products
	body := r.Body
	err2 := json.NewDecoder(body).Decode(&productDetails)
	if err2 != nil {
		resp := responseErr{StatusCode: 500, Err: "body is corupted"}
		finalResp,_ := json.Marshal(resp)

		w.WriteHeader(404)
		_, _ = w.Write(finalResp)
		return
	}

	//	send data
	productDetails.Id = id

	res, err := pHandle.productService.UpdateProductById(productDetails)
	if err != nil {
		resp := responseErr{StatusCode: 500, Err: "product not updated"}
		finalResp,_ := json.Marshal(resp)

		w.WriteHeader(500)
		_,_ = w.Write(finalResp)
		return
	}
	updatedDetails, _ :=json.Marshal(res)
	w.WriteHeader(201)
	_,_ = w.Write(updatedDetails)
}
