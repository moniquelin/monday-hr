package data

import (
	"context"
	"database/sql"
	"time"
)

type PayrollResult struct {
	ID              int64 `json:"id"`
	PayrollPeriodID int64 `json:"payroll_period_id"`
	EmployeeID      int64 `json:"employee_id"`

	BaseSalary     int64 `json:"base_salary"`
	WorkingDays    int64 `json:"working_days"`
	AttendanceDays int64 `json:"attendance_days"`
	TakeHomePay    int64 `json:"take_home_pay"`

	CreatedAt time.Time `json:"created_at"`
}

type PayrollResultModel struct {
	DB *sql.DB
}

// InsertPayrollRresult inserts the employee data and take home pay to the
// payroll_results table
func (m PayrollResultModel) InsertPayrollResult(payrollPeriodId int64, employeeId int64, baseSalary int, workingDays int, attDays int, takeHomePay int) error {
	query := `
		INSERT INTO payroll_results (payroll_period_id, employee_id, base_salary, working_days, attendance_days, take_home_pay)
		VALUES ($1,	$2, $3, $4, $5, $6)`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := m.DB.ExecContext(ctx, query, payrollPeriodId, employeeId, baseSalary, workingDays, attDays, takeHomePay)
	if err != nil {
		return err
	}

	return nil
}
