package repository

import (
	"auth_service/internal/models"
	"auth_service/utils"
	"context"
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

func NewStorage() *Storage {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		"localhost", 5432, "KinDeR", "Dimaaaa", "draft")
	var err error

	DB, err := sql.Open("postgres", psqlInfo)

	if err != nil {
		panic("coudn't connect to database")
	}

	err = DB.Ping()
	if err != nil {
		log.Fatal("Couldn't ping database:", err)
	}

	return &Storage{DB}
}

func (s *Storage) ExecTx(ctx context.Context, fn func(*sql.Tx) error) error {
	tx, err := s.db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelReadCommitted, ReadOnly: false})
	if err != nil {
		return err
	}

	err = fn(tx)

	if err == nil {
		err = tx.Commit()
	} else {
		tx.Rollback()
	}

	return err
}

func (s *Storage) CreateUserTokenTx(ctx context.Context, user models.User) (string, error) {
	token, err := utils.GenerateToken()

	if err != nil {
		return "", err
	}

	err = s.ExecTx(ctx, func(tx *sql.Tx) error {
		row := tx.QueryRowContext(ctx, "INSERT INTO users(nickname,email,password) VALUES($1, $2, $3) RETURNING id;", user.Nickname, user.Email, user.Password)

		var id int
		err = row.Scan(&id)

		if err != nil {
			return err
		}

		_, err = tx.ExecContext(ctx, "INSERT INTO email_token(user_id, token) VALUES($1, $2);", id, token)

		return err
	})

	return token, err
}

func (s *Storage) VerifyUserTx(ctx context.Context, email, token string) error {
	err := s.ExecTx(ctx, func(tx *sql.Tx) error {
		_, err := s.db.ExecContext(ctx, "UPDATE users SET verify=true WHERE email=$1;", email)
		if err != nil {
			return err
		}
		_, err = s.db.ExecContext(ctx, "DELETE FROM email_token WHERE token=$1;", token)

		return err
	})

	return err
}

func (s *Storage) CreateSession(ctx context.Context, user_id int, token string) error {
	_, err := s.db.ExecContext(ctx, "INSERT INTO sessions(user_id, token) VALUES($1, $2);", user_id, token)
	return err
}

func (s *Storage) DeleteSession(ctx context.Context, token string) error {
	_, err := s.db.ExecContext(ctx, "DELETE FROM sessions WHERE token=$1;", token)
	return err
}

func (s *Storage) FindPassword(ctx context.Context, email string) (int, string, error) {
	row := s.db.QueryRowContext(ctx, "SELECT id, password FROM users WHERE email = $1 AND verify=TRUE;", email)
	id, password := 0, ""
	err := row.Scan(&id, &password)
	return id, password, err
}

func (s *Storage) FindUserId(ctx context.Context, token string) (int, error) {

	row := s.db.QueryRowContext(ctx, "SELECT user_id FROM email_token WHERE token = $1;", token)

	var user_id int
	err := row.Scan(&user_id)

	return user_id, err
}

func (s *Storage) FindEmail(ctx context.Context, id int) (string, error) {
	row := s.db.QueryRowContext(ctx, "SELECT email FROM users WHERE id = $1;", id)

	var email string
	err := row.Scan(&email)

	return email, err
}
