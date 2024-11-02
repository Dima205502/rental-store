package controller

import (
	"context"
	"shop/models"
)

type thingManager interface {
	createThing(context.Context, models.Thing) error
	remuveThing()
	buyThing()
}

type thingGetter interface {
	getAllThings()
	getRentalThings()
	getSaleThing()
}
