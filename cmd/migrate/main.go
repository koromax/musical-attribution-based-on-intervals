package main

import (
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"interval_attribution/internal/app/config"
	"interval_attribution/internal/app/ds"
	"interval_attribution/internal/app/dsn"
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("Не удалось загрузить конфигурацию: %v", err)
	}

	db, err := gorm.Open(postgres.Open(dsn.FromConfig(cfg.Database)), &gorm.Config{})
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
