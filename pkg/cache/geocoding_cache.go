package cache

import (
	"container/list"
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"simnikah/pkg/utils"
)

// Bounded LRU cache for geocoding results with TTL per entry.
// Keys are normalized (trim + lowercase) to reduce duplicates.
type GeocodingCache struct {
	mu       sync.Mutex
	ll       *list.List
	cache    map[string]*list.Element
	capacity int
	ttl      time.Duration
	hits     int64
	misses   int64
}

type cacheEntry struct {
	key   string
	value CachedCoordinate
}

// CachedCoordinate menyimpan koordinat dengan waktu expiry
type CachedCoordinate struct {
	Latitude  float64
	Longitude float64
	CachedAt  time.Time
	ExpiresAt time.Time
}

var (
	geocodingCache *GeocodingCache
	cacheOnce      sync.Once
)

// NewGeocodingCache creates an LRU cache with capacity and ttl
func NewGeocodingCache(capacity int, ttl time.Duration) *GeocodingCache {
	gc := &GeocodingCache{
		ll:       list.New(),
		cache:    make(map[string]*list.Element),
		capacity: capacity,
		ttl:      ttl,
	}
	go gc.cleanupExpired()
	return gc
}

// GetGeocodingCache returns singleton instance (size from GEOCODING_CACHE_SIZE env)
func GetGeocodingCache() *GeocodingCache {
	cacheOnce.Do(func() {
		capDefault := 5000
		ttlDefault := 30 * 24 * time.Hour

		cap := capDefault
		if v := os.Getenv("GEOCODING_CACHE_SIZE"); v != "" {
			if n, err := strconv.Atoi(v); err == nil && n > 0 {
				cap = n
			}
		}

		ttl := ttlDefault
		if v := os.Getenv("GEOCODING_CACHE_TTL_DAYS"); v != "" {
			if n, err := strconv.Atoi(v); err == nil && n > 0 {
				ttl = time.Duration(n) * 24 * time.Hour
			}
		}

		geocodingCache = NewGeocodingCache(cap, ttl)
	})
	return geocodingCache
}

func normalizeKey(address string) string {
	return strings.TrimSpace(strings.ToLower(address))
}

// Get attempts to fetch coordinates from cache. Moves item to front on hit.
func (gc *GeocodingCache) Get(address string) (lat, lon float64, found bool) {
	key := normalizeKey(address)
	gc.mu.Lock()
	defer gc.mu.Unlock()

	ele, ok := gc.cache[key]
	if !ok {
		gc.misses++
		return 0, 0, false
	}
	entry := ele.Value.(*cacheEntry)
	if time.Now().After(entry.value.ExpiresAt) {
		// expired
		gc.ll.Remove(ele)
		delete(gc.cache, key)
		gc.misses++
		return 0, 0, false
	}
	// move to front (most recently used)
	gc.ll.MoveToFront(ele)
	gc.hits++
	return entry.value.Latitude, entry.value.Longitude, true
}

// Set stores coordinates into cache, evicting least-recently-used when over capacity
func (gc *GeocodingCache) Set(address string, lat, lon float64) {
	key := normalizeKey(address)
	gc.mu.Lock()
	defer gc.mu.Unlock()

	now := time.Now()
	if ele, ok := gc.cache[key]; ok {
		// update existing
		entry := ele.Value.(*cacheEntry)
		entry.value.Latitude = lat
		entry.value.Longitude = lon
		entry.value.CachedAt = now
		entry.value.ExpiresAt = now.Add(gc.ttl)
		gc.ll.MoveToFront(ele)
		return
	}

	entry := &cacheEntry{key: key, value: CachedCoordinate{Latitude: lat, Longitude: lon, CachedAt: now, ExpiresAt: now.Add(gc.ttl)}}
	ele := gc.ll.PushFront(entry)
	gc.cache[key] = ele

	if gc.ll.Len() > gc.capacity {
		// evict LRU
		lru := gc.ll.Back()
		if lru != nil {
			lruEntry := lru.Value.(*cacheEntry)
			delete(gc.cache, lruEntry.key)
			gc.ll.Remove(lru)
		}
	}
}

// cleanupExpired periodically removes expired entries
func (gc *GeocodingCache) cleanupExpired() {
	ticker := time.NewTicker(1 * time.Hour)
	defer ticker.Stop()
	for range ticker.C {
		gc.mu.Lock()
		now := time.Now()
		removed := 0
		for e := gc.ll.Back(); e != nil; {
			prev := e.Prev()
			entry := e.Value.(*cacheEntry)
			if now.After(entry.value.ExpiresAt) {
				delete(gc.cache, entry.key)
				gc.ll.Remove(e)
				removed++
			}
			e = prev
		}
		if removed > 0 {
			fmt.Printf("🧹 Cleaned up %d expired geocoding cache entries\n", removed)
		}
		gc.mu.Unlock()
	}
}

// Stats returns simple cache statistics
func (gc *GeocodingCache) Stats() map[string]interface{} {
	gc.mu.Lock()
	defer gc.mu.Unlock()
	return map[string]interface{}{
		"total_entries": gc.ll.Len(),
		"capacity":      gc.capacity,
		"hits":          gc.hits,
		"misses":        gc.misses,
	}
}

// GetCoordinatesFromAddressCached adalah wrapper dengan caching
func GetCoordinatesFromAddressCached(address string) (float64, float64, error) {
	cache := GetGeocodingCache()

	if lat, lon, found := cache.Get(address); found {
		fmt.Printf("🎯 Cache HIT: Geocoding untuk '%s' (%.6f, %.6f)\n", address, lat, lon)
		return lat, lon, nil
	}

	fmt.Printf("🌐 Cache MISS: Fetching geocoding untuk '%s' dari OpenStreetMap...\n", address)
	lat, lon, err := utils.GetCoordinatesFromAddress(address)
	if err != nil {
		return 0, 0, err
	}

	cache.Set(address, lat, lon)
	fmt.Printf("💾 Cached geocoding untuk '%s' (%.6f, %.6f)\n", address, lat, lon)
	return lat, lon, nil
}
