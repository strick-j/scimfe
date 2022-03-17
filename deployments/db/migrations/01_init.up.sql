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
