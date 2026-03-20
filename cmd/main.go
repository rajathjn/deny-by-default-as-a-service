package cmd

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/cors"
	"github.com/rajathjn/deny-by-default-as-a-service/internal/utils"
	"github.com/rajathjn/deny-by-default-as-a-service/internal/rate_limiter"
	"github.com/rajathjn/deny-by-default-as-a-service/internal/favicon"
)

func Server(address string) {
	// gin.SetMode(gin.ReleaseMode)
	// Create a Gin router with default middleware (logger and recovery)
	router := gin.Default()

	router.Use(cors.Default())
	router.Use(ratelimiter.Ratelimiter())

	// Default endpoint is for no
	router.GET(
		"/", 
		func(c *gin.Context) {
			// Return Response
			c.String(
				http.StatusOK,
				utils.Get_negative_reason(),
			)
	})

	// For favicon.ico
	router.GET(
		"/favicon.ico",
		func(c *gin.Context) {
			favicon_data, err := favicon.Get_favicon()
			if err != nil {
				log.Printf("Error getting favicon: %v", err)
				c.Status(http.StatusInternalServerError)
				return
			}
			c.Data(
				http.StatusOK,
				"image/x-icon",
				favicon_data,
			)
		},
	)


	// For yes
	router.NoRoute(
		func(c *gin.Context) {
			c.String(
				http.StatusOK,
				utils.Get_positive_reason(),
			)
	})
	
	log.Printf("Running the server on %s\n", address)
	// Start server on port 8080
	if err := router.Run(address); err != nil {
		log.Fatalf("Failed to run server on %s: %v", address, err)
	}
}
