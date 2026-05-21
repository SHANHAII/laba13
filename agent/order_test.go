package main

import "testing"

func TestProcessOrder(t *testing.T) {
	tests := []struct {
		name      string
		order     Order
		wantErr   bool
		wantTotal float64
	}{
		{
			name: "valid order",
			order: Order{
				ID: "ord-1", TableID: 3,
				Items: []OrderItem{
					{Name: "burger", Qty: 2, Price: 250.0},
					{Name: "cola", Qty: 1, Price: 100.0},
				},
			},
			wantErr:   false,
			wantTotal: 600.0,
		},
		{
			name:    "invalid table_id zero",
			order:   Order{ID: "ord-2", TableID: 0, Items: []OrderItem{{Name: "pizza", Qty: 1, Price: 400}}},
			wantErr: true,
		},
		{
			name:    "invalid table_id negative",
			order:   Order{ID: "ord-3", TableID: -1, Items: []OrderItem{{Name: "pizza", Qty: 1, Price: 400}}},
			wantErr: true,
		},
		{
			name:    "empty items",
			order:   Order{ID: "ord-4", TableID: 5, Items: []OrderItem{}},
			wantErr: true,
		},
		{
			name:    "invalid qty",
			order:   Order{ID: "ord-5", TableID: 2, Items: []OrderItem{{Name: "salad", Qty: 0, Price: 150}}},
			wantErr: true,
		},
		{
			name:    "negative price",
			order:   Order{ID: "ord-6", TableID: 1, Items: []OrderItem{{Name: "soup", Qty: 1, Price: -50}}},
			wantErr: true,
		},
		{
			name: "single item",
			order: Order{
				ID: "ord-7", TableID: 10,
				Items: []OrderItem{{Name: "steak", Qty: 1, Price: 1200}},
			},
			wantErr:   false,
			wantTotal: 1200.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := processOrder(tt.order)
			if (err != nil) != tt.wantErr {
				t.Fatalf("processOrder() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && result.Total != tt.wantTotal {
				t.Errorf("Total = %v, want %v", result.Total, tt.wantTotal)
			}
			if !tt.wantErr && !result.Valid {
				t.Error("Valid = false, want true")
			}
		})
	}
}
