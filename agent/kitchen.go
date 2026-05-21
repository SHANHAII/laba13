package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
)

var menuItems = map[string]bool{
	"burger":   true,
	"pizza":    true,
	"pasta":    true,
	"salad":    true,
	"soup":     true,
	"steak":    true,
	"sushi":    true,
	"sandwich": true,
}

// processKitchen проверяет наличие позиций в меню и возвращает время приготовления
func processKitchen(kt KitchenTask) (KitchenResult, error) {
	if kt.OrderID == "" {
		return KitchenResult{}, fmt.Errorf("empty order_id")
	}
	if len(kt.Items) == 0 {
		return KitchenResult{}, fmt.Errorf("no items in kitchen task")
	}
	for _, item := range kt.Items {
		if !menuItems[item.Name] {
			return KitchenResult{}, fmt.Errorf("item %q not available on menu", item.Name)
		}
	}
	return KitchenResult{
		OrderID:     kt.OrderID,
		Ready:       true,
		DurationSec: 1 + rand.Intn(5),
	}, nil
}

func handleKitchenPayload(payload string) (string, error) {
	var kt KitchenTask
	if err := json.Unmarshal([]byte(payload), &kt); err != nil {
		return "", fmt.Errorf("unmarshal kitchen task: %w", err)
	}
	res, err := processKitchen(kt)
	if err != nil {
		return "", err
	}
	out, err := json.Marshal(res)
	if err != nil {
		return "", fmt.Errorf("marshal kitchen result: %w", err)
	}
	return string(out), nil
}
