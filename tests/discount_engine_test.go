package tests

import (
	"discountengine/internal/engine"
	"discountengine/internal/models"
	"testing"
)

func TestDiscountEngine(t *testing.T) {
	// Create engine with test rules path
	discountEngine := engine.NewDiscountEngine("../config/rules.json")

	tests := []struct {
		name             string
		order            models.Order
		expectedDiscount float64
		expectedFinal    float64
	}{
		{
			name: "No discount - small order",
			order: models.Order{
				OderID:       "test-1",
				OrderTotal:   50,
				CustomerType: "regular",
			},
			expectedDiscount: 0,
			expectedFinal:    50,
		},
		{
			name: "10%  discount for order over $100 ",
			order: models.Order{
				OderID:       "test-2",
				OrderTotal:   150,
				CustomerType: "regular",
			},
			expectedDiscount: 15, // 10% of 150
			expectedFinal:    135,
		},
		{
			name: "15% discount for premium customer",
			order: models.Order{
				OderID:       "test-3",
				OrderTotal:   80,
				CustomerType: "premium",
			},
			expectedDiscount: 12, // 15% of 80
			expectedFinal:    68,
		},
		{
			name: "$20 fixed discount for order over $200",
			order: models.Order{
				OderID:       "test-4",
				OrderTotal:   250,
				CustomerType: "regular",
			},
			expectedDiscount: 20,
			expectedFinal:    230,
		},
		{
			name: "Premium customer with $250 order (multiple rules)",
			order: models.Order{
				OderID:       "test-5",
				OrderTotal:   250,
				CustomerType: "premium",
			},
			expectedDiscount: 37.5, // 15% of 250 (premium rule has higher priority)
			expectedFinal:    212.5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := discountEngine.ApplyDiscounts(tt.order)

			if result.DiscountAmount != tt.expectedDiscount {
				t.Errorf("Expected discount %v, got %v", tt.expectedDiscount, result.DiscountAmount)
			}
			if result.FinalTotal != tt.expectedFinal {
				t.Errorf("Expected final total %v, got %v", tt.expectedDiscount, result.DiscountAmount)
			}
		})
	}
}

func TestRulePriority(t *testing.T) {
	discountEngine := engine.NewDiscountEngine("../config/rules.json")

	// Test case where rule 4 (5% for regular over $50) and rule 1 (10% for order over $100) both apply
	// Rule 1 has higher priority (1) than rule 4 (also 1) but rule 1 gives higher discount

	order := models.Order{
		OderID:       "test-priority",
		OrderTotal:   150,
		CustomerType: "regular",
	}

	result := discountEngine.ApplyDiscounts(order)

	// should apply rule 1 (10% discount) because same priority but higher discount
	expectedDiscount := 15.0 // 10% of 150
	if result.DiscountAmount != expectedDiscount {
		t.Errorf("Expected discount %v, got %v", expectedDiscount, result.DiscountAmount)
	}

	if result.HighestPriorityRules.ID != "rule_1" {
		t.Errorf("Expected rule_1 to be applied, got %s", result.HighestPriorityRules.ID)
	}
}

func TestConcurrentAccess(t *testing.T) {
	discountEngine := engine.NewDiscountEngine("../config/rule.json")

	// test concurrent access to the engine
	done := make(chan bool)

	for i := 0; i < 10; i++ {
		go func(orderID string) {
			order := models.Order{
				OderID:       orderID,
				OrderTotal:   150,
				CustomerType: "regular",
			}

			_ = discountEngine.ApplyDiscounts(order)
			done <- true
		}(string(rune('A' + i)))
	}

	// Wait for all goroutines to complete
	for i := 0; i < 10; i++ {
		<-done
	}
}
