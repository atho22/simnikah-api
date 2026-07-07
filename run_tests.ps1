#!/usr/bin/env pwsh
# ============================================================
# run_tests.ps1 - SimNikah API Test Automation Script
# ============================================================
# Usage:
#   .\run_tests.ps1               -> Run all tests
#   .\run_tests.ps1 -Coverage     -> Run all tests + HTML coverage
#   .\run_tests.ps1 -Group auth   -> Run only tests matching "auth"
#   .\run_tests.ps1 -Verbose      -> Run all tests with verbose output
#   .\run_tests.ps1 -Race         -> Run with race detector
# ============================================================

param(
    [switch]$Coverage,
    [switch]$Verbose,
    [switch]$Race,
    [string]$Group = "",
    [switch]$Help
)

function Write-Green($msg)  { Write-Host $msg -ForegroundColor Green }
function Write-Red($msg)    { Write-Host $msg -ForegroundColor Red }
function Write-Yellow($msg) { Write-Host $msg -ForegroundColor Yellow }
function Write-Cyan($msg)   { Write-Host $msg -ForegroundColor Cyan }
function Write-White($msg)  { Write-Host $msg -ForegroundColor White }

if ($Help) {
    Write-Cyan "============================================================"
    Write-Cyan " SimNikah API Test Runner - Help"
    Write-Cyan "============================================================"
    Write-White "OPTIONS:"
    Write-White "  -Coverage        Generate HTML coverage report"
    Write-White "  -Verbose         Show verbose test output"
    Write-White "  -Race            Run with race condition detector"
    Write-White "  -Group <name>    Run only specific group:"
    Write-White "                   auth, catin, staff, notification, rbac,"
    Write-White "                   security, location, forward_chaining,"
    Write-White "                   dashboard, geocode, penghulu, kepala_kua"
    Write-White "  -Help            Show this help"
    exit 0
}

$timestamp = Get-Date -Format "yyyy-MM-dd HH:mm:ss"
Write-Cyan "============================================================"
Write-Cyan " SimNikah API - Automated Test Runner"
Write-Cyan " Started at: $timestamp"
Write-Cyan "============================================================"

if (-not (Get-Command go -ErrorAction SilentlyContinue)) {
    Write-Red "[ERROR] Go is not installed or not in PATH"
    exit 1
}

$goTestArgs = @("test", "./tests/...", "-count=1", "-timeout=120s")

if ($Verbose)  { $goTestArgs += "-v" }
if ($Race)     { $goTestArgs += "-race"; Write-Yellow "[INFO] Race detector ON" }
if ($Coverage) {
    $goTestArgs += "-coverprofile=coverage.out"
    $goTestArgs += "-coverpkg=./internal/..."
}

if ($Group -ne "") {
    $patterns = @{
        "auth"             = "TestRegister|TestLogin|TestGetProfile|TestHealthCheck"
        "catin"            = "TestCatin"
        "staff"            = "TestStaff"
        "notification"     = "TestNotification|TestIDOR_.*Notif"
        "rbac"             = "TestRBAC"
        "security"         = "TestJWT|TestInjection|TestHTTP|Test404"
        "location"         = "TestIDOR_.*Loc|TestLocation|TestPenghulu_UpdateCoord"
        "forward_chaining" = "TestFC"
        "dashboard"        = "TestDashboard"
        "geocode"          = "TestGeocode|TestReverseGeocode|TestSearchAddress|TestLocation_All"
        "penghulu"         = "TestPenghulu"
        "kepala_kua"       = "TestKepalaKUA"
    }
    $pattern = $patterns[$Group]
    if ($null -eq $pattern) {
        Write-Yellow "[WARN] Unknown group '$Group'. Running all tests."
        Write-Yellow "       Available groups: $($patterns.Keys -join ', ')"
    } else {
        $goTestArgs += "-run"; $goTestArgs += $pattern
        Write-Yellow "[INFO] Group: '$Group' -> pattern: $pattern"
    }
}

Write-White ""
Write-Cyan "Running: go $($goTestArgs -join ' ')"
Write-Cyan "------------------------------------------------------------"

$env:JWT_KEY = "test-secret-key-for-unit-testing-only"
$startTime = Get-Date

& go @goTestArgs
$exitCode = $LASTEXITCODE
$duration = [math]::Round(((Get-Date) - $startTime).TotalSeconds, 1)

Write-White ""
Write-Cyan "------------------------------------------------------------"
if ($exitCode -eq 0) {
    Write-Green " ALL TESTS PASSED  (${duration}s)"
} else {
    Write-Red " TESTS FAILED (exit: $exitCode, ${duration}s)"
}

if ($Coverage -and (Test-Path "coverage.out")) {
    & go tool cover -html=coverage.out -o coverage.html
    if ($LASTEXITCODE -eq 0) {
        Write-Green " Coverage: coverage.html"
        $total = (& go tool cover -func=coverage.out) | Select-String "total:"
        if ($total) { Write-Yellow " $total" }
        if ($IsWindows -or $env:OS -eq "Windows_NT") { Start-Process "coverage.html" }
    }
}
Write-Cyan "============================================================"
exit $exitCode
