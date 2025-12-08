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
		CreatedBy:  user.ID,
		UpdatedBy:  user.ID,
	}

	// Record check in
	err := app.Models.Attendance.RecordCheckIn(att)
	if err != nil {
		switch {
		// If there is already check in for the date
		case errors.Is(err, data.ErrDuplicateCheckIn):
			app.errorResponse(w, r, 409, "already recorded check in today")
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	app.writeJSON(w, http.StatusCreated, envelope{
		"attendance": att,
	}, nil)
}

// checkInHandler enables employee to record check out
func (app *Application) checkOutHandler(w http.ResponseWriter, r *http.Request) {

}
