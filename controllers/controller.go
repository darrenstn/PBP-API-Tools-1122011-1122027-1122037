package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	m "PBP-API-Tools-1122011-1122027-1122037/models"
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

func GetEmailWithContent(status string) m.Recipients {
	db := connect()
	defer db.Close()

	var result m.Recipients

	var product m.Product
	var content string
	var tempEmail string
	var email string

	content = "Rekomendasi Produk buat Pelanggan Berstatus " + status + "\n\n"

	// Channel to receive the content from the products query
	productCh := make(chan string)
	// Channel to receive the email addresses from the users query
	emailCh := make(chan string)
	// Channel to signal completion of both goroutines
	doneCh := make(chan struct{})

	// Execute the products query concurrently
	go func() {
		defer close(productCh)

		query := "SELECT id, name, price, description FROM products"
		if status == "SILVER" {
			query += " WHERE price BETWEEN 0 AND 2 LIMIT 5"
		} else if status == "GOLD" {
			query += " WHERE price BETWEEN 10 AND 20 LIMIT 5"
		}

		rows, err := db.Query(query)
		if err != nil {
			log.Println("Query Error:", err)
			return
		}
		defer rows.Close()

		for rows.Next() {
			if err := rows.Scan(&product.ID, &product.ProductName, &product.ProductPrice, &product.ProductDescription); err != nil {
				log.Println("Scan Error:", err)
				return
			}
			content += fmt.Sprintf("Product ID: %d\nName: %s\nPrice: %.2f\nProduct Description:%s\n\n", product.ID, product.ProductName, product.ProductPrice, product.ProductDescription)
		}

		// Send the content through the channel
		productCh <- content
	}()
	// Execute the users query concurrently
	go func() {
		defer close(emailCh)

		query := "SELECT email FROM users WHERE status='" + status + "'"
		rows, err := db.Query(query)
		if err != nil {
			log.Println("Query Error:", err)
			return
		}
		defer rows.Close()

		for rows.Next() {
			if err := rows.Scan(&tempEmail); err != nil {
				log.Println("Scan Error:", err)
				return
			}
			email += tempEmail + ", "
		}

		if len(email) >= 2 {
			email = email[:len(email)-2]
		} else {
			// Handle case where string length is less than 2
			email = ""
		}

		// Send the email addresses through the channel
		emailCh <- email
	}()

	// Wait for both goroutines to complete
	go func() {
		defer close(doneCh)
		<-productCh
		<-emailCh
	}()

	// Wait for both goroutines to complete
	<-doneCh

	// Set the content and email addresses in the result
	result.Content = content
	result.Email = email

	return result
}
