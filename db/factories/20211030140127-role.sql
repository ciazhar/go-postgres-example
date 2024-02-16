-- +migrate Up
CREATE TABLE role
(
    id         uuid                    NOT NULL,
    name       varchar                 NOT NULL,
    created_at timestamp DEFAULT now() NOT NULL,
    updated_at timestamp DEFAULT now() NOT NULL,
    deleted_at timestamp,
    CONSTRAINT role_pkey
        PRIMARY KEY (id)
);

-- +migrate Down
DROP TABLE IF EXISTS role CASCADE;
