package models

// Condition represents the conditions for a discount rule
type Condition struct {
	MinOrderValue *float64 `json:"min_order_value"`
	CustomerType  *string  `json:"customer_type"`
}

// Rule represents a discount rule
type Rule struct {
	ID                 string    `json:"id"`
	Description        string    `json:"description"`
	Conditions         Condition `json:"conditions"`
	DiscountPercentage *float64  `json:"discount_percentage,omitempty"`
	DiscountFixed      *float64  `json:"discount_fixed,omitempty"`
	Priority           int       `json:"priority"`
}
