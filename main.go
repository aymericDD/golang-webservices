package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type Product struct {
	ProductId		int 	`json:"productId"`
	Manufacturer	string 	`json:"manufacturer"`
	Sku				string 	`json:"sku"`
	Upc				string 	`json:"upc"`
	PricePerUnit	string 	`json:"pricePerUnit"`
	QuantityOnland	int		`json:"quantityOnHand"`
	ProductName		string 	`json:"productName"`
}

type ProductList []Product

var productList ProductList

func (products *ProductList) getNextId() int {
	var lastId int
	for _, product := range *products {
		if product.ProductId < lastId {
			continue
		}
		lastId = product.ProductId
	}

	lastId++

	return lastId
}

func findProductByID(productId int) (*Product, int) {
	for i, product := range productList {
		if product.ProductId == productId {
			return &product, i
		}
	}
	return nil, 0
}

func init() {
	productJson := `[
		{
			"productId": 1,
			"manufacturer": "Johns-Jenkins",
			"sku": "p76750ufds8",
			"pricePerUnit": "497.58",
			"quantityOnHand": 9876,
			"productName": "sticky note"
		},
		{
			"productId": 2,
			"manufacturer": "Mac-Cready",
			"sku": "x987667YU",
			"pricePerUnit": "800.58",
			"quantityOnHand": 99999876,
			"productName": "a pen"
		}
	]`
	err := json.Unmarshal([]byte(productJson), &productList)
	if err != nil {
		log.Fatal(err)
	}
}

func productHandler(w http.ResponseWriter, r *http.Request) {
	urlPathSecgments := strings.Split(r.URL.Path, "products/")
	productId, err := strconv.Atoi(urlPathSecgments[len(urlPathSecgments)-1])
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	product, listItemIndex := findProductByID(productId)
	if product == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	

	switch r.Method {
	case http.MethodGet:
		productJson, err := json.Marshal(product)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(productJson)
	case http.MethodPut:
		var updatedProduct Product
		bodyBytes, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		err = json.Unmarshal(bodyBytes, &updatedProduct)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if updatedProduct.ProductId != productId {
			w.WriteHeader(http.StatusBadGateway)
			return
		}
		product = &updatedProduct
		productList[listItemIndex] =* product
		w.WriteHeader(http.StatusOK)
		return
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func productsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		productJson, err := json.Marshal(productList)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(productJson)
	case http.MethodPost:
		var newProduct Product
		jsonBody, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Could not decode body request to json."))
			return
		}
		err = json.Unmarshal(jsonBody, &newProduct)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Could not deserialize json request to product object."))
			return
		}
		if newProduct.ProductId <= 0 {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(fmt.Sprintf("The %d productID must be an integer and greater than 0.", newProduct.ProductId)))
			return
		}
		w.WriteHeader(http.StatusCreated)
		nextId := productList.getNextId()
		newProduct.ProductId = nextId

		productList = append(productList, newProduct)
		newProductJson, err := json.Marshal(newProduct)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		w.Write(newProductJson)
	}
}

func main() {
	http.HandleFunc("/products", productsHandler)
	http.HandleFunc("/products/", productHandler)
	http.ListenAndServe(":5000", nil)
}
