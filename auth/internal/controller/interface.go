package controller

import (
	"auth_service/internal/models"
	"context"
)

type CreaterService interface {
	CreateUser(context.Context, models.User) error
}

type SessionService interface {
	DeleteSession(context.Context, string) error
	CreateSession(context.Context, models.User) (string, error)
	FindSession(context.Context, string, string) error
}

type CheckerService interface {
	CheckSend(context.Context, string) error
}
