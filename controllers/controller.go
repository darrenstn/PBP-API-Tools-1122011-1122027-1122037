package controllers

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	m "api_tools/models"
)

var ctx = context.Background()

// GetAllProducts..
func GetAllProducts(w http.ResponseWriter, r *http.Request) {
	db := connect()
	defer db.Close()

	rdb := connectRedis()

	// Check if products exist in cache
	cachedProducts, err := rdb.Get(ctx, "products").Result()
	if err == nil {
		// Products found in cache, return cached data
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(cachedProducts))
		return
	}

	query := "SELECT id, name, price, description FROM products"

	rows, err := db.Query(query)
	if err != nil {
		sendModifiedResponse(w, 400, "Query Error")
		return
	}

	var product m.Product
	var products []m.Product
	var result m.Products

	for rows.Next() {
		if err := rows.Scan(&product.ID, &product.ProductName, &product.ProductPrice, &product.ProductDescription); err != nil {
			sendModifiedResponse(w, 400, "Scan Error")
			return
		} else {
			products = append(products, product)
		}
	}

	w.Header().Set("Content-Type", "application/json")

	var response m.ProductsResponse
	response.Status = 200
	response.Message = "Success"
	result.Products = products
	response.Data = result

	// Cache the fetched products in Redis
	productsJSON, _ := json.Marshal(response)
	rdb.Set(ctx, "products", productsJSON, 1*time.Hour)

	json.NewEncoder(w).Encode(response)
}

func sendModifiedResponse(w http.ResponseWriter, stat int, msg string) {
	var response m.Response
	response.Status = stat
	response.Message = msg
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
