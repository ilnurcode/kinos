// Package middleware предоставляет Gin middleware для api-service.
package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// RateLimiter реализует token bucket rate limiter
type RateLimiter struct {
	tokens   map[string]*bucket
	mu       sync.Mutex
	rate     int           // максимум запросов
	interval time.Duration // интервал
}

type bucket struct {
	tokens     int
	lastRefill time.Time
}

// NewRateLimiter создаёт новый rate limiter
// rate — максимальное количество запросов за interval
func NewRateLimiter(rate int, interval time.Duration) *RateLimiter {
	rl := &RateLimiter{
		tokens:   make(map[string]*bucket),
		rate:     rate,
		interval: interval,
	}

	// Запускаем очистку устаревших записей
	go rl.cleanup()

	return rl
}

// Allow проверяет, разрешён ли запрос для данного ключа (IP)
func (rl *RateLimiter) Allow(key string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	b, exists := rl.tokens[key]

	if !exists {
		rl.tokens[key] = &bucket{
			tokens:     rl.rate - 1,
			lastRefill: now,
		}
		return true
	}

	// Восстанавливаем токены
	elapsed := now.Sub(b.lastRefill)
	intervals := int(elapsed / rl.interval)
	tokensToAdd := intervals * rl.rate
	if tokensToAdd > 0 {
		b.tokens += tokensToAdd
		if b.tokens > rl.rate {
			b.tokens = rl.rate
		}
		b.lastRefill = b.lastRefill.Add(time.Duration(intervals) * rl.interval)
	}

	if b.tokens <= 0 {
		return false
	}

	b.tokens--
	return true
}

// cleanup удаляет устаревшие записи каждые 5 минут
func (rl *RateLimiter) cleanup() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		rl.mu.Lock()
		now := time.Now()
		for key, b := range rl.tokens {
			if now.Sub(b.lastRefill) > rl.interval*5 {
				delete(rl.tokens, key)
			}
		}
		rl.mu.Unlock()
	}
}

// RateLimit возвращает Gin middleware, который ограничивает количество запросов
// с одного IP до rate запросов за interval.
func RateLimit(rate int, interval time.Duration) gin.HandlerFunc {
	limiter := NewRateLimiter(rate, interval)

	return func(c *gin.Context) {
		ip := c.ClientIP()
		if !limiter.Allow(ip) {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": "Слишком много запросов. Попробуйте позже.",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}
