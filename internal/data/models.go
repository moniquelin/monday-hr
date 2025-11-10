package data

import (
	"database/sql"
	"errors"
)

// App-wide errors
var (
	ErrRecordNotFound = errors.New("record not found")
)

// Models contrains all data models used in the application
type Models struct {
	Users UserModel
	// To be added:
	// Attendances AttendanceStore
	// Overtimes   OvertimeStore
	// Payrolls    PayrollStore
	// dll.
}

// Initialize all models with DB connection
func NewModels(db *sql.DB) Models {
	return Models{
		Users: UserModel{DB: db},
	}
}
