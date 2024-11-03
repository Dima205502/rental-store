package server

import (
	"log/slog"
	"net/http"
	"shop/config"
	"shop/internal/controller"
	"shop/internal/repository"
	"shop/internal/service"
	"strconv"
)

func Start(cfg *config.Config) error {
	slog.Debug("AuthManager created")

	storage := repository.NewStorage(cfg.DB)
	notifier := repository.NewNotifier(cfg)

	thingGetter := service.NewThingGetterService(storage)
	thingManager := service.NewThingManagerService(storage, notifier)

	manager := controller.NewShopManager(thingManager, thingGetter)

	http.HandleFunc("/add-thing", controller.AddThing(manager))
	http.HandleFunc("/remuve-thing", controller.RemuveThing(manager))
	http.HandleFunc("/buy-thing", controller.BuyThing(manager))
	http.HandleFunc("/show-all-things", controller.ShowAllThings(manager))
	http.HandleFunc("/show-rental-things", controller.ShowRentalThings(manager))
	http.HandleFunc("/show-sale-things", controller.ShowSaleThings(manager))

	slog.Debug("Starting the server")

	return http.ListenAndServe(cfg.Server.Host+":"+strconv.Itoa(cfg.Server.Port), nil)
}
