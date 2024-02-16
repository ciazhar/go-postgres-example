-- name: FetchRole :many
select id, name
from role
where deleted_at is null;
