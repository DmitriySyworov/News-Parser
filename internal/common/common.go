package common

import (
	"app/news-parser/configs"
	"app/news-parser/internal/custom_errors"
	"app/news-parser/pkg/send_letter"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type ResponseSuccessful struct {
	Success bool
	Data    any
}

const (
	EventClickCategory = "click_category"
	EventClickArticle  = "click_article"

	OffsetDefault = 0
	LimitDefault  = 50

	RdbTimeout = time.Second * 10

	UnixMonth = 2592000

	LengthTempCode = 6
	LengthSession  = 9
	KeySession     = "session:"
	MessageEmail   = "we sent a letter to the specified email: "
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
func ValidateOffsetAndLimit(offsetStr, limitStr string) (int, int, []custom_errors.Error) {
	var sliceError []custom_errors.Error
	var offset, limit int
	var errParseOffset, errParseLimit error
	if offsetStr != "" {
		offset, errParseOffset = strconv.Atoi(offsetStr)
		if errParseOffset != nil {
			sliceError = append(sliceError, custom_errors.Error{
				Message: custom_errors.ErrIncorrectOffset.Error(),
				Status:  http.StatusBadRequest,
			})
		}
	} else {
		offset = OffsetDefault
	}
	if limitStr != "" {
		limit, errParseLimit = strconv.Atoi(limitStr)
		if errParseLimit != nil {
			sliceError = append(sliceError, custom_errors.Error{
				Message: custom_errors.ErrIncorrectLimit.Error(),
				Status:  http.StatusBadRequest,
			})
		}
	} else {
		limit = LimitDefault
	}
	if len(sliceError) != 0 {
		return 0, 0, sliceError
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

type RequestConfirm struct {
	Code uint `json:"code" validate:"required"`
}
type ResponseAuth struct {
	Message string `json:"message"`
	JWTTemp string `json:"jwt-temp"`
}

func SendEmailLetter(userEmail string, tempCode uint, conf *configs.Configs) error {
	after := time.After(time.Second * 30)
	letter := send_letter.NewSenderLetter(conf.ApiEmail, conf.ApiPassword, conf.Address, conf.AddressHost)
	go letter.SendEmailLetter(userEmail, tempCode)
	select {
	case <-after:
		return custom_errors.ErrSendLetter
	case errSend := <-letter.ChErr:
		if errSend != nil {
			return errSend
		}
		return nil
	}
}
