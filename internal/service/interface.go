package service

import (
	"auth_service/internal/models"
	"context"
)

type RepoManager interface {
	CreateUserTokenTx(context.Context, models.User) (string, error)
	VerifyUserTx(context.Context, string, string) error
	CreateSession(context.Context, int, string) error
	DeleteSession(context.Context, string) error
}

type RepoFinder interface {
	FindPassword(context.Context, string) (int, string, error)
	FindEmail(context.Context, string) (string, error)
}
