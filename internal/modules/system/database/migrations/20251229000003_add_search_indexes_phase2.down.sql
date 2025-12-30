-- Rollback Phase 2: Drop Search Optimization Indexes

-- ============================================
-- DROP NOTIFICATIONS TABLE INDEXES
-- ============================================

-- DROP INDEX IF EXISTS idx_notifications_message_trgm;
DROP INDEX IF EXISTS idx_notifications_subject_trgm;

-- ============================================
-- DROP USERS TABLE INDEXES
-- ============================================

DROP INDEX IF EXISTS idx_users_email_verified;
DROP INDEX IF EXISTS idx_users_email_trgm;
DROP INDEX IF EXISTS idx_users_phone_trgm;
DROP INDEX IF EXISTS idx_users_name_trgm;

-- ============================================
-- DROP PRODUCTS TABLE INDEXES
-- ============================================

-- DROP INDEX IF EXISTS idx_products_fulltext;
DROP INDEX IF EXISTS idx_products_description_trgm;
DROP INDEX IF EXISTS idx_products_name_trgm;

-- ============================================
-- NOTE: pg_trgm extension is NOT dropped
-- as it might be used by other applications
-- If you want to drop it, run manually:
-- DROP EXTENSION IF EXISTS pg_trgm;
-- ============================================
