package common

import (
	"app/news-parser/internal/custom_errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

const (
	EventClickCategory = "click_category"
	EventClickArticle  = "click_article"

	OffsetDefault = 0
	LimitDefault  = 50

	RdbTimeout = time.Second * 10

	UnixMonth = 2592000
)

func DateNow() time.Time {
	now := time.Now()
	if now.Month() < 10 && now.Day() < 10 {
		date, _ := time.Parse(time.DateOnly, fmt.Sprintf("%d-0%d-0%d", now.Year(), now.Month(), now.Day()))
		return date
	} else if now.Day() < 10 {
		date, _ := time.Parse(time.DateOnly, fmt.Sprintf("%d-%d-0%d", now.Year(), now.Month(), now.Day()))
		return date
	} else if now.Month() < 10 {
		date, _ := time.Parse(time.DateOnly, fmt.Sprintf("%d-0%d-%d", now.Year(), now.Month(), now.Day()))
		return date
	} else {
		date, _ := time.Parse(time.DateOnly, fmt.Sprintf("%d-%d-%d", now.Year(), now.Month(), now.Day()))
		return date
	}
}
func ValidateOffsetAndLimit(offsetStr, limitStr string) (int, int, error) {
	var offset, limit int
	var errParseOffset, errParseLimit error
	if offsetStr != "" {
		offset, errParseOffset = strconv.Atoi(offsetStr)
	} else {
		offset = OffsetDefault
	}
	if limitStr != "" {
		limit, errParseLimit = strconv.Atoi(limitStr)
	} else {
		limit = LimitDefault
	}
	if errParseLimit != nil && errParseOffset != nil {
		return 0, 0, custom_errors.ErrIncorrectOffsetAndLimit
	} else if errParseOffset != nil {
		return 0, 0, custom_errors.ErrIncorrectOffset
	} else if errParseLimit != nil {
		return 0, 0, custom_errors.ErrIncorrectLimit
	}
	return offset, limit, nil
}
func ParseString(str string) string {
	sliceStr := strings.Fields(str)
	resStr := ""
	for _, s := range sliceStr {
		resStr += s + " "
	}
	return resStr
}
