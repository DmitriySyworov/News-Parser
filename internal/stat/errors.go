package stat

import "errors"

var (
	ErrPeriod   = errors.New("incorrect or empty period")
	ErrStatLoad = errors.New("statistic is empty or failed to load")
)
