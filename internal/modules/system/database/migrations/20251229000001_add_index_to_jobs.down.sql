-- Menghapus index jika migrasi di-rollback
DROP INDEX IF EXISTS idx_jobs_status_created_at;