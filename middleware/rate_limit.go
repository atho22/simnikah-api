package middleware

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ulule/limiter/v3"
	"github.com/ulule/limiter/v3/drivers/store/memory"
)

// RateLimiter creates a rate limiting middleware
// Default: 100 requests per minute per IP
func RateLimiter() gin.HandlerFunc {
	// Create a rate limiter store (in-memory)
	store := memory.NewStore()

	// Rate: 100 requests per minute
	rate := limiter.Rate{
		Period: 1 * time.Minute,
		Limit:  100,
	}

	// Create limiter instance
	instance := limiter.New(store, rate)

	return func(c *gin.Context) {
		// Get client IP
		clientIP := c.ClientIP()

		// Check rate limit
		context, err := instance.Get(c.Request.Context(), clientIP)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Rate limiter error",
				"error":   err.Error(),
			})
			c.Abort()
			return
		}

		// Set rate limit headers
		c.Header("X-RateLimit-Limit", fmt.Sprintf("%d", context.Limit))
		c.Header("X-RateLimit-Remaining", fmt.Sprintf("%d", context.Remaining))
		c.Header("X-RateLimit-Reset", fmt.Sprintf("%d", context.Reset))

		// Check if limit exceeded
		if context.Reached {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"success": false,
				"message": "Rate limit exceeded",
				"error":   "Terlalu banyak request. Silakan coba lagi nanti.",
				"retry_after": fmt.Sprintf("%d detik", int(time.Until(time.Unix(context.Reset, 0)).Seconds())),
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// StrictRateLimiter for sensitive endpoints (login, register)
// Rate: 5 requests per minute per IP
func StrictRateLimiter() gin.HandlerFunc {
	store := memory.NewStore()

	rate := limiter.Rate{
		Period: 1 * time.Minute,
		Limit:  5,
	}

	instance := limiter.New(store, rate)

	return func(c *gin.Context) {
		clientIP := c.ClientIP()

		context, err := instance.Get(c.Request.Context(), clientIP)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Rate limiter error",
				"error":   err.Error(),
			})
			c.Abort()
			return
		}

		c.Header("X-RateLimit-Limit", fmt.Sprintf("%d", context.Limit))
		c.Header("X-RateLimit-Remaining", fmt.Sprintf("%d", context.Remaining))
		c.Header("X-RateLimit-Reset", fmt.Sprintf("%d", context.Reset))

		if context.Reached {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"success": false,
				"message": "Terlalu banyak percobaan login",
				"error":   "Demi keamanan, Anda harus menunggu sebelum mencoba lagi.",
				"retry_after": fmt.Sprintf("%d detik", int(time.Until(time.Unix(context.Reset, 0)).Seconds())),
				"type":    "security",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

