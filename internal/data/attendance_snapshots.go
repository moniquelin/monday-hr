package data

import (
	"context"
	"database/sql"
	"time"
)

type AttendanceSnapshot struct {
	ID              int64 `json:"id"`
	PayrollPeriodID int64 `json:"payroll_period_id"`
	EmployeeID      int64 `json:"employee_id"`

	AttDate    int64 `json:"att_date"`
	CheckinAt  int64 `json:"checkin_at"`
	CheckoutAt int64 `json:"checkout_at"`

	CreatedBy int64     `json:"created_by"`
	CreatedAt time.Time `json:"created_at"`
}

type AttendanceSnapshotModel struct {
	DB *sql.DB
}

// InsertAttendanceSnapshot gets the attendance data for a payroll period and
// inserts them into attendance_snapshots table
func (m AttendanceSnapshotModel) InsertAttendanceSnapshot(payrollPeriodId int64, createdBy int64, employeeId int64, periodStartDate string, periodEndDate string) error {
	query := `
		INSERT INTO attendance_snapshots (payroll_period_id, employee_id, att_date, checkin_at, checkout_at, created_by)
		SELECT
		$1 AS payroll_period_id,
		a.employee_id,
		a.att_date,
		a.checkin_at,
		a.checkout_at,
		$2 AS created_by
		FROM attendance a
		WHERE a.employee_id = $3
		AND a.att_date BETWEEN $4 AND $5`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := m.DB.ExecContext(ctx, query, payrollPeriodId, createdBy, employeeId, periodStartDate, periodEndDate)
	if err != nil {
		return err
	}

	return nil
}

// GetAttendanceCount gets the count of attendance days from a payroll period for a specific employee
func (m AttendanceSnapshotModel) GetAttendanceDays(payrollPeriodId int64, employeeId int64) (int, error) {
	query := `
		SELECT COUNT(ps.checkin_at)
		FROM attendance_snapshots atts
		WHERE atts.payroll_period_id = $1
		AND atts.employee_id = $2`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var count int
	err := m.DB.QueryRowContext(ctx, query, payrollPeriodId, employeeId).Scan(&count)
	if err != nil {
		return count, err
	}

	return count, err
}
