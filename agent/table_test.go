package main

import "testing"

func TestProcessTable(t *testing.T) {
	tests := []struct {
		name       string
		task       TableTask
		wantErr    bool
		wantStatus string
	}{
		{
			name:       "occupy free table",
			task:       TableTask{OrderID: "ord-1", TableID: 1, Status: "occupied"},
			wantErr:    false,
			wantStatus: "occupied",
		},
		{
			name:       "reserve table",
			task:       TableTask{OrderID: "ord-2", TableID: 2, Status: "reserved"},
			wantErr:    false,
			wantStatus: "reserved",
		},
		{
			name:       "free table",
			task:       TableTask{OrderID: "ord-3", TableID: 3, Status: "free"},
			wantErr:    false,
			wantStatus: "free",
		},
		{
			name:    "invalid table_id zero",
			task:    TableTask{OrderID: "ord-4", TableID: 0, Status: "occupied"},
			wantErr: true,
		},
		{
			name:    "invalid table_id over max",
			task:    TableTask{OrderID: "ord-5", TableID: 21, Status: "occupied"},
			wantErr: true,
		},
		{
			name:    "invalid status",
			task:    TableTask{OrderID: "ord-6", TableID: 5, Status: "dirty"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			state := NewTableState()
			result, err := processTable(state, tt.task)
			if (err != nil) != tt.wantErr {
				t.Fatalf("processTable() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && result.Status != tt.wantStatus {
				t.Errorf("Status = %q, want %q", result.Status, tt.wantStatus)
			}
			if !tt.wantErr && !result.Ok {
				t.Error("Ok = false, want true")
			}
		})
	}
}

func TestProcessTableDoubleOccupy(t *testing.T) {
	state := NewTableState()

	_, err := processTable(state, TableTask{OrderID: "ord-1", TableID: 4, Status: "occupied"})
	if err != nil {
		t.Fatalf("first occupy failed: %v", err)
	}

	_, err = processTable(state, TableTask{OrderID: "ord-2", TableID: 4, Status: "occupied"})
	if err == nil {
		t.Error("expected error on double occupy, got nil")
	}
}

func TestProcessTableFreeAfterOccupy(t *testing.T) {
	state := NewTableState()

	if _, err := processTable(state, TableTask{OrderID: "ord-1", TableID: 7, Status: "occupied"}); err != nil {
		t.Fatalf("occupy failed: %v", err)
	}
	if _, err := processTable(state, TableTask{OrderID: "ord-2", TableID: 7, Status: "free"}); err != nil {
		t.Errorf("free after occupy failed: %v", err)
	}
	// теперь снова можно занять
	if _, err := processTable(state, TableTask{OrderID: "ord-3", TableID: 7, Status: "occupied"}); err != nil {
		t.Errorf("re-occupy after free failed: %v", err)
	}
}
