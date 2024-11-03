package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"shop/config"
	"shop/internal/models"
	"time"

	"github.com/IBM/sarama"
	_ "github.com/lib/pq"
)

func NewStorage(cfg config.DB) *Storage {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Dbname)
	var err error

	DB, err := sql.Open("postgres", psqlInfo)

	if err != nil {
		log.Fatal("Couldn't ping database:", err)
	}

	err = DB.Ping()
	if err != nil {
		log.Fatal("Couldn't ping database:", err)
	}

	return &Storage{DB}
}

func NewNotifier(cfg *config.Config) *Notifier {
	config := sarama.NewConfig()

	config.Producer.RequiredAcks = sarama.NoResponse
	config.Producer.Retry.Max = 0
	config.Producer.Return.Successes = true

	producer, err := sarama.NewSyncProducer(cfg.Broker_addrs, config)

	if err != nil {
		log.Fatal("Couldn't create producer:", err)
	}

	return &Notifier{producer, cfg.Topic}
}

func (s *Storage) ExecTx(ctx context.Context, fn func(*sql.Tx) error) error {
	tx, err := s.db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelReadCommitted, ReadOnly: false})
	if err != nil {
		return err
	}

	err = fn(tx)

	if err == nil {
		err = tx.Commit()
	} else {
		tx.Rollback()
	}

	return err
}

func (s *Storage) CreateThing(ctx context.Context, thing models.Thing) error {
	_, err := s.db.ExecContext(ctx, "INSERT INTO things(owner,type,description,price) VALUES($1, $2, $3, $4);", thing.Owner, thing.Type, thing.Description, thing.Price)

	return err
}

func (s *Storage) RemuveThing(ctx context.Context, nickname string, thingId int) error {
	row := s.db.QueryRowContext(ctx, "SELECT owner FROM things WHERE id=$1", thingId)

	var owner string
	if err := row.Scan(&owner); err != nil {
		return err
	}

	if owner != nickname {
		return errors.New("not enough rights")
	}

	_, err := s.db.ExecContext(ctx, "DELETE FROM things WHERE id=$1", thingId)
	return err
}

func (s *Storage) BuyThingTx(ctx context.Context, thingId int, finishTime time.Time, nickname, email string) error {
	finishTimeSQL := finishTime.Format("2006-01-02 15:04:05")

	err := s.ExecTx(ctx, func(tx *sql.Tx) error {
		_, err := s.db.ExecContext(ctx, "INSERT INTO taken_things(thing_id, buyer, email, finish_time) VALUES($1,$2,$3,$4);", thingId, nickname, email, finishTimeSQL)
		if err != nil {
			return err
		}
		_, err = s.db.ExecContext(ctx, "UPDATE things SET available=false WHERE id=$1", thingId)

		return err
	})

	return err
}

func (s *Storage) AllThings(ctx context.Context) ([]models.Thing, error) {
	var things []models.Thing

	rows, err := s.db.QueryContext(ctx, "SELECT * FROM things")

	if err != nil {
		return things, err
	}

	for rows.Next() {
		var id, price int
		var owner, typ, description string
		var available bool

		rows.Scan(&id, &owner, &typ, &description, &price, &available)
		things = append(things, models.Thing{Owner: owner, Type: typ, Description: description, Price: price, Available: available})
	}

	return things, nil
}

func (s *Storage) RentalThings(ctx context.Context, nickname string) ([]models.RentalThing, error) {
	var things []models.RentalThing

	rows, err := s.db.QueryContext(ctx, "SELECT * FROM taken_things WHERE buyer=$1;", nickname)

	if err != nil {
		return things, err
	}

	for rows.Next() {
		var thing_id int
		var buyer, email, finishTime string

		rows.Scan(&thing_id, &buyer, &email, &finishTime)

		things = append(things, models.RentalThing{ThingId: thing_id, Buyer: buyer, Email: email, FinishTime: finishTime})
	}

	return things, nil
}

func (s *Storage) SaleThings(ctx context.Context, nickname string) ([]models.Thing, error) {
	var things []models.Thing

	rows, err := s.db.QueryContext(ctx, "SELECT * FROM things WHERE owner=$1", nickname)

	if err != nil {
		return things, err
	}

	for rows.Next() {
		var id, price int
		var owner, typ, description string
		var available bool

		rows.Scan(&id, &owner, &typ, &description, &price, &available)
		things = append(things, models.Thing{Owner: owner, Type: typ, Description: description, Price: price, Available: available})
	}

	return things, nil
}

func (p *Notifier) Send(email, text string) error {

	msg := &sarama.ProducerMessage{
		Topic: p.topic,
		Value: sarama.StringEncoder(text),
		Headers: []sarama.RecordHeader{{
			Key:   []byte("Email"),
			Value: []byte(email),
		},
		},
	}

	partition, offset, err := p.producer.SendMessage(msg)

	slog.Info(fmt.Sprintf("%+v", msg), slog.Int("partition", int(partition)), slog.Int("offset", int(offset)))

	return err
}
