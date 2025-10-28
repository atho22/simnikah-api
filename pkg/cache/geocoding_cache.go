package cache

import (
	"fmt"
	"sync"
	"time"

	"simnikah/pkg/utils"
)

// GeocodingCache adalah in-memory cache untuk hasil geocoding
// Mengurangi request ke OpenStreetMap API dari 1-3 detik menjadi <1ms
type GeocodingCache struct {
	cache map[string]*CachedCoordinate
	mu    sync.RWMutex
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

// GetGeocodingCache returns singleton instance of geocoding cache
func GetGeocodingCache() *GeocodingCache {
	cacheOnce.Do(func() {
		geocodingCache = &GeocodingCache{
			cache: make(map[string]*CachedCoordinate),
		}

		// Start background cleanup goroutine (every 1 hour)
		go geocodingCache.cleanupExpired()
	})
	return geocodingCache
}

// Get mencoba ambil koordinat dari cache
func (gc *GeocodingCache) Get(address string) (lat, lon float64, found bool) {
	gc.mu.RLock()
	defer gc.mu.RUnlock()

	cached, exists := gc.cache[address]
	if !exists {
		return 0, 0, false
	}

	// Check if expired
	if time.Now().After(cached.ExpiresAt) {
		return 0, 0, false
	}

	return cached.Latitude, cached.Longitude, true
}

// Set menyimpan koordinat ke cache
// TTL default: 30 hari (alamat jarang berubah)
func (gc *GeocodingCache) Set(address string, lat, lon float64) {
	gc.mu.Lock()
	defer gc.mu.Unlock()

	now := time.Now()
	gc.cache[address] = &CachedCoordinate{
		Latitude:  lat,
		Longitude: lon,
		CachedAt:  now,
		ExpiresAt: now.Add(30 * 24 * time.Hour), // 30 days
	}
}

// cleanupExpired membersihkan cache yang sudah expired (background job)
func (gc *GeocodingCache) cleanupExpired() {
	ticker := time.NewTicker(1 * time.Hour)
	defer ticker.Stop()

	for range ticker.C {
		gc.mu.Lock()
		now := time.Now()
		count := 0

		for addr, cached := range gc.cache {
			if now.After(cached.ExpiresAt) {
				delete(gc.cache, addr)
				count++
			}
		}

		if count > 0 {
			fmt.Printf("🧹 Cleaned up %d expired geocoding cache entries\n", count)
		}
		gc.mu.Unlock()
	}
}

// Stats returns cache statistics
func (gc *GeocodingCache) Stats() map[string]interface{} {
	gc.mu.RLock()
	defer gc.mu.RUnlock()

	return map[string]interface{}{
		"total_entries": len(gc.cache),
		"cache_enabled": true,
	}
}

// GetCoordinatesFromAddressCached adalah wrapper dengan caching
// Menggunakan cache untuk menghindari request berulang ke OpenStreetMap
func GetCoordinatesFromAddressCached(address string) (float64, float64, error) {
	cache := GetGeocodingCache()

	// Try cache first
	if lat, lon, found := cache.Get(address); found {
		fmt.Printf("🎯 Cache HIT: Geocoding untuk '%s' (%.6f, %.6f)\n", address, lat, lon)
		return lat, lon, nil
	}

	// Cache miss - fetch from API
	fmt.Printf("🌐 Cache MISS: Fetching geocoding untuk '%s' dari OpenStreetMap...\n", address)
	lat, lon, err := utils.GetCoordinatesFromAddress(address)
	if err != nil {
		return 0, 0, err
	}

	// Save to cache
	cache.Set(address, lat, lon)
	fmt.Printf("💾 Cached geocoding untuk '%s' (%.6f, %.6f)\n", address, lat, lon)

	return lat, lon, nil
}
