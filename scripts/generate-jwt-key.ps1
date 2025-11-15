# Script untuk generate JWT Secret Key untuk Railway
# Jalankan dengan: .\scripts\generate-jwt-key.ps1

Write-Host "========================================" -ForegroundColor Cyan
Write-Host "  Generate JWT Secret Key" -ForegroundColor Cyan
Write-Host "========================================" -ForegroundColor Cyan
Write-Host ""

# Generate random 32-byte key dan convert ke Base64
$bytes = 1..32 | ForEach-Object { Get-Random -Minimum 0 -Maximum 256 }
$jwtKey = [Convert]::ToBase64String($bytes)

Write-Host "JWT Secret Key (32+ characters):" -ForegroundColor Green
Write-Host ""
Write-Host $jwtKey -ForegroundColor Yellow
Write-Host ""
Write-Host "========================================" -ForegroundColor Cyan
Write-Host ""
Write-Host "Cara menggunakan:" -ForegroundColor Cyan
Write-Host "1. Copy JWT key di atas" -ForegroundColor White
Write-Host "2. Buka Railway Dashboard → Variables" -ForegroundColor White
Write-Host "3. Paste ke variable JWT_KEY" -ForegroundColor White
Write-Host ""
Write-Host "Tekan Enter untuk copy ke clipboard..." -ForegroundColor Yellow
Read-Host

# Copy ke clipboard
$jwtKey | Set-Clipboard
Write-Host "✅ JWT Key sudah di-copy ke clipboard!" -ForegroundColor Green
Write-Host ""


