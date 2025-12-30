-- Phase 1: Critical Performance Indexes (ADJUSTED FOR EXISTING TABLES)
-- Orders and Products tables only

-- ============================================
-- ORDERS TABLE INDEXES
-- ============================================

-- Index for order code lookups (already has unique constraint, but verify)
CREATE INDEX IF NOT EXISTS idx_orders_order_code ON orders(order_code);

-- Index for status filtering
CREATE INDEX IF NOT EXISTS idx_orders_status ON orders(status);

-- Index for buyer filtering
CREATE INDEX IF NOT EXISTS idx_orders_buyer_id ON orders(buyer_id);

-- Index for date sorting (DESC for recent orders first)
CREATE INDEX IF NOT EXISTS idx_orders_order_date_desc ON orders(order_date DESC);

-- Composite index for common filter combination (buyer + status)
CREATE INDEX IF NOT EXISTS idx_orders_buyer_status ON orders(buyer_id, status);

-- ============================================
-- PRODUCTS TABLE INDEXES
-- ============================================

-- Index for parent product filtering
CREATE INDEX IF NOT EXISTS idx_products_parent_id ON products(parent_id);

-- Index for status filtering (ACTIVE/DRAFT)
CREATE INDEX IF NOT EXISTS idx_products_status ON products(status);

-- Index for category filtering
CREATE INDEX IF NOT EXISTS idx_products_category_slug ON products(category_slug);

-- Index for price range filtering
CREATE INDEX IF NOT EXISTS idx_products_sale_price ON products(sale_price);

-- Composite index for common filter: status + category
CREATE INDEX IF NOT EXISTS idx_products_status_category ON products(status, category_slug);

-- Composite index for parent products with status
CREATE INDEX IF NOT EXISTS idx_products_parent_status ON products(parent_id, status);

-- Partial index for parent products only (WHERE parent_id IS NULL)
CREATE INDEX IF NOT EXISTS idx_products_parent_null_status 
ON products(status, category_slug) 
WHERE parent_id IS NULL;

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
-- CATEGORIES TABLE - SLUG LOOKUP
-- ============================================

-- Index for category slug (if not already unique)
CREATE INDEX IF NOT EXISTS idx_categories_slug 
ON categories(slug);

-- ============================================
-- ANALYZE TABLES
-- ============================================

ANALYZE orders;
ANALYZE products;
ANALYZE order_items;
ANALYZE categories;
