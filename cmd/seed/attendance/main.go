package main

import (
	"fmt"
	"log"
	"math/rand/v2"
	"os"
	"time"

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

	fmt.Println("Attendance seeding started...")

	// Determine date range
	loc, _ := time.LoadLocation("Asia/Jakarta")
	seedStartDate := time.Date(2026, 01, 01, 9, 0, 0, 0, loc)
	seedEndDate := time.Date(2026, 01, 30, 9, 0, 0, 0, loc)

	// Get list of employees to record check in for
	employees, err := models.Users.GetByRole("employee")

	// Insert attendance for each date
	for d := seedStartDate; d.After(seedEndDate) == false; d = d.AddDate(0, 0, 1) {
		log.Println("--- PROCESSING DATE: ", d.String())
		// Prepare dummy check out time to be seeded to the employee attendance data
		checkOutAt := d.Add(time.Hour * 8)

		switch d.Weekday() {
		case time.Saturday, time.Sunday:
			fmt.Println("Skipping ", d, " since it is a weekend.")
		default:
			// Insert attendance for each employee
			for _, emp := range employees {
				// Determine whether employee checked in or not using random mechanism
				randomNumber := rand.IntN(100)
				if randomNumber <= 90 {
					// If random integer is less than 90, then
					// employee does check in on that date
					// but no check out
					log.Println("Recording check in on ", d, " for employeee with ID:", emp.ID)
					empAttendance := data.Attendance{
						EmployeeID: emp.ID,
						AttDate:    d.Format("2006-01-02"),
						CheckInAt:  d.In(loc),
						CreatedBy:  1,
						UpdatedBy:  1,
					}
					err = models.Attendance.RecordCheckIn(&empAttendance)
					if err != nil {
						log.Println("Error recording check in on ", d, " for employeee with ID:", emp.ID)
						log.Println(err)
					}

					if randomNumber <= 70 {
						// If random integer is less than 70, then
						// employee does check in on that date
						// and also checks out
						log.Println("Recording check out for the same employee...")
						empAttendance.CheckOutAt = &checkOutAt
						err = models.Attendance.RecordCheckOut(&empAttendance)
					}

				}
			}
		}
	}
	fmt.Println("Completed seeding.")
}
