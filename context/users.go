package context

import (
	"context"

	"github.com/silasburger/lenslocked/models"
)

type userKey string

const (
	uKey userKey = "user"
)

func WithUser(ctx context.Context, user *models.User) context.Context {
	return context.WithValue(ctx, uKey, user)
}

func User(ctx context.Context) *models.User {
	val := ctx.Value(uKey)
	user, ok := val.(*models.User)
	if !ok {
		// The most likely case is that nothing was ever stored in the context,
		// so it doesn't have a type of *models.User. It is also possible that
		// other code in this package wrote an invalid value using the user key,
		// so it is important to review code changes in tshis package.
		return nil
	}
	return user
}
