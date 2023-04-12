package handlers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/MichaelD3289/go_coffee_shop/data"
	"github.com/gorilla/mux"
)

type Products struct {
	logger *log.Logger
}

func NewProducts(logger *log.Logger) *Products {
	return &Products{logger}
}

func (p *Products) GetProducts(res http.ResponseWriter, req *http.Request) {
	allProducts := data.GetProducts()
	err := allProducts.ToJSON(res)
	if err != nil {
		http.Error(res, "Unable to marshal json", http.StatusInternalServerError)
	}
}

func (p *Products) AddProduct(res http.ResponseWriter, req *http.Request) {
	p.logger.Println("Handle POST")

	product := req.Context().Value(KeyProduct{}).(data.Product)

	data.AddProduct(&product)
}

func (p *Products) UpdateProducts(res http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(res, "Unable to convert id", http.StatusBadRequest)
		return
	}

	p.logger.Println("Handle PUT", id)

	product := req.Context().Value(KeyProduct{}).(data.Product)

	err = data.UpdateProduct(id, &product)
	if err == data.ErrProductNotFound {
		http.Error(res, "Product not found", http.StatusNotFound)
		return
	}

	if err != nil {
		http.Error(res, "Product not found", http.StatusInternalServerError)
		return
	}
}

type KeyProduct struct {}

func (p Products) MiddlewareProductValidation(next http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		product := data.Product{}

		err := product.FromJSON(req.Body)
		if err != nil {
			p.logger.Println("[ERROR] deserializing product", err)
			http.Error(res, "Error Reading Product", http.StatusBadRequest)
			return
		}

		// Validate product
		err = product.Validate()
		if err != nil {
			p.logger.Println("[ERROR] validating product", err)
			http.Error(
				res, 
				fmt.Sprintf("Unable validation product: %s", err),
				http.StatusBadRequest,
			)
			return
		}

		ctx := context.WithValue(req.Context(), KeyProduct{}, product)

		reqCopy := req.WithContext(ctx)

		next.ServeHTTP(res, reqCopy)
	})
}