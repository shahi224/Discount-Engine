package main

import (
	"discountengine/handlers"
	"discountengine/internal/engine"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize discount engine
	discountEngine := engine.NewDiscountEngine("config/rules.json")

	// Initialize HTTP handler
	discountHandler := handlers.NewDiscountHandler(discountEngine)

	// create gin router
	r := gin.Default()

	// API routes
	r.GET("/health", discountHandler.Healthcheck)
	r.POST("/discount", discountHandler.CalculateDiscount)
	r.GET("/rules", discountHandler.GetRules)
	r.POST("/reload-rules", discountHandler.ReloadRules)

	// start server
	fmt.Println("Discount engine server starting on port 8080...")
	if err := r.Run(":8080"); err != nil {
		log.Fatal("failed to start server", err)
	}
}
