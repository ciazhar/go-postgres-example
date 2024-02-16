-- +migrate Up
CREATE EXTENSION if not exists "uuid-ossp" with schema extension;
CREATE EXTENSION if not exists pgcrypto WITH SCHEMA extension;
CREATE EXTENSION if not exists pg_trgm WITH SCHEMA extension;
CREATE EXTENSION if not exists postgis WITH SCHEMA extension;
-- +migrate Down

