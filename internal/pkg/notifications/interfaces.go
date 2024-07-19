package notifications

import (
	"birthdayNotification/internal/models"
)

type NotificationUsecase interface {
	SendTodayNotification() error
	SendCongratulations() error
}
type NotificationRepository interface {
	SelectNotificationToday() ([]*models.Notification, error)
	SelectBirthdayPersonsToday() ([]*models.BirthdayPerson, error)
}
