package repo

import (
	"birthdayNotification/internal/models"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
)

type Repository struct {
	db     *sql.DB
	logger *logrus.Logger
}

func NewRepository(db *sql.DB, log *logrus.Logger) *Repository {
	return &Repository{
		db:     db,
		logger: log,
	}
}

func (r *Repository) CreateUser(ctx context.Context, user *models.User, userData *models.Employee) (int, error) {
	var userID int
	insertEmployee := `INSERT INTO employees(name, surname, birth_date) VALUES($1, $2, $3) RETURNING id`
	insertUser := `INSERT INTO users (email, password, id_employee) VALUES ($1, $2, $3)`
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		r.logger.Errorf("Error starting transaction: %v", err)
		return 0, errors.New("error starting transaction")
	}
	if err := tx.QueryRowContext(ctx, insertEmployee, userData.Name, userData.Surname, userData.BirthDate).Scan(&userID); err != nil {
		tx.Rollback()
		r.logger.Errorf("Error inserting user: %v", err)
		return 0, errors.New("error creating user")
	}
	if _, err := tx.ExecContext(ctx, insertUser, user.Email, user.Password, userID); err != nil {
		tx.Rollback()
		r.logger.Errorf("Error inserting employee: %v", err)
		return 0, errors.New("error creating user")
	}
	err = tx.Commit()
	if err != nil {
		r.logger.Errorf("Error committing transaction: %v", err)
		return 0, errors.New("error committing user")
	}
	return userID, nil
}

func (r *Repository) CheckUser(ctx context.Context, user *models.User) (int, error) {
	var userID int
	var password string
	query := `SELECT id, password FROM users WHERE email = $1`
	if err := r.db.QueryRowContext(ctx, query, user.Email).Scan(&userID, &password); err != nil {
		r.logger.Errorf("Error checking user: %v", err)
		return 0, errors.New("user not found")
	}
	if password != user.Password {
		r.logger.Error("Error checking user: wrong password")
		return 0, errors.New("invalid password")
	}
	return userID, nil

}
func (r *Repository) CreateSubscription(ctx context.Context, from, to int) error {
	query := `INSERT INTO subscriptions (id_from, id_to) VALUES ($1, $2)`
	fmt.Println(from, to)
	_, err := r.db.ExecContext(ctx, query, from, to)
	if err != nil {
		r.logger.Errorf("Error inserting subscription: %v", err)
		return errors.New("error inserting subscription")
	}
	return err
}
func (r *Repository) CancelSubscription(ctx context.Context, from, to int) error {
	query := `DELETE FROM subscriptions WHERE id_from = $1 AND id_to = $2`
	_, err := r.db.ExecContext(ctx, query, from, to)
	if err != nil {
		r.logger.Errorf("Error cancelling subscription: %v", err)
	}
	return err
}
func (r *Repository) GetEmploees(ctx context.Context) ([]models.Employee, error) {
	query := `SELECT id, name, surname, birth_date FROM employees`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		r.logger.Errorf("Error getting employees: %v", err)
		return nil, err
	}
	defer rows.Close()

	var employees []models.Employee
	for rows.Next() {
		var employee models.Employee
		if err := rows.Scan(&employee.Id, &employee.Name, &employee.Surname, &employee.BirthDate); err != nil {
			r.logger.Errorf("Error getting employees: %v", err)
			return nil, err
		}
		employees = append(employees, employee)
	}
	if err := rows.Err(); err != nil {
		r.logger.Errorf("Error getting employees: %v", err)
		return nil, err
	}
	return employees, nil
}
