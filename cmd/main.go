package cmd

import (
	"context"
	"log"
	"math/rand/v2"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/rajathjn/deny-by-default-as-a-service/internal/favicon"
	"github.com/rajathjn/deny-by-default-as-a-service/internal/help"
	"github.com/rajathjn/deny-by-default-as-a-service/internal/rate_limiter"
	"github.com/rajathjn/deny-by-default-as-a-service/internal/utils"
)

type jsonResponse struct {
	Reason string `json:"reason"`
	Type   string `json:"type"`
}

func respondWithReason(c *gin.Context, reason, reasonType string) {
	if utils.WantsJSON(c) {
		c.JSON(
			http.StatusOK,
			jsonResponse{
				Reason: reason,
				Type:   reasonType,
			},
		)
		return
	}
	c.String(http.StatusOK, reason)
}

func Server(address string) {
	gin.SetMode(gin.ReleaseMode)
	// Create a Gin router with default middleware (logger and recovery)
	router := gin.Default()

	router.SetTrustedProxies(nil) // Disable trusted proxies to get real client IPs
	router.Use(cors.Default())
	router.Use(ratelimiter.RateLimiter())

	// Health check endpoint
	router.GET(
		"/health",
		func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"status": "ok"})
		},
	)

	// Help endpoint
	router.GET(
		"/help", 
		help.Help,
	)

	// Default endpoint is for no
	router.GET(
		"/",
		func(c *gin.Context) {
			respondWithReason(c, utils.GetNegativeReason(), "no")
		},
	)

	// Explicit /no endpoint
	router.GET(
		"/no",
		func(c *gin.Context) {
			respondWithReason(c, utils.GetNegativeReason(), "no")
		},
	)

	// Explicit /yes endpoint
	router.GET(
		"/yes",
		func(c *gin.Context) {
			respondWithReason(c, utils.GetPositiveReason(), "yes")
		},
	)

	// Random yes or no
	router.GET(
		"/random",
		func(c *gin.Context) {
			if rand.IntN(2) == 0 {
				respondWithReason(c, utils.GetNegativeReason(), "no")
			} else {
				respondWithReason(c, utils.GetPositiveReason(), "yes")
			}
		},
	)

	// For favicon.ico
	router.GET(
		"/favicon.ico",
		func(c *gin.Context) {
			faviconData, err := favicon.GetFavicon()
			if err != nil {
				log.Printf("Error getting favicon: %v", err)
				c.Status(http.StatusInternalServerError)
				return
			}
			c.Data(
				http.StatusOK,
				"image/x-icon",
				faviconData,
			)
		},
	)

	// Catch-all returns a positive reason
	router.NoRoute(
		func(c *gin.Context) {
			respondWithReason(c, utils.GetPositiveReason(), "yes")
		},
	)

	// Graceful shutdown
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	srv := &http.Server{
		Addr:    address,
		Handler: router,
		// Below recommended timeouts help prevent Slowloris attacks and improve server resilience
		ReadTimeout:       5 * time.Second,
		ReadHeaderTimeout: 2 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       120 * time.Second,
	}

	go func() {
		log.Printf("Running the server on %s\n", address)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to run server on %s: %v", address, err)
		}
	}()

	<-ctx.Done()
	log.Println("Shutting down server...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited")
}
