package cache

import (
	"testing"
)

func BenchmarkGeocachingCacheGet(b *testing.B) {
	cache := GetGeocodingCache()

	// Populate cache
	cache.Set("Test Address 123", 40.7128, -74.0060)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cache.Get("Test Address 123")
	}
}

func BenchmarkGeocachingCacheSet(b *testing.B) {
	cache := GetGeocodingCache()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cache.Set("Test Address", 40.7128, -74.0060)
	}
}

func TestGeocachingCache(t *testing.T) {
	cache := GetGeocodingCache()

	// Test Set and Get
	testAddr := "Jl. Test No. 123"
	testLat := -6.2088
	testLon := 106.8456

	cache.Set(testAddr, testLat, testLon)

	lat, lon, found := cache.Get(testAddr)
	if !found {
		t.Error("Expected to find cached address")
	}

	if lat != testLat || lon != testLon {
		t.Errorf("Expected (%f, %f), got (%f, %f)", testLat, testLon, lat, lon)
	}

	// Test cache miss
	_, _, found = cache.Get("Non-existent address")
	if found {
		t.Error("Expected cache miss for non-existent address")
	}
}
