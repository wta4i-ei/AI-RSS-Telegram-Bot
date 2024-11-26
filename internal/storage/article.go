package storage

import (
	"AI-RSS-Telegram-Bot/internal/model"
	"context"
	"database/sql"
	"time"
)

type ArticlePostgresStorage struct {
	db *sql.DB
}

type dbArticleWithPriority struct {
	ID             int64          `db:"a_id"`
	SourcePriority int64          `db:"s_priority"`
	SourceID       int64          `db:"s_id"`
	Title          string         `db:"a_title"`
	Link           string         `db:"a_link"`
	Summary        sql.NullString `db:"a_summary"`
	PublishedAt    time.Time      `db:"a_published_at"`
	PostedAt       sql.NullTime   `db:"a_posted_at"`
	CreatedAt      time.Time      `db:"a_created_at"`
}

func NewArticleStorage(db *sql.DB) *ArticlePostgresStorage {
	return &ArticlePostgresStorage{db: db}
}

func (s *ArticlePostgresStorage) Store(ctx context.Context, article model.Article) error {
	_, err := s.db.ExecContext(
		ctx,
		`INSERT INTO articles (source_id, title, link, summary, published_at)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT DO NOTHING;`,
		article.SourceID,
		article.Title,
		article.Link,
		article.Summary,
		article.PublishedAt,
	)
	return err
}

func (s *ArticlePostgresStorage) AllNotPosted(ctx context.Context, since time.Time, limit uint64) ([]model.Article, error) {
	var articles []dbArticleWithPriority

	rows, err := s.db.QueryContext(
		ctx,
		`SELECT 
			a.id AS a_id, 
			s.priority AS s_priority,
			s.id AS s_id,
			a.title AS a_title,
			a.link AS a_link,
			a.summary AS a_summary,
			a.published_at AS a_published_at,
			a.posted_at AS a_posted_at,
			a.created_at AS a_created_at
		FROM articles a JOIN sources s ON s.id = a.source_id
		WHERE a.posted_at IS NULL 
			AND a.published_at >= $1::timestamp
		ORDER BY a.created_at DESC, s_priority DESC LIMIT $2;`,
		since.UTC().Format(time.RFC3339),
		limit,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var article dbArticleWithPriority
		if err := rows.Scan(&article.ID, &article.SourcePriority,
			&article.SourceID, &article.Title, &article.Link,
			&article.Summary, &article.PublishedAt,
			&article.PostedAt, &article.CreatedAt); err != nil {
			return nil, err
		}
		articles = append(articles, article)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	result := make([]model.Article, len(articles))
	for i, article := range articles {
		result[i] = model.Article{
			ID:          article.ID,
			SourceID:    article.SourceID,
			Title:       article.Title,
			Link:        article.Link,
			Summary:     article.Summary.String,
			PublishedAt: article.PublishedAt,
			CreatedAt:   article.CreatedAt,
		}
	}
	return result, nil
}

func (s *ArticlePostgresStorage) MarkAsPosted(ctx context.Context, article model.Article) error {
	_, err := s.db.ExecContext(
		ctx,
		`UPDATE articles SET posted_at = $1::timestamp WHERE id = $2;`,
		time.Now().UTC(),
		article.ID,
	)
	return err
}
