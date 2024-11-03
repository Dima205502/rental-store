package service

type ThingManagerService struct {
	thingManagerRepo thingManagerRepo
	notifier         sender
}

type ThingGetterService struct {
	thingGetterRepo thingGetterRepo
}
