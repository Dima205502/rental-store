package service

import (
	"context"
	"shop/internal/models"
	"time"
)

type thingManagerRepo interface {
	CreateThing(context.Context, models.Thing) error
	RemuveThing(context.Context, string, int) error
	BuyThingTx(context.Context, int, time.Time, string, string) error
}

type thingGetterRepo interface {
	AllThings(context.Context) ([]models.Thing, error)
	RentalThings(context.Context, string) ([]models.RentalThing, error)
	SaleThings(context.Context, string) ([]models.Thing, error)
}

type sender interface {
	Send(string, string) error
}
