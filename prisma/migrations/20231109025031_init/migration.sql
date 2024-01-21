-- Grant permissions to relevant roles on future created resources
ALTER DEFAULT PRIVILEGES GRANT SELECT ON TABLES TO donations_ro;

ALTER DEFAULT PRIVILEGES GRANT INSERT, UPDATE, DELETE ON TABLES TO donations_rw;
ALTER DEFAULT PRIVILEGES GRANT USAGE ON SEQUENCES TO donations_rw;
ALTER DEFAULT PRIVILEGES GRANT EXECUTE ON FUNCTIONS TO donations_rw;

ALTER DEFAULT PRIVILEGES GRANT DROP ON TABLES TO donations_maintenance;
