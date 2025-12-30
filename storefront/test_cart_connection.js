/**
 * Test Script for Cart API Connection
 * 
 * This script tests the add to cart functionality and its connection to the backend.
 * 
 * Usage:
 *   node test_cart_connection.js
 * 
 * Prerequisites:
 *   - Backend server must be running on http://localhost:8080
 *   - Redis must be running
 *   - Valid JWT token (get from login or set in environment variable)
 */

const API_BASE_URL = process.env.API_URL || 'http://localhost:8080/api/v1';
const ACCESS_TOKEN = process.env.ACCESS_TOKEN || '';

// Colors for console output
const colors = {
    reset: '\x1b[0m',
    green: '\x1b[32m',
    red: '\x1b[31m',
    yellow: '\x1b[33m',
    blue: '\x1b[34m',
    cyan: '\x1b[36m'
};

function log(message, color = colors.reset) {
    console.log(`${color}${message}${colors.reset}`);
}

function logTest(name, passed, details = {}) {
    const icon = passed ? '✅' : '❌';
    const statusColor = passed ? colors.green : colors.red;
    log(`${icon} ${name}`, statusColor);
    if (Object.keys(details).length > 0) {
        log(`   Details: ${JSON.stringify(details, null, 2)}`, colors.cyan);
    }
    return passed;
}

async function makeRequest(method, endpoint, data = null) {
    const url = `${API_BASE_URL}${endpoint}`;

    const options = {
        method: method,
        headers: {
            'Content-Type': 'application/json',
        },
    };

    if (ACCESS_TOKEN) {
        options.headers['Authorization'] = `Bearer ${ACCESS_TOKEN}`;
    }

    if (data) {
        options.body = JSON.stringify(data);
    }

    try {
        const response = await fetch(url, options);
        const responseData = await response.json();
        return {
            ok: response.ok,
            status: response.status,
            data: responseData
        };
    } catch (error) {
        return {
            ok: false,
            status: 0,
            data: { error: error.message }
        };
    }
}

// Test 1: Backend Connection
async function testBackendConnection() {
    log('\n' + '='.repeat(60), colors.blue);
    log('TEST 1: Backend Connection', colors.blue);
    log('='.repeat(60), colors.blue);

    const result = await makeRequest('GET', '/products?limit=1');

    const passed = logTest(
        'Backend is accessible',
        result.ok,
        {
            status: result.status,
            endpoint: '/products?limit=1'
        }
    );

    return passed;
}

// Test 2: Add to Cart
async function testAddToCart(productId = 1, quantity = 1, size = 'M', color = 'Red') {
    log('\n' + '='.repeat(60), colors.blue);
    log('TEST 2: Add to Cart', colors.blue);
    log('='.repeat(60), colors.blue);

    const result = await makeRequest('POST', '/carts', {
        product_id: productId,
        quantity: quantity,
        size: size,
        color: color,
        sku: `SKU-${productId}-${size}-${color}`
    });

    const passed = logTest(
        'Add to cart endpoint works',
        result.ok,
        {
            status: result.status,
            message: result.data.message,
            data: result.data
        }
    );

    return passed;
}

// Test 3: Get Cart
async function testGetCart() {
    log('\n' + '='.repeat(60), colors.blue);
    log('TEST 3: Get Cart', colors.blue);
    log('='.repeat(60), colors.blue);

    const result = await makeRequest('GET', '/carts');

    const passed = logTest(
        'Get cart endpoint works',
        result.ok,
        {
            status: result.status,
            itemCount: result.data?.data?.length || 0,
            data: result.data
        }
    );

    return passed;
}

// Test 4: Remove from Cart
async function testRemoveFromCart(productId = 1) {
    log('\n' + '='.repeat(60), colors.blue);
    log('TEST 4: Remove from Cart', colors.blue);
    log('='.repeat(60), colors.blue);

    const result = await makeRequest('DELETE', `/carts?product_id=${productId}`);

    const passed = logTest(
        'Remove from cart endpoint works',
        result.ok,
        {
            status: result.status,
            message: result.data.message,
            data: result.data
        }
    );

    return passed;
}

// Test 5: Clear Cart
async function testClearCart() {
    log('\n' + '='.repeat(60), colors.blue);
    log('TEST 5: Clear Cart', colors.blue);
    log('='.repeat(60), colors.blue);

    const result = await makeRequest('DELETE', '/carts/all');

    const passed = logTest(
        'Clear cart endpoint works',
        result.ok,
        {
            status: result.status,
            message: result.data.message,
            data: result.data
        }
    );

    return passed;
}

// Test 6: Full Workflow
async function testFullWorkflow(productId = 1, quantity = 1, size = 'M', color = 'Red') {
    log('\n' + '='.repeat(60), colors.blue);
    log('TEST 6: Full Workflow (Add -> Get -> Remove)', colors.blue);
    log('='.repeat(60), colors.blue);

    const steps = [];
    let allPassed = true;

    // Step 1: Add to cart
    log('\nStep 1: Adding item to cart...', colors.yellow);
    const addResult = await makeRequest('POST', '/carts', {
        product_id: productId,
        quantity: quantity,
        size: size,
        color: color,
        sku: `SKU-${productId}-${size}-${color}`
    });
    const addPassed = logTest('Add to cart', addResult.ok, { status: addResult.status });
    steps.push({ step: 'Add to Cart', passed: addPassed, status: addResult.status });
    allPassed = allPassed && addPassed;

    // Wait a bit for Redis to update
    await new Promise(resolve => setTimeout(resolve, 100));

    // Step 2: Get cart
    log('\nStep 2: Retrieving cart...', colors.yellow);
    const getResult = await makeRequest('GET', '/carts');
    const getPassed = logTest('Get cart', getResult.ok, {
        status: getResult.status,
        itemCount: getResult.data?.data?.length || 0
    });
    steps.push({ step: 'Get Cart', passed: getPassed, status: getResult.status });
    allPassed = allPassed && getPassed;

    // Step 3: Remove from cart
    log('\nStep 3: Removing item from cart...', colors.yellow);
    const removeResult = await makeRequest('DELETE', `/carts?product_id=${productId}`);
    const removePassed = logTest('Remove from cart', removeResult.ok, { status: removeResult.status });
    steps.push({ step: 'Remove from Cart', passed: removePassed, status: removeResult.status });
    allPassed = allPassed && removePassed;

    log('\n' + '-'.repeat(60), colors.cyan);
    log('Workflow Summary:', colors.cyan);
    steps.forEach((step, index) => {
        const icon = step.passed ? '✅' : '❌';
        log(`  ${index + 1}. ${step.step}: ${icon}`, step.passed ? colors.green : colors.red);
    });

    return allPassed;
}

// Test 7: Test Multiple Items
async function testMultipleItems() {
    log('\n' + '='.repeat(60), colors.blue);
    log('TEST 7: Add Multiple Items to Cart', colors.blue);
    log('='.repeat(60), colors.blue);

    const items = [
        { product_id: 1, quantity: 2, size: 'M', color: 'Red' },
        { product_id: 2, quantity: 1, size: 'L', color: 'Blue' },
        { product_id: 3, quantity: 3, size: 'S', color: 'Green' }
    ];

    let allPassed = true;

    for (const item of items) {
        log(`\nAdding product ${item.product_id}...`, colors.yellow);
        const result = await makeRequest('POST', '/carts', {
            product_id: item.product_id,
            quantity: item.quantity,
            size: item.size,
            color: item.color,
            sku: `SKU-${item.product_id}-${item.size}-${item.color}`
        });

        const passed = logTest(
            `Add product ${item.product_id} to cart`,
            result.ok,
            { status: result.status }
        );
        allPassed = allPassed && passed;

        // Wait a bit between requests
        await new Promise(resolve => setTimeout(resolve, 100));
    }

    // Get cart to verify all items
    log('\nVerifying cart contents...', colors.yellow);
    const getResult = await makeRequest('GET', '/carts');
    const getPassed = logTest(
        'Verify cart has multiple items',
        getResult.ok && getResult.data?.data?.length >= items.length,
        {
            expected: items.length,
            actual: getResult.data?.data?.length || 0
        }
    );
    allPassed = allPassed && getPassed;

    return allPassed;
}

// Test 8: Test Update Quantity
async function testUpdateQuantity(productId = 1) {
    log('\n' + '='.repeat(60), colors.blue);
    log('TEST 8: Update Item Quantity', colors.blue);
    log('='.repeat(60), colors.blue);

    // First add item
    log('\nAdding initial item...', colors.yellow);
    await makeRequest('POST', '/carts', {
        product_id: productId,
        quantity: 1,
        size: 'M',
        color: 'Red',
        sku: `SKU-${productId}-M-Red`
    });

    await new Promise(resolve => setTimeout(resolve, 100));

    // Add same item again (should update quantity)
    log('\nAdding same item again (should update quantity)...', colors.yellow);
    const result = await makeRequest('POST', '/carts', {
        product_id: productId,
        quantity: 2,
        size: 'M',
        color: 'Red',
        sku: `SKU-${productId}-M-Red`
    });

    // Get cart to verify quantity updated
    await new Promise(resolve => setTimeout(resolve, 100));
    const getResult = await makeRequest('GET', '/carts');

    const cartItem = getResult.data?.data?.find(item => item.id === productId);
    const quantity = cartItem?.quantity || 0;

    const passed = logTest(
        'Quantity updated correctly',
        quantity === 3, // 1 + 2 = 3
        {
            expected: 3,
            actual: quantity
        }
    );

    // Cleanup
    await makeRequest('DELETE', `/carts?product_id=${productId}`);

    return passed;
}

// Main test runner
async function runAllTests() {
    log('\n' + '='.repeat(60), colors.cyan);
    log('🧪 CART API CONNECTION TEST SUITE', colors.cyan);
    log('='.repeat(60), colors.cyan);
    log(`API URL: ${API_BASE_URL}`, colors.cyan);
    log(`Access Token: ${ACCESS_TOKEN ? 'Set' : 'Not Set'}`, colors.cyan);
    log('='.repeat(60), colors.cyan);

    const results = [];

    try {
        results.push(await testBackendConnection());
        results.push(await testAddToCart());
        results.push(await testGetCart());
        results.push(await testRemoveFromCart());
        results.push(await testClearCart());
        results.push(await testFullWorkflow());
        results.push(await testMultipleItems());
        results.push(await testUpdateQuantity());

        // Summary
        const passed = results.filter(r => r).length;
        const failed = results.length - passed;
        const successRate = ((passed / results.length) * 100).toFixed(2);

        log('\n' + '='.repeat(60), colors.cyan);
        log('📊 TEST SUMMARY', colors.cyan);
        log('='.repeat(60), colors.cyan);
        log(`Total Tests: ${results.length}`, colors.cyan);
        log(`Passed: ${passed}`, colors.green);
        log(`Failed: ${failed}`, colors.red);
        log(`Success Rate: ${successRate}%`, colors.cyan);
        log('='.repeat(60), colors.cyan);

        if (failed === 0) {
            log('\n🎉 All tests passed! Cart API is working correctly.', colors.green);
        } else {
            log('\n⚠️  Some tests failed. Please check the results above.', colors.yellow);
        }

    } catch (error) {
        log(`\n❌ Error running tests: ${error.message}`, colors.red);
        console.error(error);
    }
}

// Run tests
runAllTests().catch(console.error);
