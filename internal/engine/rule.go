package engine

import "discountengine/internal/models"

// RuleEvaluator evaluates if a rule applies to an order
type RuleEvaluator struct{}

// NewRuleEvaluator creates a new RuleEvaluator
func NewRuleEvaluator() *RuleEvaluator {
	return &RuleEvaluator{}
}

// Evaluate checks if all conditios of a rule are satisfied
func (re *RuleEvaluator) Evaluate(rule models.Rule, order models.Order) bool {
	// check min_order_value condition
	if rule.Conditions.MinOrderValue != nil {
		if order.OrderTotal < *rule.Conditions.MinOrderValue {
			return false
		}
	}

	// check customer_type conditon
	if rule.Conditions.CustomerType != nil {
		if order.CustomerType != *rule.Conditions.CustomerType {
			return false
		}
	}
	return true
}

// CalculateDiscount calculate the discount amount for a rule
func (re *RuleEvaluator) CalculateDiscount(rule models.Rule, orderTotal float64) float64 {
	if rule.DiscountPercentage != nil {
		return orderTotal * (*rule.DiscountPercentage / 100)
	}

	if rule.DiscountFixed != nil {
		return *rule.DiscountFixed
	}
	return 0
}
