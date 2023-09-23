CREATE DATABASE IF NOT EXISTS donationsdb;
CREATE SCHEMA IF NOT EXISTS donationsdb.donations;

-- Creating the roles
CREATE ROLE IF NOT EXISTS donations_rw;
CREATE ROLE IF NOT EXISTS donations_migrators;
CREATE ROLE IF NOT EXISTS donations_maintenance;

-- Setting up donations_rw role
GRANT CONNECT ON DATABASE donationsdb TO donations_rw;
GRANT USAGE ON SCHEMA donationsdb.donations TO donations_rw;

-- Setting up donations_migrators role
GRANT donations_rw TO donations_migrators;
GRANT CREATE ON DATABASE donationsdb TO donations_migrators;
GRANT CREATE ON SCHEMA donationsdb.donations TO donations_migrators;

-- Setting up donations_maintenance role
GRANT donations_rw TO donations_maintenance;
GRANT CREATE ON DATABASE donationsdb TO donations_maintenance;
GRANT SYSTEM VIEWACTIVITY, VIEWCLUSTERMETADATA TO donations_maintenance;

-- Creating application user
CREATE USER IF NOT EXISTS donation_mgmt_app LOGIN PASSWORD NULL;
GRANT donations_rw TO donation_mgmt_app;

-- Creating migration user
CREATE USER IF NOT EXISTS donation_mgmt_migrator LOGIN PASSWORD NULL;
GRANT donations_migrators TO donation_mgmt_migrator;

-- Creating maintenance user
CREATE USER IF NOT EXISTS donation_mgmt_maintenance LOGIN PASSWORD NULL;
GRANT donations_maintenance TO donation_mgmt_maintenance;
