package common

import (
	"app/news-parser/configs"
	"app/news-parser/internal/custom_errors"
	"app/news-parser/pkg/send_letter"
	"math/rand/v2"
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
	EventClickCategory         = "click_category"
	EventClickArticle          = "click_article"
	EventCreateUserArticle     = "create_article"
	EventUpdateUserArticle     = "update_article"
	EventSoftDeleteUserArticle = "soft_delete_article"
	EventHardDeleteUserArticle = "hard_delete_article"
	EventRecoveryUserArticle   = "recovery_article"

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
func SendRequest(url string) (*http.Response, error) {
	poolUserAgents := []string{
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/124.0.0.0 Safari/537.36",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 14_4_1) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.4.1 Safari/605.1.15",
		"Mozilla/5.0 (X11; Linux x64; rv:125.0) Gecko/20100101 Firefox/125.0",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/124.0.0.0 Safari/537.36 Edg/124.0.0.0",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/122.0.0.0 Safari/537.36 OPR/108.0.0.0",
		"Mozilla/5.0 (Linux; Android 14; Pixel 8) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/124.0.6367.82 Mobile Safari/537.36",
		"Mozilla/5.0 (iPhone; CPU iPhone OS 17_4_1 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.4.1 Mobile/15E148 Safari/605.1.15",
		"Mozilla/5.0 (Linux; Android 14; SAMSUNG SM-S928B) AppleWebKit/537.36 (KHTML, like Gecko) SamsungBrowser/24.0 Chrome/115.0.0.0 Mobile Safari/537.36",
	}
	request, errReq := http.NewRequest(http.MethodGet, url, nil)
	if errReq != nil {
		return nil, errReq
	}
	randomIndex := rand.IntN(len(poolUserAgents))
	request.Header.Set("User-Agent", poolUserAgents[randomIndex])
	request.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8")
	request.Header.Set("Accept-Language", "ru,en;q=0.9,en-US;q=0.8")
	request.Header.Set("Connection", "keep-alive")
	client := http.Client{
		Timeout: 5 * time.Second,
	}
	response, errResp := client.Do(request)
	if errResp != nil {
		return nil, errResp
	}
	return response, nil
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
	if offset < 0 {
		sliceError = append(sliceError, custom_errors.Error{
			Message: custom_errors.ErrNegativeOffset.Error(),
			Status:  http.StatusBadRequest,
		})
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
	if limit < 0 {
		sliceError = append(sliceError, custom_errors.Error{
			Message: custom_errors.ErrNegativeLimit.Error(),
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
