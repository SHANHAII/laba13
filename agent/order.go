package main

import (
	"encoding/json"
	"fmt"
)

// processOrder валидирует заказ и считает итоговую сумму
func processOrder(o Order) (OrderResult, error) {
	if o.TableID <= 0 {
		return OrderResult{}, fmt.Errorf("invalid table_id: %d", o.TableID)
	}
	if len(o.Items) == 0 {
		return OrderResult{}, fmt.Errorf("order %q has no items", o.ID)
	}

	total := 0.0
	for _, item := range o.Items {
		if item.Qty <= 0 {
			return OrderResult{}, fmt.Errorf("item %q has invalid qty: %d", item.Name, item.Qty)
		}
		if item.Price < 0 {
			return OrderResult{}, fmt.Errorf("item %q has negative price", item.Name)
		}
		total += float64(item.Qty) * item.Price
	}

	return OrderResult{
		OrderID: o.ID,
		TableID: o.TableID,
		Total:   total,
		Valid:   true,
	}, nil
}

func handleOrderPayload(payload string) (string, error) {
	var o Order
	if err := json.Unmarshal([]byte(payload), &o); err != nil {
		return "", fmt.Errorf("unmarshal order: %w", err)
	}
	res, err := processOrder(o)
	if err != nil {
		return "", err
	}
	out, err := json.Marshal(res)
	if err != nil {
		return "", fmt.Errorf("marshal order result: %w", err)
	}
	return string(out), nil
}
