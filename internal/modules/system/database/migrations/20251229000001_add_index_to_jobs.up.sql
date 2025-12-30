-- Menambahkan Composite Index pada kolom status dan created_at
-- Ini akan mempercepat query: WHERE status = 'pending' ORDER BY created_at ASC
CREATE INDEX IF NOT EXISTS idx_jobs_status_created_at ON jobs (status, created_at);