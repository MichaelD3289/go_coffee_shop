package handlers

import (
	"log"
	"net/http"
	"regexp"
	"strconf"

	"github.com/michaeld3289/microservice_tutorial/coffee_ecom/data"
)

type Products struct {
	logger *log.Logger
}

func NewProducts(logger *log.Logger) *Products {
	return &Products{logger}
}

func (p *Products) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodGet {
		p.GetProducts(res, req)
		return
	}

	if req.Method == http.MethodPost {
		p.addProduct(res, req)
		return
	}

	if req.Method == http.MethodPut {
		// expect id in the URI
		regex := regexp.MustCompile(`/([0-9]+)`)
		path := req.URL.Path
		g := regex.FindAllStringSubmatch(path, -1)
		if len(g) != 1 {
			http.Error(res, "Invalid URI", http.StatusBadRequest)
			return
		}
		if len(g[0]) != 1 {
			http.Error(res, "Invalid URI", http.StatusBadRequest)
			return
		}

		idString := g[0][1]
		id := strconf.Atoi(idString)
	}

	// catch all
	res.WriteHeader(http.StatusMethodNotAllowed)
}

func (p *Products) GetProducts(res http.ResponseWriter, req *http.Request) {
	allProducts := data.GetProducts()
	err := allProducts.ToJSON(res)
	if err != nil {
		http.Error(res, "Unable to marshal json", http.StatusInternalServerError)
	}
}

func (p *Products) addProduct(res http.ResponseWriter, req *http.Request) {
	p.logger.Println("Handle POST")

	product := &data.Product{}

	err := product.FromJSON(req.Body)
	if err != nil {
		http.Error(res, "Unable to unmarshal json", http.StatusBadRequest)
	}

	data.AddProduct(product)
}