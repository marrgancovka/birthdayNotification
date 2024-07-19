package delivery

import (
	"birthdayNotification/internal/middleware"
	"birthdayNotification/internal/models"
	"birthdayNotification/internal/pkg/employees/usecase"
	"birthdayNotification/internal/utils"
	"errors"
	"github.com/sirupsen/logrus"
	"net/http"
)

type Handler struct {
	u      usecase.Usecase
	logger *logrus.Logger
}

func NewHandler(u *usecase.Usecase, logger *logrus.Logger) *Handler {
	return &Handler{
		u:      *u,
		logger: logger,
	}
}

// SignUp godoc
// @Summary Sign up a new user
// @Description Create a new user and employee profile
// @Tags auth
// @Accept json
// @Produce json
// @Param signup body models.SignUpRequest true "Sign Up Request"
// @Success 201 {string} string "Token"
// @Failure 400 {object} string "Invalid request"
// @Router /auth/signup [post]
func (h *Handler) SignUp(w http.ResponseWriter, r *http.Request) {
	user := &models.User{}
	employee := &models.Employee{}
	signup := models.SignUpRequest{}
	if err := utils.ReadRequestData(r, &signup); err != nil {
		h.logger.Error(err)
		utils.WriteError(w, http.StatusBadRequest, errors.New("invalid request"))
		return
	}
	user.Password = signup.Password
	user.Email = signup.Email
	employee.Surname = signup.Surname
	employee.Name = signup.Name
	employee.BirthDate = signup.BirthDate

	token, exp, err := h.u.SignUp(r.Context(), user, employee)
	if err != nil {
		h.logger.Error(err)
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}
	h.logger.Info("User signed up successfully")
	http.SetCookie(w, utils.TokenCookie(middleware.CookieName, token, exp))
	utils.WriteJSON(w, http.StatusCreated, token)
}

// SignIn godoc
// @Summary Sign in a user
// @Description Authenticate a user and return a token
// @Tags auth
// @Accept json
// @Produce json
// @Param user body models.User true "User credentials"
// @Success 200 {string} string "Token"
// @Failure 400 {object} string "Invalid request"
// @Router /auth/signin [post]
func (h *Handler) SignIn(w http.ResponseWriter, r *http.Request) {
	user := &models.User{}
	if err := utils.ReadRequestData(r, user); err != nil {
		utils.WriteError(w, http.StatusBadRequest, errors.New("invalid request"))
		h.logger.Error(err)
		return
	}
	token, exp, err := h.u.SignIn(r.Context(), user)
	if err != nil {
		h.logger.Error(err)
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}
	h.logger.Info("User signed in successfully")
	http.SetCookie(w, utils.TokenCookie(middleware.CookieName, token, exp))
	utils.WriteJSON(w, http.StatusOK, token)
}

// Subscribe godoc
// @Summary Subscribe to employee notifications
// @Description Subscribe to notifications for an employee's birthday
// @Security ApiKeyAuth
// @Tags employee
// @Accept json
// @Produce json
// @Param subscription body models.SubscriptionRequest true "Subscription details"
// @Success 200 {string} string "subscription subscribed"
// @Failure 400 {object} string "Invalid request"
// @Router /employee/subscribe [post]
func (h *Handler) Subscribe(w http.ResponseWriter, r *http.Request) {
	subscription := &models.Subscription{}
	from, ok := r.Context().Value(middleware.CookieName).(int)
	if !ok {
		h.logger.Error("Failed to get cookie from")
		utils.WriteError(w, http.StatusBadRequest, errors.New("no cookie value"))
		return
	}

	if err := utils.ReadRequestData(r, &subscription); err != nil {
		h.logger.Error(err)
		utils.WriteError(w, http.StatusBadRequest, errors.New("invalid request"))
		return
	}
	subscription.IdFrom = from
	if err := h.u.Subscribe(r.Context(), subscription.IdFrom, subscription.IdTo); err != nil {
		h.logger.Error(err)
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}
	h.logger.Info("Subscription subscribed")
	utils.WriteJSON(w, http.StatusOK, "subscription subscribed")
}

// Unsubscribe godoc
// @Summary Unsubscribe from employee notifications
// @Description Unsubscribe from notifications for an employee's birthday
// @Security ApiKeyAuth
// @Tags employee
// @Accept json
// @Produce json
// @Param subscription body models.SubscriptionRequest true "Subscription details"
// @Success 200 {string} string "subscription unsubscribed"
// @Failure 400 {object} string "Invalid request"
// @Router /employee/unsubscribe [post]
func (h *Handler) Unsubscribe(w http.ResponseWriter, r *http.Request) {
	subscription := &models.Subscription{}
	from, ok := r.Context().Value(middleware.CookieName).(int)
	if !ok {
		h.logger.Error("Failed to get cookie from")
		utils.WriteError(w, http.StatusBadRequest, errors.New("no cookie value"))
		return
	}

	if err := utils.ReadRequestData(r, &subscription); err != nil {
		h.logger.Error(err)
		utils.WriteError(w, http.StatusBadRequest, errors.New("invalid request"))
		return
	}
	subscription.IdFrom = from
	if err := h.u.Unsubscribe(r.Context(), subscription.IdFrom, subscription.IdTo); err != nil {
		h.logger.Error(err)
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}
	h.logger.Info("Subscription unsubscribed")
	utils.WriteJSON(w, http.StatusOK, "subscription unsubscribed")
}

// GetEmployeesList godoc
// @Summary Get list of employees
// @Description Retrieve a list of all employees
// @Security ApiKeyAuth
// @Tags employee
// @Produce json
// @Success 200 {array} models.Employee
// @Failure 500 {object} string "Internal server error"
// @Router /employee/list [get]
func (h *Handler) GetEmployeesList(w http.ResponseWriter, r *http.Request) {
	employees, err := h.u.EmployeesList(r.Context())
	if err != nil {
		h.logger.Error(err)
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, employees)
}
