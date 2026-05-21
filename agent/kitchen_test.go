package main

import "testing"

func TestProcessKitchen(t *testing.T) {
	tests := []struct {
		name    string
		task    KitchenTask
		wantErr bool
	}{
		{
			name:    "valid single item",
			task:    KitchenTask{OrderID: "ord-1", Items: []OrderItem{{Name: "burger", Qty: 1, Price: 250}}},
			wantErr: false,
		},
		{
			name: "valid multiple items",
			task: KitchenTask{
				OrderID: "ord-2",
				Items: []OrderItem{
					{Name: "pizza", Qty: 1, Price: 500},
					{Name: "salad", Qty: 2, Price: 200},
				},
			},
			wantErr: false,
		},
		{
			name:    "item not on menu",
			task:    KitchenTask{OrderID: "ord-3", Items: []OrderItem{{Name: "unicorn_dish", Qty: 1, Price: 999}}},
			wantErr: true,
		},
		{
			name:    "empty order_id",
			task:    KitchenTask{OrderID: "", Items: []OrderItem{{Name: "pizza", Qty: 1, Price: 400}}},
			wantErr: true,
		},
		{
			name:    "empty items",
			task:    KitchenTask{OrderID: "ord-5", Items: []OrderItem{}},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := processKitchen(tt.task)
			if (err != nil) != tt.wantErr {
				t.Fatalf("processKitchen() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr {
				if !result.Ready {
					t.Error("Ready = false, want true")
				}
				if result.DurationSec < 1 || result.DurationSec > 5 {
					t.Errorf("DurationSec = %d, want [1..5]", result.DurationSec)
				}
				if result.OrderID != tt.task.OrderID {
					t.Errorf("OrderID = %q, want %q", result.OrderID, tt.task.OrderID)
				}
			}
		})
	}
}
