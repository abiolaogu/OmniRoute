// Package main provides the inventory management service entry point.
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
	port := getEnv("SERVER_PORT", "8103")
	router := gin.Default()

	router.GET("/health", func(c *gin.Context) { c.JSON(200, gin.H{"status": "healthy"}) })
	router.GET("/ready", func(c *gin.Context) { c.JSON(200, gin.H{"status": "ready"}) })

	// Stock Levels
	stock := router.Group("/stock")
	{
		stock.GET("", getStockLevelsHandler)
		stock.GET("/product/:productId", getProductStockHandler)
		stock.GET("/warehouse/:warehouseId", getWarehouseStockHandler)
		stock.POST("/adjust", adjustStockHandler)
		stock.POST("/transfer", transferStockHandler)
		stock.GET("/low", getLowStockHandler)
		stock.GET("/out-of-stock", getOutOfStockHandler)
	}

	// Reservations
	reservations := router.Group("/reservations")
	{
		reservations.GET("", listReservationsHandler)
		reservations.POST("", createReservationHandler)
		reservations.GET("/:id", getReservationHandler)
		reservations.DELETE("/:id", cancelReservationHandler)
		reservations.POST("/:id/confirm", confirmReservationHandler)
	}

	// Stock Movements
	movements := router.Group("/movements")
	{
		movements.GET("", listMovementsHandler)
		movements.GET("/:id", getMovementHandler)
		movements.POST("/receive", receiveStockHandler)
		movements.POST("/dispatch", dispatchStockHandler)
	}

	// Batches
	batches := router.Group("/batches")
	{
		batches.GET("", listBatchesHandler)
		batches.POST("", createBatchHandler)
		batches.GET("/:id", getBatchHandler)
		batches.GET("/expiring", getExpiringBatchesHandler)
	}

	// Stock Takes / Cycle Counts
	stocktakes := router.Group("/stocktakes")
	{
		stocktakes.GET("", listStocktakesHandler)
		stocktakes.POST("", createStocktakeHandler)
		stocktakes.GET("/:id", getStocktakeHandler)
		stocktakes.POST("/:id/count", submitCountHandler)
		stocktakes.POST("/:id/complete", completeStocktakeHandler)
		stocktakes.GET("/:id/variances", getVariancesHandler)
	}

	// Hasura Actions
	actions := router.Group("/actions")
	{
		actions.POST("/check-availability", checkAvailabilityAction)
		actions.POST("/reserve-stock", reserveStockAction)
		actions.POST("/get-stock-level", getStockLevelAction)
	}

	srv := &http.Server{Addr: ":" + port, Handler: router}
	go func() {
		log.Printf("Inventory service starting on port %s", port)
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
func getStockLevelsHandler(c *gin.Context)     { c.JSON(200, gin.H{"stock": []interface{}{}}) }
func getProductStockHandler(c *gin.Context)    { c.JSON(200, gin.H{"quantity": 100}) }
func getWarehouseStockHandler(c *gin.Context)  { c.JSON(200, gin.H{"items": []interface{}{}}) }
func adjustStockHandler(c *gin.Context)        { c.JSON(200, gin.H{"message": "adjusted"}) }
func transferStockHandler(c *gin.Context)      { c.JSON(200, gin.H{"transfer_id": "id"}) }
func getLowStockHandler(c *gin.Context)        { c.JSON(200, gin.H{"items": []interface{}{}}) }
func getOutOfStockHandler(c *gin.Context)      { c.JSON(200, gin.H{"items": []interface{}{}}) }
func listReservationsHandler(c *gin.Context)   { c.JSON(200, gin.H{"reservations": []interface{}{}}) }
func createReservationHandler(c *gin.Context)  { c.JSON(201, gin.H{"id": "res-id"}) }
func getReservationHandler(c *gin.Context)     { c.JSON(200, gin.H{"id": c.Param("id")}) }
func cancelReservationHandler(c *gin.Context)  { c.JSON(200, gin.H{"message": "cancelled"}) }
func confirmReservationHandler(c *gin.Context) { c.JSON(200, gin.H{"message": "confirmed"}) }
func listMovementsHandler(c *gin.Context)      { c.JSON(200, gin.H{"movements": []interface{}{}}) }
func getMovementHandler(c *gin.Context)        { c.JSON(200, gin.H{"id": c.Param("id")}) }
func receiveStockHandler(c *gin.Context)       { c.JSON(200, gin.H{"received": true}) }
func dispatchStockHandler(c *gin.Context)      { c.JSON(200, gin.H{"dispatched": true}) }
func listBatchesHandler(c *gin.Context)        { c.JSON(200, gin.H{"batches": []interface{}{}}) }
func createBatchHandler(c *gin.Context)        { c.JSON(201, gin.H{"id": "batch-id"}) }
func getBatchHandler(c *gin.Context)           { c.JSON(200, gin.H{"id": c.Param("id")}) }
func getExpiringBatchesHandler(c *gin.Context) { c.JSON(200, gin.H{"batches": []interface{}{}}) }
func listStocktakesHandler(c *gin.Context)     { c.JSON(200, gin.H{"stocktakes": []interface{}{}}) }
func createStocktakeHandler(c *gin.Context)    { c.JSON(201, gin.H{"id": "stocktake-id"}) }
func getStocktakeHandler(c *gin.Context)       { c.JSON(200, gin.H{"id": c.Param("id")}) }
func submitCountHandler(c *gin.Context)        { c.JSON(200, gin.H{"message": "submitted"}) }
func completeStocktakeHandler(c *gin.Context)  { c.JSON(200, gin.H{"message": "completed"}) }
func getVariancesHandler(c *gin.Context)       { c.JSON(200, gin.H{"variances": []interface{}{}}) }
func checkAvailabilityAction(c *gin.Context)   { c.JSON(200, gin.H{"available": true}) }
func reserveStockAction(c *gin.Context)        { c.JSON(200, gin.H{"reservation_id": "id"}) }
func getStockLevelAction(c *gin.Context)       { c.JSON(200, gin.H{"level": 100}) }
