package q

import (
	"encoding/json"
	"log"

	nsq "github.com/bitly/go-nsq"
	"github.com/o3labs/transaction-log/models"
)

type QService struct {
	Producer *nsq.Producer
}

type QServiceInterface interface {
	PushTransaction(tx models.TransactionLog) error
	Stop()
}

func Shared() *QService {
	producer, err := NewProducer()
	if err != nil {
		log.Printf("error NewProducer %v", err)
		return nil
	}
	return &QService{Producer: producer}
}

var _ QServiceInterface = (*QService)(nil)

func NewProducer() (*nsq.Producer, error) {
	nsqConfig := nsq.NewConfig()
	return nsq.NewProducer("127.0.0.1:4150", nsqConfig)
}

func (q *QService) Stop() {
	q.Producer.Stop()
}

func (q *QService) PushTransaction(tx models.TransactionLog) error {
	jsonData, err := json.Marshal(tx)
	if err != nil {
		return err
	}
	q.Producer.Publish("neo.network", jsonData)
	return nil
}
