package repo

import (
	"birthdayNotification/internal/models"
	"database/sql"
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

func (repo *Repository) SelectNotificationToday() ([]*models.Notification, error) {
	selectQuery := `SELECT
    u1.email,
    STRING_AGG(e2.name || ' ' || e2.surname, ', ') AS subscribed_to,
    COUNT(e2.id) AS birthday_count
FROM
    subscriptions s
        JOIN
    users u1 ON s.id_from = u1.id
        JOIN
    users u2 ON s.id_to = u2.id
        JOIN
    employees e2 ON u2.id_employee = e2.id
WHERE
    DATE_PART('day', e2.birth_date) = DATE_PART('day', CURRENT_DATE)
  AND DATE_PART('month', e2.birth_date) = DATE_PART('month', CURRENT_DATE)
GROUP BY
    u1.email;
`
	notifications := make([]*models.Notification, 0)
	rows, err := repo.db.Query(selectQuery)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		notification := &models.Notification{}
		if err = rows.Scan(&notification.Email, &notification.NameBirthdayPerson, &notification.Count); err != nil {
			repo.logger.Error(err)
			return nil, err
		}
		notifications = append(notifications, notification)
	}

	return notifications, nil
}

func (repo *Repository) SelectBirthdayPersonsToday() ([]*models.BirthdayPerson, error) {
	selectQuery := `SELECT e.name, e.surname, u.email FROM employees e
	JOIN users AS u ON e.id = u.id_employee
	WHERE DATE_PART('day', e.birth_date) = DATE_PART('day', CURRENT_DATE)
  AND DATE_PART('month', e.birth_date) = DATE_PART('month', CURRENT_DATE);`
	rows, err := repo.db.Query(selectQuery)
	if err != nil {
		repo.logger.Error("SelectBirthdayPersonsToday ", err.Error())
		return nil, err
	}
	defer rows.Close()

	birthdayPersons := make([]*models.BirthdayPerson, 0)
	for rows.Next() {
		birthdayPerson := &models.BirthdayPerson{}
		if err = rows.Scan(&birthdayPerson.FirstName, &birthdayPerson.LastName, &birthdayPerson.EmailAddress); err != nil {
			repo.logger.Error("SelectBirthdayPersonsToday ", err.Error())
			return nil, err
		}
		birthdayPersons = append(birthdayPersons, birthdayPerson)
	}
	repo.logger.Info("success select today birthday persons")
	return birthdayPersons, nil
}
