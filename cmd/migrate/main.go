package main

import (
	"log"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"interval_attribution/internal/app/ds"
	"interval_attribution/internal/app/dsn"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("Предупреждение: .env файл не найден, используются стандартные переменные")
	}

	db, err := gorm.Open(postgres.Open(dsn.FromEnv()), &gorm.Config{})
	if err != nil {
		log.Fatalf("Ошибка подключения к БД: %v", err)
	}

	err = db.AutoMigrate(
		&ds.User{},
		&ds.Composer{},
		&ds.ComposerLike{},
	)
	if err != nil {
		log.Fatalf("Ошибка выполнения миграции: %v", err)
	}

	log.Println("Миграция успешно завершена")
}
