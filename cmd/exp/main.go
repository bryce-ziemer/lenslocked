package main

import (
	"context"
	"fmt"
	"strings"
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
	value := ctx.Value(favoriteColorKey) // value is of type any
	// fmt.Println(strings.HasPrefix(value, "b")) // WOULD COMPLAIN - value is of type any, not string

	intValue, ok := value.(int)

	if !ok {
		fmt.Println("Not an int !")
	} else {
		fmt.Println(intValue)
		fmt.Println(intValue + 4)
	}

	strValue, ok := value.(string)
	// fmt.Println(strings.HasPrefix(value, "b")) // WOULD COMPLAIN - value is of type any, not string

	if !ok {
		fmt.Println("Not an int !")
	} else {
		fmt.Println(strValue)
		fmt.Println(strings.HasPrefix(strValue, "b"))
	}
}
