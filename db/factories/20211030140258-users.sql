-- +migrate Up
CREATE TABLE users
(
    id           uuid      DEFAULT uuid_generate_v4() NOT NULL,
    name         varchar                                        NOT NULL,
    username     varchar,
    email        varchar,
    phone_number varchar                                        NOT NULL,
    password     varchar                                        NOT NULL,
    fcm_token    varchar,
    role_id      uuid                                           NOT NULL,
    created_at   timestamp DEFAULT now()                        NOT NULL,
    updated_at   timestamp DEFAULT now()                        NOT NULL,
    deleted_at   timestamp,
    PRIMARY KEY (id)
);
CREATE INDEX users_role_id
    ON users (role_id);
ALTER TABLE users
    ADD CONSTRAINT FKusers973872 FOREIGN KEY (role_id) REFERENCES role (id);

-- +migrate Down
ALTER TABLE users
    DROP CONSTRAINT FKusers973872;
DROP TABLE IF EXISTS users CASCADE;
