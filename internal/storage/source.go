package storage

import (
	"AI-RSS-Telegram-Bot/internal/model"
	"context"
	"database/sql"
	"time"
)

type dbSource struct {
	ID        int64     `db:"id"`
	Name      string    `db:"name"`
	FeedURL   string    `db:"feed_url"`
	Priority  int       `db:"priority"`
	CreatedAt time.Time `db:"created_at"`
}

type SourcePostgresStorage struct {
	db *sql.DB
}

func NewSourceStorage(db *sql.DB) *SourcePostgresStorage {
	return &SourcePostgresStorage{db: db}
}

func (s *SourcePostgresStorage) Sources(ctx context.Context) ([]model.Source, error) {
	rows, err := s.db.QueryContext(ctx, `SELECT * FROM sources`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sources []model.Source
	for rows.Next() {
		var source dbSource
		if err := rows.Scan(&source.ID, &source.Name, &source.FeedURL, &source.Priority, &source.CreatedAt); err != nil {
			return nil, err
		}
		sources = append(sources, model.Source(source))
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return sources, nil
}

func (s *SourcePostgresStorage) SourceByID(ctx context.Context, id int64) (*model.Source, error) {
	var source dbSource

	err := s.db.QueryRowContext(ctx, `SELECT * FROM sources WHERE id = $1`, id).Scan(&source.ID, &source.Name, &source.FeedURL, &source.Priority, &source.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return (*model.Source)(&source), nil
}

func (s *SourcePostgresStorage) Add(ctx context.Context, source model.Source) (int64, error) {
	var id int64

	err := s.db.QueryRowContext(
		ctx,
		`INSERT INTO sources (name, feed_url, priority) VALUES ($1, $2, $3) RETURNING id;`,
		source.Name, source.FeedURL, source.Priority,
	).Scan(&id)

	if err != nil {
		return 0, err
	}

	return id, nil
}

func (s *SourcePostgresStorage) SetPriority(ctx context.Context, id int64, priority int) error {
	_, err := s.db.ExecContext(ctx,
		`UPDATE sources SET priority = $1 WHERE id = $2`,
		priority, id)
	return err

}

func (s *SourcePostgresStorage) Delete(ctx context.Context, id int64) error {
	_, err := s.db.ExecContext(ctx,
		`DELETE FROM sources WHERE id = $1`, id)
	return err
}
