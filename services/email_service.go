package services

import (
	"fmt"

	"github.com/go-mail/mail/v2"
)

type EmailService struct {
	dialer *mail.Dialer
	from   string
}

func NewEmailService(host string, port int, user, password string) *EmailService {
	d := mail.NewDialer(host, port, user, password)
	return &EmailService{
		dialer: d,
		from:   user,
	}
}

func (s *EmailService) SendPaymentNotification(to string, amount float64) error {
	m := mail.NewMessage()
	m.SetHeader("From", s.from)
	m.SetHeader("To", to)
	m.SetHeader("Subject", "Платёж проведён")

	body := fmt.Sprintf(`
		<h1>Спасибо за оплату!</h1>
		<p>Сумма: <strong>%.2f RUB</strong></p>
		<small>Это автоматическое уведомление</small>
	`, amount)
	m.SetBody("text/html", body)

	return s.dialer.DialAndSend(m)
}

func (s *EmailService) SendCreditReminder(to string, amount float64, dueDate string) error {
	m := mail.NewMessage()
	m.SetHeader("From", s.from)
	m.SetHeader("To", to)
	m.SetHeader("Subject", "Напоминание о платеже по кредиту")

	body := fmt.Sprintf(`
		<h1>Напоминание о платеже</h1>
		<p>Сумма платежа: <strong>%.2f RUB</strong></p>
		<p>Дата: %s</p>
		<small>Пожалуйста, пополните счёт</small>
	`, amount, dueDate)
	m.SetBody("text/html", body)

	return s.dialer.DialAndSend(m)
}
