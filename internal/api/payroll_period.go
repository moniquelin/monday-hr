package api

import (
	"errors"
	"net/http"
	"time"

	"github.com/moniquelin/monday-hr/internal/data"
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

	// Initialize new payroll period
	payrollPeriod := data.PayrollPeriod{
		StartDate: startDate.Format("2006-01-02"),
		EndDate:   endDate.Format("2006-01-02"),
		CreatedBy: user.ID,
		UpdatedBy: user.ID,
	}

	// Insert payroll period
	err = app.Models.PayrollPeriod.Insert(payrollPeriod)
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
