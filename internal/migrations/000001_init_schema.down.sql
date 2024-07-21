BEGIN;

-- Rename the users table to users_archive for archival purposes
ALTER TABLE users RENAME TO users_archive;

-- Drop the ai queries table if it's no longer needed
DROP TABLE IF EXISTS ais CASCADE;

-- Assuming you want to drop the original users_view if it's meant to be replaced or no longer needed
DROP VIEW IF EXISTS users_view CASCADE;

-- Create a read-only view for the users_archive table
CREATE OR REPLACE VIEW users_archive_view AS
SELECT * FROM users_archive;

-- Assuming the existence of a moderators table or view that also needs to be handled
-- DROP TABLE IF EXISTS moderators CASCADE;
-- DROP VIEW IF EXISTS moderators_view CASCADE;

COMMIT;

-- Optional: Cleaning up by dropping the archived table if it's not needed
-- Be very cautious with this step; ensure that you have backups or that this is indeed desired
-- DROP TABLE IF EXISTS users_archive CASCADE;

-- VACUUM and ANALYZE (Consider your database workload and maintenance windows before executing)
-- VACUUM (VERBOSE, ANALYZE) users_archive;

