// Package main provides the fleet management service entry point.
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
	port := getEnv("SERVER_PORT", "8106")
	router := gin.Default()

	router.GET("/health", func(c *gin.Context) { c.JSON(200, gin.H{"status": "healthy"}) })
	router.GET("/ready", func(c *gin.Context) { c.JSON(200, gin.H{"status": "ready"}) })

	// Vehicles
	vehicles := router.Group("/vehicles")
	{
		vehicles.GET("", listVehiclesHandler)
		vehicles.POST("", createVehicleHandler)
		vehicles.GET("/:id", getVehicleHandler)
		vehicles.PUT("/:id", updateVehicleHandler)
		vehicles.DELETE("/:id", deleteVehicleHandler)
		vehicles.GET("/:id/location", getVehicleLocationHandler)
		vehicles.GET("/:id/telemetry", getVehicleTelemetryHandler)
		vehicles.GET("/:id/trips", getVehicleTripsHandler)
		vehicles.GET("/:id/maintenance", getVehicleMaintenanceHandler)
		vehicles.POST("/:id/assign", assignVehicleHandler)
		vehicles.POST("/:id/unassign", unassignVehicleHandler)
	}

	// Drivers
	drivers := router.Group("/drivers")
	{
		drivers.GET("", listDriversHandler)
		drivers.POST("", createDriverHandler)
		drivers.GET("/:id", getDriverHandler)
		drivers.PUT("/:id", updateDriverHandler)
		drivers.GET("/:id/score", getDriverScoreHandler)
		drivers.GET("/:id/trips", getDriverTripsHandler)
		drivers.GET("/:id/violations", getDriverViolationsHandler)
	}

	// Telemetry
	telemetry := router.Group("/telemetry")
	{
		telemetry.POST("/ingest", ingestTelemetryHandler)
		telemetry.GET("/stream", streamTelemetryHandler)
		telemetry.GET("/history", getTelemetryHistoryHandler)
	}

	// Maintenance
	maintenance := router.Group("/maintenance")
	{
		maintenance.GET("", listMaintenanceHandler)
		maintenance.POST("", createMaintenanceHandler)
		maintenance.GET("/:id", getMaintenanceHandler)
		maintenance.PUT("/:id", updateMaintenanceHandler)
		maintenance.POST("/:id/complete", completeMaintenanceHandler)
		maintenance.GET("/schedule", getMaintenanceScheduleHandler)
		maintenance.GET("/overdue", getOverdueMaintenanceHandler)
	}

	// Fuel
	fuel := router.Group("/fuel")
	{
		fuel.GET("", listFuelRecordsHandler)
		fuel.POST("", addFuelRecordHandler)
		fuel.GET("/consumption", getFuelConsumptionHandler)
		fuel.GET("/anomalies", getFuelAnomaliesHandler)
	}

	// Geofencing
	geofences := router.Group("/geofences")
	{
		geofences.GET("", listGeofencesHandler)
		geofences.POST("", createGeofenceHandler)
		geofences.GET("/:id", getGeofenceHandler)
		geofences.DELETE("/:id", deleteGeofenceHandler)
		geofences.GET("/violations", getGeofenceViolationsHandler)
	}

	// Trips
	trips := router.Group("/trips")
	{
		trips.GET("", listTripsHandler)
		trips.GET("/:id", getTripHandler)
		trips.GET("/:id/route", getTripRouteHandler)
		trips.GET("/:id/events", getTripEventsHandler)
	}

	// Hasura Actions
	actions := router.Group("/actions")
	{
		actions.POST("/get-vehicle-status", getVehicleStatusAction)
		actions.POST("/track-vehicle", trackVehicleAction)
		actions.POST("/get-fleet-summary", getFleetSummaryAction)
	}

	srv := &http.Server{Addr: ":" + port, Handler: router}
	go func() {
		log.Printf("Fleet service starting on port %s", port)
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
func listVehiclesHandler(c *gin.Context)           { c.JSON(200, gin.H{"vehicles": []interface{}{}}) }
func createVehicleHandler(c *gin.Context)          { c.JSON(201, gin.H{"id": "veh-id"}) }
func getVehicleHandler(c *gin.Context)             { c.JSON(200, gin.H{"id": c.Param("id")}) }
func updateVehicleHandler(c *gin.Context)          { c.JSON(200, gin.H{"message": "updated"}) }
func deleteVehicleHandler(c *gin.Context)          { c.JSON(200, gin.H{"message": "deleted"}) }
func getVehicleLocationHandler(c *gin.Context)     { c.JSON(200, gin.H{"lat": 0, "lng": 0}) }
func getVehicleTelemetryHandler(c *gin.Context)    { c.JSON(200, gin.H{"telemetry": nil}) }
func getVehicleTripsHandler(c *gin.Context)        { c.JSON(200, gin.H{"trips": []interface{}{}}) }
func getVehicleMaintenanceHandler(c *gin.Context)  { c.JSON(200, gin.H{"records": []interface{}{}}) }
func assignVehicleHandler(c *gin.Context)          { c.JSON(200, gin.H{"message": "assigned"}) }
func unassignVehicleHandler(c *gin.Context)        { c.JSON(200, gin.H{"message": "unassigned"}) }
func listDriversHandler(c *gin.Context)            { c.JSON(200, gin.H{"drivers": []interface{}{}}) }
func createDriverHandler(c *gin.Context)           { c.JSON(201, gin.H{"id": "drv-id"}) }
func getDriverHandler(c *gin.Context)              { c.JSON(200, gin.H{"id": c.Param("id")}) }
func updateDriverHandler(c *gin.Context)           { c.JSON(200, gin.H{"message": "updated"}) }
func getDriverScoreHandler(c *gin.Context)         { c.JSON(200, gin.H{"score": 85}) }
func getDriverTripsHandler(c *gin.Context)         { c.JSON(200, gin.H{"trips": []interface{}{}}) }
func getDriverViolationsHandler(c *gin.Context)    { c.JSON(200, gin.H{"violations": []interface{}{}}) }
func ingestTelemetryHandler(c *gin.Context)        { c.JSON(200, gin.H{"ingested": true}) }
func streamTelemetryHandler(c *gin.Context)        { c.JSON(200, gin.H{"stream": "ws://..."}) }
func getTelemetryHistoryHandler(c *gin.Context)    { c.JSON(200, gin.H{"history": []interface{}{}}) }
func listMaintenanceHandler(c *gin.Context)        { c.JSON(200, gin.H{"records": []interface{}{}}) }
func createMaintenanceHandler(c *gin.Context)      { c.JSON(201, gin.H{"id": "maint-id"}) }
func getMaintenanceHandler(c *gin.Context)         { c.JSON(200, gin.H{"id": c.Param("id")}) }
func updateMaintenanceHandler(c *gin.Context)      { c.JSON(200, gin.H{"message": "updated"}) }
func completeMaintenanceHandler(c *gin.Context)    { c.JSON(200, gin.H{"message": "completed"}) }
func getMaintenanceScheduleHandler(c *gin.Context) { c.JSON(200, gin.H{"schedule": []interface{}{}}) }
func getOverdueMaintenanceHandler(c *gin.Context)  { c.JSON(200, gin.H{"overdue": []interface{}{}}) }
func listFuelRecordsHandler(c *gin.Context)        { c.JSON(200, gin.H{"records": []interface{}{}}) }
func addFuelRecordHandler(c *gin.Context)          { c.JSON(201, gin.H{"id": "fuel-id"}) }
func getFuelConsumptionHandler(c *gin.Context)     { c.JSON(200, gin.H{"consumption": nil}) }
func getFuelAnomaliesHandler(c *gin.Context)       { c.JSON(200, gin.H{"anomalies": []interface{}{}}) }
func listGeofencesHandler(c *gin.Context)          { c.JSON(200, gin.H{"geofences": []interface{}{}}) }
func createGeofenceHandler(c *gin.Context)         { c.JSON(201, gin.H{"id": "geo-id"}) }
func getGeofenceHandler(c *gin.Context)            { c.JSON(200, gin.H{"id": c.Param("id")}) }
func deleteGeofenceHandler(c *gin.Context)         { c.JSON(200, gin.H{"message": "deleted"}) }
func getGeofenceViolationsHandler(c *gin.Context)  { c.JSON(200, gin.H{"violations": []interface{}{}}) }
func listTripsHandler(c *gin.Context)              { c.JSON(200, gin.H{"trips": []interface{}{}}) }
func getTripHandler(c *gin.Context)                { c.JSON(200, gin.H{"id": c.Param("id")}) }
func getTripRouteHandler(c *gin.Context)           { c.JSON(200, gin.H{"route": nil}) }
func getTripEventsHandler(c *gin.Context)          { c.JSON(200, gin.H{"events": []interface{}{}}) }
func getVehicleStatusAction(c *gin.Context)        { c.JSON(200, gin.H{"status": "active"}) }
func trackVehicleAction(c *gin.Context)            { c.JSON(200, gin.H{"location": nil}) }
func getFleetSummaryAction(c *gin.Context)         { c.JSON(200, gin.H{"summary": nil}) }
