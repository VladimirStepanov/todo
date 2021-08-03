//https://medium.com/glottery/sending-emails-with-go-golang-and-gmail-39bc20423cf0

package service

import (
	"fmt"
	"net/smtp"
	"strings"

	"github.com/VladimirStepanov/todo-app/internal/models"
)

type MailService struct {
	Email    string
	Password string
	Domain   string
}

func (ms *MailService) SendConfirmationsEmail(user *models.User) error {
	from := ms.Email
	pass := ms.Password
	to := user.Email
	server := "smtp.gmail.com"
	port := "587"

	msg := "From: " + from + "\n" +
		"To: " + to + "\n" +
		"Subject: Email conficmation\n\n" +
		fmt.Sprintf(
			"Confirm your email: http://%s/auth/confirm/%s",
			ms.Domain, user.ActivatedLink,
		)

	return smtp.SendMail(strings.Join([]string{server, port}, ":"),
		smtp.PlainAuth("", from, pass, server),
		from, []string{to}, []byte(msg))
}

func NewMailService(Email, Password, Domain string) models.MailService {
	return &MailService{
		Email:    Email,
		Password: Password,
		Domain:   Domain,
	}
}
