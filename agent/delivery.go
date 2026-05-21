package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
)

var waiterPool = []string{"Алексей", "Мария", "Дмитрий", "Елена", "Иван"}

// processDelivery назначает официанта и подтверждает доставку
func processDelivery(dt DeliveryTask) (DeliveryResult, error) {
	if dt.OrderID == "" {
		return DeliveryResult{}, fmt.Errorf("empty order_id")
	}
	if dt.TableID <= 0 {
		return DeliveryResult{}, fmt.Errorf("invalid table_id: %d", dt.TableID)
	}
	return DeliveryResult{
		OrderID: dt.OrderID,
		Waiter:  waiterPool[rand.Intn(len(waiterPool))],
		Done:    true,
	}, nil
}

func handleDeliveryPayload(payload string) (string, error) {
	var dt DeliveryTask
	if err := json.Unmarshal([]byte(payload), &dt); err != nil {
		return "", fmt.Errorf("unmarshal delivery task: %w", err)
	}
	res, err := processDelivery(dt)
	if err != nil {
		return "", err
	}
	out, err := json.Marshal(res)
	if err != nil {
		return "", fmt.Errorf("marshal delivery result: %w", err)
	}
	return string(out), nil
}
