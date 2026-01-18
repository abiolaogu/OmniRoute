// Package main provides the customer management service entry point.
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
	port := getEnv("SERVER_PORT", "8104")
	router := gin.Default()

	router.GET("/health", func(c *gin.Context) { c.JSON(200, gin.H{"status": "healthy"}) })
	router.GET("/ready", func(c *gin.Context) { c.JSON(200, gin.H{"status": "ready"}) })

	// Customers
	customers := router.Group("/customers")
	{
		customers.GET("", listCustomersHandler)
		customers.POST("", createCustomerHandler)
		customers.GET("/:id", getCustomerHandler)
		customers.PUT("/:id", updateCustomerHandler)
		customers.DELETE("/:id", deleteCustomerHandler)
		customers.GET("/:id/orders", getCustomerOrdersHandler)
		customers.GET("/:id/addresses", getCustomerAddressesHandler)
		customers.POST("/:id/addresses", addCustomerAddressHandler)
		customers.GET("/:id/credit", getCustomerCreditHandler)
		customers.PUT("/:id/credit", updateCustomerCreditHandler)
		customers.GET("/:id/segments", getCustomerSegmentsHandler)
		customers.GET("/:id/analytics", getCustomerAnalyticsHandler)
	}

	// Segments
	segments := router.Group("/segments")
	{
		segments.GET("", listSegmentsHandler)
		segments.POST("", createSegmentHandler)
		segments.GET("/:id", getSegmentHandler)
		segments.PUT("/:id", updateSegmentHandler)
		segments.DELETE("/:id", deleteSegmentHandler)
		segments.GET("/:id/customers", getSegmentCustomersHandler)
		segments.POST("/:id/recalculate", recalculateSegmentHandler)
	}

	// Loyalty
	loyalty := router.Group("/loyalty")
	{
		loyalty.GET("/programs", listLoyaltyProgramsHandler)
		loyalty.POST("/programs", createLoyaltyProgramHandler)
		loyalty.GET("/members/:customerId", getLoyaltyMemberHandler)
		loyalty.POST("/members/:customerId/points/earn", earnPointsHandler)
		loyalty.POST("/members/:customerId/points/redeem", redeemPointsHandler)
		loyalty.GET("/members/:customerId/transactions", getLoyaltyTransactionsHandler)
		loyalty.GET("/tiers", listLoyaltyTiersHandler)
	}

	// Tags
	tags := router.Group("/tags")
	{
		tags.GET("", listTagsHandler)
		tags.POST("", createTagHandler)
		tags.DELETE("/:id", deleteTagHandler)
		tags.POST("/customers/:customerId", addTagToCustomerHandler)
		tags.DELETE("/customers/:customerId/:tagId", removeTagFromCustomerHandler)
	}

	// Hasura Actions
	actions := router.Group("/actions")
	{
		actions.POST("/get-customer-360", getCustomer360Action)
		actions.POST("/check-credit-limit", checkCreditLimitAction)
		actions.POST("/calculate-loyalty-tier", calculateLoyaltyTierAction)
	}

	srv := &http.Server{Addr: ":" + port, Handler: router}
	go func() {
		log.Printf("Customer service starting on port %s", port)
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
func listCustomersHandler(c *gin.Context)        { c.JSON(200, gin.H{"customers": []interface{}{}}) }
func createCustomerHandler(c *gin.Context)       { c.JSON(201, gin.H{"id": "cust-id"}) }
func getCustomerHandler(c *gin.Context)          { c.JSON(200, gin.H{"id": c.Param("id")}) }
func updateCustomerHandler(c *gin.Context)       { c.JSON(200, gin.H{"message": "updated"}) }
func deleteCustomerHandler(c *gin.Context)       { c.JSON(200, gin.H{"message": "deleted"}) }
func getCustomerOrdersHandler(c *gin.Context)    { c.JSON(200, gin.H{"orders": []interface{}{}}) }
func getCustomerAddressesHandler(c *gin.Context) { c.JSON(200, gin.H{"addresses": []interface{}{}}) }
func addCustomerAddressHandler(c *gin.Context)   { c.JSON(201, gin.H{"id": "addr-id"}) }
func getCustomerCreditHandler(c *gin.Context)    { c.JSON(200, gin.H{"limit": 10000, "used": 0}) }
func updateCustomerCreditHandler(c *gin.Context) { c.JSON(200, gin.H{"message": "updated"}) }
func getCustomerSegmentsHandler(c *gin.Context)  { c.JSON(200, gin.H{"segments": []interface{}{}}) }
func getCustomerAnalyticsHandler(c *gin.Context) { c.JSON(200, gin.H{"analytics": nil}) }
func listSegmentsHandler(c *gin.Context)         { c.JSON(200, gin.H{"segments": []interface{}{}}) }
func createSegmentHandler(c *gin.Context)        { c.JSON(201, gin.H{"id": "seg-id"}) }
func getSegmentHandler(c *gin.Context)           { c.JSON(200, gin.H{"id": c.Param("id")}) }
func updateSegmentHandler(c *gin.Context)        { c.JSON(200, gin.H{"message": "updated"}) }
func deleteSegmentHandler(c *gin.Context)        { c.JSON(200, gin.H{"message": "deleted"}) }
func getSegmentCustomersHandler(c *gin.Context)  { c.JSON(200, gin.H{"customers": []interface{}{}}) }
func recalculateSegmentHandler(c *gin.Context)   { c.JSON(200, gin.H{"count": 0}) }
func listLoyaltyProgramsHandler(c *gin.Context)  { c.JSON(200, gin.H{"programs": []interface{}{}}) }
func createLoyaltyProgramHandler(c *gin.Context) { c.JSON(201, gin.H{"id": "prog-id"}) }
func getLoyaltyMemberHandler(c *gin.Context)     { c.JSON(200, gin.H{"points": 0, "tier": "bronze"}) }
func earnPointsHandler(c *gin.Context)           { c.JSON(200, gin.H{"new_balance": 100}) }
func redeemPointsHandler(c *gin.Context)         { c.JSON(200, gin.H{"new_balance": 50}) }
func getLoyaltyTransactionsHandler(c *gin.Context) {
	c.JSON(200, gin.H{"transactions": []interface{}{}})
}
func listLoyaltyTiersHandler(c *gin.Context)      { c.JSON(200, gin.H{"tiers": []interface{}{}}) }
func listTagsHandler(c *gin.Context)              { c.JSON(200, gin.H{"tags": []interface{}{}}) }
func createTagHandler(c *gin.Context)             { c.JSON(201, gin.H{"id": "tag-id"}) }
func deleteTagHandler(c *gin.Context)             { c.JSON(200, gin.H{"message": "deleted"}) }
func addTagToCustomerHandler(c *gin.Context)      { c.JSON(200, gin.H{"message": "added"}) }
func removeTagFromCustomerHandler(c *gin.Context) { c.JSON(200, gin.H{"message": "removed"}) }
func getCustomer360Action(c *gin.Context)         { c.JSON(200, gin.H{"customer": nil}) }
func checkCreditLimitAction(c *gin.Context)       { c.JSON(200, gin.H{"approved": true}) }
func calculateLoyaltyTierAction(c *gin.Context)   { c.JSON(200, gin.H{"tier": "gold"}) }
