package data

import (
	"database/sql"
)

// Models contrains all data models used in the application
type Models struct {
	Users         UserModel
	Attendance    AttendanceModel
	PayrollPeriod PayrollPeriodModel
	// To be added:
	// Payrolls    PayrollModel
	// dll.
}

// Initialize all models with DB connection
func NewModels(db *sql.DB) Models {
	return Models{
		Users:         UserModel{DB: db},
		Attendance:    AttendanceModel{DB: db},
		PayrollPeriod: PayrollPeriodModel{DB: db},
	}
}
