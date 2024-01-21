package main

import (
	"enrichment-user-info/internal/app"
	"enrichment-user-info/internal/config"
	"enrichment-user-info/internal/handler"
	"enrichment-user-info/internal/lib/logger"
	"enrichment-user-info/internal/repository/postgres"
	"enrichment-user-info/internal/service"
	"fmt"
	"github.com/gin-gonic/gin"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	cfg := config.MustLoadConfig()

	log := initLogger()

	log.Info("Config loaded",
		slog.String("port", cfg.App.Port),
		slog.String("env", cfg.App.Env),
	)

	gin.SetMode(cfg.App.Env)

	db, err := repository.ConnectDb(repository.Config{
		Host:     cfg.DB.Host,
		Port:     cfg.DB.Port,
		Username: cfg.DB.Username,
		Password: cfg.DB.Password,
		DBname:   cfg.DB.DBname,
		SSL:      cfg.DB.SSL,
	})
	if err != nil {
		panic(fmt.Sprintf("failed to connect to db: %s", err))
	}

	repo := repository.NewRepository(db)
	services := service.NewService(repo, log)
	handlers := handler.NewHandler(services, log)
	application := app.NewApp(log)

	go application.Run(handlers.InitRoutes(), cfg.App.Port)

	stop := make(chan os.Signal)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	<-stop

	application.Stop()
}

func initLogger() *slog.Logger {
	opts := logger.PrettyHandlerOptions{
		SlogOpts: slog.HandlerOptions{
			Level: slog.LevelDebug,
		},
	}
	handler := logger.NewPrettyHandler(os.Stdout, opts)
	log := slog.New(handler)

	return log
}
