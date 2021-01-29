package middlewares

import (
	"fmt"
	"math"
	"time"

	"github.com/gin-gonic/gin"
)

//Benchmark ---
func Benchmark() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		elapsed := time.Since(start)
		fmt.Printf("Request took %v milliseconds\n", float64(elapsed.Nanoseconds())/math.Pow(float64(10), float64(6)))
	}
}
