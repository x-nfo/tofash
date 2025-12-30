-- Phase 2: Search Optimization with Trigram Indexes (ADJUSTED)
-- Enable fuzzy search for products and users only

-- ============================================
-- ENABLE TRIGRAM EXTENSION
-- ============================================

CREATE EXTENSION IF NOT EXISTS pg_trgm;

-- ============================================
-- PRODUCTS TABLE - FULL TEXT SEARCH
-- ============================================

-- Trigram index for product name fuzzy search
CREATE INDEX IF NOT EXISTS idx_products_name_trgm 
ON products USING gin(name gin_trgm_ops);

-- Trigram index for product description fuzzy search
CREATE INDEX IF NOT EXISTS idx_products_description_trgm 
ON products USING gin(description gin_trgm_ops);

-- ============================================
-- USERS TABLE - SEARCH OPTIMIZATION
-- ============================================

-- Trigram index for user name search
CREATE INDEX IF NOT EXISTS idx_users_name_trgm 
ON users USING gin(name gin_trgm_ops);

-- Trigram index for phone number search (if phone column exists)
CREATE INDEX IF NOT EXISTS idx_users_phone_trgm 
ON users USING gin(phone gin_trgm_ops);

-- Trigram index for email search (for ILIKE queries)
CREATE INDEX IF NOT EXISTS idx_users_email_trgm 
ON users USING gin(email gin_trgm_ops);

-- Composite index for login optimization
CREATE INDEX IF NOT EXISTS idx_users_email_verified 
ON users(email, is_verified);

-- ============================================
-- ANALYZE TABLES
-- ============================================

ANALYZE products;
ANALYZE users;
