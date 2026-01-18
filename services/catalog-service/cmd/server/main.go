// Package main provides the catalog service entry point.
// This service manages product catalogs, categories, variants, and media.
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
	port := getEnv("SERVER_PORT", "8101")

	router := gin.Default()

	// Health endpoints
	router.GET("/health", healthHandler)
	router.GET("/ready", readyHandler)

	// Categories
	categories := router.Group("/categories")
	{
		categories.GET("", listCategoriesHandler)
		categories.POST("", createCategoryHandler)
		categories.GET("/:id", getCategoryHandler)
		categories.PUT("/:id", updateCategoryHandler)
		categories.DELETE("/:id", deleteCategoryHandler)
		categories.GET("/:id/products", getCategoryProductsHandler)
		categories.GET("/tree", getCategoryTreeHandler)
	}

	// Products
	products := router.Group("/products")
	{
		products.GET("", listProductsHandler)
		products.POST("", createProductHandler)
		products.GET("/:id", getProductHandler)
		products.PUT("/:id", updateProductHandler)
		products.DELETE("/:id", deleteProductHandler)
		products.GET("/:id/variants", getProductVariantsHandler)
		products.POST("/:id/variants", createProductVariantHandler)
		products.GET("/:id/media", getProductMediaHandler)
		products.POST("/:id/media", uploadProductMediaHandler)
		products.DELETE("/:id/media/:mediaId", deleteProductMediaHandler)
		products.GET("/:id/related", getRelatedProductsHandler)
		products.POST("/:id/duplicate", duplicateProductHandler)
	}

	// Variants
	variants := router.Group("/variants")
	{
		variants.GET("/:id", getVariantHandler)
		variants.PUT("/:id", updateVariantHandler)
		variants.DELETE("/:id", deleteVariantHandler)
	}

	// Brands
	brands := router.Group("/brands")
	{
		brands.GET("", listBrandsHandler)
		brands.POST("", createBrandHandler)
		brands.GET("/:id", getBrandHandler)
		brands.PUT("/:id", updateBrandHandler)
		brands.DELETE("/:id", deleteBrandHandler)
	}

	// Attributes
	attributes := router.Group("/attributes")
	{
		attributes.GET("", listAttributesHandler)
		attributes.POST("", createAttributeHandler)
		attributes.GET("/:id", getAttributeHandler)
		attributes.PUT("/:id", updateAttributeHandler)
		attributes.DELETE("/:id", deleteAttributeHandler)
	}

	// Bulk operations
	bulk := router.Group("/bulk")
	{
		bulk.POST("/import", bulkImportHandler)
		bulk.GET("/export", bulkExportHandler)
		bulk.POST("/update", bulkUpdateHandler)
	}

	// Search
	router.GET("/search", searchProductsHandler)

	// Hasura Actions
	actions := router.Group("/actions")
	{
		actions.POST("/get-product-details", getProductDetailsAction)
		actions.POST("/search-products", searchProductsAction)
		actions.POST("/get-category-tree", getCategoryTreeAction)
	}

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: router,
	}

	go func() {
		log.Printf("Catalog service starting on port %s", port)
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
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// Handler stubs
func healthHandler(c *gin.Context)               { c.JSON(200, gin.H{"status": "healthy"}) }
func readyHandler(c *gin.Context)                { c.JSON(200, gin.H{"status": "ready"}) }
func listCategoriesHandler(c *gin.Context)       { c.JSON(200, gin.H{"categories": []string{}}) }
func createCategoryHandler(c *gin.Context)       { c.JSON(201, gin.H{"id": "new-id"}) }
func getCategoryHandler(c *gin.Context)          { c.JSON(200, gin.H{"id": c.Param("id")}) }
func updateCategoryHandler(c *gin.Context)       { c.JSON(200, gin.H{"message": "updated"}) }
func deleteCategoryHandler(c *gin.Context)       { c.JSON(200, gin.H{"message": "deleted"}) }
func getCategoryProductsHandler(c *gin.Context)  { c.JSON(200, gin.H{"products": []string{}}) }
func getCategoryTreeHandler(c *gin.Context)      { c.JSON(200, gin.H{"tree": []string{}}) }
func listProductsHandler(c *gin.Context)         { c.JSON(200, gin.H{"products": []string{}}) }
func createProductHandler(c *gin.Context)        { c.JSON(201, gin.H{"id": "new-id"}) }
func getProductHandler(c *gin.Context)           { c.JSON(200, gin.H{"id": c.Param("id")}) }
func updateProductHandler(c *gin.Context)        { c.JSON(200, gin.H{"message": "updated"}) }
func deleteProductHandler(c *gin.Context)        { c.JSON(200, gin.H{"message": "deleted"}) }
func getProductVariantsHandler(c *gin.Context)   { c.JSON(200, gin.H{"variants": []string{}}) }
func createProductVariantHandler(c *gin.Context) { c.JSON(201, gin.H{"id": "new-id"}) }
func getProductMediaHandler(c *gin.Context)      { c.JSON(200, gin.H{"media": []string{}}) }
func uploadProductMediaHandler(c *gin.Context)   { c.JSON(201, gin.H{"id": "new-id"}) }
func deleteProductMediaHandler(c *gin.Context)   { c.JSON(200, gin.H{"message": "deleted"}) }
func getRelatedProductsHandler(c *gin.Context)   { c.JSON(200, gin.H{"products": []string{}}) }
func duplicateProductHandler(c *gin.Context)     { c.JSON(201, gin.H{"id": "new-id"}) }
func getVariantHandler(c *gin.Context)           { c.JSON(200, gin.H{"id": c.Param("id")}) }
func updateVariantHandler(c *gin.Context)        { c.JSON(200, gin.H{"message": "updated"}) }
func deleteVariantHandler(c *gin.Context)        { c.JSON(200, gin.H{"message": "deleted"}) }
func listBrandsHandler(c *gin.Context)           { c.JSON(200, gin.H{"brands": []string{}}) }
func createBrandHandler(c *gin.Context)          { c.JSON(201, gin.H{"id": "new-id"}) }
func getBrandHandler(c *gin.Context)             { c.JSON(200, gin.H{"id": c.Param("id")}) }
func updateBrandHandler(c *gin.Context)          { c.JSON(200, gin.H{"message": "updated"}) }
func deleteBrandHandler(c *gin.Context)          { c.JSON(200, gin.H{"message": "deleted"}) }
func listAttributesHandler(c *gin.Context)       { c.JSON(200, gin.H{"attributes": []string{}}) }
func createAttributeHandler(c *gin.Context)      { c.JSON(201, gin.H{"id": "new-id"}) }
func getAttributeHandler(c *gin.Context)         { c.JSON(200, gin.H{"id": c.Param("id")}) }
func updateAttributeHandler(c *gin.Context)      { c.JSON(200, gin.H{"message": "updated"}) }
func deleteAttributeHandler(c *gin.Context)      { c.JSON(200, gin.H{"message": "deleted"}) }
func bulkImportHandler(c *gin.Context)           { c.JSON(200, gin.H{"imported": 0}) }
func bulkExportHandler(c *gin.Context)           { c.JSON(200, gin.H{"url": "export-url"}) }
func bulkUpdateHandler(c *gin.Context)           { c.JSON(200, gin.H{"updated": 0}) }
func searchProductsHandler(c *gin.Context)       { c.JSON(200, gin.H{"results": []string{}}) }
func getProductDetailsAction(c *gin.Context)     { c.JSON(200, gin.H{"product": nil}) }
func searchProductsAction(c *gin.Context)        { c.JSON(200, gin.H{"results": []string{}}) }
func getCategoryTreeAction(c *gin.Context)       { c.JSON(200, gin.H{"tree": []string{}}) }
