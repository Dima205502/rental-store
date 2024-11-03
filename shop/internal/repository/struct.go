package repository

import (
	"database/sql"

	"github.com/IBM/sarama"
)

type Storage struct {
	db *sql.DB
}

type Notifier struct {
	producer sarama.SyncProducer
	topic    string
}
