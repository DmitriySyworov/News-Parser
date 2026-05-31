package common

import (
	"app/news-parser/internal/custom_errors"
	"app/news-parser/internal/response"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const (
	OffsetDefault = 0
	LimitDefault  = 50

	RdbTimeout = time.Second * 10

	Day = time.Hour * 24

	UnixMonth = 2592000

	LengthTempCode = 6
	LengthSession  = 9
	KeySession     = "session:"
	MessageEmail   = "we sent a letter to the specified email: "
)

type StatDataUserArticle struct {
	UserUUID string
	Number   int
}

func DateNow() time.Time {
	now := time.Now().Format(time.DateOnly)
	date, _ := time.Parse(time.DateOnly, now)
	return date
}
func ValidateOffsetAndLimit(offsetStr, limitStr string) (int, int, []response.Error) {
	var sliceError []response.Error
	var offset, limit int
	var errParseOffset, errParseLimit error
	if offsetStr != "" {
		offset, errParseOffset = strconv.Atoi(offsetStr)
		if errParseOffset != nil {
			sliceError = append(sliceError, response.Error{
				Message: custom_errors.ErrIncorrectOffset.Error(),
				Status:  http.StatusBadRequest,
			})
		}
	} else {
		offset = OffsetDefault
	}
	if offset < 0 {
		sliceError = append(sliceError, response.Error{
			Message: custom_errors.ErrIncorrectOffset.Error(),
			Status:  http.StatusBadRequest,
		})
	}
	if limitStr != "" {
		limit, errParseLimit = strconv.Atoi(limitStr)
		if errParseLimit != nil {
			sliceError = append(sliceError, response.Error{
				Message: custom_errors.ErrIncorrectLimit.Error(),
				Status:  http.StatusBadRequest,
			})
		}
	} else {
		limit = LimitDefault
	}
	if limit < 0 {
		sliceError = append(sliceError, response.Error{
			Message: custom_errors.ErrIncorrectLimit.Error(),
			Status:  http.StatusBadRequest,
		})
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
