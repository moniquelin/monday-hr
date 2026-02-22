package api

import (
	"fmt"
	"log"
	"math/rand/v2"
	"net/http"
	"time"

	"github.com/moniquelin/monday-hr/internal/data"
	"github.com/moniquelin/monday-hr/internal/validator"
)

// seedUsers seeds 1 Super Admin, 1 Admin, and 100 Employees into the users table
func (app *Application) seedUsers(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Start seeding users...")
	// 1. CREATE SUPER ADMIN / SEED ADMIN
	superAdmin := &data.User{
		Role:      "super_admin",
		Name:      "Super Admin",
		Email:     "superadmin@example.com",
		Salary:    0,
		CreatedBy: 0,
		UpdatedBy: 0,
	}

	err := superAdmin.Password.Set("Password123!")
	if err != nil {
		fmt.Println("Error setting password!")
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.Models.Users.Insert(superAdmin)
	if err != nil {
		fmt.Println("Error inserting super admin!")
		app.serverErrorResponse(w, r, err)
		return
	}

	fmt.Println("Created 1 super admin user")

	// 2. CREATE ADMIN USER (ordinary admin, not super admin)
	admin := &data.User{
		Role:      "admin",
		Name:      "Admin 1",
		Email:     "admin1@example.com",
		Salary:    0,
		CreatedBy: superAdmin.ID,
		UpdatedBy: superAdmin.ID,
	}

	err = admin.Password.Set("Password123!")
	if err != nil {
		fmt.Println("Error setting password!")
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.Models.Users.Insert(admin)
	if err != nil {
		fmt.Println("Error inserting admin!")
		app.serverErrorResponse(w, r, err)
		return
	}

	fmt.Println("Created 1 admin user")

	// 3. CREATE 100 EMPLOYEES
	for i := 1; i <= 100; i++ {
		u := &data.User{
			Role:      "employee",
			Name:      fmt.Sprintf("Employee %d", i),
			Email:     fmt.Sprintf("employee%d@example.com", i),
			Salary:    int64(5000000 + i*10000),
			CreatedBy: admin.ID,
			UpdatedBy: admin.ID,
		}

		err = u.Password.Set("Password123!")
		if err != nil {
			fmt.Println("Error setting password!")
			app.serverErrorResponse(w, r, err)
			return
		}

		err = app.Models.Users.Insert(u)
		if err != nil {
			fmt.Println("Error inserting employee!")
			app.serverErrorResponse(w, r, err)
			return
		}
	}
	fmt.Println("Completed seeding 100 employees")
	app.writeJSON(w, http.StatusCreated, envelope{
		"message": "Completed seeding 100 employees",
	}, nil)
}

// seedAttendance randomly inserts attendance for every employees
// between two dates, with 70% probability of employee checking in and out,
// 20% probability of employee only checking in,  and
// 10% probability of employee being absent
func (app *Application) seedAttendance(w http.ResponseWriter, r *http.Request) {
	var input struct {
		StartDate string `json:"start_date"`
		EndDate   string `json:"end_date"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	// Parse dates
	startDate, err := time.Parse("2006-01-02", input.StartDate)
	if err != nil {
		app.failedValidationResponse(w, r, map[string]string{
			"start_date": "must be a valid date (YYYY-MM-DD)",
		})
		return
	}

	endDate, err := time.Parse("2006-01-02", input.EndDate)
	if err != nil {
		app.failedValidationResponse(w, r, map[string]string{
			"end_date": "must be a valid date (YYYY-MM-DD)",
		})
		return
	}

	// Validation
	v := validator.New()

	validator.ValidateDate(v, &startDate, "start_date")
	validator.ValidateDate(v, &endDate, "end_date")

	// Domain rule: end_date >= start_date
	v.Check(!endDate.Before(startDate), "end_date", "must be on or after start_date")

	if !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	fmt.Println("Start seeding attendance...")

	// Determine date range
	loc, _ := time.LoadLocation("Asia/Jakarta")
	seedStartDate := time.Date(startDate.Year(), startDate.Month(), startDate.Day(), 9, 0, 0, 0, loc)
	seedEndDate := time.Date(endDate.Year(), endDate.Month(), endDate.Day(), 9, 0, 0, 0, loc)

	// Get list of employees to record check in for
	employees, err := app.Models.Users.GetByRole("employee")

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
					err = app.Models.Attendance.RecordCheckIn(&empAttendance)
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
						err = app.Models.Attendance.RecordCheckOut(&empAttendance)
						if err != nil {
							log.Println("Error recording check out on ", d, " for employeee with ID:", emp.ID)
							log.Println(err)
						}
					}

				}
			}
		}
	}
	fmt.Println("Completed seeding attendance")
	app.writeJSON(w, http.StatusCreated, envelope{
		"message": "Completed seeding attendance",
	}, nil)
}
