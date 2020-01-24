package repository

import (
	"database/sql"
	"github.com/Marshality/tech-db/service"
	"github.com/Marshality/tech-db/tools/queries"
)

type ServiceRepository struct {
	db *sql.DB
}

func NewServiceRepository(conn *sql.DB) service.Repository {
	return &ServiceRepository{
		db: conn,
	}
}

func (sr *ServiceRepository) Status() (uint64, uint64, uint64, uint64, error) {
	var forumsStatus, usersStatus, threadsStatus, postsStatus uint64

	err := sr.db.QueryRow(queries.StatusQuery).Scan(&forumsStatus, &threadsStatus, &postsStatus, &usersStatus)

	return forumsStatus, threadsStatus, postsStatus, usersStatus, err
}

func (sr *ServiceRepository) Clear() error {
	_, err := sr.db.Exec(queries.ClearQuery)

	return err
}
