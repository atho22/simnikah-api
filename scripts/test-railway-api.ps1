# Script untuk test Railway API endpoints
# Usage: .\scripts\test-railway-api.ps1 -Url "https://your-app.up.railway.app"

param(
    [Parameter(Mandatory=$true)]
    [string]$Url
)

Write-Host "========================================" -ForegroundColor Cyan
Write-Host "  Test Railway API Endpoints" -ForegroundColor Cyan
Write-Host "========================================" -ForegroundColor Cyan
Write-Host ""
Write-Host "Testing API: $Url" -ForegroundColor Yellow
Write-Host ""

# Test 1: Health Check
Write-Host "1. Testing Health Endpoint..." -ForegroundColor Green
try {
    $healthResponse = Invoke-RestMethod -Uri "$Url/health" -Method Get
    Write-Host "   ✅ Health Check: SUCCESS" -ForegroundColor Green
    Write-Host "   Status: $($healthResponse.status)" -ForegroundColor White
    Write-Host "   Service: $($healthResponse.service)" -ForegroundColor White
} catch {
    Write-Host "   ❌ Health Check: FAILED" -ForegroundColor Red
    Write-Host "   Error: $($_.Exception.Message)" -ForegroundColor Red
}
Write-Host ""

# Test 2: Register User
Write-Host "2. Testing Register Endpoint..." -ForegroundColor Green
$registerData = @{
    username = "testuser_$(Get-Random -Minimum 1000 -Maximum 9999)"
    email = "test_$(Get-Random -Minimum 1000 -Maximum 9999)@example.com"
    password = "test123"
    nama = "Test User"
    role = "user_biasa"
} | ConvertTo-Json

try {
    $registerResponse = Invoke-RestMethod -Uri "$Url/register" -Method Post -Body $registerData -ContentType "application/json"
    Write-Host "   ✅ Register: SUCCESS" -ForegroundColor Green
    Write-Host "   User ID: $($registerResponse.user.user_id)" -ForegroundColor White
    Write-Host "   Username: $($registerResponse.user.username)" -ForegroundColor White
    $testUsername = $registerData | ConvertFrom-Json | Select-Object -ExpandProperty username
    $testPassword = "test123"
} catch {
    Write-Host "   ❌ Register: FAILED" -ForegroundColor Red
    Write-Host "   Error: $($_.Exception.Message)" -ForegroundColor Red
    $testUsername = $null
    $testPassword = $null
}
Write-Host ""

# Test 3: Login
if ($testUsername) {
    Write-Host "3. Testing Login Endpoint..." -ForegroundColor Green
    $loginData = @{
        username = $testUsername
        password = $testPassword
    } | ConvertTo-Json

    try {
        $loginResponse = Invoke-RestMethod -Uri "$Url/login" -Method Post -Body $loginData -ContentType "application/json"
        Write-Host "   ✅ Login: SUCCESS" -ForegroundColor Green
        Write-Host "   Token: $($loginResponse.token.Substring(0, 50))..." -ForegroundColor White
        Write-Host "   User: $($loginResponse.user.nama)" -ForegroundColor White
    } catch {
        Write-Host "   ❌ Login: FAILED" -ForegroundColor Red
        Write-Host "   Error: $($_.Exception.Message)" -ForegroundColor Red
    }
    Write-Host ""
}

# Summary
Write-Host "========================================" -ForegroundColor Cyan
Write-Host "  Test Summary" -ForegroundColor Cyan
Write-Host "========================================" -ForegroundColor Cyan
Write-Host ""
Write-Host "API URL: $Url" -ForegroundColor Yellow
Write-Host ""
Write-Host "✅ Test selesai!" -ForegroundColor Green
Write-Host ""


