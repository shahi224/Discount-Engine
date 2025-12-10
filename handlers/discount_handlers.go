package handlers

import (
	"discountengine/internal/engine"
	"discountengine/internal/models"
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
)

// DiscountHandler handles HTTP requests for discount calculations
type DiscountHandler struct {
	engine *engine.DiscountEngine
}

// NewDiscountHandler creates a new DiscountHandler
func NewDiscountHandler(engine *engine.DiscountEngine) *DiscountHandler {
	return &DiscountHandler{
		engine: engine,
	}
}

// CalculateDiscount handles POST / discount
func (h *DiscountHandler) CalculateDiscount(c *gin.Context) {
	var order models.Order

	if err := c.ShouldBindJSON(&order); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// Validate required fields
	if order.OderID == "" || order.OrderTotal <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "order_id and order_total (positive) are required"})
		return
	}

	// Apply discounts
	result := h.engine.ApplyDiscounts(order)

	c.JSON(http.StatusOK, result)
}

// ReloadRules handles POST /reload-rules (for admin purposes)
func (h *DiscountHandler) ReloadRules(c *gin.Context) {
	if err := h.engine.ReloadRules(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to reload rules:" + err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Rules reload successfully"})
}

// GetRules handels GET /rules (for debugging)
func (h *DiscountHandler) GetRules(c *gin.Context) {
	rules := h.engine.GetRules()

	// Pretty print JSON
	jsonData, err := json.MarshalIndent(rules, "", "  ")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to serialize rules"})
		return
	}
	c.Data(http.StatusOK, "application/json", jsonData)
}

// Healthcheck handles GET /health
func (h *DiscountHandler) Healthcheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "healthy", "service": "discount-engine"})
}
