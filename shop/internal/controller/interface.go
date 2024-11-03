package controller

import (
	"context"
	"shop/internal/models"
)

type thingManager interface {
	CreateThing(context.Context, models.Thing) error
	RemuveThing(context.Context, string, int) error
	BuyThing(context.Context, models.BuyThingRequest, string, string) error
}

type thingGetter interface {
	GetAllThings(context.Context) ([]models.Thing, error)
	GetRentalThings(context.Context, string) ([]models.RentalThing, error)
	GetSaleThings(context.Context, string) ([]models.Thing, error)
}
