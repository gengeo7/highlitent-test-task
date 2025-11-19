package main

import (
	"fmt"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/gengeo7/highlitent/config"
	"github.com/gengeo7/highlitent/controllers/answers"
	"github.com/gengeo7/highlitent/controllers/questions"
	"github.com/gengeo7/highlitent/logger"
	"github.com/gengeo7/highlitent/middleware"
	"github.com/gengeo7/highlitent/storage/gormdb"
)

func main() {
	err := config.Initialize()
	if err != nil {
		fmt.Println(err)
		return
	}

	logger.Init()
	defer func() {
		if err := recover(); err != nil {
			logger.Error("PANIC", "error", err, "stacktrace", string(debug.Stack()))
		}
	}()

	db := gormdb.NewDb()
	dsnConfig := gormdb.DsnConfig{
		Host:     config.Conf.PostgresHost,
		Port:     config.Conf.PostgresPort,
		User:     config.Conf.PostgresUser,
		Password: config.Conf.PostgresPassword,
		Database: config.Conf.PostgresDatabase,
	}
	err = db.Open(&dsnConfig)
	if err != nil {
		logger.Error(err.Error())
		return
	}
	err = db.Migrate(config.Conf.MigrationPath)
	if err != nil {
		logger.Error(err.Error())
		return
	}

	mux := http.NewServeMux()
	answersController := answers.NewAnswersController(db)
	answersController.RegisterController(mux)
	questionsController := questions.NewQuestionsController(db)
	questionsController.RegisterController(mux)

	handler := middleware.Log(
		middleware.TimeElapsed(
			middleware.Recoverer(mux),
		),
	)

	server := http.Server{
		Addr:         fmt.Sprintf("%s:%d", config.Conf.Host, config.Conf.Port),
		Handler:      handler,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  600 * time.Second,
	}

	logger.Info("Ready to go...")
	if err := server.ListenAndServe(); err != nil {
		logger.Error("Startup error", "error", err)
		return
	}
}
