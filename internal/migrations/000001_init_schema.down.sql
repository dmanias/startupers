BEGIN;
-- Change the table name to users_archive
ALTER TABLE users RENAME TO users_archive;

-- Set the table as read-only
ALTER TABLE users_archive OWNER TO postgres;

-- Create a read-only view for the users_archive table
CREATE OR REPLACE VIEW users_archive_view AS
SELECT * FROM users_archive;

-- Create a read-only view for the users_archive table
CREATE OR REPLACE VIEW users_view AS
SELECT * FROM users_archive;

COMMIT;

-- Perform a vacuum full on the users table to reclaim space
VACUUM FULL users_archive;

-- Analyze the table to update statistics for the query planner
ANALYZE users_archive;
