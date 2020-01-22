package repository

import (
	"database/sql"
	"github.com/Marshality/tech-db/models"
	. "github.com/Marshality/tech-db/tools"
	"github.com/Marshality/tech-db/tools/queries"
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

func (ur *UserRepository) SelectWhere(nickname, email string) ([]*models.User, error) {
	var users []*models.User

	rows, err := ur.db.Query(queries.SelectUsersWhereNicknameOrEmail, nickname, email)

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

func (ur *UserRepository) Insert(user *models.User) error {
	return ur.db.QueryRow(queries.InsertIntoUsers,
		user.Nickname,
		user.Fullname,
		user.Email,
		user.About,
	).Scan(&user.ID)
}

func (ur *UserRepository) SelectByNickname(nickname string) (*models.User, error) {
	u := &models.User{}

	if err := ur.db.QueryRow(queries.SelectUserWhereNickname, nickname).Scan(
		&u.ID, &u.Nickname, &u.Fullname, &u.Email, &u.About); err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}

		return nil, err
	}

	return u, nil
}

func (ur *UserRepository) SelectByEmail(email string) (*models.User, error) {
	u := &models.User{}

	if err := ur.db.QueryRow(queries.SelectUserWhereEmail, email).Scan(
		&u.ID, &u.Nickname, &u.Fullname, &u.Email, &u.About); err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}

		return nil, err
	}

	return u, nil
}

func (ur *UserRepository) Update(u *models.User) error {
	res, err := ur.db.Exec(queries.UpdateUsers, u.About, u.Fullname, u.Email, u.Nickname)

	if err != nil {
		return err
	}

	count, err := res.RowsAffected()

	if err != nil {
		return err
	}

	if count == 0 {
		return ErrNotFound
	}

	return nil
}
