package repository

import (
	"database/sql"

	"github.com/CVWO-Backend/internal/models"
)

type DatabaseRepo interface {
	Connection() *sql.DB
	AllThreads() ([]*models.Thread, error)
}