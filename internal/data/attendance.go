package data

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/lib/pq"
)

// Error for duplicate emails
var (
	ErrDuplicateCheckIn  = errors.New("employee have already checked in on the date")
	ErrDuplicateCheckOut = errors.New("employee have already checked out on the date")
)

// Attendance struct represents attendance data of one date
type Attendance struct {
	ID         int64      `json:"id"`
	EmployeeID int64      `json:"employee_id"`
	AttDate    string     `json:"att_date"`
	CheckInAt  time.Time  `json:"checkin_at"`
	CheckOutAt *time.Time `json:"checkout_at,omitempty"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
	CreatedBy  int64      `json:"created_by"`
	UpdatedBy  int64      `json:"updated_by"`
}

// AttendanceModel struct wraps the connection pool
type AttendanceModel struct {
	DB *sql.DB
}

// Record new employee check-in in the database
func (m AttendanceModel) RecordCheckIn(attendance *Attendance) error {
	query := `
		INSERT INTO attendance (employee_id, att_date, checkin_at, created_by, updated_by)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at, updated_at`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// If the table already contains a check in record for this employee at the check in date,
	// there will be a violation of the UNIQUE constraint, since attendance on the same day should
	// count as one
	err := m.DB.QueryRowContext(ctx, query,
		attendance.EmployeeID,
		attendance.AttDate,
		&attendance.CheckInAt,
		attendance.CreatedBy,
		attendance.UpdatedBy,
	).Scan(&attendance.ID, &attendance.CreatedAt, &attendance.UpdatedAt)
	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) {
			if pqErr.Code == "23505" { // unique_violation
				return ErrDuplicateCheckIn
			}
		}
		return err
	}

	return nil
}

// Get attendance data from the database
func (m AttendanceModel) Get(employeeId int64, date string) (*Attendance, error) {
	query := `
        SELECT id, employee_id, att_date, checkin_at, checkout_at, created_at, created_by, updated_at, updated_by
        FROM attendance
        WHERE employee_id = $1 AND att_date = $2`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var attendance Attendance

	err := m.DB.QueryRowContext(ctx, query, employeeId, date).Scan(
		&attendance.ID,
		&attendance.EmployeeID,
		&attendance.AttDate,
		&attendance.CheckInAt,
		&attendance.CheckOutAt,
		&attendance.CreatedAt,
		&attendance.CreatedBy,
		&attendance.UpdatedAt,
		&attendance.UpdatedBy,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrRecordNotFound
		}
		return nil, err
	}

	return &attendance, nil
}

// Record new employee check-out in the database
func (m AttendanceModel) RecordCheckOut(attendance *Attendance) error {
	query := `
		UPDATE attendance 
		SET updated_by = $1, checkout_at = $2, updated_at = now()
		WHERE employee_id = $3 AND att_date = $4
	`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	result, err := m.DB.ExecContext(ctx, query,
		attendance.UpdatedBy,
		attendance.CheckOutAt,
		attendance.EmployeeID,
		attendance.AttDate,
	)

	// Check if UPDATE applies to any row
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	// If there is no rows affected → no attendance data
	if rows == 0 {
		return ErrRecordNotFound
	}

	return nil
}
