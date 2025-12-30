-- Verification Queries for Performance Indexes
-- Run these queries to verify indexes are created and working

-- ============================================
-- 1. CHECK ALL INDEXES
-- ============================================

SELECT 
    schemaname,
    tablename,
    indexname,
    indexdef
FROM pg_indexes
WHERE schemaname = 'public'
ORDER BY tablename, indexname;

-- ============================================
-- 2. CHECK SPECIFIC TABLE INDEXES
-- ============================================

-- Orders table
SELECT indexname, indexdef 
FROM pg_indexes 
WHERE tablename = 'orders'
ORDER BY indexname;

-- Products table
SELECT indexname, indexdef 
FROM pg_indexes 
WHERE tablename = 'products'
ORDER BY indexname;

-- Payments table
SELECT indexname, indexdef 
FROM pg_indexes 
WHERE tablename = 'payments'
ORDER BY indexname;

-- Users table
SELECT indexname, indexdef 
FROM pg_indexes 
WHERE tablename = 'users'
ORDER BY indexname;

-- Notifications table
SELECT indexname, indexdef 
FROM pg_indexes 
WHERE tablename = 'notifications'
ORDER BY indexname;

-- ============================================
-- 3. CHECK EXTENSION
-- ============================================

-- Check if pg_trgm extension is enabled
SELECT * FROM pg_extension WHERE extname = 'pg_trgm';

-- ============================================
-- 4. INDEX SIZES
-- ============================================

SELECT 
    tablename,
    indexname,
    pg_size_pretty(pg_relation_size(indexname::regclass)) AS index_size,
    pg_size_pretty(pg_total_relation_size(tablename::regclass)) AS table_total_size
FROM pg_indexes
WHERE schemaname = 'public'
ORDER BY pg_relation_size(indexname::regclass) DESC;

-- ============================================
-- 5. TABLE STATISTICS
-- ============================================

SELECT 
    schemaname,
    tablename,
    pg_size_pretty(pg_total_relation_size(schemaname||'.'||tablename)) AS total_size,
    pg_size_pretty(pg_relation_size(schemaname||'.'||tablename)) AS table_size,
    pg_size_pretty(pg_total_relation_size(schemaname||'.'||tablename) - pg_relation_size(schemaname||'.'||tablename)) AS indexes_size
FROM pg_tables
WHERE schemaname = 'public'
ORDER BY pg_total_relation_size(schemaname||'.'||tablename) DESC;

-- ============================================
-- 6. VERIFY SPECIFIC INDEXES EXIST
-- ============================================

-- Check critical indexes
SELECT 
    CASE 
        WHEN EXISTS (SELECT 1 FROM pg_indexes WHERE indexname = 'idx_orders_buyer_id') 
        THEN '✓ idx_orders_buyer_id exists'
        ELSE '✗ idx_orders_buyer_id MISSING'
    END AS orders_buyer_id,
    
    CASE 
        WHEN EXISTS (SELECT 1 FROM pg_indexes WHERE indexname = 'idx_products_status_category') 
        THEN '✓ idx_products_status_category exists'
        ELSE '✗ idx_products_status_category MISSING'
    END AS products_status_category,
    
    CASE 
        WHEN EXISTS (SELECT 1 FROM pg_indexes WHERE indexname = 'idx_products_name_trgm') 
        THEN '✓ idx_products_name_trgm exists'
        ELSE '✗ idx_products_name_trgm MISSING'
    END AS products_name_trgm,
    
    CASE 
        WHEN EXISTS (SELECT 1 FROM pg_indexes WHERE indexname = 'idx_payments_user_id') 
        THEN '✓ idx_payments_user_id exists'
        ELSE '✗ idx_payments_user_id MISSING'
    END AS payments_user_id,
    
    CASE 
        WHEN EXISTS (SELECT 1 FROM pg_indexes WHERE indexname = 'idx_notifications_unread') 
        THEN '✓ idx_notifications_unread exists'
        ELSE '✗ idx_notifications_unread MISSING'
    END AS notifications_unread;

-- ============================================
-- 7. INDEX USAGE STATISTICS
-- ============================================

SELECT 
    schemaname,
    tablename,
    indexname,
    idx_scan as index_scans,
    idx_tup_read as tuples_read,
    idx_tup_fetch as tuples_fetched
FROM pg_stat_user_indexes
WHERE schemaname = 'public'
ORDER BY idx_scan DESC;

-- ============================================
-- 8. UNUSED INDEXES (Run after some time)
-- ============================================

SELECT 
    schemaname,
    tablename,
    indexname,
    idx_scan,
    pg_size_pretty(pg_relation_size(indexrelid)) as index_size
FROM pg_stat_user_indexes
WHERE schemaname = 'public'
  AND idx_scan = 0
  AND indexrelname NOT LIKE '%pkey'
ORDER BY pg_relation_size(indexrelid) DESC;
