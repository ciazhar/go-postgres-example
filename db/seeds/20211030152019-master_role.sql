-- +migrate Up
insert into role(id, name, created_at, updated_at)
values (gen_random_uuid(), 'User', now(), now());
insert into role(id, name, created_at, updated_at)
values (gen_random_uuid(), 'Admin', now(), now());
-- +migrate Down
