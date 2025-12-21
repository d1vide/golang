package repo

import (
	"context"
	"errors"
	"fmt"
	"time"

	"example.com/notes-api/internal/core"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	defaultPageSize = 20
)

type NoteDBRepo struct {
	db *pgxpool.Pool
}

func NewNoteDBRepo(pool *pgxpool.Pool) *NoteDBRepo {
	return &NoteDBRepo{db: pool}
}

type PaginationParams struct {
	Limit             int
	LastSeenID        int64
	LastSeenCreatedAt time.Time
}

func (r *NoteDBRepo) Create(n core.Note) (*core.Note, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
		INSERT INTO notes (title, content)
		VALUES ($1, $2)
		RETURNING id, created_at
	`

	var id int64
	var createdAt time.Time
	err := r.db.QueryRow(ctx, query, n.Title, n.Content).Scan(&id, &createdAt)
	if err != nil {
		return nil, fmt.Errorf("failed to create note: %w", err)
	}

	n.ID = id
	n.CreatedAt = createdAt
	return &n, nil
}

func (r *NoteDBRepo) GetAll() ([]core.Note, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
		SELECT id, title, content, created_at
		FROM notes
		ORDER BY created_at DESC, id DESC
	`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get all notes: %w", err)
	}
	defer rows.Close()

	return r.scanNotes(rows)
}

func (r *NoteDBRepo) GetPaginated(limit int, lastSeenID int64, lastSeenCreatedAt time.Time) ([]core.Note, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if limit <= 0 {
		limit = defaultPageSize
	}

	var rows pgx.Rows
	var err error

	if lastSeenID == 0 && lastSeenCreatedAt.IsZero() {
		query := `
			SELECT id, title, content, created_at
			FROM notes
			ORDER BY created_at DESC, id DESC
			LIMIT $1
		`
		rows, err = r.db.Query(ctx, query, limit)
	} else {
		query := `
			SELECT id, title, content, created_at
			FROM notes
			WHERE (created_at, id) < ($2, $1)
			ORDER BY created_at DESC, id DESC
			LIMIT $3
		`
		rows, err = r.db.Query(ctx, query, lastSeenID, lastSeenCreatedAt, limit)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to get paginated notes: %w", err)
	}
	defer rows.Close()

	return r.scanNotes(rows)
}

func (r *NoteDBRepo) GetByIDs(ids []int64) (map[int64]*core.Note, error) {
	if len(ids) == 0 {
		return make(map[int64]*core.Note), nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
		SELECT id, title, content, created_at
		FROM notes
		WHERE id = ANY($1)
		ORDER BY created_at DESC, id DESC
	`

	rows, err := r.db.Query(ctx, query, ids)
	if err != nil {
		return nil, fmt.Errorf("failed to get notes by ids: %w", err)
	}
	defer rows.Close()

	notes := make(map[int64]*core.Note)
	for rows.Next() {
		var note core.Note
		err := rows.Scan(&note.ID, &note.Title, &note.Content, &note.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan note: %w", err)
		}
		notes[note.ID] = &note
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	return notes, nil
}

func (r *NoteDBRepo) GetByID(id int64) (*core.Note, error) {
	notes, err := r.GetByIDs([]int64{id})
	if err != nil {
		return nil, err
	}

	note, exists := notes[id]
	if !exists {
		return nil, ErrNotFound
	}

	return note, nil
}

func (r *NoteDBRepo) Update(id int64, title, content string) (*core.Note, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
		UPDATE notes
		SET title = $2, content = $3
		WHERE id = $1
		RETURNING id, title, content, created_at
	`

	var note core.Note
	err := r.db.QueryRow(ctx, query, id, title, content).Scan(
		&note.ID,
		&note.Title,
		&note.Content,
		&note.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("failed to update note: %w", err)
	}

	return &note, nil
}

func (r *NoteDBRepo) Delete(id int64) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
		DELETE FROM notes
		WHERE id = $1
	`

	result, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete note: %w", err)
	}

	if result.RowsAffected() == 0 {
		return ErrNotFound
	}

	return nil
}

func (r *NoteDBRepo) SearchByTitle(searchTerm string, limit int) ([]core.Note, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if limit <= 0 {
		limit = defaultPageSize
	}

	query := `
		SELECT id, title, content, created_at
		FROM notes
		WHERE to_tsvector('simple', title) @@ plainto_tsquery('simple', $1)
		ORDER BY created_at DESC, id DESC
		LIMIT $2
	`

	rows, err := r.db.Query(ctx, query, searchTerm, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to search notes: %w", err)
	}
	defer rows.Close()

	return r.scanNotes(rows)
}

func (r *NoteDBRepo) Close() {
	r.db.Close()
}

func (r *NoteDBRepo) scanNotes(rows pgx.Rows) ([]core.Note, error) {
	notes := make([]core.Note, 0)
	for rows.Next() {
		var note core.Note
		err := rows.Scan(&note.ID, &note.Title, &note.Content, &note.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan note: %w", err)
		}
		notes = append(notes, note)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	return notes, nil
}

func (r *NoteDBRepo) GetTotalCount() (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var count int64
	err := r.db.QueryRow(ctx, "SELECT COUNT(*) FROM notes").Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to get total count: %w", err)
	}

	return count, nil
}

func (r *NoteDBRepo) GetPaginatedOffset(limit int, offset int64) ([]core.Note, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if limit <= 0 {
		limit = defaultPageSize
	}

	query := `
		SELECT id, title, content, created_at
		FROM notes
		ORDER BY created_at DESC, id DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := r.db.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get paginated notes with offset: %w", err)
	}
	defer rows.Close()

	return r.scanNotes(rows)
}
