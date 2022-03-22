-- Support to auto-generate UUIDs (aka GUIDs)
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Users table
--
-- Email limit is 254 chars according to RFC 5321, don't ask me who need such long mail.
-- Password is encrypted using bcrypt, which is always has 60 chars.
CREATE TABLE "users"
(
    "id"       uuid PRIMARY KEY    NOT NULL DEFAULT uuid_generate_v4(),
    "email"    VARCHAR(254) UNIQUE NOT NULL,
    "name"     VARCHAR(64)         NOT NULL,
    "password" CHAR(60)            NOT NULL
);

-- Pamusers table
-- 
-- Used to store information retrieved from PAM SCIM Server
CREATE TABLE "pamusers"
(
    "id"            INT PRIMARY KEY NOT NULL GENERATED ALWAYS AS IDENTITY,
    "username"      varchar(100),
    "displayname"   varchar(100),
    "usertype"      varchar(50),
    "active"        bool,
    "user_id"       INT             NOT NULL,
    "entitlements"  text[],
    "schemas"       text[]
);

-- Name table
-- 
-- References Pamusers
-- Used to store name information for users retrieved from PAM SCIM Server
CREATE TABLE "name"
(
    "name_id"       INT PRIMARY KEY NOT NULL GENERATED ALWAYS AS IDENTITY,
    "user_id"       INT             NOT NULL,
    "givenname"     varchar(100),
    "familyname"    varchar(100),
    CONSTRAINT fk_user
        FOREIGN KEY(user_id)
            REFERENCES pamusers(user_id)
            ON DELETE CASCADE
);

-- Meta table
-- 
-- References Pamusers
-- Used to store meta information for users retrieved from PAM SCIM Server
CREATE TABLE "meta"
(
    "meta_id"       INT PRIMARY KEY NOT NULL GENERATED ALWAYS AS IDENTITY,
    "user_id"       INT,
    "resourceType"  varchar(100),
    "created"       timestamp,
    "lastModified"  timestamp,
    "location"      varchar(200),
    PRIMARY KEY(meta_id),
    CONSTRAINT fk_user
        FOREIGN KEY(user_id)
            REFERENCES pamusers(user_id)
            ON DELETE CASCADE
);

-- Auth table
-- 
-- Used to store auth information for Oauth2 Token
CREATE TABLE "auth"
(
    "id"            INT PRIMARY KEY NOT NULL GENERATED ALWAYS AS IDENTITY,
    "access_token"  text            NOT NULL,
    "token_type"    varchar(50)     NOT NULL,
    "expiry"        timestamp       NOT NULL
);
