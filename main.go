package main

import (
	"net/http"

	"github.com/aymericdd/golan-webservices/product"
)

const apiBasePath = "/api"

func main() {
	product.SetupRoutes(apiBasePath)
	http.ListenAndServe(":5000", nil)
}
