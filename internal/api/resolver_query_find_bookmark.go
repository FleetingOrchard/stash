package api

import (
	"context"

	"github.com/stashapp/stash/pkg/models"
)

func (r *queryResolver) AllBookmarks(ctx context.Context) (ret []*models.Bookmark, err error) {
	if err := r.withTxn(ctx, func(ctx context.Context) error {
		ret, err = r.repository.Bookmark.All(ctx)
		return err
	}); err != nil {
		return nil, err
	}
	return ret, nil
}
