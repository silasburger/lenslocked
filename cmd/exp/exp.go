package main

import (
	stdctx "context"
	"fmt"

	"github.com/silasburger/lenslocked/context"
	"github.com/silasburger/lenslocked/models"
)

func main() {
	ctx := stdctx.Background()

	user := models.User{
		Email: "jon@calhoun.io",
	}

	ctx = context.WithUser(ctx, &user)
	newUser := context.User(ctx)
	fmt.Println(newUser.Email)
}
