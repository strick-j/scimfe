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

-- Pamusers table
-- 
-- Used to store information retrieved from PAM SCIM Server
CREATE TABLE IF NOT EXISTS pamusers
(
    "id" INT PRIMARY KEY NOT NULL GENERATED ALWAYS AS IDENTITY,
    "username" VARCHAR(100),
    "displayname" VARCHAR(100),
    "usertype" VARCHAR(50),
    "active" BOOL,
    "user_id" INT UNIQUE NOT NULL,
    "entitlements" TEXT[],
    "schemas" TEXT[]
);

-- Name table
-- 
-- References Pamusers
-- Used to store name information for users retrieved from PAM SCIM Server
CREATE TABLE IF NOT EXISTS name
(
    "name_id" INT PRIMARY KEY NOT NULL GENERATED ALWAYS AS IDENTITY,
    "user_id"  INT NOT NULL,
    "givenname" VARCHAR(64),
    "familyname" VARCHAR(64),
    CONSTRAINT fk_user
        FOREIGN KEY(user_id)
            REFERENCES pamusers(user_id)
            ON DELETE CASCADE
);

-- Meta table
-- 
-- References Pamusers
-- Used to store meta information for users retrieved from PAM SCIM Server
CREATE TABLE IF NOT EXISTS meta
(
    "meta_id" INT PRIMARY KEY NOT NULL GENERATED ALWAYS AS IDENTITY,
    "user_id" INT,
    "resourceType" VARCHAR(100),
    "created" TIMESTAMP,
    "lastModified" TIMESTAMP,
    "location" VARCHAR(200),
    PRIMARY KEY(meta_id),
    CONSTRAINT fk_user
        FOREIGN KEY(user_id)
            REFERENCES pamusers(user_id)
            ON DELETE CASCADE
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
