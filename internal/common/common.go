package common

import (
	"fmt"
	"time"
)

const (
	EventClickCategory = "click_category"
	EventClickArticle  = "click_article"

	RdbTimeout = time.Second * 10
)

func DateNow() time.Time{
	now := time.Now()
	if now.Month() < 10 {
		date, _ := time.Parse(time.DateTime,  fmt.Sprintf("%d-0%d-%d", now.Year(), now.Month(), now.Day()))
		return date
	} else if now.Day() < 10 {
		date, _ :=  time.Parse(time.DateOnly, fmt.Sprintf("%d-%d-0%d", now.Year(), now.Month(), now.Day()))
		return date
	} else if now.Day() < 10 && now.Month() < 10 {
		date, _ := time.Parse(time.DateOnly, fmt.Sprintf("%d-0%d-0%d", now.Year(), now.Month(), now.Day()))
		return date
	} else {
		date, _ :=  time.Parse(time.DateOnly, fmt.Sprintf("%d-%d-%d", now.Year(), now.Month(), now.Day()))
		return date
	}
}
