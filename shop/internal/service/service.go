package service

import (
	"context"
	"fmt"
	"log/slog"
	"shop/internal/models"
	"shop/internal/repository"
	"time"
)

func NewThingManagerService(storage *repository.Storage, notifier *repository.Notifier) *ThingManagerService {
	return &ThingManagerService{thingManagerRepo: storage, notifier: notifier}
}

func NewThingGetterService(storage *repository.Storage) *ThingGetterService {
	return &ThingGetterService{thingGetterRepo: storage}
}

func (t *ThingManagerService) CreateThing(ctx context.Context, thing models.Thing) error {
	return t.thingManagerRepo.CreateThing(ctx, thing)
}

func (t *ThingManagerService) RemuveThing(ctx context.Context, nickname string, thingId int) error {
	return t.thingManagerRepo.RemuveThing(ctx, nickname, thingId)
}

func (t *ThingManagerService) BuyThing(ctx context.Context, buyThng models.BuyThingRequest, nickname, email string) error {
	thingId := buyThng.ThingId

	months := buyThng.Months
	days := buyThng.Days
	hours := buyThng.Hours

	finishTime := time.Now().AddDate(0, months, days).Add(time.Duration(hours) * time.Hour)

	err := t.thingManagerRepo.BuyThingTx(ctx, thingId, finishTime, nickname, email)

	if err != nil {
		slog.Error("BuyThing", slog.String("error", err.Error()))
		return err
	}

	go func() {
		msg := fmt.Sprintf("Subject: Purchase Notification\nYou rented an thing with the code %d, the rental deadline is %s", thingId, finishTime.Format("2006-01-02 15:04:05"))
		slog.Info("Before Send", slog.String("text", msg))
		err := t.notifier.Send(email, msg)

		if err != nil {
			slog.Error("BuyThing", slog.String("place", "Send"), slog.String("error", err.Error()))
		}
	}()

	return err
}

func (t *ThingGetterService) GetAllThings(ctx context.Context) ([]models.Thing, error) {
	return t.thingGetterRepo.AllThings(ctx)
}

func (t *ThingGetterService) GetRentalThings(ctx context.Context, nickname string) ([]models.RentalThing, error) {
	return t.thingGetterRepo.RentalThings(ctx, nickname)
}

func (t *ThingGetterService) GetSaleThings(ctx context.Context, nickname string) ([]models.Thing, error) {
	return t.thingGetterRepo.SaleThings(ctx, nickname)
}
