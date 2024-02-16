-- name: ValidateRegister :one
select result::varchar
from validate_register(@id::uuid, @phone_number::varchar, @role_id::uuid);

-- name: Register :one
insert into users(name, phone_number, role_id, password, fcm_token, created_at, updated_at)
values (@name::varchar, @phone_number::varchar,
        (select id from role where role.name = 'User'), @password::varchar,
        @fcm_token::varchar, now(),
        now()) returning id::uuid, name::varchar, phone_number::varchar, role_id::uuid, password::varchar,created_at::date, updated_at::date;

-- name: FetchUser :many
select u.id,
       coalesce(u.name, '')       as name,
       coalesce(u.username, '')   as username,
       coalesce(u.email, '')      as email,
       json_build_object(
               'id', r.id,
               'name', coalesce(r.name, ''),
               'created_at', r.created_at,
               'updated_at', r.updated_at
           )                      as role
from users u
         join role r on r.id = u.role_id and r.deleted_at is null
where u.deleted_at is null limit $1
offset $2;

-- name: CountUser :one
select count(*) as count
from users u
    join role r
on r.id = u.role_id and r.deleted_at is null
where u.deleted_at is null;

-- name: FetchUserByPhone :one
select u.id,
       u.name,
       r.id                    as role_id,
       r.name                  as role_name,
       phone_number,
       password,
       coalesce(fcm_token, '') as fcm_token
from users u
         join role r on u.role_id = r.id
where phone_number = @phone_number;

-- name: UpdateUser :exec
update users
set name         = (case when @name:: varchar = '' then name else @name:: varchar end),
    phone_number = (case when @phone_number:: varchar = '' then phone_number else @phone_number:: varchar end),
    password     = (case when @password:: varchar = '' then password else @password:: varchar end),
    fcm_token    = (case when @fcm_token:: varchar = '' then fcm_token else @fcm_token:: varchar end),
    role_id      = (case
                        when @role_id::uuid = '00000000-0000-0000-0000-000000000000' then role_id
                        else @role_id::uuid end),
    updated_at   = now()
where id = @id::uuid;

-- name: GetUserById :one
select u.id, u.name, r.name as role_name, phone_number, password
from users u
         join role r on r.id = u.role_id and r.deleted_at is null
where u.deleted_at is null
  and u.id = @id::uuid;

-- name: ForgotPassword :exec
update users
set password = @password::varchar
where phone_number = @phone_number:: varchar
  and @password = @re_password:: varchar;

-- name: FindByPhone :one
select id
from users
where phone_number = @phone_number;