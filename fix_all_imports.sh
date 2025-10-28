#!/bin/bash

# Fix all import paths after restructure

echo "🔧 Fixing all import paths..."

# Find all Go files (excluding vendor and go modules)
GO_FILES=$(find . -name "*.go" -type f -not -path "./go/*" -not -path "./vendor/*")

# Count total files
TOTAL=$(echo "$GO_FILES" | wc -l)
echo "Found $TOTAL Go files to update"
echo ""

# Fix imports one by one
for file in $GO_FILES; do
    echo "Processing: $file"
    
    # Handlers
    sed -i 's|"simnikah/catin"|"simnikah/internal/handlers/catin"|g' "$file"
    sed -i 's|"simnikah/staff"|"simnikah/internal/handlers/staff"|g' "$file"
    sed -i 's|"simnikah/penghulu"|"simnikah/internal/handlers/penghulu"|g' "$file"
    sed -i 's|"simnikah/kepala_kua"|"simnikah/internal/handlers/kepala_kua"|g' "$file"
    sed -i 's|"simnikah/notification"|"simnikah/internal/handlers/notification"|g' "$file"
    
    # Models
    sed -i 's|"simnikah/structs"|"simnikah/internal/models"|g' "$file"
    
    # Middleware
    sed -i 's|"simnikah/middleware"|"simnikah/internal/middleware"|g' "$file"
    
    # Services
    sed -i 's|"simnikah/services"|"simnikah/internal/services"|g' "$file"
    
    # Helper → pkg/* (more complex, need to check usage)
    # We'll replace with pkg/utils first, then manually fix specific cases
    sed -i 's|"simnikah/helper"|"simnikah/pkg/utils"|g' "$file"
done

echo ""
echo "✅ Import paths updated!"
echo ""
echo "⚠️  Additional manual fixes needed:"
echo "  - Check files that use validation.go → might need pkg/validator"
echo "  - Check files that use bcrypt.go → might need pkg/crypto"
echo "  - Check files that use geocoding_cache.go → might need pkg/cache"
echo ""
echo "Run 'go build ./cmd/api' to test"

