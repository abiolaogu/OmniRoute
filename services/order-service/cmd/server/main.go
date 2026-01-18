// Package main provides the order management service entry point.
package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	port := getEnv("SERVER_PORT", "8102")
	router := gin.Default()

	router.GET("/health", func(c *gin.Context) { c.JSON(200, gin.H{"status": "healthy"}) })
	router.GET("/ready", func(c *gin.Context) { c.JSON(200, gin.H{"status": "ready"}) })

	// Orders
	orders := router.Group("/orders")
	{
		orders.GET("", listOrdersHandler)
		orders.POST("", createOrderHandler)
		orders.GET("/:id", getOrderHandler)
		orders.PUT("/:id", updateOrderHandler)
		orders.POST("/:id/cancel", cancelOrderHandler)
		orders.POST("/:id/confirm", confirmOrderHandler)
		orders.POST("/:id/ship", shipOrderHandler)
		orders.POST("/:id/complete", completeOrderHandler)
		orders.GET("/:id/history", getOrderHistoryHandler)
		orders.GET("/:id/items", getOrderItemsHandler)
		orders.POST("/:id/items", addOrderItemHandler)
		orders.DELETE("/:id/items/:itemId", removeOrderItemHandler)
		orders.POST("/:id/split", splitOrderHandler)
		orders.POST("/:id/merge", mergeOrdersHandler)
	}

	// Quotes
	quotes := router.Group("/quotes")
	{
		quotes.GET("", listQuotesHandler)
		quotes.POST("", createQuoteHandler)
		quotes.GET("/:id", getQuoteHandler)
		quotes.POST("/:id/convert", convertQuoteToOrderHandler)
	}

	// Returns
	returns := router.Group("/returns")
	{
		returns.GET("", listReturnsHandler)
		returns.POST("", createReturnHandler)
		returns.GET("/:id", getReturnHandler)
		returns.POST("/:id/approve", approveReturnHandler)
		returns.POST("/:id/reject", rejectReturnHandler)
		returns.POST("/:id/process", processReturnHandler)
	}

	// Shipments
	shipments := router.Group("/shipments")
	{
		shipments.GET("", listShipmentsHandler)
		shipments.POST("", createShipmentHandler)
		shipments.GET("/:id", getShipmentHandler)
		shipments.PUT("/:id/tracking", updateTrackingHandler)
		shipments.POST("/:id/deliver", markDeliveredHandler)
	}

	// Hasura Actions
	actions := router.Group("/actions")
	{
		actions.POST("/create-order", createOrderAction)
		actions.POST("/get-order-status", getOrderStatusAction)
		actions.POST("/calculate-order-total", calculateOrderTotalAction)
	}

	srv := &http.Server{Addr: ":" + port, Handler: router}
	go func() {
		log.Printf("Order service starting on port %s", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	srv.Shutdown(ctx)
}

func getEnv(key, defaultValue string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return defaultValue
}

// Handler stubs
func listOrdersHandler(c *gin.Context)          { c.JSON(200, gin.H{"orders": []interface{}{}}) }
func createOrderHandler(c *gin.Context)         { c.JSON(201, gin.H{"id": "order-id"}) }
func getOrderHandler(c *gin.Context)            { c.JSON(200, gin.H{"id": c.Param("id")}) }
func updateOrderHandler(c *gin.Context)         { c.JSON(200, gin.H{"message": "updated"}) }
func cancelOrderHandler(c *gin.Context)         { c.JSON(200, gin.H{"message": "cancelled"}) }
func confirmOrderHandler(c *gin.Context)        { c.JSON(200, gin.H{"message": "confirmed"}) }
func shipOrderHandler(c *gin.Context)           { c.JSON(200, gin.H{"message": "shipped"}) }
func completeOrderHandler(c *gin.Context)       { c.JSON(200, gin.H{"message": "completed"}) }
func getOrderHistoryHandler(c *gin.Context)     { c.JSON(200, gin.H{"history": []interface{}{}}) }
func getOrderItemsHandler(c *gin.Context)       { c.JSON(200, gin.H{"items": []interface{}{}}) }
func addOrderItemHandler(c *gin.Context)        { c.JSON(201, gin.H{"id": "item-id"}) }
func removeOrderItemHandler(c *gin.Context)     { c.JSON(200, gin.H{"message": "removed"}) }
func splitOrderHandler(c *gin.Context)          { c.JSON(200, gin.H{"orders": []string{}}) }
func mergeOrdersHandler(c *gin.Context)         { c.JSON(200, gin.H{"order_id": "merged"}) }
func listQuotesHandler(c *gin.Context)          { c.JSON(200, gin.H{"quotes": []interface{}{}}) }
func createQuoteHandler(c *gin.Context)         { c.JSON(201, gin.H{"id": "quote-id"}) }
func getQuoteHandler(c *gin.Context)            { c.JSON(200, gin.H{"id": c.Param("id")}) }
func convertQuoteToOrderHandler(c *gin.Context) { c.JSON(200, gin.H{"order_id": "new-order"}) }
func listReturnsHandler(c *gin.Context)         { c.JSON(200, gin.H{"returns": []interface{}{}}) }
func createReturnHandler(c *gin.Context)        { c.JSON(201, gin.H{"id": "return-id"}) }
func getReturnHandler(c *gin.Context)           { c.JSON(200, gin.H{"id": c.Param("id")}) }
func approveReturnHandler(c *gin.Context)       { c.JSON(200, gin.H{"message": "approved"}) }
func rejectReturnHandler(c *gin.Context)        { c.JSON(200, gin.H{"message": "rejected"}) }
func processReturnHandler(c *gin.Context)       { c.JSON(200, gin.H{"message": "processed"}) }
func listShipmentsHandler(c *gin.Context)       { c.JSON(200, gin.H{"shipments": []interface{}{}}) }
func createShipmentHandler(c *gin.Context)      { c.JSON(201, gin.H{"id": "shipment-id"}) }
func getShipmentHandler(c *gin.Context)         { c.JSON(200, gin.H{"id": c.Param("id")}) }
func updateTrackingHandler(c *gin.Context)      { c.JSON(200, gin.H{"message": "updated"}) }
func markDeliveredHandler(c *gin.Context)       { c.JSON(200, gin.H{"message": "delivered"}) }
func createOrderAction(c *gin.Context)          { c.JSON(200, gin.H{"order_id": "id"}) }
func getOrderStatusAction(c *gin.Context)       { c.JSON(200, gin.H{"status": "pending"}) }
func calculateOrderTotalAction(c *gin.Context)  { c.JSON(200, gin.H{"total": 0}) }
