package validator

import "time"

// ValidateDate checks if the date is present and has valid YYYY-MM-DD format.
func ValidateDate(v *Validator, date *time.Time, field string) {
	v.Check(date != nil, field, "must be provided")
	v.Check(!date.IsZero(), field, "must be a valid date")
}
