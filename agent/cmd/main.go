package main

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/alexvancasper/TunnelBroker/agent/internal/broker"
	"github.com/alexvancasper/TunnelBroker/agent/internal/doer"
	formatter "github.com/fabienm/go-logrus-formatters"
	"github.com/rabbitmq/amqp091-go"
	"github.com/sirupsen/logrus"
)

func main() {

	//Initialize Logging connections
	var MyLogger = logrus.New()

	gelfFmt := formatter.NewGelf("agent")
	MyLogger.SetFormatter(gelfFmt)
	MyLogger.SetOutput(os.Stdout)
	loglevel, err := logrus.ParseLevel("debug")
	if err != nil {
		MyLogger.WithField("function", "main").Fatalf("error %v", err)
	}
	MyLogger.SetLevel(loglevel)

	// Initialize message broker connection
	m, err := broker.MsgBrokerInit(os.Getenv("BROKER_CONN"), os.Getenv("QUEUENAME"))
	if err != nil {
		MyLogger.Fatalf("Message broker error init: %s", err)
	}
	defer m.Close()

	MyLogger.Debug("Message broker connected")

	msgs, err := m.AddRegisterConsumer()
	if err != nil {
		MyLogger.Fatalf("AddRegisterConsumer error: %s", err)
	}
	var wg sync.WaitGroup
	ctx, ctxCancel := context.WithCancel(context.Background())

	MyLogger.Info(" [*] Waiting for messages. To exit press CTRL+C")
	wg.Add(1)
	go Listener(ctx, &wg, msgs, MyLogger)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
	MyLogger.Info("Graceful shutdown")
	ctxCancel()

}

func Listener(ctx context.Context, wg *sync.WaitGroup, msgs <-chan amqp091.Delivery, log *logrus.Logger) {
	defer wg.Done()
	l := log.WithField("function", "Listener")
	h := doer.Handler{
		Log: log,
	}
	for {
		time.Sleep(100 * time.Millisecond)
		select {
		case <-ctx.Done():
			l.Debug("Context closed")
			return
		case msg := <-msgs:
			l.WithField("message type", msg.Type).WithField("body", msg.Body).Infof("Received message from queue")
			switch msg.Type {
			case string(broker.ADD):
				wg.Add(1)
				go h.AddTunnel(wg, msg.Body)
			case string(broker.DELETE):
				wg.Add(1)
				go h.DeleteTunnel(wg, msg.Body)
			}
		}
	}
}
