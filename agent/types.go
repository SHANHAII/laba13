package main

// Task — универсальная обёртка задачи для NATS
type Task struct {
	ID      string `json:"id"`
	Type    string `json:"type"`
	Payload string `json:"payload"`
}

// Result — универсальный ответ агента
type Result struct {
	TaskID  string `json:"task_id"`
	Success bool   `json:"success"`
	Output  string `json:"output"`
	Error   string `json:"error,omitempty"`
}

// OrderItem — позиция в заказе
type OrderItem struct {
	Name  string  `json:"name"`
	Qty   int     `json:"qty"`
	Price float64 `json:"price"`
}

// Order — входные данные агента приёма заказа
type Order struct {
	ID           string      `json:"id"`
	TableID      int         `json:"table_id"`
	Items        []OrderItem `json:"items"`
	CustomerName string      `json:"customer_name"`
}

// OrderResult — итог валидации заказа
type OrderResult struct {
	OrderID string  `json:"order_id"`
	TableID int     `json:"table_id"`
	Total   float64 `json:"total"`
	Valid   bool    `json:"valid"`
}

// KitchenTask — задача для агента кухни
type KitchenTask struct {
	OrderID string      `json:"order_id"`
	Items   []OrderItem `json:"items"`
}

// KitchenResult — результат приготовления
type KitchenResult struct {
	OrderID     string `json:"order_id"`
	Ready       bool   `json:"ready"`
	DurationSec int    `json:"duration_sec"`
}

// TableTask — задача для агента управления столами
type TableTask struct {
	OrderID string `json:"order_id"`
	TableID int    `json:"table_id"`
	Status  string `json:"status"` // occupied | free | reserved
}

// TableResult — результат обновления статуса стола
type TableResult struct {
	OrderID string `json:"order_id"`
	TableID int    `json:"table_id"`
	Status  string `json:"status"`
	Ok      bool   `json:"ok"`
}

// DeliveryTask — задача для агента доставки
type DeliveryTask struct {
	OrderID string `json:"order_id"`
	TableID int    `json:"table_id"`
}

// DeliveryResult — результат доставки заказа
type DeliveryResult struct {
	OrderID string `json:"order_id"`
	Waiter  string `json:"waiter"`
	Done    bool   `json:"done"`
}
