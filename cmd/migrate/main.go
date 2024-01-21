package main

import (
	"enrichment-user-info/internal/config"
	"enrichment-user-info/internal/repository/postgres"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

func main() {
	cfg := config.MustLoadConfig()

	db, err := repository.ConnectDb(repository.Config{
		Host:     cfg.DB.Host,
		Port:     cfg.DB.Port,
		Username: cfg.DB.Username,
		Password: cfg.DB.Password,
		DBname:   cfg.DB.DBname,
		SSL:      cfg.DB.SSL,
	})
	if err != nil {
		fmt.Println("Ошибка подключения к базе данных:", err)
		return
	}
	fmt.Println(cfg.DB.Host, cfg.DB.Port, cfg.DB.Username, cfg.DB.Password, cfg.DB.DBname, cfg.DB.SSL)

	driver, err := postgres.WithInstance(db.DB, &postgres.Config{})
	if err != nil {
		fmt.Println("Ошибка создания экземпляра драйвера базы данных:", err)
		return
	}

	// Путь к директории с миграциями
	migrationsDir := "./migrations"

	// Создание экземпляра объекта Migrate
	m, err := migrate.NewWithDatabaseInstance(fmt.Sprintf("file://%s", migrationsDir), "postgres", driver)
	if err != nil {
		fmt.Println("Ошибка создания экземпляра объекта Migrate:", err)
		return
	}

	// Применение всех миграций
	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		fmt.Println("Ошибка применения миграций:", err)
		return
	}

	fmt.Println("Все миграции успешно применены")
}
