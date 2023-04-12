package data

import (
	"encoding/json"
	"fmt"
	"io"
	"regexp"
	"time"

	"github.com/go-playground/validator/v10"
)

type Product struct {
	ID int `json:"id"`
	Name string `json:"name" validate:"required"`
	Description string `json:"description"`
	Price float32 `json:"price" validate:"gt=0"`
	SKU string `json:"sku" validate:"required,sku"`
	CreatedOn string `json:"-"`
	UpdatedOn string `json:"-"`
	DeletedOn string `json:"-"`
}

func (products *Product) FromJSON(reader io.Reader) error {
	err := json.NewDecoder(reader)
	return err.Decode(products)
}

func validateSKU(fl validator.FieldLevel) bool {
	// sku is of format abc-abc-abcd
	re := regexp.MustCompile(`[a-z]+-[a-z]+-[a-z]+`)
	matches := re.FindAllString(fl.Field().String(), -1)

	if len(matches) != 1 {
		return false
	}

	return true
}

func (products *Product) Validate() error {
	validate := validator.New()
	validate.RegisterValidation("sku", validateSKU)
	return validate.Struct(products)
}

type Products []*Product

func (products *Products) ToJSON(writer io.Writer) error {
	encoder := json.NewEncoder(writer)
	return encoder.Encode(products)
}

func GetProducts() Products {
	return productList
}

func UpdateProduct(id int, p *Product) error {
	_, position, error := findProduct(id)
	if error != nil {
		return error
	}

	p.ID = id
	productList[position] = p

	return nil
}

var ErrProductNotFound = fmt.Errorf("Product not found")

func findProduct(id int) (*Product, int, error) {
	for i, p :=  range productList {
		if p.ID == id {
			return p, i, nil
		}
	}

	return nil, -1, ErrProductNotFound
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