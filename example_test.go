package urlshort

import (
	"fmt"
	"github.com/GearFramework/urlshort/internal/app"
	"github.com/GearFramework/urlshort/internal/pkg/auth"
)

// ExampleGetID example of generated unique user ID
func ExampleGetID() {
	gen := app.UserGenID{}
	userID := gen.GetID()
	fmt.Println(userID)
}

// ExampleBuildJWT example create token by user ID
func ExampleBuildJWT() {
	gen := app.UserGenID{}
	userID := gen.GetID()
	token, err := auth.BuildJWT(userID)
	if err != nil {
		fmt.Printf("jwt error: %s\n", err)
	}
	fmt.Println(token)
}

// ExampleGetUserIDFromJWT example get user from token
func ExampleGetUserIDFromJWT() {
	gen := app.UserGenID{}
	userID := gen.GetID()
	token, err := auth.BuildJWT(userID)
	if err != nil {
		fmt.Printf("jwt error: %s\n", err)
	}
	fmt.Println(token)
	userID = auth.GetUserIDFromJWT(token)
	fmt.Println(userID)
}
