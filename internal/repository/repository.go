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
		res, err := tx.ExecContext(ctx, "INSERT INTO users(login,email,password) VALUES($1, $2, $3)", user.Email, user.Login, user.Password)
		if err != nil {
			return err
		}

		user_id, err := res.LastInsertId()

		if err != nil {
			return err
		}

		_, err = tx.ExecContext(ctx, "INSERT INTO email_token(user_id, token)", user_id, token)

		return err
	})

	return token, err
}

func (s *Storage) VerifyUserTx(ctx context.Context, email, token string) error {
	err := s.ExecTx(ctx, func(tx *sql.Tx) error {
		_, err := s.db.ExecContext(ctx, "UPDATE users SET verify=true WHERE email=$1", email)
		if err != nil {
			return err
		}
		_, err = s.db.ExecContext(ctx, "DELETE FROM email_token WHERE token=$1", token)

		return err
	})

	return err
}

func (s *Storage) CreateSession(ctx context.Context, user_id int, token string) error {
	_, err := s.db.ExecContext(ctx, "INSERT INTO sessions(email, token) VALUES($1, $2)", user_id, token)
	return err
}

func (s *Storage) DeleteSession(ctx context.Context, token string) error {
	_, err := s.db.ExecContext(ctx, "DELETE FROM sessions WHERE token=$1", token)
	return err
}

func (s *Storage) FindPassword(ctx context.Context, email string) (int, string, error) {
	row := s.db.QueryRowContext(ctx, "SELECT id, password FROM users WHERE email=$1 AND verify=true", email)
	id, password := 0, ""
	err := row.Scan(&id, &password)
	return id, password, err
}

func (s *Storage) FindEmail(ctx context.Context, token string) (string, error) {
	row := s.db.QueryRowContext(ctx, "SELECT email FROM email_token WHERE token=$1", token)

	var email string
	err := row.Scan(&email)
	if err != nil {
		return "", err
	}

	return email, err
}
