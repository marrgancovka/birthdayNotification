package employees

import (
	"birthdayNotification/internal/models"
	"context"
	"time"
)

type EmployeeUsecase interface {
	SignUp(ctx context.Context, user *models.User, userData *models.Employee) (string, time.Time, error)
	SignIn(ctx context.Context, user *models.User) (string, time.Time, error)
	Subscribe(ctx context.Context, from, to int) error
	Unsubscribe(ctx context.Context, from, to int) error
	EmployeesList(ctx context.Context) ([]models.Employee, error)
}

type EmployeeRepository interface {
	CreateUser(ctx context.Context, user *models.User, userData *models.Employee) (int, error)
	CheckUser(ctx context.Context, user *models.User) (int, error)
	CreateSubscription(ctx context.Context, from, to int) error
	CancelSubscription(ctx context.Context, from, to int) error
	GetEmploees(ctx context.Context) ([]models.Employee, error)
}
