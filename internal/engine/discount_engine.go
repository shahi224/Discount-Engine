package engine

import (
	"discountengine/internal/models"
	"discountengine/internal/utils"
	"sort"
	"sync"
)

// DiscountEngine is the main engine that applies discount rules
type DiscountEngine struct {
	configLoader  *utils.ConfigLoader
	ruleEvaluator *RuleEvaluator
	mu            sync.RWMutex
}

// NewDiscountEngine creates a new Discountengine
func NewDiscountEngine(rulesPath string) *DiscountEngine {
	configLoader := utils.NewConfigLoader(rulesPath)

	// load initial rules
	if err := configLoader.LoadRules(); err != nil {
		panic("failed to load discount rules: " + err.Error())
	}

	return &DiscountEngine{
		configLoader:  configLoader,
		ruleEvaluator: NewRuleEvaluator(),
	}
}

// Applydiscounts applies all applicable discounts to an order
func (de *DiscountEngine) ApplyDiscounts(order models.Order) models.DiscountResponse {
	de.mu.RLock()
	defer de.mu.RUnlock()

	rules := de.configLoader.GetRules()

	var applicableRules []models.Rule
	var ruleDiscounts []struct {
		rule     models.Rule
		discount float64
	}

	// Find all applicable rules
	for _, rule := range rules {
		if de.ruleEvaluator.Evaluate(rule, order) {
			applicableRules = append(applicableRules, rule)

			discount := de.ruleEvaluator.CalculateDiscount(rule, order.OrderTotal)

			ruleDiscounts = append(ruleDiscounts, struct {
				rule     models.Rule
				discount float64
			}{rule, discount})
		}
	}

	// if no rules apply return original total
	if len(applicableRules) == 0 {
		return models.DiscountResponse{
			OrderID:              order.OderID,
			OriginalTotal:        order.OrderTotal,
			DiscountAmount:       0,
			FinalTotal:           order.OrderTotal,
			AppliedRules:         []models.Rule{},
			HighestPriorityRules: models.Rule{},
		}
	}
	// Sort applicable rules by priority(descending) and discount amount(descending)
	sort.Slice(ruleDiscounts, func(i, j int) bool {
		if ruleDiscounts[i].rule.Priority != ruleDiscounts[j].rule.Priority {
			return ruleDiscounts[i].rule.Priority > ruleDiscounts[j].rule.Priority
		}
		return ruleDiscounts[i].discount > ruleDiscounts[j].discount
	})

	// Apply the best discount (highest priority, then highest discount)
	bestRule := ruleDiscounts[0].rule
	bestDiscount := ruleDiscounts[0].discount

	finalTotal := order.OrderTotal - bestDiscount
	if finalTotal < 0 {
		finalTotal = 0
	}

	return models.DiscountResponse{
		OrderID:              order.OderID,
		OriginalTotal:        order.OrderTotal,
		DiscountAmount:       bestDiscount,
		FinalTotal:           finalTotal,
		AppliedRules:         applicableRules,
		HighestPriorityRules: bestRule,
	}

}

// ReloadRules reloads rules from configuration file
func (de *DiscountEngine) ReloadRules() error {
	de.mu.Lock()
	defer de.mu.Unlock()

	return de.configLoader.ReloadRules()
}

// GetRules returns current rules (for debugging/testing)
func (de *DiscountEngine) GetRules() []models.Rule {
	return de.configLoader.GetRules()
}
