package repository

import (
	"fmt"

	"interval_attribution/internal/app/ds"
)

func (r *Repository) CreateUser(login, password string) (ds.User, error) {
	var count int64
	r.db.Model(&ds.User{}).Where("login = ?", login).Count(&count)
	if count > 0 {
		return ds.User{}, fmt.Errorf("пользователь с таким логином уже существует")
	}

	user := ds.User{
		Login:       login,
		Password:    password,
		IsModerator: false,
	}

	err := r.db.Create(&user).Error
	if err != nil {
		return ds.User{}, fmt.Errorf("ошибка создания пользователя: %w", err)
	}

	return user, nil
}
