package api

import (
	"context"
	"strconv"

	"github.com/stashapp/stash/pkg/models"
)

func (r *mutationResolver) getBookmark(ctx context.Context, id int) (ret *models.Bookmark, err error) {
	if err := r.withTxn(ctx, func(ctx context.Context) error {
		ret, err = r.repository.Bookmark.Find(ctx, id)
		return err
	}); err != nil {
		return nil, err
	}

	return ret, nil
}

func (r *mutationResolver) BookmarkCreate(ctx context.Context, input BookmarkCreateInput) (*models.Bookmark, error) {
	newBookmark := models.Bookmark{
		URL:  input.URL,
		Name: input.Name,
	}

	var b *models.Bookmark
	if err := r.withTxn(ctx, func(ctx context.Context) error {
		qb := r.repository.Bookmark

		var err error
		pos, err := qb.GetMaxPosition(ctx)
		if err != nil {
			pos = 0
		}

		newBookmark.Position = pos + 1
		b, err = qb.Create(ctx, newBookmark)
		if err != nil {
			return err
		}
		return nil
	}); err != nil {
		return nil, err
	}

	return r.getBookmark(ctx, b.ID)
}

func (r *mutationResolver) BookmarkUpdate(ctx context.Context, input BookmarkUpdateInput) (*models.Bookmark, error) {
	bookmarkID, err := strconv.Atoi(input.ID)
	if err != nil {
		return nil, err
	}

	newBookmark := models.Bookmark{
		ID:       bookmarkID,
		URL:      input.URL,
		Name:     input.Name,
		Position: input.Position,
	}

	var b *models.Bookmark
	if err := r.withTxn(ctx, func(ctx context.Context) error {
		var err error
		qb := r.repository.Bookmark

		b, err = qb.UpdateFull(ctx, newBookmark)
		if err != nil {
			return err
		}
		return nil
	}); err != nil {
		return nil, err
	}

	return r.getBookmark(ctx, b.ID)
}

func (r *mutationResolver) BookmarkDestroy(ctx context.Context, input BookmarkDestroyInput) (bool, error) {
	bookmarkID, err := strconv.Atoi(input.ID)
	if err != nil {
		return false, err
	}

	if err := r.withTxn(ctx, func(ctx context.Context) error {
		qb := r.repository.Bookmark
		err := qb.Destroy(ctx, bookmarkID)
		if err != nil {
			return err
		}
		return nil
	}); err != nil {
		return false, err
	}

	return true, nil
}
