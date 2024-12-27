CREATE DATABASE donationsdb;
CREATE SCHEMA donations;

-- Creating the roles
CREATE ROLE donations_ro;
CREATE ROLE donations_rw;
CREATE ROLE donations_maintenance;

-- Setting up donations_ro role
GRANT CONNECT ON DATABASE donationsdb TO GROUP donations_ro;
GRANT USAGE ON SCHEMA donations TO GROUP donations_ro;

-- Setting up donations_rw role
GRANT donations_ro TO donations_rw;

-- Setting up donations_maintenance role
GRANT donations_rw TO donations_maintenance;
GRANT CREATE ON DATABASE donationsdb TO GROUP donations_maintenance;
GRANT donations_rw TO donations_maintenance;

-- Creating application user
CREATE USER donation_mgmt_app WITH LOGIN PASSWORD 'yq2REWv0iD8nepOe1BFskFwPDgn69mFbjt2q3hzmB8THCLtteXKHMws1teMKLIu7';
GRANT donations_rw TO donation_mgmt_app;

-- Creating migration user
CREATE USER donation_mgmt_migrator WITH LOGIN PASSWORD 'MvGfJjzmxzLOYQjhF9i1fvq9dQemZmGvJdFVIbAsb37nopG4gR3GE4D4nOf3xWvX' CREATEDB; -- Migrator receives the CREATEDB option to handle Prisma's ShadowDB;
GRANT donations_rw TO donation_mgmt_migrator;
GRANT CREATE ON DATABASE donationsdb TO donation_mgmt_migrator;
GRANT CREATE ON SCHEMA donations TO donation_mgmt_migrator;

-- Creating maintenance user
CREATE USER donation_mgmt_maintenance LOGIN PASSWORD 'NaG8Y9Jo1t7ziNCPNL1RpsUYJSbjPbGoU8RVjoINn1iSHR2xUamqRpfcoge428r1';
GRANT donations_maintenance TO donation_mgmt_maintenance;
