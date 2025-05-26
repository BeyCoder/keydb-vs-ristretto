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
const CacheIterations = 100000

func TestRedisCache(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	db := 0

	rdb := redis.NewClient(&redis.Options{
		Addr:     config.Values.KeyDB.Address,
		Password: config.Values.KeyDB.Password,
		DB:       db,
	})
	if err := rdb.Ping(ctx).Err(); err != nil {
		panic(err)
	}
	defer rdb.Close()
	// test of writing 10k plus keys
	for i := range CacheIterations {
		rdb.Set(t.Context(), fmt.Sprintf("%s%d", KeyBase, i), fmt.Sprintf("%d", i), 0)
	}
	// test of getting 10k plus keys
	for i := range CacheIterations {
		rdb.Get(t.Context(), fmt.Sprintf("%s%d", KeyBase, i)).Result()
	}
}

func TestRistrettoCache(t *testing.T) {
	cache, err := ristretto.NewCache(&ristretto.Config[string, string]{
		NumCounters: 1e7,     // number of keys to track frequency of (10M).
		MaxCost:     1 << 30, // 1GB.
		BufferItems: 64,      // number of keys per Get buffer.
	})
	if err != nil {
	}
	defer cache.Close()
	// writing
	for i := range CacheIterations {
		cache.Set(fmt.Sprintf("%s%d", KeyBase, i), fmt.Sprintf("%d", i), 1)
	}
	// reading
	for i := range CacheIterations {
		cache.Get(fmt.Sprintf("%s%d", KeyBase, i))
	}
}
