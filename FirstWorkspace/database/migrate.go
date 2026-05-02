package database

import (
	"errors"
	"log"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func RunMigrations(dsn string) {
	m, err := migrate.New(
		"file://migrations",
		"postgres://"+dsn,
	)

	if err != nil {
		log.Fatal("Ошибка создании мигратора: ", err)
	}

	if err := m.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			log.Println("Новых миграций нет!")
			return
		}
		log.Fatal(err)
	}

	log.Println("Миграции выполнены!")
}
