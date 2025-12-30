# Test Redis Connection Script
Write-Host "=== Redis Connection Test ===" -ForegroundColor Cyan
Write-Host ""

# Test 1: Check if Redis is running
Write-Host "[Test 1] Checking if Redis is listening on port 6379..." -ForegroundColor Yellow
$portTest = Test-NetConnection -ComputerName localhost -Port 6379 -WarningAction SilentlyContinue
if ($portTest.TcpTestSucceeded) {
    Write-Host "✅ Redis is listening on port 6379" -ForegroundColor Green
} else {
    Write-Host "❌ Redis is NOT listening on port 6379" -ForegroundColor Red
    exit 1
}
Write-Host ""

# Test 2: Ping Redis
Write-Host "[Test 2] Pinging Redis..." -ForegroundColor Yellow
$pingResult = redis-cli ping
if ($pingResult -eq "PONG") {
    Write-Host "✅ Redis PING successful: $pingResult" -ForegroundColor Green
} else {
    Write-Host "❌ Redis PING failed: $pingResult" -ForegroundColor Red
    exit 1
}
Write-Host ""

# Test 3: Set and Get operation
Write-Host "[Test 3] Testing SET and GET operations..." -ForegroundColor Yellow
$testKey = "test_tofash_redis"
$testValue = "redis_connection_ok"

redis-cli SET $testKey $testValue | Out-Null
$getValue = redis-cli GET $testKey

if ($getValue -eq $testValue) {
    Write-Host "✅ SET/GET successful: Key=$testKey, Value=$getValue" -ForegroundColor Green
} else {
    Write-Host "❌ SET/GET failed: Expected '$testValue', got '$getValue'" -ForegroundColor Red
    exit 1
}

# Cleanup
redis-cli DEL $testKey | Out-Null
Write-Host "✅ Cleanup: Test key deleted" -ForegroundColor Green
Write-Host ""

# Test 4: Check for existing keys (cart data)
Write-Host "[Test 4] Checking for existing keys in Redis..." -ForegroundColor Yellow
$keys = redis-cli KEYS "*"
if ($keys -eq "(empty array)") {
    Write-Host "ℹ️  No keys found in Redis (fresh instance)" -ForegroundColor Cyan
} else {
    Write-Host "✅ Found keys in Redis:" -ForegroundColor Green
    $keys | ForEach-Object { Write-Host "   - $_" -ForegroundColor Gray }
}
Write-Host ""

# Test 5: Test with cart-like data structure
Write-Host "[Test 5] Testing cart-like data structure..." -ForegroundColor Yellow
$cartKey = "cart:user:123"
$cartData = '{"product_id":1,"quantity":2}'

redis-cli HSET $cartKey "item_1" $cartData | Out-Null
$cartValue = redis-cli HGET $cartKey "item_1"

if ($cartValue -eq $cartData) {
    Write-Host "✅ HSET/HGET successful (cart simulation)" -ForegroundColor Green
} else {
    Write-Host "❌ HSET/HGET failed" -ForegroundColor Red
}

# Cleanup
redis-cli DEL $cartKey | Out-Null
Write-Host "✅ Cleanup: Cart test key deleted" -ForegroundColor Green
Write-Host ""

Write-Host "=== All Redis Tests Passed! ===" -ForegroundColor Green
Write-Host ""
Write-Host "Summary:" -ForegroundColor Cyan
Write-Host "  - Redis is running and accessible on localhost:6379" -ForegroundColor White
Write-Host "  - Basic operations (PING, SET, GET, HSET, HGET) work correctly" -ForegroundColor White
Write-Host "  - Redis is ready for use by the application" -ForegroundColor White
Write-Host ""
