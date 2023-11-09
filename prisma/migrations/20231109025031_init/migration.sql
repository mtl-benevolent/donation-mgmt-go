-- Grant permissions to relevant roles on future created resources
ALTER DEFAULT PRIVILEGES FOR ROLE donations_migrators GRANT SELECT ON TABLES TO donations_ro;

ALTER DEFAULT PRIVILEGES FOR ROLE donations_migrators GRANT INSERT, UPDATE, DELETE ON TABLES TO donations_rw;
ALTER DEFAULT PRIVILEGES FOR ROLE donations_migrators GRANT USAGE ON SEQUENCES TO donations_rw;
ALTER DEFAULT PRIVILEGES FOR ROLE donations_migrators GRANT EXECUTE ON FUNCTIONS TO donations_rw;

ALTER DEFAULT PRIVILEGES FOR ROLE donations_migrators GRANT DROP ON TABLES TO donations_maintenance;
