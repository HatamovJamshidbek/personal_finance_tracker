package rabbitMq

import (
	pbb "budgeting_service/genproto/budgeting_service"
	"budgeting_service/pkg/logger"
	"budgeting_service/storage"
	"context"
	"encoding/json"
	"fmt"
	"time"
)

func SetUpConsumers(log logger.ILogger, mongoStore storage.IStorage) {
	queues := []string{"transaction_created", "budget_updated", "goal_progress_updated", "notification_created"}

	for _, queueName := range queues {
		go func(queueName string) {
			for {
				time.Sleep(10 * time.Second)

				readerMq, err := NewRabbitMqConsumerImpl("amqp://user:1111@rabbitmq:5672/", queueName, mongoStore)
				if err != nil {
					log.Error("Error while creating RabbitMQ consumer for "+queueName, logger.Error(err))
					time.Sleep(5 * time.Second) // Retry after a delay
					continue
				}
				defer func() {
					//if err := readerMq.Close(); err != nil {
					//	log.Error("Error while closing RabbitMQ consumer for "+queueName, logger.Error(err))
					//}
				}()

				messageChan := make(chan []byte)
				errorChan := make(chan error)

				go func() {
					err := readerMq.ConsumeMessages(func(message []byte) {
						messageChan <- message
					})
					if err != nil {
						errorChan <- err
					}
				}()

				for {
					select {
					case message := <-messageChan:
						ProcessMessage(queueName, message, log, mongoStore)
					case err := <-errorChan:
						log.Error("Error while consuming messages for "+queueName, logger.Error(err))
						time.Sleep(5 * time.Second) // Retry after a delay
					case <-time.After(30 * time.Second): // Handle timeout
						log.Warn("No messages received for " + queueName + " in the last 30 seconds")
					}
				}
			}
		}(queueName)
	}
}

func ProcessMessage(queueName string, message []byte, log logger.ILogger, mongoStore storage.IStorage) {
	switch queueName {
	case "transaction_created":
		var txnMsg pbb.TransactionRequest
		if err := json.Unmarshal(message, &txnMsg); err != nil {
			log.Error("Error unmarshalling transaction_created message", logger.Error(err))
			return
		}
		fmt.Printf("Processing transaction: %+v\n", txnMsg)
		_, err := mongoStore.Transactions().CreateTransaction(context.TODO(), &txnMsg)
		if err != nil {
			log.Error("Error saving transaction to storage", logger.Error(err))
			return
		}

	case "budget_updated":
		var budgetMsg pbb.Budget
		if err := json.Unmarshal(message, &budgetMsg); err != nil {
			log.Error("Error unmarshalling budget_updated message", logger.Error(err))
			return
		}
		fmt.Printf("Processing budget update: %+v\n", budgetMsg)
		_, err := mongoStore.Budgets().UpdateBudget(context.TODO(), &budgetMsg)
		if err != nil {
			log.Error("Error updating budget in storage", logger.Error(err))
			return
		}

	case "goal_progress_updated":
		var goalMsg pbb.Goal
		if err := json.Unmarshal(message, &goalMsg); err != nil {
			log.Error("Error unmarshalling goal_progress_updated message", logger.Error(err))
			return
		}
		fmt.Printf("Processing goal progress update: %+v\n", goalMsg)
		_, err := mongoStore.Goals().UpdateGoal(context.TODO(), &goalMsg)
		if err != nil {
			log.Error("Error updating goal in storage", logger.Error(err))
			return
		}

	default:
		log.Warn("Unknown queue name: " + queueName)
	}
}
