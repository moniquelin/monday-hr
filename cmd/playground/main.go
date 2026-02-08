package main

import (
	"fmt"

	"github.com/moniquelin/monday-hr/internal/domain"
)

func main() {
	// dsn := os.Getenv("MONDAY_HR_DB_DSN")
	// if dsn == "" {
	// 	log.Fatal("MONDAY_HR_DB_DSN is not set")
	// }

	// db, err := database.OpenDB(dsn)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// defer db.Close()
	// loc, err := time.LoadLocation("Asia/Jakarta")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	recipientSalary := 5010000
	takeHomePay := domain.CalculateTakeHomePay(18, 22, recipientSalary)
	fmt.Println(takeHomePay)
}
