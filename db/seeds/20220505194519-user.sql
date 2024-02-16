
-- +migrate Up
insert into users(name, phone_number, password, fcm_token, role_id)
select name,
       '1',
       '$2a$14$azQyCPnm/fU7rEx/ss3M9OuLLvJjSnnpUG.B2kREAU8hOuEQnAV62',
       '',
       id
from role
where name = 'User';
insert into users(name, phone_number, password, fcm_token, role_id)
select name,
       '2',
       '$2a$14$azQyCPnm/fU7rEx/ss3M9OuLLvJjSnnpUG.B2kREAU8hOuEQnAV62',
       '',
       id
from role
where name = 'Admin';

-- +migrate Down
