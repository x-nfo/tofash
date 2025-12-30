-- Rollback Phase 3: Drop Notification and Additional Indexes

-- ============================================
-- DROP CATEGORIES TABLE INDEXES
-- ============================================

DROP INDEX IF EXISTS idx_categories_slug;

-- ============================================
-- DROP PAYMENT LOGS TABLE INDEXES
-- ============================================

DROP INDEX IF EXISTS idx_payment_logs_payment_created;
DROP INDEX IF EXISTS idx_payment_logs_status;
DROP INDEX IF EXISTS idx_payment_logs_payment_id;

-- ============================================
-- DROP ORDER ITEMS TABLE INDEXES
-- ============================================

DROP INDEX IF EXISTS idx_order_items_product_id;
DROP INDEX IF EXISTS idx_order_items_order_id;

-- ============================================
-- DROP NOTIFICATIONS TABLE INDEXES
-- ============================================

DROP INDEX IF EXISTS idx_notifications_receiver_status;
DROP INDEX IF EXISTS idx_notifications_read;
DROP INDEX IF EXISTS idx_notifications_unread;
DROP INDEX IF EXISTS idx_notifications_status;
DROP INDEX IF EXISTS idx_notifications_receiver_id;
