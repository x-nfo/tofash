-- Rollback Phase 1: Drop Critical Performance Indexes

-- ============================================
-- DROP ORDERS TABLE INDEXES
-- ============================================

DROP INDEX IF EXISTS idx_orders_buyer_status;
DROP INDEX IF EXISTS idx_orders_order_date_desc;
DROP INDEX IF EXISTS idx_orders_buyer_id;
DROP INDEX IF EXISTS idx_orders_status;
DROP INDEX IF EXISTS idx_orders_order_code;

-- ============================================
-- DROP PRODUCTS TABLE INDEXES
-- ============================================

DROP INDEX IF EXISTS idx_products_parent_null_status;
DROP INDEX IF EXISTS idx_products_parent_status;
DROP INDEX IF EXISTS idx_products_status_category;
DROP INDEX IF EXISTS idx_products_sale_price;
DROP INDEX IF EXISTS idx_products_category_slug;
DROP INDEX IF EXISTS idx_products_status;
DROP INDEX IF EXISTS idx_products_parent_id;

-- ============================================
-- DROP PAYMENTS TABLE INDEXES
-- ============================================

DROP INDEX IF EXISTS idx_payments_user_created;
DROP INDEX IF EXISTS idx_payments_created_at_desc;
DROP INDEX IF EXISTS idx_payments_order_id;
DROP INDEX IF EXISTS idx_payments_user_id;
DROP INDEX IF EXISTS idx_payments_payment_status;
DROP INDEX IF EXISTS idx_payments_payment_method;
