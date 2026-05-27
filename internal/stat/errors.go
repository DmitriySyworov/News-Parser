package stat

import "errors"

var (
	ErrIncorrectDate = errors.New("the 'date' must be a YYYY-MM-DD")
	ErrPeriod       = errors.New("incorrect or empty period")
	ErrStatNotFound = errors.New("statistic not found or empty")
)
