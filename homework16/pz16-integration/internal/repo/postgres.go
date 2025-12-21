package repo

import (
	"context"
	"database/sql"
	"errors"
	"pz16/internal/models"
)

type NoteRepo struct{ DB *sql.DB }

func (r NoteRepo) Create(ctx context.Context, n *models.Note) error {
	return r.DB.QueryRowContext(ctx,
		`INSERT INTO notes(title, content) VALUES($1,$2) RETURNING id`,
		n.Title, n.Content,
	).Scan(&n.ID)
}

func (r NoteRepo) Get(ctx context.Context, id int64) (models.Note, error) {
	var n models.Note
	err := r.DB.QueryRowContext(ctx,
		`SELECT id, title, content, created_at, updated_at FROM notes WHERE id=$1`, id,
	).Scan(&n.ID, &n.Title, &n.Content, &n.CreatedAt, &n.UpdatedAt)
	if err == sql.ErrNoRows {
		return models.Note{}, errors.New("not found")
	}
	return n, err
}
