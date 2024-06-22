package middleware

import (
	"expvar"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func MetricsMiddleware() gin.HandlerFunc {
	// Initialize the new expvar variables.
	totalRequestsReceived := expvar.NewInt("total_requests_received")
	totalResponsesSent := expvar.NewInt("total_responses_sent")
	totalProcessingTimeMicroseconds := expvar.NewInt("total_processing_time_us")
	totalResponsesSentByStatus := expvar.NewMap("total_responses_sent_by_status")

	return func(c *gin.Context) {
		// Increment the total requests received counter.
		totalRequestsReceived.Add(1)

		// Start timer
		start := time.Now()

		// Process the request
		c.Next()

		// Calculate duration
		duration := time.Since(start)

		// Increment the total responses sent counter.
		totalResponsesSent.Add(1)

		// Add the processing time in microseconds.
		totalProcessingTimeMicroseconds.Add(duration.Microseconds())

		// Increment the counter for the specific status code.
		totalResponsesSentByStatus.Add(strconv.Itoa(c.Writer.Status()), 1)
	}
}
