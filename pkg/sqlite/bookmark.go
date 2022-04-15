package sqlite

import (
	"context"
	"database/sql"
	"errors"

	"github.com/stashapp/stash/pkg/models"
)

const bookmarkTable = "bookmarks"
const bookmarkPositionColumn = "position"

type bookmarkQueryBuilder struct {
	repository
}

var BookmarkReaderWriter = &bookmarkQueryBuilder{
	repository{
		tableName: bookmarkTable,
		idColumn:  idColumn,
	},
}

func (qb *bookmarkQueryBuilder) Create(ctx context.Context, newObject models.Bookmark) (*models.Bookmark, error) {
	var ret models.Bookmark
	if err := qb.insertObject(ctx, newObject, &ret); err != nil {
		return nil, err
	}

	return &ret, nil
}

func (qb *bookmarkQueryBuilder) UpdateFull(ctx context.Context, updatedObject models.Bookmark) (*models.Bookmark, error) {
	/* TODO: Work out if position updates work without this or not
	currentObject, err := qb.Find(updatedObject.ID)

	if err != nil {
		return nil, err
	}

	if currentObject.Position != updatedObject.Position {
		if currentObject.Position < updatedObject.Position {
		} else {
		}
	}
	*/

	const partial = false
	if err := qb.update(ctx, updatedObject.ID, updatedObject, partial); err != nil {
		return nil, err
	}

	return qb.Find(ctx, updatedObject.ID)
}

func (qb *bookmarkQueryBuilder) Find(ctx context.Context, id int) (*models.Bookmark, error) {
	var ret models.Bookmark
	if err := qb.getByID(ctx, id, &ret); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &ret, nil
}

func (qb *bookmarkQueryBuilder) Destroy(ctx context.Context, id int) error {
	return qb.destroyExisting(ctx, []int{id})
}

func (qb *bookmarkQueryBuilder) Count(ctx context.Context) (int, error) {
	return qb.runCountQuery(ctx, qb.buildCountQuery("SELECT bookmarks.id FROM bookmarks"), nil)
}

func (qb *bookmarkQueryBuilder) All(ctx context.Context) ([]*models.Bookmark, error) {
	return qb.queryBookmarks(ctx, selectAll(bookmarkTable)+qb.getDefaultBookmarkSort(), nil)
}

func (qb *bookmarkQueryBuilder) GetMaxPosition(ctx context.Context) (int, error) {
	return qb.runMaxQuery(ctx, "SELECT MAX(bookmarks.position) as max FROM bookmarks", nil)
}

func (qb *bookmarkQueryBuilder) queryBookmarks(ctx context.Context, query string, args []interface{}) ([]*models.Bookmark, error) {
	var ret models.Bookmarks
	if err := qb.query(ctx, query, args, &ret); err != nil {
		return nil, err
	}

	return []*models.Bookmark(ret), nil
}

func (qb *bookmarkQueryBuilder) getDefaultBookmarkSort() string {
	return getSort(bookmarkPositionColumn, "ASC", bookmarkTable)
}
