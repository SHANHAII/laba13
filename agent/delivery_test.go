package main

import "testing"

func TestProcessDelivery(t *testing.T) {
	tests := []struct {
		name    string
		task    DeliveryTask
		wantErr bool
	}{
		{
			name:    "valid delivery",
			task:    DeliveryTask{OrderID: "ord-1", TableID: 4},
			wantErr: false,
		},
		{
			name:    "empty order_id",
			task:    DeliveryTask{OrderID: "", TableID: 4},
			wantErr: true,
		},
		{
			name:    "invalid table_id zero",
			task:    DeliveryTask{OrderID: "ord-2", TableID: 0},
			wantErr: true,
		},
		{
			name:    "invalid table_id negative",
			task:    DeliveryTask{OrderID: "ord-3", TableID: -5},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := processDelivery(tt.task)
			if (err != nil) != tt.wantErr {
				t.Fatalf("processDelivery() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr {
				if !result.Done {
					t.Error("Done = false, want true")
				}
				if result.Waiter == "" {
					t.Error("Waiter is empty")
				}
				if result.OrderID != tt.task.OrderID {
					t.Errorf("OrderID = %q, want %q", result.OrderID, tt.task.OrderID)
				}
			}
		})
	}
}
