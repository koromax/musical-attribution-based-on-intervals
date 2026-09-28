package handler

import "interval_attribution/internal/app/repository"

type Handler struct {
	Repo *repository.Repository
}

func NewHandler(repo *repository.Repository) *Handler {
	return &Handler{
		Repo: repo,
	}
}
