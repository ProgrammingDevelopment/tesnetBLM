package services

import (
	"context"
	"log"
	"sync"
	"sync/atomic"

	"github.com/redis/go-redis/v9"
)

type RedisService struct {
	Client    *redis.Client
	connected bool
	// Fallback in-memory counter when Redis is unavailable
	memoryQuota int64
	mu          sync.Mutex
}

var ctx = context.Background()

func NewRedisService() *RedisService {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	service := &RedisService{
		Client:      client,
		connected:   false,
		memoryQuota: 5000, // Default quota
	}

	_, err := client.Ping(ctx).Result()
	if err != nil {
		log.Printf("⚠️  Redis not connected: %v", err)
		log.Println("📦 Using IN-MEMORY fallback mode (quota: 5000)")
		service.connected = false
	} else {
		log.Println("✅ Redis connected")
		service.connected = true
		// Initialize quota for testing
		client.SetNX(ctx, "ticket_quota", 5000, 0)
	}

	return service
}

func (s *RedisService) AtomicDecreaseQuota() (int64, error) {
	if s.connected {
		// Use Redis DECR (atomic)
		val, err := s.Client.Decr(ctx, "ticket_quota").Result()
		return val, err
	}

	// Fallback: use atomic in-memory counter
	newVal := atomic.AddInt64(&s.memoryQuota, -1)
	return newVal, nil
}

func (s *RedisService) GetQuota() (int64, error) {
	if s.connected {
		val, err := s.Client.Get(ctx, "ticket_quota").Int64()
		if err == redis.Nil {
			return 0, nil
		}
		return val, err
	}

	// Fallback: return in-memory quota
	return atomic.LoadInt64(&s.memoryQuota), nil
}

func (s *RedisService) IsConnected() bool {
	return s.connected
}
