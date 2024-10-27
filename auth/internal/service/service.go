package service

import (
	"auth_service/internal/models"
	"auth_service/internal/repository"
	"auth_service/utils"
	"context"
	"errors"
	"log/slog"
)

func NewUserManager(storage *repository.Storage) *UserManager {
	return &UserManager{createrRepo: storage}
}

func NewSessionManager(storage *repository.Storage) *SessionManager {
	return &SessionManager{sessionRepo: storage, finderRepo: storage}
}

func NewEmailManager(storage *repository.Storage) *EmailManager {
	return &EmailManager{finderRepo: storage, userTokenRepo: storage}
}

func (u UserManager) CreateUser(ctx context.Context, user models.User) error {
	hashedPassword, err := utils.Hashed(user.Password)

	if err != nil {
		slog.Error("Service layer", slog.String("place", "Hashed"), slog.String("error", err.Error()))
		return err
	}

	user.Password = hashedPassword

	token, err := u.createrRepo.CreateUserTokenTx(ctx, user)

	if err != nil {
		slog.Error("Service layer", slog.String("place", "CreateUserTokenTx"), slog.String("error", err.Error()))
		return err
	}

	msg := "Subject: Verify Email\nClick on the link to confirm your email\nhttp://localhost:8080/verify-email?token=" + token

	err = utils.Send(ctx, user.Email, msg)

	if err != nil {
		slog.Error("Service layer", slog.String("place", "Send"), slog.String("error", err.Error()))
		return err
	}

	return nil
}

func (s *SessionManager) CreateSession(ctx context.Context, user models.User) (string, error) {
	id, password, err := s.finderRepo.FindPassword(ctx, user.Email)

	if err != nil {
		slog.Error("Service layer", slog.String("place", "FindPassword"), slog.String("error", err.Error()))
		return "", err
	}

	err = utils.ComparePassword(password, user.Password)

	if err != nil {
		return "", errors.New("wrong password")
	}

	token, err := utils.GenerateToken()

	if err != nil {
		slog.Error("Service layer", slog.String("place", "GenerateToken"), slog.String("error", err.Error()))
		return "", err
	}

	err = s.sessionRepo.CreateSession(ctx, id, token)
	if err != nil {
		slog.Error("Service layer", slog.String("place", "CreateSession"), slog.String("error", err.Error()))
		return "", err
	}

	return token, err
}

func (s *SessionManager) DeleteSession(ctx context.Context, token string) error {
	slog.Debug(token)
	err := s.sessionRepo.DeleteSession(ctx, token)
	if err != nil {
		slog.Error("Service layer", slog.String("place", "DeleteSession"), slog.String("error", err.Error()))
	}
	return err
}

func (e *EmailManager) CheckSend(ctx context.Context, token string) error {
	user_id, err := e.finderRepo.FindUserId(ctx, token)

	if err != nil {
		slog.Error("Service layer", slog.String("place", "FindUserId"), slog.String("error", err.Error()))
		return err
	}

	email, err := e.finderRepo.FindEmail(ctx, user_id)

	if err != nil {
		slog.Error("Service layer", slog.String("place", "FindUserId"), slog.String("error", err.Error()))
		return err
	}

	err = e.userTokenRepo.VerifyUserTx(ctx, email, token)

	if err != nil {
		slog.Error("Service layer", slog.String("place", "VerifyUserTx"), slog.String("error", err.Error()))
		return err
	}

	msg := "Subject: Create Account\nYou have successfully registered!"

	err = utils.Send(ctx, email, msg)

	if err != nil {
		slog.Error("Service layer", slog.String("place", "Send"), slog.String("error", err.Error()))
		return err
	}

	return err
}
