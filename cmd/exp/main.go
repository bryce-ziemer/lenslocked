package main

import (
	stdctx "context"
	"fmt"

	"bryce-ziemer/github.com/lenslocked/context"
	"bryce-ziemer/github.com/lenslocked/models"
)

// Best practice so no one else can accidently collide (more unique, because linked to type)
// should not be exportable, local to package (lower case starting letter)
type ctxKey string

// also do not export keys (gurantee code within this package can access these keys/types)
const (
	favoriteColorKey ctxKey = "favorite-color"
)

func main() {
	ctx := stdctx.Background()

	user := models.User{
		Email: "me@email.com",
	}

	ctx = context.WithUser(ctx, &user)

	retrivedUser := context.User(ctx)
	fmt.Println(retrivedUser.Email)
}
