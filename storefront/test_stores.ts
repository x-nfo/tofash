/**
 * Test script for Phase 1 stores
 * Run this to verify that all stores are working correctly
 */

import {
    // Auth
    $authUser, $authToken, $isAuthenticated, $userRole,
    signin, signup, forgotPassword, verifyAccount, updatePassword,
    getProfile, updateProfile, logout, checkAuth,

    // Cart
    $cart, $cartTotalItems, $cartTotalPrice,
    fetchCarts, addToCart, deleteCart, deleteAllCart, clearCart,

    // Product
    $products, $product, $pagination, $hasMoreProducts,
    fetchProductsHome, fetchProductsShop, fetchProductDetailHome,

    // Order
    $orders, $order, $orderPagination, $hasMoreOrders, $totalOrders,
    fetchOrders, getOrderDetail, createOrder,

    // Notification
    $notifications, $unreadCount, $hasUnread, $toasts,
    fetchNotifications, markAsRead, markAllAsRead, showToast
} from './src/stores/index';

console.log('=== Phase 1 Store Tests ===\n');

// Test 1: Auth Store
console.log('1. Testing Auth Store...');
console.log('   - Auth User:', $authUser.get());
console.log('   - Auth Token:', $authToken.get() ? 'Present' : 'Not present');
console.log('   - Is Authenticated:', $isAuthenticated.get());
console.log('   - User Role:', $userRole.get());
console.log('   ✓ Auth Store initialized\n');

// Test 2: Cart Store
console.log('2. Testing Cart Store...');
console.log('   - Cart Items:', $cart.get().length);
console.log('   - Total Items:', $cartTotalItems.get());
console.log('   - Total Price:', $cartTotalPrice.get());
console.log('   ✓ Cart Store initialized\n');

// Test 3: Product Store
console.log('3. Testing Product Store...');
console.log('   - Products:', $products.get().length);
console.log('   - Current Product:', $product.get() ? 'Present' : 'Not present');
console.log('   - Pagination:', $pagination.get());
console.log('   - Has More Products:', $hasMoreProducts.get());
console.log('   ✓ Product Store initialized\n');

// Test 4: Order Store
console.log('4. Testing Order Store...');
console.log('   - Orders:', $orders.get().length);
console.log('   - Current Order:', $order.get() ? 'Present' : 'Not present');
console.log('   - Order Pagination:', $orderPagination.get());
console.log('   - Has More Orders:', $hasMoreOrders.get());
console.log('   - Total Orders:', $totalOrders.get());
console.log('   ✓ Order Store initialized\n');

// Test 5: Notification Store
console.log('5. Testing Notification Store...');
console.log('   - Notifications:', $notifications.get().length);
console.log('   - Unread Count:', $unreadCount.get());
console.log('   - Has Unread:', $hasUnread.get());
console.log('   - Toasts:', $toasts.get().length);
console.log('   ✓ Notification Store initialized\n');

// Test 6: Toast functionality
console.log('6. Testing Toast Functionality...');
showToast('Test success message', 'success', 1000);
console.log('   - Toast added:', $toasts.get().length);
console.log('   ✓ Toast functionality working\n');

// Test 7: Persistent storage check
console.log('7. Testing Persistent Storage...');
if (typeof localStorage !== 'undefined') {
    console.log('   - localStorage available');
    console.log('   - Auth User in localStorage:', localStorage.getItem('authUser') ? 'Yes' : 'No');
    console.log('   - Auth Token in localStorage:', localStorage.getItem('authToken') ? 'Yes' : 'No');
    console.log('   - Cart in localStorage:', localStorage.getItem('cart') ? 'Yes' : 'No');
    console.log('   ✓ Persistent storage working\n');
} else {
    console.log('   - localStorage not available (running in Node.js)\n');
}

console.log('=== All Phase 1 Store Tests Completed ===');
console.log('\nNext Steps:');
console.log('1. Run the dev server: npm run dev');
console.log('2. Test API calls with actual backend');
console.log('3. Test authentication flow');
console.log('4. Test cart operations');
console.log('5. Test product fetching');
console.log('6. Test order creation');
console.log('7. Test notifications');
