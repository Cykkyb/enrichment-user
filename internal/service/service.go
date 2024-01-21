package service

import (
	"enrichment-user-info/internal/entity"
	"enrichment-user-info/internal/repository/postgres"
	"log/slog"
)

type Enrichment interface {
	EnrichUser(*entity.User) error
	CreateUser(*entity.User) error
	UpdateUser(*entity.User) error
	GetUser(int) (*entity.User, error)
	GetUsersWithFilter(*entity.UserFilter) ([]entity.User, error)
	DeleteUser(int) error
}

type Service struct {
	Enrichment
}

func NewService(repo *repository.Repository, log *slog.Logger) *Service {
	return &Service{
		Enrichment: NewEnrichmentService(repo.Enrichment, log),
	}
}
