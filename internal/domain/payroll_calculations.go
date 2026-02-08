package domain

import (
	"time"
)

func CountWorkingDays(start, end time.Time) (int, error) {
	businessDays := 0

	for d := start; !d.After(end); d = d.AddDate(0, 0, 1) {
		weekday := d.Weekday()
		if weekday != time.Saturday && weekday != time.Sunday {
			businessDays++
		}
	}

	return businessDays, nil
}

func CalculateTakeHomePay(attendanceDays, workingDays, baseSalary int) int {
	takeHomePay := float64(attendanceDays) / float64(workingDays) * float64(baseSalary)
	return int(takeHomePay)
}
