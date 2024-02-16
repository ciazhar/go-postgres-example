-- name: IsTokenExceedLimitByIssuer :one
select count(*) > 3 as is_limit
from token
where issuer = @issuer::uuid
  and deleted_at is null;

-- name: DeleteTokenByIssuer :exec
update token
set deleted_at = now()
where issuer = @issuer::uuid
  and deleted_at is null;

-- name: CreateToken :exec
insert into token (issuer, expires_at, subject, issued_at)
values (@issuer::uuid, @expires_at::timestamp, @subject::varchar, @issued_at::timestamp);
