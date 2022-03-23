-- Internal Tables -----------------------------------------------------------------------------------------------

-- Support to auto-generate UUIDs (aka GUIDs)
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Users table
--
-- Email limit is 254 chars according to RFC 5321, don't ask me who need such long mail.
-- Password is encrypted using bcrypt, which is always has 60 chars.
CREATE TABLE IF NOT EXISTS users
(
    "id" uuid PRIMARY KEY NOT NULL DEFAULT uuid_generate_v4(),
    "email" VARCHAR(254) UNIQUE NOT NULL,
    "name" VARCHAR(64) NOT NULL,
    "password" CHAR(60) NOT NULL
);

-- Actions table
--
-- Stores completed actions and their result
-- Ex: action: add, resourcetype: Group, Success: True
CREATE TABLE IF NOT EXISTS actions
(
    "id" uuid PRIMARY KEY NOT NULL DEFAULT uuid_generate_v4(),
    "action" VARCHAR(254) UNIQUE NOT NULL,
    "resourceType" VARCHAR(64) NOT NULL,
    "success" BOOL NOT NULL,
    "time" CURRENT_TIMESTAMP
);

-- Auth table
-- 
-- Used to store auth information for Oauth2 Token
CREATE TABLE IF NOT EXISTS auth
(
    "id" INT PRIMARY KEY NOT NULL GENERATED ALWAYS AS IDENTITY,
    "access_token" TEXT NOT NULL,
    "token_type" VARCHAR(50) NOT NULL,
    "expiry" TIMESTAMP NOT NULL
);

-- SCIM PAM Tables -----------------------------------------------------------------------------------------------

-- The tables below are utilized to store information retrieved from the PAM SCIM Server
-- All table fields are based on the SCIM (RFC7643) Schema, specifically those utilized by the CyberArk
-- PAM Solution.
-- SCIM 2.0 Schema: https://datatracker.ietf.org/doc/html/rfc7643
-- CyberArk PAM Schema: https://identity-developer.cyberark.com/docs/manage-pam-objects-with-scim-endpoints

-- SCIM PAM User Tables ------------------------------------------------------------------------------------------

-- Pamuser table
-- 
-- Used to store pam user information retrieved from PAM SCIM Server
CREATE TABLE IF NOT EXISTS pamuser
(
    "id" INT PRIMARY KEY UNIQUE NOT NULL,
    "username" VARCHAR(64),
    "displayname" VARCHAR(64),
    "nickname" VARCHAR(64),
    "profileUrl" VARCHAR(64),
    "title" VARCHAR(64),
    "userType" VARCHAR(64),
    "locale" VARCHAR(64),
    "timezone" VARCHAR(16),
    "active" BOOL,
    "entitlements" TEXT[],
    "schemas" TEXT[]
);

-- Pamuser_name table
-- 
-- References Pamusers
-- Used to store pam user name information for users retrieved from PAM SCIM Server
CREATE TABLE IF NOT EXISTS pamuser_name
(
    "name_id" INT PRIMARY KEY NOT NULL GENERATED ALWAYS AS IDENTITY,
    "id"  INT NOT NULL,
    "givenname" VARCHAR(64),
    "middlename" VARCHAR(64),
    "familyname" VARCHAR(64),
    "formatted" VARCHAR(64),
    "honorificPrefix" VARCHAR(8),
    "honorificSuffix" VARCHAR(8),
    CONSTRAINT fk_user
        FOREIGN KEY(id)
            REFERENCES pamuser(id)
            ON DELETE CASCADE
);

-- Pamuser_emails table
-- 
-- References Pamuser
-- Used to store pam user email information for users retrieved from PAM SCIM Server
CREATE TABLE IF NOT EXISTS pamuser_emails
(
    "email_id" INT PRIMARY KEY NOT NULL DEFAULT uuid_generate_v4(),
    "id" INT,
    "name" VARCHAR(64),
    "primary" BOOL,
    "display" VARCHAR(64),
    "value" VARCHAR(254),
    "$ref",
    CONSTRAINT fk_user
        FOREIGN KEY(id)
            REFERENCES pamuser(id)
            ON DELETE CASCADE
);

-- Pamuser_phonenumbers table
-- 
-- References Pamuser
-- Used to store pam user email information for users retrieved from PAM SCIM Server
CREATE TABLE IF NOT EXISTS pamuser_phonenumbers
(
    "phonenumber_id" INT PRIMARY KEY NOT NULL DEFAULT uuid_generate_v4(),
    "id" INT,
    "name" VARCHAR(64),
    "primary" BOOL,
    "display" VARCHAR(64),
    "value" VARCHAR(32),
    "ref" TEXT,
    CONSTRAINT fk_user
        FOREIGN KEY(id)
            REFERENCES pamuser(id)
            ON DELETE CASCADE
);

-- Pamuser_groups table
-- 
-- References Pamuser / Pamgroup
-- Used to store pam user group information for users retrieved from PAM SCIM Server
CREATE TABLE IF NOT EXISTS pamuser_groups
(
    "id" INT,
    "type" VARCHAR(64),
    "display" VARCHAR(64),
    "value" INT,
    "ref" TEXT,
    CONSTRAINT fk_user
        FOREIGN KEY(id)
            REFERENCES pamuser(id)
            ON DELETE CASCADE
    CONSTRAINT fk_group
        FOREIGN KEY(value)
            REFERENCES pamgroup(id)
            ON DELETE CASCADE
);

-- Pamuser_meta table
-- 
-- References Pamusers
-- Used to store pam user meta information for users retrieved from PAM SCIM Server
CREATE TABLE IF NOT EXISTS pamuser_meta
(
    "meta_id" INT PRIMARY KEY NOT NULL GENERATED ALWAYS AS IDENTITY,
    "user_id" INT,
    "resourceType" VARCHAR(100),
    "created" TIMESTAMP,
    "lastModified" TIMESTAMP,
    "location" VARCHAR(200),
    CONSTRAINT fk_user
        FOREIGN KEY(user_id)
            REFERENCES pamuser(user_id)
            ON DELETE CASCADE
);


-- SCIM PAM Group Tables ------------------------------------------------------------------

-- Pamgroup table
-- 
-- Used to store group information retrieved from PAM SCIM Server
CREATE TABLE IF NOT EXISTS pamgroup
(
    "id" INT PRIMARY KEY UNIQUE NOT NULL,
    "displayname" VARCHAR(100),
    "external_id" TEXT,
    "entitlements" TEXT[],
    "schemas" TEXT[]
);

-- Pamgroup_members table
-- 
-- References Pamgroup
-- Used to store group member information for users retrieved from PAM SCIM Server
CREATE TABLE IF NOT EXISTS pamgroup_members
(
    "value" INT NOT NULL,
    "display" VARCHAR(64),
    "ref" TEXT,
    CONSTRAINT fk_user
        FOREIGN KEY(value)
            REFERENCES pamuser(id)
            ON DELETE CASCADE
);

-- Pamgroup_meta table
-- 
-- References Pamgroup
-- Used to store meta information for pam groups retrieved from PAM SCIM Server
CREATE TABLE IF NOT EXISTS pamgroup_meta
(
    "id" INT NOT NULL,
    "resourceType" VARCHAR(100),
    "created" TIMESTAMP,
    "lastModified" TIMESTAMP,
    "location" VARCHAR(200),
    CONSTRAINT fk_group
        FOREIGN KEY(id)
            REFERENCES pamgroup(id)
            ON DELETE CASCADE
);



