-- +migrate Up
-- +migrate StatementBegin
CREATE OR REPLACE function validate_register(
    v_id uuid default '00000000-0000-0000-0000-000000000000',
    v_phone_number varchar default '',
    v_role_id uuid default '00000000-0000-0000-0000-000000000000'
)
    RETURNS TABLE
            (
                result varchar
            )
    LANGUAGE sql
as
$function$
with cte_phone_number_unique as (
    select count(*) = 0 as correct
    from users
    where case
              when v_phone_number <> '' then phone_number = v_phone_number
                  and deleted_at is null
              else false end
      and id <> v_id
)
   , cte_role_id as (
    select case
               when v_role_id = '00000000-0000-0000-0000-000000000000' then true
               else count(*) = 1 end as correct
    from role
    where id = v_role_id
)
select case
           when cpnu.correct = false then 'phone number must unique'
           when cri.correct = false then 'role id not found'
           else 'Validated'
           end
           as result
from cte_phone_number_unique cpnu
         cross join cte_role_id cri
$function$;
-- +migrate StatementEnd

-- +migrate Down
