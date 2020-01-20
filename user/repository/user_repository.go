package repository

import (
	"database/sql"
	"github.com/Marshality/tech-db/models"
	"github.com/Marshality/tech-db/user"
	"log"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(conn *sql.DB) user.Repository {
	return &UserRepository{
		db: conn,
	}
}

func (ur *UserRepository) GetUsersWhere(nickname, email string) ([]*models.User, error) {
	var users []*models.User

	rows, err := ur.db.Query("SELECT id, nickname, fullname, email, about " +
		"FROM users WHERE nickname = $1 OR email = $2", nickname, email)

	if err != nil {
		return nil, err
	}

	defer func() {
		if err := rows.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	for rows.Next() {
		u := &models.User{}

		if err := rows.Scan(&u.ID, &u.Nickname, &u.Fullname, &u.Email, &u.About); err != nil {
			return nil, err
		}

		users = append(users, u)
	}

	return users, nil
}

func (ur *UserRepository) Create(user *models.User) error {
	return ur.db.QueryRow("INSERT INTO users (nickname, fullname, email, about) VALUES ($1, $2, $3, $4) RETURNING id",
		user.Nickname,
		user.Fullname,
		user.Email,
		user.About,
	).Scan(&user.ID)
}
