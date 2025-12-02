package main

import (
	"fmt"
	"log"
	"os"

	"github.com/moniquelin/monday-hr/internal/data"
	"github.com/moniquelin/monday-hr/internal/database"
)

func main() {
	dsn := os.Getenv("MONDAY_HR_DB_DSN")
	if dsn == "" {
		log.Fatal("MONDAY_HR_DB_DSN is not set")
	}

	db, err := database.OpenDB(dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	models := data.NewModels(db)

	// 1. CREATE ADMIN USER
	admin := &data.User{
		IsAdmin:   true,
		Name:      "Admin",
		Email:     "admin@example.com",
		Salary:    0,
		CreatedBy: 0,
		UpdatedBy: 0,
	}

	err = admin.Password.Set("Password123!")
	if err != nil {
		log.Fatal(err)
	}

	err = models.Users.Insert(admin)
	if err != nil {
		log.Fatal("error inserting admin:", err)
	}

	fmt.Println("Created admin user with ID:", admin.ID)

	// 2. CREATE 100 EMPLOYEES
	for i := 1; i <= 100; i++ {
		u := &data.User{
			IsAdmin:   false,
			Name:      fmt.Sprintf("Employee %d", i),
			Email:     fmt.Sprintf("employee%d@example.com", i),
			Salary:    int64(5000000 + i*10000),
			CreatedBy: admin.ID,
			UpdatedBy: admin.ID,
		}

		err = u.Password.Set("Password123!")
		if err != nil {
			log.Fatal(err)
		}

		err = models.Users.Insert(u)
		if err != nil {
			log.Fatal("error inserting employee:", err)
		}
	}

	fmt.Println("Completed seeding 100 employees.")
}
