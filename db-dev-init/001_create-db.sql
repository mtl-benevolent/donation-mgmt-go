CREATE DATABASE donationsdb;
CREATE SCHEMA donationsdb.donations;

-- Creating the roles
CREATE ROLE donations_ro;
CREATE ROLE donations_rw;
CREATE ROLE donations_maintenance;

-- Setting up donations_ro role
GRANT CONNECT ON DATABASE donationsdb TO donations_ro;
GRANT USAGE ON SCHEMA donationsdb.donations TO donations_ro;

-- Setting up donations_rw role
GRANT donations_ro TO donations_rw;

-- Setting up donations_maintenance role
GRANT donations_rw TO donations_maintenance;
GRANT CREATE ON DATABASE donationsdb TO donations_maintenance;
GRANT SYSTEM VIEWACTIVITY, VIEWCLUSTERMETADATA TO donations_maintenance;

-- Creating application user
CREATE USER donation_mgmt_app WITH LOGIN PASSWORD NULL;
GRANT donations_rw TO donation_mgmt_app;

-- Creating migration user
CREATE USER donation_mgmt_migrator WITH LOGIN PASSWORD NULL CREATEDB; -- Migrator receives the CREATEDB option to handle Prisma's ShadowDB;
GRANT donations_rw TO donation_mgmt_migrator;
GRANT CREATE ON DATABASE donationsdb TO donation_mgmt_migrator;
GRANT CREATE ON SCHEMA donationsdb.donations TO donation_mgmt_migrator;

-- Creating maintenance user
CREATE USER donation_mgmt_maintenance LOGIN PASSWORD NULL;
GRANT donations_maintenance TO donation_mgmt_maintenance;
