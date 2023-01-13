package main

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	handler2 "layeredArchitecture/handler/product"
	"layeredArchitecture/service/product"
	"layeredArchitecture/store/brand"
	"layeredArchitecture/store/products"
	"net/http"
)

func main() {
	db, err := sql.Open("mysql", "root:Shiva_9121@(127.0.0.1)/catalog")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	// store layer
	prod := products.New(db)
	brand := brand.New(db)

	// service layer
	service := product.New(prod, brand)

	// delivery layer
	handler := handler2.New(service)

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/product", handler.GetById).Methods("GET")
	router.HandleFunc("/product", handler.Create).Methods("POST")
	router.HandleFunc("/product/{id}", handler.Delete).Methods("DELETE")
	router.HandleFunc("/product/{id}", handler.Update).Methods("PUT")

	_ = http.ListenAndServe(":8080", router)
}
