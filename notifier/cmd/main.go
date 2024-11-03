package main

import (
	"fmt"
	"log"
	"log/slog"
	"notifier/config"
	"notifier/internal"
	"os"
	"os/signal"
	"syscall"

	"github.com/IBM/sarama"
)

func main() {
	cfg := config.Init()

	internal.Init(cfg)

	consumer, err := sarama.NewConsumer(cfg.Broker_addrs, nil)
	if err != nil {
		log.Fatalf("Failed to create consumer: %v", err)
	}

	defer consumer.Close()

	partitions, err := consumer.Partitions(cfg.Topic)

	if err != nil {
		log.Fatalf("Failed to consume partition: %v", err)
	}

	for _, partition := range partitions {
		partitionConsumer, err := consumer.ConsumePartition(cfg.Topic, partition, sarama.OffsetNewest)

		if err != nil {
			log.Fatalf("Failed to start consuming partition %d: %v", partition, err)
		}

		go func(pc sarama.PartitionConsumer) {
			for {
				select {
				case msg := <-pc.Messages():
					text := string(msg.Value)

					var email string
					for _, header := range msg.Headers {
						if string(header.Key) == "Email" {
							email = string(header.Value)
							break
						}
					}

					fmt.Printf("Message = %+v\n", msg)

					if err := internal.Send(email, text); err != nil {
						slog.Error("Send", slog.String("err", err.Error()))
					}

				case err := <-pc.Errors():
					log.Printf("Ошибка при получении сообщения: %v", err)
				}
			}
		}(partitionConsumer)
	}

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	<-sigs
}
