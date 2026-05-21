package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/nats-io/nats.go"
)

func main() {
	agentType := flag.String("type", "order", "agent type: order|kitchen|table|delivery")
	natsURL := flag.String("nats", nats.DefaultURL, "NATS server URL")
	logFilePath := flag.String("log", "", "additional log file path (stdout always enabled)")
	flag.Parse()

	writers := []io.Writer{os.Stdout}
	if *logFilePath != "" {
		f, err := os.OpenFile(*logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Fatalf("open log file %q: %v", *logFilePath, err)
		}
		defer f.Close()
		writers = append(writers, f)
	}
	logger := log.New(io.MultiWriter(writers...), fmt.Sprintf("[%s] ", *agentType), log.LstdFlags)

	nc, err := nats.Connect(*natsURL)
	if err != nil {
		logger.Fatalf("ERROR NATS connect to %s: %v", *natsURL, err)
	}
	defer nc.Close()
	logger.Printf("INFO connected to NATS at %s", *natsURL)

	counter := &Counter{}
	tableState := NewTableState()

	var subject string
	var handler func(string) (string, error)

	switch *agentType {
	case "order":
		subject = "restaurant.order.new"
		handler = handleOrderPayload
	case "kitchen":
		subject = "restaurant.kitchen.new"
		handler = handleKitchenPayload
	case "table":
		subject = "restaurant.table.update"
		handler = func(p string) (string, error) { return handleTablePayload(tableState, p) }
	case "delivery":
		subject = "restaurant.delivery.new"
		handler = handleDeliveryPayload
	default:
		logger.Fatalf("ERROR unknown agent type: %q", *agentType)
	}

	// QueueSubscribe обеспечивает балансировку нагрузки между экземплярами одного типа агента
	queueGroup := *agentType + "-workers"
	_, err = nc.QueueSubscribe(subject, queueGroup, func(m *nats.Msg) {
		var task Task
		if unmarshalErr := json.Unmarshal(m.Data, &task); unmarshalErr != nil {
			logger.Printf("ERROR unmarshal task: %v", unmarshalErr)
			return
		}
		logger.Printf("INFO received task %s (type=%s)", task.ID, task.Type)

		output, procErr := handler(task.Payload)

		res := Result{TaskID: task.ID, Success: procErr == nil, Output: output}
		if procErr != nil {
			res.Error = procErr.Error()
			logger.Printf("ERROR task %s failed: %v", task.ID, procErr)
		} else {
			counter.Inc()
			logger.Printf("INFO task %s completed, total=%d", task.ID, counter.Get())
		}

		data, _ := json.Marshal(res)
		if pubErr := nc.Publish("restaurant.tasks.completed", data); pubErr != nil {
			logger.Printf("ERROR publish result for task %s: %v", task.ID, pubErr)
		}
	})
	if err != nil {
		logger.Fatalf("ERROR subscribe to %s: %v", subject, err)
	}

	logger.Printf("INFO %s agent started, queue=%s, subject=%s", *agentType, queueGroup, subject)
	select {}
}
