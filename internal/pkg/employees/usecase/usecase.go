package usecase

import (
	"birthdayNotification/internal/models"
	"birthdayNotification/internal/pkg/employees/repo"
	"birthdayNotification/internal/utils"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"github.com/sirupsen/logrus"
	"time"
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

func (u *Usecase) SignUp(ctx context.Context, user *models.User, userData *models.Employee) (string, time.Time, error) {
	user.Password = hashString(user.Password)
	id, err := u.repo.CreateUser(ctx, user, userData)
	if err != nil {
		return "", time.Now(), err
	}
	token, exp, err := utils.GenerateToken(id)
	if err != nil {
		return "", time.Now(), err
	}
	return token, exp, nil
}
func (u *Usecase) SignIn(ctx context.Context, user *models.User) (string, time.Time, error) {
	user.Password = hashString(user.Password)
	id, err := u.repo.CheckUser(ctx, user)
	if err != nil {
		return "", time.Now(), err
	}
	token, exp, err := utils.GenerateToken(id)
	if err != nil {
		return "", time.Now(), err
	}
	return token, exp, nil
}
func (u *Usecase) Subscribe(ctx context.Context, from, to int) error {
	if err := u.repo.CreateSubscription(ctx, from, to); err != nil {
		u.logger.Errorf("Failed to create subscription: %v", err)
		return err
	}
	return nil
}
func (u *Usecase) Unsubscribe(ctx context.Context, from, to int) error {
	if err := u.repo.CancelSubscription(ctx, from, to); err != nil {
		u.logger.Errorf("Failed to canceld subscription: %v", err)
		return err
	}
	return nil

}
func (u *Usecase) EmployeesList(ctx context.Context) ([]models.Employee, error) {
	employees, err := u.repo.GetEmploees(ctx)
	if err != nil {
		u.logger.Errorf("Failed to get employees: %v", err)
		return nil, err
	}
	return employees, nil
}

func hashString(input string) string {
	hasher := sha256.New()
	hasher.Write([]byte(input))

	hashedBytes := hasher.Sum(nil)
	hashedString := hex.EncodeToString(hashedBytes)

	return hashedString
}
