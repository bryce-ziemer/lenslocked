package main

import (
	"context"
	"fmt"
)

// Best practice so no one else can accidently collide (more unique, because linked to type)
// should not be exportable, local to package (lower case starting letter)
type ctxKey string

// also do not export keys (gurantee code within this package can access these keys/types)
const (
	favoriteColorKey ctxKey = "favorite-color"
)

func main() {
	ctx := context.Background()
	ctx = context.WithValue(ctx, favoriteColorKey, "blue")
	value := ctx.Value(favoriteColorKey)
	fmt.Println(value)
}
