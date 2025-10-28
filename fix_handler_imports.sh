#!/bin/bash

echo "🔧 Fixing handler imports and helper references..."

# Fix imports in staff.go
echo "Fixing staff.go..."
sed -i '/^import (/,/^)/ {
    /simnikah\/pkg\/utils/!{
        /^)$/i\
	"simnikah/pkg/utils"\
	"simnikah/pkg/crypto"\
	"simnikah/pkg/validator"
    }
}' internal/handlers/staff/staff.go

# Replace helper calls
sed -i 's/helper\.HashPassword/crypto.HashPassword/g' internal/handlers/staff/staff.go

# Fix imports in catin/daftar.go
echo "Fixing catin/daftar.go..."
sed -i '/^import (/,/^)/ {
    /simnikah\/pkg\/utils/!{
        /^)$/i\
	"simnikah/pkg/utils"\
	"simnikah/pkg/cache"\
	"simnikah/pkg/validator"
    }
}' internal/handlers/catin/daftar.go

# Replace helper calls in catin/daftar.go
sed -i 's/helper\.GetCoordinatesFromAddressCached/cache.GetCoordinatesFromAddressCached/g' internal/handlers/catin/daftar.go
sed -i 's/helper\.CheckValidValue/validator.CheckValidValue/g' internal/handlers/catin/daftar.go
sed -i 's/helper\.ValidateParentFields/validator.ValidateParentFields/g' internal/handlers/catin/daftar.go
sed -i 's/helper\.ValidatePersonFields/validator.ValidatePersonFields/g' internal/handlers/catin/daftar.go

# Fix imports in catin/location.go
echo "Fixing catin/location.go..."
sed -i '/^import (/,/^)/ {
    /simnikah\/pkg\/cache/!{
        /^)$/i\
	"simnikah/pkg/cache"
    }
}' internal/handlers/catin/location.go

# Replace helper calls in location.go
sed -i 's/helper\.GetCoordinatesFromAddressCached/cache.GetCoordinatesFromAddressCached/g' internal/handlers/catin/location.go

echo "✅ Handler imports fixed!"
echo ""
echo "Test build: go build -o bin/simnikah-api cmd/api/main.go"

