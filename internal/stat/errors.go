package stat

import "errors"

var (
	ErrIncorrectDate = errors.New("the 'date' must be a YYYY-MM-DD")
	ErrStatNotFound  = errors.New("statistic not found or empty")
)
