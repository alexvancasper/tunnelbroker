package main

import (
	"context"
	"os"
	"os/signal"
	"sync"

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

	addMsgs, err := m.AddRegisterConsumer()
	if err != nil {
		MyLogger.Fatalf("AddRegisterConsumer error: %s", err)
	}
	delMsgs, err := m.DeleteRegisterConsumer()
	if err != nil {
		MyLogger.Fatalf("DeleteRegisterConsumer error: %s", err)
	}

	var wg sync.WaitGroup
	ctx, ctxCancel := context.WithCancel(context.Background())

	MyLogger.Info(" [*] Waiting for messages. To exit press CTRL+C")
	wg.Add(1)
	go Listener(ctx, addMsgs, delMsgs, MyLogger, &wg)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
	MyLogger.Info("Graceful shutdown")
	ctxCancel()

}

func Listener(ctx context.Context, addMsgs, delMsgs <-chan amqp091.Delivery, log *logrus.Logger, wg *sync.WaitGroup) {
	defer wg.Done()
	l := log.WithField("function", "Listener")
	h := doer.Handler{
		Log: log,
	}
	for {
		select {
		case <-ctx.Done():
			l.Debug("Context closed")
			return
		case addmsg := <-addMsgs:
			l.Debugf("Received from ADD queue a message: %s", addmsg.Body)
			wg.Add(1)
			go h.AddTunnel(wg, addmsg.Body)
		case delmsg := <-delMsgs:
			l.Debugf("Received from DELETE queue a message: %s", delmsg.Body)
			wg.Add(1)
			go h.DeleteTunnel(wg, delmsg.Body)
		}
	}

}