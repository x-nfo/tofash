$ErrorActionPreference = "Stop"

function Invoke-Api {
    param(
        [string]$Url,
        [string]$Method = "GET",
        [hashtable]$Body = @{},
        [string]$Token = ""
    )
    
    $Headers = @{ "Content-Type" = "application/json" }
    if ($Token) { $Headers["Authorization"] = "Bearer $Token" }

    $JsonBody = $Body | ConvertTo-Json -Depth 10
    
    Write-Host "[$Method] $Url"
    try {
        if ($Method -eq "GET") {
            $Response = Invoke-RestMethod -Uri $Url -Method $Method -Headers $Headers
        }
        else {
            $Response = Invoke-RestMethod -Uri $Url -Method $Method -Headers $Headers -Body $JsonBody
        }
        return $Response
    }
    catch {
        Write-Error $_.Exception.Message
        if ($_.Exception.Response) {
            $Stream = $_.Exception.Response.GetResponseStream()
            $Reader = New-Object System.IO.StreamReader($Stream)
            Write-Host "Error Body: $($Reader.ReadToEnd())" -ForegroundColor Red
        }
        exit 1
    }
}

# 1. Login
Write-Host "`n--- 1. Login ---"
$LoginRes = Invoke-Api -Url "http://localhost:8080/api/v1/login" -Method "POST" -Body @{
    email    = "superadmin@mail.com"
    password = "admin123"
}
$Token = $LoginRes.token
Write-Host "Token obtained"

# 2. Check Categories (Need one for product)
Write-Host "`n--- 2. Get Categories ---"
$Cats = Invoke-Api -Url "http://localhost:8080/api/v1/categories" -Method "GET" 
# Assuming there is at least one category seeded or empty list.
# For this test we need a product. Since we don't have a "Create Product" endpoint exposed easily (it might require admin or specific role),
# let's check if we can Create Product. The handler `productH.Create` (if enabled) or we assume database seeding seeded some products?
# The code doesn't show product seeding in the visible files (admin_seed/role_seed). 
# Wait, I checked `admin_seed.go` and `role_seed.go`. Is there a `product_seed.go`?
# If no products traverse, we can't add to cart.

# Let's try to verify if there are products.
$Products = Invoke-Api -Url "http://localhost:8080/api/v1/products" -Method "GET"
$ProductID = 0

if ($Products.data.Count -gt 0) {
    $ProductID = $Products.data[0].id
    Write-Host "Found existing product ID: $ProductID"
}
else {
    Write-Host "No products found. Skipping Cart/Order test."
    exit
}

# 3. Add to Cart
Write-Host "`n--- 3. Add to Cart ---"
Invoke-Api -Url "http://localhost:8080/api/v1/carts" -Method "POST" -Token $Token -Body @{
    product_id = $ProductID
    quantity   = 2
    size       = "M"
    color      = "Black"
    sku        = "SKU-TEST-001"
}
Write-Host "Added to Cart"

# 4. Get Cart
Write-Host "`n--- 4. Get Cart ---"
$Cart = Invoke-Api -Url "http://localhost:8080/api/v1/carts" -Method "GET" -Token $Token
Write-Host "Cart Items: $($Cart.data.Count)"

# 5. Create Order
Write-Host "`n--- 5. Create Order ---"
# Need address ID? Getting User Profile first?
# The Order Create endpoint usually takes address/payment info.
# Let's inspect the `CreateOrder` handler requirement briefly or just try sending minimal body.
# `orderH.CreateOrder`
Invoke-Api -Url "http://localhost:8080/api/v1/orders" -Method "POST" -Token $Token -Body @{
    address_id        = 1
    payment_method_id = 1
    couriers_id       = 1
    total_amount      = 50000 
    shipping_type     = "Regular"
    payment_type      = "Bank Transfer"
    order_time        = "12:00"
    order_date        = "2025-01-01"
    order_details     = @(
        @{
            product_id = $ProductID
            quantity   = 2
            size       = "M"
            color      = "Black"
            sku        = "SKU-TEST-001"
        }
    )
}
Write-Host "Order Created"
