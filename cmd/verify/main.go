package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

const baseURL = "http://localhost:8080/api/v1"

func main() {
	// 1. Login
	fmt.Println("1. Logging in...")
	token, err := login("superadmin@mail.com", "admin123")
	if err != nil {
		fmt.Printf("Login failing: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Login Success! Token: %s...\n", token[:20])

	// 2. Get Products (to find ID)
	fmt.Println("\n2. Getting Products...")
	productID, err := getFirstProductID()
	if err != nil {
		fmt.Printf("Get Products failed: %v\n", err)
		// Don't exit, maybe we can't test cart but login worked.
		os.Exit(1)
	}
	fmt.Printf("Found Product ID: %d\n", productID)

	// 3. Add to Cart
	fmt.Println("\n3. Adding to Cart...")
	err = addToCart(token, productID, 2)
	if err != nil {
		fmt.Printf("Add to Cart failed: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("Add to Cart Success!")

	// 4. Get Cart
	fmt.Println("\n4. Verifying Cart...")
	count, err := getCartCount(token)
	if err != nil {
		fmt.Printf("Get Cart failed: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Cart has %d items\n", count)

	// 5. Create Order
	fmt.Println("\n5. Creating Order...")
	err = createOrder(token, productID)
	if err != nil {
		fmt.Printf("Create Order failed: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("Create Order Success!")
}

func login(email, password string) (string, error) {
	data := map[string]string{"email": email, "password": password}
	resp, err := post("/login", data, "")
	if err != nil {
		return "", err
	}
	var res struct {
		Data struct {
			AccessToken string `json:"access_token"`
		} `json:"data"`
	}
	if err := json.Unmarshal(resp, &res); err != nil {
		return "", err
	}
	if res.Data.AccessToken == "" {
		return "", fmt.Errorf("no token in response: %s", string(resp))
	}
	return res.Data.AccessToken, nil
}

func getFirstProductID() (int64, error) {
	resp, err := get("/products", "")
	if err != nil {
		return 0, err
	}
	var res struct {
		Data []struct {
			ID int64 `json:"id"`
		} `json:"data"`
	}
	if err := json.Unmarshal(resp, &res); err != nil {
		return 0, err
	}
	if len(res.Data) == 0 {
		return 0, fmt.Errorf("no products found")
	}
	return res.Data[0].ID, nil
}

func addToCart(token string, productID int64, qty int) error {
	data := map[string]interface{}{"product_id": productID, "quantity": qty}
	_, err := post("/carts", data, token)
	return err
}

func getCartCount(token string) (int, error) {
	resp, err := get("/carts", token)
	if err != nil {
		return 0, err
	}
	var res struct {
		Data []interface{} `json:"data"`
	}
	if err := json.Unmarshal(resp, &res); err != nil {
		return 0, err
	}
	return len(res.Data), nil
}

func createOrder(token string, productID int64) error {
	// Payload for CreateOrderRequest
	data := map[string]interface{}{
		"buyer_id":      1,
		"order_date":    "2025-01-01",
		"total_amount":  50000,
		"shipping_type": "JNE Regular",
		"payment_type":  "Bank Transfer",
		"order_time":    "12:00:00",
		"remarks":       "Test Order",
		"order_details": []map[string]interface{}{
			{
				"product_id": productID,
				"quantity":   2,
			},
		},
	}
	_, err := post("/orders", data, token)
	return err
}

func post(endpoint string, body interface{}, token string) ([]byte, error) {
	jsonBody, _ := json.Marshal(body)
	req, _ := http.NewRequest("POST", baseURL+endpoint, bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	b, _ := ioutil.ReadAll(resp.Body)
	if resp.StatusCode >= 400 {
		return b, fmt.Errorf("statue %d: %s", resp.StatusCode, string(b))
	}
	return b, nil
}

func get(endpoint string, token string) ([]byte, error) {
	req, _ := http.NewRequest("GET", baseURL+endpoint, nil)
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	b, _ := ioutil.ReadAll(resp.Body)
	if resp.StatusCode >= 400 {
		return b, fmt.Errorf("statue %d: %s", resp.StatusCode, string(b))
	}
	return b, nil
}
