-- name: CreateArticle :one
INSERT INTO articles (name) VALUES (:name) RETURNING *;

-- name: GetArticle :one
SELECT * FROM articles WHERE id = ?;

-- name: ListArticles :many
SELECT * FROM articles ORDER BY id;

-- name: UpdateArticle :one
UPDATE articles SET name = :name WHERE id = :id RETURNING *;

-- name: DeleteArticle :exec
DELETE FROM articles WHERE id = :id;
