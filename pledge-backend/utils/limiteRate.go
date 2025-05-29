package utils

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"sync"
	"time"
)

// RateLimiter 令牌桶限流器
type RateLimiter struct {
	capacity  int64            // 桶容量
	rate      int64            // 令牌生成速率 (个/秒)
	tokens    map[string]int64 // 客户端令牌数量
	lastToken map[string]int64 // 客户端上次令牌时间 (纳秒)
	mu        sync.RWMutex     // 读写锁
}

// NewRateLimiter 创建新的限流器
func NewRateLimiter(capacity, rate int64) *RateLimiter {
	return &RateLimiter{
		capacity:  capacity,
		rate:      rate,
		tokens:    make(map[string]int64),
		lastToken: make(map[string]int64),
	}
}

// Allow 判断请求是否允许
func (l *RateLimiter) Allow(clientIP string) bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	now := time.Now().UnixNano()

	// 初始化或获取客户端的令牌信息
	tokens, exists := l.tokens[clientIP]
	if !exists {
		tokens = l.capacity - 1
		l.tokens[clientIP] = tokens
		l.lastToken[clientIP] = now
		return true
	}

	last := l.lastToken[clientIP]
	// 计算从上次请求到现在应生成的令牌数
	elapsed := now - last
	generatedTokens := elapsed * l.rate / 1e9 // 转换为秒

	// 更新令牌数，但不超过容量
	tokens = tokens + generatedTokens
	if tokens > l.capacity {
		tokens = l.capacity
	}

	// 保存更新后的令牌数和时间
	l.tokens[clientIP] = tokens
	l.lastToken[clientIP] = now

	// 判断是否有足够的令牌
	if tokens >= 1 {
		l.tokens[clientIP] = tokens - 1
		return true
	}

	return false
}

// RateLimitMiddleware 创建限流中间件
func RateLimitMiddleware(capacity, rate int64) gin.HandlerFunc {
	limiter := NewRateLimiter(capacity, rate)
	return func(c *gin.Context) {
		clientIP := c.ClientIP()
		if !limiter.Allow(clientIP) {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"error": "请求频率过高，请稍后再试",
			})
			return
		}
		c.Next()
	}
}
