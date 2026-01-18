// Package main provides the analytics service entry point.
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
	port := getEnv("SERVER_PORT", "8105")
	router := gin.Default()

	router.GET("/health", func(c *gin.Context) { c.JSON(200, gin.H{"status": "healthy"}) })
	router.GET("/ready", func(c *gin.Context) { c.JSON(200, gin.H{"status": "ready"}) })

	// Dashboards
	dashboards := router.Group("/dashboards")
	{
		dashboards.GET("/executive", executiveDashboardHandler)
		dashboards.GET("/sales", salesDashboardHandler)
		dashboards.GET("/operations", operationsDashboardHandler)
		dashboards.GET("/finance", financeDashboardHandler)
		dashboards.GET("/logistics", logisticsDashboardHandler)
	}

	// Reports
	reports := router.Group("/reports")
	{
		reports.GET("", listReportsHandler)
		reports.POST("", createReportHandler)
		reports.GET("/:id", getReportHandler)
		reports.GET("/:id/download", downloadReportHandler)
		reports.POST("/:id/schedule", scheduleReportHandler)
		reports.GET("/sales/daily", dailySalesReportHandler)
		reports.GET("/sales/weekly", weeklySalesReportHandler)
		reports.GET("/sales/monthly", monthlySalesReportHandler)
		reports.GET("/inventory/status", inventoryStatusReportHandler)
		reports.GET("/customers/acquisition", customerAcquisitionReportHandler)
		reports.GET("/workers/performance", workerPerformanceReportHandler)
	}

	// Metrics
	metrics := router.Group("/metrics")
	{
		metrics.GET("/gmv", gmvMetricsHandler)
		metrics.GET("/orders", orderMetricsHandler)
		metrics.GET("/customers", customerMetricsHandler)
		metrics.GET("/products", productMetricsHandler)
		metrics.GET("/workers", workerMetricsHandler)
		metrics.GET("/delivery", deliveryMetricsHandler)
		metrics.GET("/revenue", revenueMetricsHandler)
		metrics.GET("/custom", customMetricsHandler)
	}

	// KPIs
	kpis := router.Group("/kpis")
	{
		kpis.GET("", listKPIsHandler)
		kpis.GET("/:id", getKPIHandler)
		kpis.GET("/trends", kpiTrendsHandler)
		kpis.GET("/targets", kpiTargetsHandler)
	}

	// Cohorts
	cohorts := router.Group("/cohorts")
	{
		cohorts.GET("", listCohortsHandler)
		cohorts.POST("", createCohortHandler)
		cohorts.GET("/:id", getCohortHandler)
		cohorts.GET("/:id/retention", cohortRetentionHandler)
	}

	// Funnels
	funnels := router.Group("/funnels")
	{
		funnels.GET("", listFunnelsHandler)
		funnels.POST("", createFunnelHandler)
		funnels.GET("/:id", getFunnelHandler)
		funnels.GET("/:id/conversion", funnelConversionHandler)
	}

	// Alerts
	alerts := router.Group("/alerts")
	{
		alerts.GET("", listAlertsHandler)
		alerts.POST("", createAlertHandler)
		alerts.PUT("/:id", updateAlertHandler)
		alerts.DELETE("/:id", deleteAlertHandler)
	}

	// Hasura Actions
	actions := router.Group("/actions")
	{
		actions.POST("/get-dashboard-data", getDashboardDataAction)
		actions.POST("/run-report", runReportAction)
		actions.POST("/get-kpi-summary", getKPISummaryAction)
	}

	srv := &http.Server{Addr: ":" + port, Handler: router}
	go func() {
		log.Printf("Analytics service starting on port %s", port)
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
func executiveDashboardHandler(c *gin.Context)        { c.JSON(200, gin.H{"metrics": nil}) }
func salesDashboardHandler(c *gin.Context)            { c.JSON(200, gin.H{"metrics": nil}) }
func operationsDashboardHandler(c *gin.Context)       { c.JSON(200, gin.H{"metrics": nil}) }
func financeDashboardHandler(c *gin.Context)          { c.JSON(200, gin.H{"metrics": nil}) }
func logisticsDashboardHandler(c *gin.Context)        { c.JSON(200, gin.H{"metrics": nil}) }
func listReportsHandler(c *gin.Context)               { c.JSON(200, gin.H{"reports": []interface{}{}}) }
func createReportHandler(c *gin.Context)              { c.JSON(201, gin.H{"id": "report-id"}) }
func getReportHandler(c *gin.Context)                 { c.JSON(200, gin.H{"id": c.Param("id")}) }
func downloadReportHandler(c *gin.Context)            { c.JSON(200, gin.H{"url": "download-url"}) }
func scheduleReportHandler(c *gin.Context)            { c.JSON(200, gin.H{"scheduled": true}) }
func dailySalesReportHandler(c *gin.Context)          { c.JSON(200, gin.H{"data": nil}) }
func weeklySalesReportHandler(c *gin.Context)         { c.JSON(200, gin.H{"data": nil}) }
func monthlySalesReportHandler(c *gin.Context)        { c.JSON(200, gin.H{"data": nil}) }
func inventoryStatusReportHandler(c *gin.Context)     { c.JSON(200, gin.H{"data": nil}) }
func customerAcquisitionReportHandler(c *gin.Context) { c.JSON(200, gin.H{"data": nil}) }
func workerPerformanceReportHandler(c *gin.Context)   { c.JSON(200, gin.H{"data": nil}) }
func gmvMetricsHandler(c *gin.Context)                { c.JSON(200, gin.H{"gmv": 0}) }
func orderMetricsHandler(c *gin.Context)              { c.JSON(200, gin.H{"count": 0}) }
func customerMetricsHandler(c *gin.Context)           { c.JSON(200, gin.H{"count": 0}) }
func productMetricsHandler(c *gin.Context)            { c.JSON(200, gin.H{"count": 0}) }
func workerMetricsHandler(c *gin.Context)             { c.JSON(200, gin.H{"count": 0}) }
func deliveryMetricsHandler(c *gin.Context)           { c.JSON(200, gin.H{"count": 0}) }
func revenueMetricsHandler(c *gin.Context)            { c.JSON(200, gin.H{"revenue": 0}) }
func customMetricsHandler(c *gin.Context)             { c.JSON(200, gin.H{"data": nil}) }
func listKPIsHandler(c *gin.Context)                  { c.JSON(200, gin.H{"kpis": []interface{}{}}) }
func getKPIHandler(c *gin.Context)                    { c.JSON(200, gin.H{"id": c.Param("id")}) }
func kpiTrendsHandler(c *gin.Context)                 { c.JSON(200, gin.H{"trends": nil}) }
func kpiTargetsHandler(c *gin.Context)                { c.JSON(200, gin.H{"targets": nil}) }
func listCohortsHandler(c *gin.Context)               { c.JSON(200, gin.H{"cohorts": []interface{}{}}) }
func createCohortHandler(c *gin.Context)              { c.JSON(201, gin.H{"id": "cohort-id"}) }
func getCohortHandler(c *gin.Context)                 { c.JSON(200, gin.H{"id": c.Param("id")}) }
func cohortRetentionHandler(c *gin.Context)           { c.JSON(200, gin.H{"retention": nil}) }
func listFunnelsHandler(c *gin.Context)               { c.JSON(200, gin.H{"funnels": []interface{}{}}) }
func createFunnelHandler(c *gin.Context)              { c.JSON(201, gin.H{"id": "funnel-id"}) }
func getFunnelHandler(c *gin.Context)                 { c.JSON(200, gin.H{"id": c.Param("id")}) }
func funnelConversionHandler(c *gin.Context)          { c.JSON(200, gin.H{"conversion": nil}) }
func listAlertsHandler(c *gin.Context)                { c.JSON(200, gin.H{"alerts": []interface{}{}}) }
func createAlertHandler(c *gin.Context)               { c.JSON(201, gin.H{"id": "alert-id"}) }
func updateAlertHandler(c *gin.Context)               { c.JSON(200, gin.H{"message": "updated"}) }
func deleteAlertHandler(c *gin.Context)               { c.JSON(200, gin.H{"message": "deleted"}) }
func getDashboardDataAction(c *gin.Context)           { c.JSON(200, gin.H{"data": nil}) }
func runReportAction(c *gin.Context)                  { c.JSON(200, gin.H{"report_id": "id"}) }
func getKPISummaryAction(c *gin.Context)              { c.JSON(200, gin.H{"summary": nil}) }
