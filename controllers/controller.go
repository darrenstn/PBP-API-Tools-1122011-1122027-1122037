package controllers

import (
	"encoding/json"
	"net/http"

	m "api_tools/models"
)

// GetAllRooms..
func GetAllProducts(w http.ResponseWriter, r *http.Request) {
	db := connect()
	defer db.Close()

	query := "SELECT * FROM products"

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
	json.NewEncoder(w).Encode(response)
}

func sendModifiedResponse(w http.ResponseWriter, stat int, msg string) {
	var response m.Response
	response.Status = stat
	response.Message = msg
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
