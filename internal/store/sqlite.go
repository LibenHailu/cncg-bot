package store

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	_ "modernc.org/sqlite"
)

type Store struct{ DB *sql.DB }

type Item struct {
	ID           int64
	Source       string
	Title        string
	URL          string
	Summary      string
	PublishedAt  time.Time
	Tags         string // comma-separated
	Hash         string // sha256(url+title)
	Score        float64
	Posted       bool
}

func Open(path string) (*Store, error) {
	db, err := sql.Open("sqlite", path)
	if err != nil { return nil, err }
	s := &Store{DB: db}
	if err := s.migrate(); err != nil { return nil, err }
	return s, nil
}

func (s *Store) migrate() error {
	_, err := s.DB.Exec(`
CREATE TABLE IF NOT EXISTS items (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	source TEXT NOT NULL,
	title TEXT NOT NULL,
	url TEXT NOT NULL,
	summary TEXT NOT NULL,
	published_at TIMESTAMP NOT NULL,
	tags TEXT,
	hash TEXT NOT NULL UNIQUE,
	score REAL NOT NULL DEFAULT 0,
	posted INTEGER NOT NULL DEFAULT 0
);
CREATE TABLE IF NOT EXISTS errors (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	when_ts TIMESTAMP NOT NULL,
	component TEXT NOT NULL,
	message TEXT NOT NULL
);
CREATE INDEX IF NOT EXISTS idx_items_posted ON items(posted);
`)
	return err
}

func Hash(url, title string) string {
	h := sha256.Sum256([]byte(url + "|" + title))
	return hex.EncodeToString(h[:])
}

func (s *Store) InsertIfNew(ctx context.Context, it Item) (bool, error) {
	fmt.Println("Inserting item:", it.Score, it.Posted)
	_, err := s.DB.ExecContext(ctx, `
INSERT INTO items (source,title,url,summary,published_at,tags,hash,score,posted)
VALUES (?,?,?,?,?,?,?,?,0)
ON CONFLICT(hash) DO NOTHING
`, it.Source, it.Title, it.URL, it.Summary, it.PublishedAt, it.Tags, it.Hash, it.Score)
	if err != nil { return false, err }
	// Check if it exists now
	var cnt int
	if err := s.DB.QueryRowContext(ctx, `SELECT COUNT(1) FROM items WHERE hash=?`, it.Hash).Scan(&cnt); err != nil { return false, err }
	return cnt == 1, nil
}

func (s *Store) NextUnposted(ctx context.Context, minScore float64, limit int) ([]Item, error) {
	fmt.Println("Fetching next unposted items with min score:", minScore, "limit:", limit)
	rows, err := s.DB.QueryContext(ctx, `
SELECT id,source,title,url,summary,published_at,tags,hash,score,posted
FROM items
WHERE posted=0 AND score >= ?
ORDER BY score DESC, published_at DESC
LIMIT ?`, minScore, limit)
	if err != nil { return nil, err }
	defer rows.Close()
	var out []Item
	for rows.Next() {
		var it Item
		if err := rows.Scan(&it.ID,&it.Source,&it.Title,&it.URL,&it.Summary,&it.PublishedAt,&it.Tags,&it.Hash,&it.Score,&it.Posted); err != nil { return nil, err }
		out = append(out, it)
	}
	fmt.Println(out)
	return out, rows.Err()
}

func (s *Store) MarkPosted(ctx context.Context, id int64) error {
	_, err := s.DB.ExecContext(ctx, `UPDATE items SET posted=1 WHERE id=?`, id)
	return err
}

func (s *Store) LogError(ctx context.Context, component, msg string) {
	_, _ = s.DB.ExecContext(ctx, `INSERT INTO errors(when_ts,component,message) VALUES (?,?,?)`,
		time.Now().UTC(), component, msg)
}

var ErrNotFound = errors.New("not found")
