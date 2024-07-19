package usecase

import (
	"birthdayNotification/internal/pkg/notifications/repo"
	"fmt"
	"github.com/sirupsen/logrus"
	"gopkg.in/gomail.v2"
	"os"
)

type Usecase struct {
	repo   repo.Repository
	logger *logrus.Logger
}

func NewUsecase(repo *repo.Repository, logger *logrus.Logger) *Usecase {
	return &Usecase{
		repo:   *repo,
		logger: logger,
	}
}

func (u *Usecase) SendTodayNotification() error {
	u.logger.Info("sending today notification")
	notifications, err := u.repo.SelectNotificationToday()
	if err != nil {
		u.logger.Error(err)
		return err
	}
	for _, notification := range notifications {
		u.sendEmail(notification.Email, notification.NameBirthdayPerson, notification.Count)
	}
	return nil
}

func (u *Usecase) sendEmail(emailAddress, birthdayPersons string, count int) {
	email := gomail.NewMessage()
	email.SetHeader("From", os.Getenv("EMAIL_USER"))
	email.SetHeader("To", emailAddress)
	email.SetHeader("Subject", "Напоминание о дне рождении!")
	if count > 1 {
		email.SetBody("text/html", fmt.Sprintf("Сегодня свой день рождения отмечают %v! Не забудьте поздравить их!", birthdayPersons))
	} else {
		email.SetBody("text/html", fmt.Sprintf("Сегодня свой день рождения отмечает %v! Не забудьте поздравить!", birthdayPersons))
	}
	d := gomail.NewDialer(os.Getenv("EMAIL_HOST"), 465, os.Getenv("EMAIL_USER"), os.Getenv("EMAIL_PASSWORD"))
	if err := d.DialAndSend(email); err != nil {
		u.logger.Error("Ошибка в отправке сообщения " + emailAddress + ": " + err.Error())
	}
	u.logger.Info("Отправлено сообщение для " + emailAddress)
}

func (u *Usecase) SendCongratulations() error {
	persons, err := u.repo.SelectBirthdayPersonsToday()
	if err != nil {
		return err
	}
	for _, person := range persons {
		u.sendCongratulation(person.FirstName, person.LastName, person.EmailAddress)
	}
	return nil
}

func (u *Usecase) sendCongratulation(name, surname, emailAddress string) {
	email := gomail.NewMessage()
	email.SetHeader("From", os.Getenv("EMAIL_USER"))
	email.SetHeader("To", emailAddress)
	email.SetHeader("Subject", "Поздравляем с днем рождения!")

	email.SetBody("text/html", fmt.Sprintf("%v %v, поздравляем вас с днем рождения!", name, surname))
	d := gomail.NewDialer(os.Getenv("EMAIL_HOST"), 465, os.Getenv("EMAIL_USER"), os.Getenv("EMAIL_PASSWORD"))
	if err := d.DialAndSend(email); err != nil {
		u.logger.Error("Ошибка в отправке поздравления " + emailAddress + ": " + err.Error())
	}
	u.logger.Info("Отправлено поздравление для " + emailAddress)
}
