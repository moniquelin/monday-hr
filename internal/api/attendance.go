package api

import (
	"errors"
	"net/http"
	"time"

	"github.com/moniquelin/monday-hr/internal/data"
)

// checkInHandler enables employee to record check in
func (app *Application) checkInHandler(w http.ResponseWriter, r *http.Request) {
	user := app.contextGetUser(r)

	// Determine attendance date in WIB
	loc, _ := time.LoadLocation("Asia/Jakarta")
	attDateWIB := time.Now().In(loc)

	// Validate if date is not weekend
	if attDateWIB.Weekday() == time.Saturday || attDateWIB.Weekday() == time.Sunday {
		app.errorResponse(w, r, 422, "cannot check in on the weekend")
		return
	}

	att := &data.Attendance{
		EmployeeID: user.ID,
		AttDate:    attDateWIB.Format("2006-01-02"),
		CheckInAt:  attDateWIB,
		CreatedBy:  user.ID,
		UpdatedBy:  user.ID,
	}

	// Record check in
	err := app.Models.Attendance.RecordCheckIn(att)
	if err != nil {
		switch {
		// If there is already check in for the date
		case errors.Is(err, data.ErrDuplicateCheckIn):
			app.errorResponse(w, r, 409, err)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	app.writeJSON(w, http.StatusCreated, envelope{
		"message":    "checked-in successfully",
		"attendance": att,
	}, nil)
}

// checkOutHandler enables employee to record check out
func (app *Application) checkOutHandler(w http.ResponseWriter, r *http.Request) {
	user := app.contextGetUser(r)

	// Determine attendance date in WIB
	loc, _ := time.LoadLocation("Asia/Jakarta")
	attDateWIB := time.Now().In(loc)

	// Validate if date is not weekend
	if attDateWIB.Weekday() == time.Saturday || attDateWIB.Weekday() == time.Sunday {
		app.errorResponse(w, r, 422, "cannot check out on the weekend")
		return
	}

	// Check if there is already an attendance data
	att, err := app.Models.Attendance.Get(user.ID, attDateWIB.Format("2006-01-02"))
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.errorResponse(w, r, 404, "no check-in data for the date")
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	if att.CheckOutAt != nil {
		// If there is already check in for the date
		app.errorResponse(w, r, 409, "already checked out for the date")
		return
	}

	att = &data.Attendance{
		EmployeeID: user.ID,
		CheckOutAt: &attDateWIB,
		AttDate:    attDateWIB.Format("2006-01-02"),
		UpdatedBy:  user.ID,
	}

	// Record check out
	err = app.Models.Attendance.RecordCheckOut(att)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.errorResponse(w, r, 404, "no check-in data for the date")
		default:
			app.serverErrorResponse(w, r, err)

		}
		return
	}
	app.writeJSON(w, http.StatusCreated, envelope{
		"message":    "checked-out successfully",
		"attendance": att,
	}, nil)
}
