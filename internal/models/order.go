package models

// Order represents an order from customer
type Order struct {
	OderID       string  `json:"order_id"`
	OrderTotal   float64 `json:"order_total"`
	CustomerType string  `json:"customer_type"`
}

// DiscountResponse represent the respone after applying discounts
type DiscountResponse struct {
	OrderID              string  `json:"order_id"`
	OriginalTotal        float64 `json:"original_total"`
	DiscountAmount       float64 `json:"discount_amount"`
	FinalTotal           float64 `json:"final_total"`
	AppliedRules         []Rule  `json:"applied_rules"`
	HighestPriorityRules Rule    `json:"highest_priority_rule"`
}
