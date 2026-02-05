package main

import (
	"fmt"
	"log"
	"os"

	"github.com/moniquelin/monday-hr/internal/data"
	"github.com/moniquelin/monday-hr/internal/database"
)

func main() {
	// Setup logger
	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	// Get DB connection string (sama kayak di main.go)
	dsn := os.Getenv("MONDAY_HR_DB_DSN")
	if dsn == "" {
		log.Fatal("missing MONDAY_HR_DB_DSN environment variable")
	}

	// Connect ke DB
	db, err := database.OpenDB(dsn)
	if err != nil {
		logger.Fatal(err)
	}
	defer db.Close()

	logger.Println("database connection pool established")

	// Initialize models (sama kayak di app)
	models := data.NewModels(db)

	// ==========================================
	// AREA TESTING - Edit bagian ini sesuka hati
	// ==========================================

	// Test GetByRole
	fmt.Println("\n=== Testing GetByRole ===")
	users, err := models.Users.GetByRole("employee")
	if err != nil {
		log.Printf("Error: %v", err)
	} else {
		fmt.Printf("%", users)
		// fmt.Printf("Found %d users with role 'admin'\n", len(users))
		// for i, user := range users {
		// 	fmt.Printf("%d. ID: %d, Name: %s, Email: %s, Salary: %.2f\n",
		// 		i+1, user.ID, user.Name, user.Email, user.Salary)
		// }
	}

	// Test function lain
	// fmt.Println("\n=== Testing GetByID ===")
	// user, err := models.Users.GetByID(1)
	// if err != nil {
	//     log.Printf("Error: %v", err)
	// } else {
	//     fmt.Printf("User: %+v\n", user)
	// }

	// Atau test apapun yang lu mau...
	// fmt.Println("\n=== Testing Insert ===")
	// newUser := &data.User{
	//     Name:  "Test User",
	//     Email: "test@example.com",
	//     Role:  "employee",
	// }
	// err = models.Users.Insert(newUser)
	// if err != nil {
	//     log.Printf("Error: %v", err)
	// } else {
	//     fmt.Printf("New user created with ID: %d\n", newUser.ID)
	// }
}
