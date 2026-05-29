package send_letter

import (
	"app/news-parser/internal/custom_errors"
	"fmt"
	"net/smtp"
	"time"

	"github.com/jordan-wright/email"
)

type SenderLetter struct {
	ChErr       chan error
	ApiEmail    string
	ApiPassword string
	Address     string
	AddressHost string
}

func NewSenderLetter(apiEmail, apiPassword, address, addressHost string) *SenderLetter {
	return &SenderLetter{
		ChErr:       make(chan error),
		ApiEmail:    apiEmail,
		ApiPassword: apiPassword,
		Address:     address,
		AddressHost: addressHost,
	}
}
func (l *SenderLetter) send(userEmail string, tempCode uint) {
	defer close(l.ChErr)
	e := email.NewEmail()
	e.From = l.ApiEmail
	e.To = []string{userEmail}
	e.Text = []byte(fmt.Sprintf("your code %d for authorization in the news parsing service", tempCode))
	errSend := e.Send(l.AddressHost, smtp.PlainAuth("", l.ApiEmail, l.ApiPassword, l.Address))
	if errSend != nil {
		l.ChErr <- errSend
	}
	l.ChErr <- nil
}
func (l *SenderLetter) SendEmailLetter(userEmail string, tempCode uint) error {
	after := time.After(time.Second * 30)
	go l.send(userEmail, tempCode)
	select {
	case <-after:
		return custom_errors.ErrSendLetter
	case errSend := <-l.ChErr:
		if errSend != nil {
			return errSend
		}
		return nil
	}
}
