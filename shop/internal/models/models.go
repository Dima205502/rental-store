package models

type Thing struct {
	Owner       string `json:"owner"`
	Type        string `json:"type" validate:"required"`
	Description string `json:"description"`
	Price       int    `json:"price" validate:"required"`
	Available   bool   `json:"available"`
}

type RentalThing struct {
	ThingId    int    `json:"thing_id"`
	Buyer      string `json:"buyer"`
	Email      string `json:"email"`
	FinishTime string `json:"finish_time"`
}

type ThingId struct {
	ThingId int `json:"thing_id" validate:"required"`
}

type BuyThingRequest struct {
	ThingId      int `json:"thing_id" validate:"required"`
	TimeInterval `json:"time_interval" validate:"required"`
}

type TimeInterval struct {
	Months int `json:"months"`
	Days   int `json:"days"`
	Hours  int `json:"hours"`
}
