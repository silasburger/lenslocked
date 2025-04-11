package context

import (
	"context"

	"github.com/silasburger/lenslocked/models"
)

type galleryKey string

const (
	gKey galleryKey = "gallery"
)

func WithGallery(ctx context.Context, gallery *models.Gallery) context.Context {
	return context.WithValue(ctx, gKey, gallery)
}

func Gallery(ctx context.Context) *models.Gallery {
	val := ctx.Value(gKey)
	gallery, ok := val.(*models.Gallery)
	if !ok {
		// The most likely case is that nothing was ever stored in the context,
		// so it doesn't have a type of *models.User. It is also possible that
		// other code in this package wrote an invalid value using the gallery key,
		// so it is important to review code changes in this package.
		return nil
	}
	return gallery
}
