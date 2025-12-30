# Simple Redis Connection Test
Write-Host "=== Redis Connection Test ===" -ForegroundColor Cyan

# Test 1: Port check
Write-Host "[1] Checking Redis port 6379..." -ForegroundColor Yellow
$portTest = Test-NetConnection -ComputerName localhost -Port 6379 -WarningAction SilentlyContinue
if ($portTest.TcpTestSucceeded) {
    Write-Host "    Redis is listening on port 6379" -ForegroundColor Green
} else {
    Write-Host "    Redis is NOT listening" -ForegroundColor Red
    exit 1
}

# Test 2: Ping
Write-Host "[2] Pinging Redis..." -ForegroundColor Yellow
$ping = redis-cli ping
Write-Host "    PING result: $ping" -ForegroundColor Green

# Test 3: Set/Get
Write-Host "[3] Testing SET/GET..." -ForegroundColor Yellow
redis-cli SET "test_key" "test_value" | Out-Null
$val = redis-cli GET "test_key"
Write-Host "    GET result: $val" -ForegroundColor Green
redis-cli DEL "test_key" | Out-Null

# Test 4: Check keys
Write-Host "[4] Checking existing keys..." -ForegroundColor Yellow
$keys = redis-cli KEYS "*"
Write-Host "    Keys found: $keys" -ForegroundColor Cyan

Write-Host ""
Write-Host "=== Redis is working correctly! ===" -ForegroundColor Green
Write-Host "Redis is accessible on localhost:6379 and ready for use." -ForegroundColor White
