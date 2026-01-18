// Package main provides the market intelligence service entry point.
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
	port := getEnv("SERVER_PORT", "8110")
	router := gin.Default()

	router.GET("/health", func(c *gin.Context) { c.JSON(200, gin.H{"status": "healthy"}) })
	router.GET("/ready", func(c *gin.Context) { c.JSON(200, gin.H{"status": "ready"}) })

	// Price Intelligence
	prices := router.Group("/prices")
	{
		prices.GET("/compare", comparePricesHandler)
		prices.GET("/trends", priceTrendsHandler)
		prices.GET("/alerts", priceAlertsHandler)
		prices.GET("/competitor/:productId", competitorPricesHandler)
		prices.POST("/optimize", priceOptimizationHandler)
	}

	// Market Analysis
	market := router.Group("/market")
	{
		market.GET("/size", marketSizeHandler)
		market.GET("/share", marketShareHandler)
		market.GET("/growth", marketGrowthHandler)
		market.GET("/segments", marketSegmentsHandler)
		market.GET("/opportunities", marketOpportunitiesHandler)
	}

	// Competitor Intelligence
	competitors := router.Group("/competitors")
	{
		competitors.GET("", listCompetitorsHandler)
		competitors.GET("/:id", getCompetitorHandler)
		competitors.GET("/:id/products", getCompetitorProductsHandler)
		competitors.GET("/:id/pricing", getCompetitorPricingHandler)
		competitors.GET("/activity", competitorActivityHandler)
	}

	// Consumer Insights
	insights := router.Group("/insights")
	{
		insights.GET("/trends", consumerTrendsHandler)
		insights.GET("/preferences", consumerPreferencesHandler)
		insights.GET("/sentiment", sentimentAnalysisHandler)
		insights.GET("/demographics", demographicsHandler)
	}

	// Benchmarking
	benchmarks := router.Group("/benchmarks")
	{
		benchmarks.GET("/industry", industryBenchmarksHandler)
		benchmarks.GET("/performance", performanceBenchmarksHandler)
		benchmarks.GET("/pricing", pricingBenchmarksHandler)
	}

	// Hasura Actions
	actions := router.Group("/actions")
	{
		actions.POST("/get-market-summary", getMarketSummaryAction)
		actions.POST("/analyze-opportunity", analyzeOpportunityAction)
	}

	srv := &http.Server{Addr: ":" + port, Handler: router}
	go func() {
		log.Printf("Market Intel service starting on port %s", port)
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
func comparePricesHandler(c *gin.Context)         { c.JSON(200, gin.H{"comparison": nil}) }
func priceTrendsHandler(c *gin.Context)           { c.JSON(200, gin.H{"trends": nil}) }
func priceAlertsHandler(c *gin.Context)           { c.JSON(200, gin.H{"alerts": []interface{}{}}) }
func competitorPricesHandler(c *gin.Context)      { c.JSON(200, gin.H{"prices": nil}) }
func priceOptimizationHandler(c *gin.Context)     { c.JSON(200, gin.H{"suggestions": nil}) }
func marketSizeHandler(c *gin.Context)            { c.JSON(200, gin.H{"size_usd": 1000000000}) }
func marketShareHandler(c *gin.Context)           { c.JSON(200, gin.H{"share": 0.15}) }
func marketGrowthHandler(c *gin.Context)          { c.JSON(200, gin.H{"growth_rate": 0.12}) }
func marketSegmentsHandler(c *gin.Context)        { c.JSON(200, gin.H{"segments": []interface{}{}}) }
func marketOpportunitiesHandler(c *gin.Context)   { c.JSON(200, gin.H{"opportunities": []interface{}{}}) }
func listCompetitorsHandler(c *gin.Context)       { c.JSON(200, gin.H{"competitors": []interface{}{}}) }
func getCompetitorHandler(c *gin.Context)         { c.JSON(200, gin.H{"id": c.Param("id")}) }
func getCompetitorProductsHandler(c *gin.Context) { c.JSON(200, gin.H{"products": []interface{}{}}) }
func getCompetitorPricingHandler(c *gin.Context)  { c.JSON(200, gin.H{"pricing": nil}) }
func competitorActivityHandler(c *gin.Context)    { c.JSON(200, gin.H{"activity": []interface{}{}}) }
func consumerTrendsHandler(c *gin.Context)        { c.JSON(200, gin.H{"trends": nil}) }
func consumerPreferencesHandler(c *gin.Context)   { c.JSON(200, gin.H{"preferences": nil}) }
func sentimentAnalysisHandler(c *gin.Context)     { c.JSON(200, gin.H{"sentiment": nil}) }
func demographicsHandler(c *gin.Context)          { c.JSON(200, gin.H{"demographics": nil}) }
func industryBenchmarksHandler(c *gin.Context)    { c.JSON(200, gin.H{"benchmarks": nil}) }
func performanceBenchmarksHandler(c *gin.Context) { c.JSON(200, gin.H{"benchmarks": nil}) }
func pricingBenchmarksHandler(c *gin.Context)     { c.JSON(200, gin.H{"benchmarks": nil}) }
func getMarketSummaryAction(c *gin.Context)       { c.JSON(200, gin.H{"summary": nil}) }
func analyzeOpportunityAction(c *gin.Context)     { c.JSON(200, gin.H{"analysis": nil}) }
