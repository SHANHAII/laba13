package main

import (
	"encoding/json"
	"fmt"
	"sync"
)

const maxTableID = 20

var validTableStatuses = map[string]bool{
	"occupied": true,
	"free":     true,
	"reserved": true,
}

// TableState хранит статусы столов; потокобезопасна
type TableState struct {
	mu     sync.Mutex
	tables map[int]string
}

func NewTableState() *TableState {
	return &TableState{tables: make(map[int]string)}
}

// processTable обновляет статус стола с проверкой бизнес-правил
func processTable(s *TableState, tt TableTask) (TableResult, error) {
	if tt.TableID <= 0 || tt.TableID > maxTableID {
		return TableResult{}, fmt.Errorf("table_id %d out of range [1..%d]", tt.TableID, maxTableID)
	}
	if !validTableStatuses[tt.Status] {
		return TableResult{}, fmt.Errorf("invalid status %q", tt.Status)
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if s.tables[tt.TableID] == "occupied" && tt.Status == "occupied" {
		return TableResult{}, fmt.Errorf("table %d is already occupied", tt.TableID)
	}
	s.tables[tt.TableID] = tt.Status

	return TableResult{
		OrderID: tt.OrderID,
		TableID: tt.TableID,
		Status:  tt.Status,
		Ok:      true,
	}, nil
}

func handleTablePayload(s *TableState, payload string) (string, error) {
	var tt TableTask
	if err := json.Unmarshal([]byte(payload), &tt); err != nil {
		return "", fmt.Errorf("unmarshal table task: %w", err)
	}
	res, err := processTable(s, tt)
	if err != nil {
		return "", err
	}
	out, err := json.Marshal(res)
	if err != nil {
		return "", fmt.Errorf("marshal table result: %w", err)
	}
	return string(out), nil
}
