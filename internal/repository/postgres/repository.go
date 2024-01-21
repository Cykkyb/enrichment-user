package repository

import (
	"enrichment-user-info/internal/entity"
	"github.com/jmoiron/sqlx"
)

type Enrichment interface {
	CreateUser(user *entity.User) error
	GetUser(id int) (*entity.User, error)
	GetUsersWithFilter(filter *entity.UserFilter) ([]entity.User, error)
	UpdateUser(*entity.User) error
	DeleteUser(id int) error
}

type Repository struct {
	Enrichment
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Enrichment: NewEnrichmentPostgres(db),
	}
}
