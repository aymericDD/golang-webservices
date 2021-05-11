package main

import (
	"net/http"

	"github.com/aymericdd/golan-webservices/database"
	"github.com/aymericdd/golan-webservices/product"
	"github.com/aymericdd/golan-webservices/receipt"
	_ "github.com/go-sql-driver/mysql"
)

const apiBasePath = "/api"

func main() {
	database.SetupDatabase()
	receipt.SetupRoutes(apiBasePath)
	product.SetupRoutes(apiBasePath)
	http.ListenAndServe(":5000", nil)
}
