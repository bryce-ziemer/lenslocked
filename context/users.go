package context

import (
	"bryce-ziemer/github.com/lenslocked/models"
	"context"
)

type key string

const (
	userKey = "key"
)

func WithUser(ctx context.Context, user *models.User) context.Context {
	return context.WithValue(ctx, userKey, user)
}

func User(ctx context.Context) *models.User {
	val := ctx.Value(userKey)

	user, ok := val.(*models.User)
	if !ok {
		// The most likely case is that nothign was ever stored in the context,
		// so it does not have a type of *models.User. It is also possibel that
		// other code in this package wrote an invalid value using the user key.
		return nil // nil gets converted into a *models.User
	}

	return user
}
