package app

import (
	"birthdayNotification/internal/middleware"
	"birthdayNotification/internal/pkg/employees/delivery"
	repoE "birthdayNotification/internal/pkg/employees/repo"
	usecaseE "birthdayNotification/internal/pkg/employees/usecase"
	repoN "birthdayNotification/internal/pkg/notifications/repo"
	usecaseN "birthdayNotification/internal/pkg/notifications/usecase"
	"context"
	"database/sql"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	httpSwagger "github.com/swaggo/http-swagger"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type App struct {
	logger *logrus.Logger
}

func NewApp(logger *logrus.Logger) *App {
	return &App{
		logger: logger,
	}
}

func (app *App) Start() error {
	app.logger.Infof("starting server")

	_ = godotenv.Load(".env")
	db, err := sql.Open("postgres", fmt.Sprintf("postgres://%v:%v@%v:%v/%v?sslmode=disable",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASS"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME")))
	if err != nil {
		app.logger.Error("error connecting to database: " + err.Error())
		return err
	}
	if err = db.Ping(); err != nil {
		app.logger.Error("failed to ping database" + err.Error())
		return err
	}
	defer db.Close()

	r := mux.NewRouter().PathPrefix("/api").Subrouter()

	repoEmployees := repoE.NewRepository(db, app.logger)
	ucEmployyes := usecaseE.NewUsecase(repoEmployees, app.logger)
	hEmployees := delivery.NewHandler(ucEmployyes, app.logger)

	r.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	auth := r.PathPrefix("/auth").Subrouter()
	auth.HandleFunc("/signup", hEmployees.SignUp).Methods(http.MethodPost)
	auth.HandleFunc("/signin", hEmployees.SignIn).Methods(http.MethodPost)

	employee := r.PathPrefix("/employee").Subrouter()
	employee.Handle("/subscribe", middleware.Auth(http.HandlerFunc(hEmployees.Subscribe), app.logger)).Methods(http.MethodPost)
	employee.Handle("/unsubscribe", middleware.Auth(http.HandlerFunc(hEmployees.Unsubscribe), app.logger)).Methods(http.MethodPost)
	employee.Handle("/list", middleware.Auth(http.HandlerFunc(hEmployees.GetEmployeesList), app.logger)).Methods(http.MethodGet)

	repoNotification := repoN.NewRepository(db, app.logger)
	usecaseNotification := usecaseN.NewUsecase(repoNotification, app.logger)

	srv := http.Server{
		Addr:         ":8080",
		Handler:      r,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		app.logger.Info("Start server on ", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			app.logger.Error("Error in listen: ", err.Error())
		}
	}()
	go func() {
		for {
			if err = usecaseNotification.SendTodayNotification(); err != nil {
				app.logger.Error("Error in sending today notification: ", err.Error())
			}
			time.Sleep(24 * time.Hour)
		}
	}()
	go func() {
		for {
			if err = usecaseNotification.SendCongratulations(); err != nil {
				app.logger.Error("Error in sending today congratulations: ", err.Error())
			}
			time.Sleep(24 * time.Hour)
		}
	}()

	sig := <-signalCh
	app.logger.Info("Received signal: ", sig)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		app.logger.Fatal("Server shutdown failed: ", err.Error())
	}
	return nil
}
