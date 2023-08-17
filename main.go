package main

import (
	"bufio"
	"database/sql"
	"fmt"
	"log"
	"strconv"
	"time"

	"golang.org/x/crypto/bcrypt"
)

func RegisterUser(db *sql.DB, scanner *bufio.Scanner) {
	fmt.Println("Register")
	fmt.Println("--------")

	fmt.Print("Enter username: ")
	scanner.Scan()
	username := scanner.Text()

	fmt.Print("Enter email: ")
	scanner.Scan()
	email := scanner.Text()

	fmt.Print("Enter password: ")
	scanner.Scan()
	password := scanner.Text()

	// Meng-hash password sebelum menyimpannya
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal(err)
	}

	createdAt := time.Now()

	// Insert data ke tabel Users
	_, err = db.Exec("INSERT INTO Users (username, email, password, created_at) VALUES (?, ?, ?, ?)",
		username, email, hashedPassword, createdAt)
	if err != nil {
		log.Fatal(err)
	}

	// Insert data ke tabel Users
	_, err = db.Exec("INSERT INTO Users (username, email, password, created_at) VALUES (?, ?, ?, ?)",
		username, email, hashedPassword, createdAt)
	if err != nil {
		log.Fatal(err)
	}

	// Ambil user_id yang baru saja di-generate
	var userID int
	err = db.QueryRow("SELECT LAST_INSERT_ID()").Scan(&userID)
	if err != nil {
		log.Fatal(err)
	}

	// Insert data ke tabel user_profiles dan menghubungkannya dengan user_id
	_, err = db.Exec("INSERT INTO user_profiles (user_id, full_name) VALUES (?, ?)",
		userID, "Default Full Name")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("User registered successfully!")
}

func LoginUser(db *sql.DB, scanner *bufio.Scanner) {
	fmt.Println("Login")
	fmt.Println("-----")

	fmt.Print("Enter username: ")
	scanner.Scan()
	username := scanner.Text()

	fmt.Print("Enter password: ")
	scanner.Scan()
	password := scanner.Text()

	var storedPassword []byte
	var userID int
	err := db.QueryRow("SELECT user_id, password FROM Users WHERE username = ?", username).Scan(&userID, &storedPassword)
	if err != nil {
		log.Fatal(err)
	}

	// Memeriksa apakah password cocok dengan hashed password yang disimpan
	if err := bcrypt.CompareHashAndPassword(storedPassword, []byte(password)); err != nil {
		log.Fatal("Login failed: Incorrect username or password")
	}

	fmt.Println("Login successful!")
}

func ListLaptops(db *sql.DB) {
	fmt.Println("List Laptops")
	fmt.Println("------------")

	rows, err := db.Query("SELECT laptop_id, brand, model, price, stock_quantity FROM Laptops")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	fmt.Println("Laptop list:")
	for rows.Next() {
		var laptopID int
		var brand, model string
		var price float64
		var stockQuantity int

		err := rows.Scan(&laptopID, &brand, &model, &price, &stockQuantity)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("ID: %d, Brand: %s, Model: %s, Price: %.2f, Stock: %d\n", laptopID, brand, model, price, stockQuantity)
	}
}

func BuyLaptop(db *sql.DB, scanner *bufio.Scanner) {
	fmt.Println("Buy Laptop")
	fmt.Println("----------")

	fmt.Print("Enter user ID: ")
	scanner.Scan()
	userIDStr := scanner.Text()
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		log.Fatal(err)
	}

	ListLaptops(db)

	fmt.Print("Enter laptop ID to buy: ")
	scanner.Scan()
	laptopIDStr := scanner.Text()
	laptopID, err := strconv.Atoi(laptopIDStr)
	if err != nil {
		log.Fatal(err)
	}

	var laptopPrice float64
	err = db.QueryRow("SELECT price FROM Laptops WHERE laptop_id = ?", laptopID).Scan(&laptopPrice)
	if err != nil {
		log.Fatal(err)
	}

	// Insert data ke tabel Orders
	orderDate := time.Now()
	_, err = db.Exec("INSERT INTO Orders (user_id, order_date, total_amount) VALUES (?, ?, ?)",
		userID, orderDate, laptopPrice)
	if err != nil {
		log.Fatal(err)
	}

	// Ambil order_id yang baru saja di-generate
	var orderID int
	err = db.QueryRow("SELECT LAST_INSERT_ID()").Scan(&orderID)
	if err != nil {
		log.Fatal(err)
	}

	// Insert data ke tabel order_items
	quantity := 1 // Jumlah item yang dibeli, bisa disesuaikan
	subtotal := laptopPrice
	_, err = db.Exec("INSERT INTO order_items (order_id, laptop_id, quantity, subtotal) VALUES (?, ?, ?, ?)",
		orderID, laptopID, quantity, subtotal)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Laptop purchased successfully!")
}

func EditUser(db *sql.DB, scanner *bufio.Scanner) {
	fmt.Println("Edit User")
	fmt.Println("---------")

	fmt.Print("Enter user ID: ")
	scanner.Scan()
	userIDStr := scanner.Text()
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Print("Enter new full name: ")
	scanner.Scan()
	fullName := scanner.Text()

	fmt.Print("Enter new address: ")
	scanner.Scan()
	address := scanner.Text()

	fmt.Print("Enter new phone number: ")
	scanner.Scan()
	phoneNumber := scanner.Text()

	fmt.Print("Enter new birthdate (YYYY-MM-DD): ")
	scanner.Scan()
	birthdateStr := scanner.Text()
	birthdate, err := time.Parse("2006-01-02", birthdateStr)
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec("UPDATE user_profiles SET full_name = ?, address = ?, phone_number = ?, birthdate = ? WHERE user_id = ?",
		fullName, address, phoneNumber, birthdate, userID)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("User information updated successfully!")
}

func DeleteUser(db *sql.DB, scanner *bufio.Scanner) {
	fmt.Println("Delete User")
	fmt.Println("-----------")

	fmt.Print("Enter user ID: ")
	scanner.Scan()
	userIDStr := scanner.Text()
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec("DELETE FROM Users WHERE user_id = ?", userID)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("User deleted successfully!")
}

func PrintUserReport(db *sql.DB, scanner *bufio.Scanner) {
	rows, err := db.Query("SELECT user_id, username, email, created_at FROM users")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	fmt.Println("User Report:")
	fmt.Println("-----------------------------------------------------------------")
	fmt.Printf("| %-8s | %-15s | %-30s | %-19s |\n", "User ID", "Username", "Email", "Created At")
	fmt.Println("-----------------------------------------------------------------")
	for rows.Next() {
		var userID int
		var username, email, createdAt string
		err := rows.Scan(&userID, &username, &email, &createdAt)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("| %-8d | %-15s | %-30s | %-19s |\n", userID, username, email, createdAt)
	}
	fmt.Println("-----------------------------------------------------------------")
}

func PrintOrderReport(db *sql.DB, scanner *bufio.Scanner) {
	rows, err := db.Query("SELECT orders.order_id, users.username, orders.order_date, orders.total_amount FROM orders INNER JOIN users ON orders.user_id = users.user_id")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	fmt.Println("\nOrder Report:")
	fmt.Println("------------------------------------------------------------------------------------------------")
	fmt.Printf("| %-9s | %-15s | %-20s | %-10s |\n", "Order ID", "Username", "Order Date", "Total Amount")
	fmt.Println("------------------------------------------------------------------------------------------------")
	for rows.Next() {
		var orderID int
		var username, orderDate string
		var totalAmount float64
		err := rows.Scan(&orderID, &username, &orderDate, &totalAmount)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("| %-9d | %-15s | %-20s | %-10.2f |\n", orderID, username, orderDate, totalAmount)
	}
	fmt.Println("------------------------------------------------------------------------------------------------")
}

func PrintStockLaptopReport(db *sql.DB, scanner *bufio.Scanner) {
	rows, err := db.Query("SELECT laptop_id, brand, model, stock_quantity FROM laptops")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	fmt.Println("\nStock Laptop Report:")
	fmt.Println("--------------------------------------------------")
	fmt.Printf("| %-10s | %-15s | %-20s | %-13s |\n", "Laptop ID", "Brand", "Model", "Stock Quantity")
	fmt.Println("--------------------------------------------------")
	for rows.Next() {
		var laptopID int
		var brand, model string
		var stockQuantity int
		err := rows.Scan(&laptopID, &brand, &model, &stockQuantity)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("| %-10d | %-15s | %-20s | %-13d |\n", laptopID, brand, model, stockQuantity)
	}
	fmt.Println("--------------------------------------------------")
}
