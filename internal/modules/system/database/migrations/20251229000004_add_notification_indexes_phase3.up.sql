-- Phase 3: Notifications and Additional Indexes
-- Specialized indexes for notification queries

-- ============================================
-- NOTIFICATIONS TABLE - ADVANCED INDEXES
-- ============================================

-- Index for receiver filtering
CREATE INDEX IF NOT EXISTS idx_notifications_receiver_id 
ON notifications(receiver_id);

-- Index for status filtering
CREATE INDEX IF NOT EXISTS idx_notifications_status 
ON notifications(status);

-- Partial index for unread notifications (most common query)
CREATE INDEX IF NOT EXISTS idx_notifications_unread 
ON notifications(receiver_id, sent_at DESC) 
WHERE read_at IS NULL;

-- Partial index for read notifications
CREATE INDEX IF NOT EXISTS idx_notifications_read 
ON notifications(receiver_id, read_at DESC) 
WHERE read_at IS NOT NULL;

-- Composite index for notification listing with status
CREATE INDEX IF NOT EXISTS idx_notifications_receiver_status 
ON notifications(receiver_id, status, sent_at DESC);

-- ============================================
-- ORDER ITEMS TABLE - FOREIGN KEY OPTIMIZATION
-- ============================================

-- Index for order_id foreign key lookups
CREATE INDEX IF NOT EXISTS idx_order_items_order_id 
ON order_items(order_id);

-- Index for product_id lookups (for inventory checks)
CREATE INDEX IF NOT EXISTS idx_order_items_product_id 
ON order_items(product_id);

-- ============================================
-- PAYMENT LOGS TABLE - AUDIT TRAIL
-- ============================================

-- Index for payment_id foreign key
CREATE INDEX IF NOT EXISTS idx_payment_logs_payment_id 
ON payment_logs(payment_id);

-- Index for status filtering in logs
CREATE INDEX IF NOT EXISTS idx_payment_logs_status 
ON payment_logs(status);

-- Composite index for payment audit trail
CREATE INDEX IF NOT EXISTS idx_payment_logs_payment_created 
ON payment_logs(payment_id, created_at DESC);

-- ============================================
-- CATEGORIES TABLE - SLUG LOOKUP
-- ============================================

-- Index for category slug (if not already unique)
CREATE INDEX IF NOT EXISTS idx_categories_slug 
ON categories(slug);

-- ============================================
-- ANALYZE TABLES
-- ============================================

ANALYZE notifications;
ANALYZE order_items;
ANALYZE payment_logs;
ANALYZE categories;
