-- +migrate Up
CREATE TABLE token
(
    id         uuid      DEFAULT           uuid_generate_v4() NOT NULL,
    issuer     uuid                                           NOT NULL,
    expires_at timestamp                                      NOT NULL,
    subject    varchar                                        NOT NULL,
    issued_at  timestamp                                      NOT NULL,
    created_at timestamp DEFAULT now()                        NOT NULL,
    updated_at timestamp DEFAULT now()                        NOT NULL,
    deleted_at timestamp,
    PRIMARY KEY (id)
);

-- +migrate Down
