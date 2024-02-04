// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0
// source: queries.sql

package db

import (
	"context"
)

const getPageBySlug = `-- name: GetPageBySlug :one
SELECT page_id, title, slug, content, created_at, updated_at FROM pages
WHERE slug = $1 LIMIT 1
`

func (q *Queries) GetPageBySlug(ctx context.Context, slug string) (Page, error) {
	row := q.db.QueryRow(ctx, getPageBySlug, slug)
	var i Page
	err := row.Scan(
		&i.PageID,
		&i.Title,
		&i.Slug,
		&i.Content,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const savePage = `-- name: SavePage :exec
INSERT INTO pages (slug, title, content)
VALUES ($1, $2, $3)
ON CONFLICT (slug)
DO UPDATE SET title = $2, content = $3
RETURNING page_id, title, slug, content, created_at, updated_at
`

type SavePageParams struct {
	Slug    string `json:"slug"`
	Title   string `json:"title"`
	Content string `json:"content"`
}

func (q *Queries) SavePage(ctx context.Context, arg SavePageParams) error {
	_, err := q.db.Exec(ctx, savePage, arg.Slug, arg.Title, arg.Content)
	return err
}