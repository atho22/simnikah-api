#!/bin/bash

# Test Address Autocomplete API
# Pastikan server running dan Anda sudah login

echo "================================"
echo "TEST ADDRESS AUTOCOMPLETE API"
echo "================================"
echo ""

# Ganti dengan token JWT Anda
TOKEN="your_jwt_token_here"

echo "1. Test pencarian 'Banjarmasin'"
echo "--------------------------------"
curl -s -X GET "http://localhost:8080/simnikah/location/search?q=Banjarmasin" \
  -H "Authorization: Bearer ${TOKEN}" | jq

echo ""
echo ""

echo "2. Test pencarian 'Jl. Pangeran Antasari'"
echo "----------------------------------------"
curl -s -X GET "http://localhost:8080/simnikah/location/search?q=Jl.%20Pangeran%20Antasari" \
  -H "Authorization: Bearer ${TOKEN}" | jq

echo ""
echo ""

echo "3. Test pencarian query pendek (harus error)"
echo "-------------------------------------------"
curl -s -X GET "http://localhost:8080/simnikah/location/search?q=Ba" \
  -H "Authorization: Bearer ${TOKEN}" | jq

echo ""
echo ""

echo "4. Test pencarian 'Kalimantan Selatan'"
echo "-------------------------------------"
curl -s -X GET "http://localhost:8080/simnikah/location/search?q=Kalimantan%20Selatan" \
  -H "Authorization: Bearer ${TOKEN}" | jq

echo ""
echo "================================"
echo "Test selesai!"
echo "================================"


