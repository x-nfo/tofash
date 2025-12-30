# Migration Guide - Performance Indexes

## 📋 Overview

Migration files telah dibuat untuk menerapkan performance indexes yang direkomendasikan dari SQL Performance Audit. Migrations dibagi menjadi 3 fase untuk memudahkan rollout bertahap.

---

## 📁 Migration Files

### Phase 1: Critical Indexes (PRIORITY: HIGH)
- **File**: `20251229000002_add_performance_indexes_phase1.up.sql`
- **Rollback**: `20251229000002_add_performance_indexes_phase1.down.sql`
- **Tables**: `orders`, `products`, `payments`
- **Impact**: 8-10x faster queries untuk product listing, order filtering, payment history

### Phase 2: Search Optimization (PRIORITY: MEDIUM)
- **File**: `20251229000003_add_search_indexes_phase2.up.sql`
- **Rollback**: `20251229000003_add_search_indexes_phase2.down.sql`
- **Tables**: `products`, `users`, `notifications`
- **Requires**: PostgreSQL `pg_trgm` extension
- **Impact**: Fuzzy search 10x lebih cepat

### Phase 3: Supporting Tables (PRIORITY: LOW)
- **File**: `20251229000004_add_notification_indexes_phase3.up.sql`
- **Rollback**: `20251229000004_add_notification_indexes_phase3.down.sql`
- **Tables**: `notifications`, `order_items`, `payment_logs`, `categories`
- **Impact**: Optimasi untuk notification queries dan foreign key lookups

---

## 🚀 How to Apply Migrations

### Option 1: Using Migration Tool (Recommended)

Jika aplikasi menggunakan migration runner (seperti golang-migrate atau gorm AutoMigrate):

```bash
# Check migration status
migrate -path internal/modules/system/database/migrations -database "postgresql://user:pass@localhost:5432/tofash_db?sslmode=disable" version

# Apply all pending migrations
migrate -path internal/modules/system/database/migrations -database "postgresql://user:pass@localhost:5432/tofash_db?sslmode=disable" up

# Apply specific version
migrate -path internal/modules/system/database/migrations -database "postgresql://user:pass@localhost:5432/tofash_db?sslmode=disable" goto 20251229000002
```

### Option 2: Manual Application (Development)

```bash
# Connect to database
psql -U your_user -d tofash_db

# Apply Phase 1
\i internal/modules/system/database/migrations/20251229000002_add_performance_indexes_phase1.up.sql

# Apply Phase 2
\i internal/modules/system/database/migrations/20251229000003_add_search_indexes_phase2.up.sql

# Apply Phase 3
\i internal/modules/system/database/migrations/20251229000004_add_notification_indexes_phase3.up.sql
```

### Option 3: Production Deployment (CONCURRENT)

Untuk production, gunakan `CREATE INDEX CONCURRENTLY` untuk menghindari table locks:

```sql
-- Example: Create index without blocking writes
CREATE INDEX CONCURRENTLY idx_products_name_trgm 
ON products USING gin(name gin_trgm_ops);
```

> ⚠️ **WARNING**: `CONCURRENTLY` tidak bisa dijalankan dalam transaction block. Modify migration files jika diperlukan.

---

## ✅ Verification Steps

### 1. Check if Migrations Applied

```sql
-- List all indexes
SELECT 
    schemaname,
    tablename,
    indexname,
    indexdef
FROM pg_indexes
WHERE schemaname = 'public'
ORDER BY tablename, indexname;
```

### 2. Verify Specific Indexes

```sql
-- Check products table indexes
SELECT indexname, indexdef 
FROM pg_indexes 
WHERE tablename = 'products';

-- Check if pg_trgm extension is enabled
SELECT * FROM pg_extension WHERE extname = 'pg_trgm';
```

### 3. Check Index Sizes

```sql
SELECT 
    tablename,
    indexname,
    pg_size_pretty(pg_relation_size(indexname::regclass)) AS index_size
FROM pg_indexes
WHERE schemaname = 'public'
ORDER BY pg_relation_size(indexname::regclass) DESC;
```

### 4. Test Query Performance

```sql
-- Before: Should show Seq Scan
EXPLAIN ANALYZE
SELECT * FROM products 
WHERE parent_id IS NULL 
  AND status = 'ACTIVE' 
  AND name ILIKE '%shirt%'
ORDER BY created_at DESC
LIMIT 20;

-- After: Should show Index Scan
-- Look for "Index Scan using idx_products_..." in output
```

---

## 📊 Expected Results

### Index Count by Table

| Table | Indexes Before | Indexes After | New Indexes |
|-------|---------------|---------------|-------------|
| `orders` | ~2 | ~7 | +5 |
| `products` | ~2 | ~10 | +8 |
| `payments` | ~1 | ~7 | +6 |
| `users` | ~3 | ~7 | +4 |
| `notifications` | ~1 | ~6 | +5 |
| `order_items` | ~1 | ~3 | +2 |
| `payment_logs` | ~1 | ~4 | +3 |
| `categories` | ~1 | ~2 | +1 |

**Total New Indexes**: ~34

### Storage Impact

Estimated additional storage for indexes:
- Small dataset (<10K products): ~50-100 MB
- Medium dataset (10K-100K products): ~200-500 MB
- Large dataset (>100K products): ~1-2 GB

---

## 🔄 Rollback Instructions

Jika terjadi masalah, rollback migrations dalam urutan terbalik:

```bash
# Rollback Phase 3
psql -U your_user -d tofash_db < internal/modules/system/database/migrations/20251229000004_add_notification_indexes_phase3.down.sql

# Rollback Phase 2
psql -U your_user -d tofash_db < internal/modules/system/database/migrations/20251229000003_add_search_indexes_phase2.down.sql

# Rollback Phase 1
psql -U your_user -d tofash_db < internal/modules/system/database/migrations/20251229000002_add_performance_indexes_phase1.down.sql
```

Or using migrate tool:
```bash
migrate -path internal/modules/system/database/migrations -database "postgresql://..." down 3
```

---

## ⚠️ Important Notes

1. **Backup First**: Always backup database before applying migrations
2. **Test in Development**: Apply and test in dev environment first
3. **Monitor Performance**: Watch for slow index creation on large tables
4. **Disk Space**: Ensure sufficient disk space for indexes
5. **pg_trgm Extension**: Phase 2 requires PostgreSQL 9.1+ with pg_trgm extension

---

## 🧪 Performance Testing

### Before Applying Indexes

```bash
# Run benchmark queries
psql -U your_user -d tofash_db -f test_queries_before.sql
```

### After Applying Indexes

```bash
# Run same queries and compare
psql -U your_user -d tofash_db -f test_queries_after.sql
```

### Sample Test Queries

```sql
-- Product search
\timing on
SELECT * FROM products 
WHERE parent_id IS NULL 
  AND status = 'ACTIVE' 
  AND name ILIKE '%shirt%'
ORDER BY created_at DESC
LIMIT 20;

-- Order listing by buyer
SELECT * FROM orders 
WHERE buyer_id = 123 
  AND status = 'Pending'
ORDER BY order_date DESC
LIMIT 50;

-- User search
SELECT * FROM users 
WHERE name ILIKE '%john%' 
   OR email ILIKE '%john%'
LIMIT 20;
```

---

## 📞 Support

Jika mengalami masalah:
1. Check PostgreSQL logs: `/var/log/postgresql/postgresql-*.log`
2. Verify index creation: `SELECT * FROM pg_stat_progress_create_index;`
3. Check for blocking queries: `SELECT * FROM pg_stat_activity WHERE wait_event_type = 'Lock';`

---

## ✨ Next Steps

1. ✅ Apply Phase 1 migrations (critical)
2. ✅ Verify indexes created successfully
3. ✅ Test query performance improvements
4. ✅ Apply Phase 2 migrations (search optimization)
5. ✅ Apply Phase 3 migrations (supporting tables)
6. ✅ Monitor application performance
7. ✅ Update audit report with actual results
