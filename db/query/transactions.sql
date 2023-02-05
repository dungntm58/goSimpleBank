-- name: CreateTransaction :one
INSERT INTO transactions (
  from_account_id,
  to_account_id,
  amount
) VALUES (
  $1, $2, $3
) RETURNING *;

-- name: GetTransaction :one
SELECT * FROM transactions
WHERE id = $1 LIMIT 1;

-- name: ListTransactions :many
SELECT * FROM transactions
ORDER BY id
LIMIT $1
OFFSET $2;

-- name: DeleteTransaction :exec
DELETE FROM transactions
WHERE id = $1;