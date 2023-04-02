package data

import (
	"encoding/json"
	"io"
	"time"
)

type Product struct {
	ID int `json:"id"`
	Name string `json:"name"`
	Description string `json:"description"`
	Price float32 `json:"price"`
	SKU string `json:"sku"`
	CreatedOn string `json:"-"`
	UpdatedOn string `json:"-"`
	DeletedOn string `json:"-"`
}

func (products *Product) FromJSON(reader io.Reader) error {
	err := json.NewDecoder(reader)
	return err.Decode(products)
}

type Products []*Product

func (products *Products) ToJSON(writer io.Writer) error {
	encoder := json.NewEncoder(writer)
	return encoder.Encode(products)
}

func GetProducts() Products {
	return productList
}

func AddProduct(product *Product) {
	product.ID = getNextID()
	productList = append(productList, product)
}

func getNextID() int {
	lp := productList[len(productList) - 1]
	return lp.ID + 1
}

var productList = []*Product{
	{
		ID: 1,
		Name: "Latte",
		Description: "Frothy milky coffee",
		Price: 2.45,
		SKU: "abc323",
		CreatedOn: time.Now().UTC().String(),
		UpdatedOn: time.Now().UTC().String(),
	},
	{
		ID: 2,
		Name: "Espresso",
		Description: "Short and  strong coffee without milk",
		Price: 1.99,
		SKU: "jkl456",
		CreatedOn: time.Now().UTC().String(),
		UpdatedOn: time.Now().UTC().String(),
	},
}