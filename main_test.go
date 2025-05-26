package main

import (
	"BeyCoder/keydb-vs-ristretto/config"
	"context"
	"fmt"
	"github.com/dgraph-io/ristretto/v2"
	"github.com/redis/go-redis/v9"
	"testing"
	"time"
)

const KeyBase = "key3"
const CacheIterations = 100_000

func init() {
	config.LoadConfig()
}

func TestRedisCache(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	rdb := redis.NewClient(&redis.Options{
		Network:  "unix",
		Addr:     config.Values.KeyDB.Address,
		Password: config.Values.KeyDB.Password,
		DB:       0,
	})
	defer rdb.Close()

	start := time.Now()

	for i := 0; i < CacheIterations; i++ {
		key := fmt.Sprintf("%s%d", KeyBase, i)
		value := fmt.Sprintf("%d", i)
		if err := rdb.Set(ctx, key, value, 0).Err(); err != nil {
			t.Fatalf("Redis Set error at i=%d: %v", i, err)
		}
	}
	for i := 0; i < CacheIterations; i++ {
		key := fmt.Sprintf("%s%d", KeyBase, i)
		_, err := rdb.Get(ctx, key).Result()
		if err != nil && err != redis.Nil {
			t.Fatalf("Redis Get error at i=%d: %v", i, err)
		}
	}

	t.Logf("Redis test finished in %s", time.Since(start))
}

func TestRistrettoCache(t *testing.T) {
	cache, err := ristretto.NewCache(&ristretto.Config[string, string]{
		NumCounters: 1e7,
		MaxCost:     1 << 30,
		BufferItems: 64,
	})
	if err != nil {
		t.Fatalf("Failed to create Ristretto cache: %v", err)
	}
	defer cache.Close()

	start := time.Now()

	for i := 0; i < CacheIterations; i++ {
		cache.Set(fmt.Sprintf("%s%d", KeyBase, i), fmt.Sprintf("%d", i), 1)
	}
	cache.Wait()

	hits := 0
	for i := 0; i < CacheIterations; i++ {
		if _, ok := cache.Get(fmt.Sprintf("%s%d", KeyBase, i)); ok {
			hits++
		}
	}

	t.Logf("Ristretto test finished in %s", time.Since(start))
	t.Logf("Ristretto hits: %d / %d (%.2f%%)", hits, CacheIterations, float64(hits)*100/CacheIterations)
}
