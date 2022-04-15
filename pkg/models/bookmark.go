package models

import "context"

type BookmarkReader interface {
	Find(ctx context.Context, id int) (*Bookmark, error)
	All(ctx context.Context) ([]*Bookmark, error)
	GetMaxPosition(ctx context.Context) (int, error)
}

type BookmarkWriter interface {
	Create(ctx context.Context, newBookmark Bookmark) (*Bookmark, error)
	UpdateFull(ctx context.Context, updateBookmark Bookmark) (*Bookmark, error)
	Destroy(ctx context.Context, id int) error
}

type BookmarkReaderWriter interface {
	BookmarkReader
	BookmarkWriter
}
