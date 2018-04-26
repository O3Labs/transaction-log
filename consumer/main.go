package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	nsq "github.com/bitly/go-nsq"
	"github.com/o3labs/neo-utils/neoutils/neorpc"
	"github.com/o3labs/transaction-log/consumer/transaction"
	"github.com/o3labs/transaction-log/services"
)

var maxInFlight = 200

type NoopNSQLogger struct{}

// Output allows us to implement the nsq.Logger interface
func (l *NoopNSQLogger) Output(calldepth int, s string) error {
	return nil
}
func main() {

	nsqConfig := nsq.NewConfig()

	client := neorpc.NewClient("http://seed1.o3node.org:10332")
	saveTransactionQ, _ := nsq.NewConsumer("neo.network", "tx", nsqConfig)
	saveTransactionQ.ChangeMaxInFlight(maxInFlight)
	saveTransactionQ.AddConcurrentHandlers(&transaction.MessageHandler{
		Service:      services.NewTransactionService(),
		NEORPCClient: client,
	}, 20)

	saveTransactionQ.SetLogger(
		&NoopNSQLogger{},
		nsq.LogLevelError,
	)

	if err := saveTransactionQ.ConnectToNSQLookupd("127.0.0.1:4161"); err != nil {
		log.Panicf("saveQ : Could not connect")
	}

	log.Printf("Consumer is running: %v", "")
	shutdown := make(chan os.Signal, 2)
	signal.Notify(shutdown, syscall.SIGINT)
	for {
		select {
		case <-saveTransactionQ.StopChan:
			return // uh oh consumer disconnected. Time to quit.
		case <-shutdown:
			// Synchronously drain the queue before falling out of main
			saveTransactionQ.Stop()
		}
	}
}
