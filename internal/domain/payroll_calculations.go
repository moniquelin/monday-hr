package domain

import (
	"time"
)

func CountWorkingDays(startDate, endDate string) (int, error) {
	const layout = "2006-01-02"

	start, err := time.Parse(layout, startDate)
	if err != nil {
		return 0, err
	}

	end, err := time.Parse(layout, endDate)
	if err != nil {
		return 0, err
	}

	// Guard: start > end
	if start.After(end) {
		return 0, nil
	}

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
	return (attendanceDays / workingDays) * baseSalary
}
