package api

import (
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/moniquelin/monday-hr/internal/data"
	"github.com/moniquelin/monday-hr/internal/domain"
	"github.com/moniquelin/monday-hr/internal/validator"
)

func (app *Application) createPayrollPeriodHandler(w http.ResponseWriter, r *http.Request) {
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

	// Get user from context
	user := app.contextGetUser(r)

	// Insert payroll period
	err = app.Models.PayrollPeriod.Insert(startDate, endDate, user.ID)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrPayrollPeriodOverlap):
			app.errorResponse(w, r, http.StatusConflict, err)
		case errors.Is(err, data.ErrPayrollPeriodDateOrder):
			app.errorResponse(w, r, http.StatusConflict, err)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	app.writeJSON(w, http.StatusCreated, envelope{
		"message": "payroll period created successfully",
	}, nil)
}

func (app *Application) runPayrollHandler(w http.ResponseWriter, r *http.Request) {
	// Get current user information
	user := app.contextGetUser(r)

	var input struct {
		PayrollPeriodId int64 `json:"id"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	// Get the current payroll period information
	log.Println("Processing payroll period with ID: ", input.PayrollPeriodId)
	p, err := app.Models.PayrollPeriod.Get(input.PayrollPeriodId)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	// Check the status and raise error if payroll has been processed
	log.Println("Check if payroll period has been processed...")
	if p.Status == "processed" {
		app.errorResponse(w, r, 409, "payroll period has already been processed")
		return
	}

	// Count number of working days in the payroll period
	log.Println("Counting the number of working days in the payroll period...")
	workingDays, err := domain.CountWorkingDays(p.StartDate, p.EndDate)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	// Get all employees (non-admin) to process the payroll for
	log.Println("Getting employees for the payroll...")
	recipientList, err := app.Models.Users.GetByRole("employee")
	if err != nil {
		app.serverErrorResponse(w, r, err)
		log.Println("Error getting employees!")
		return
	}

	// Generate payroll result for each employee
	for _, recipient := range recipientList {
		// Generate payroll result for each employee
		log.Println("-- Processing Employee ID: ", recipient.ID)
		// Snapshot employee attendance data into attendance_snapshots table
		log.Println("Snapshotting employee attendance data...")
		err = app.Models.AttendanceSnapshot.InsertAttendanceSnapshot(input.PayrollPeriodId, user.ID, recipient.ID, p.StartDate, p.EndDate)
		if err != nil {
			app.serverErrorResponse(w, r, err)
			log.Println("Error Inserting Attendance Snapshot!")
			return
		}
		// Get attendance days count
		log.Println("Getting attendance days count...")
		attDays, err := app.Models.AttendanceSnapshot.GetAttendanceDays(input.PayrollPeriodId, recipient.ID)
		if err != nil {
			app.serverErrorResponse(w, r, err)
			log.Println("Error Getting Attendance Count!")
			return
		}

		// Get employee salary
		log.Println("Getting employee salary...")
		recipientSalary := int(recipient.Salary)

		// Calculate employee take home pay
		log.Println("Calculating employee take home pay...")
		takeHomePay := domain.CalculateTakeHomePay(attDays, workingDays, recipientSalary)

		// Insert employee data and take home pay to payroll_results table
		log.Println("Saving payroll results...")
		err = app.Models.PayrollResult.InsertPayrollResult(p.ID, recipient.ID, recipientSalary, workingDays, attDays, takeHomePay)
		if err != nil {
			app.serverErrorResponse(w, r, err)
			log.Println("Error Inserting THP!")
			return
		}
	}

	// Change payroll period status from draft to processed
	err = app.Models.PayrollPeriod.MarkAsProcessed(input.PayrollPeriodId, user.ID)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	log.Println("Successfully run payroll period!")
	app.writeJSON(w, http.StatusCreated, envelope{
		"message": "successfully run payroll period",
	}, nil)
}
