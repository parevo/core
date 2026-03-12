package main

import (
	"context"
	"fmt"

	"github.com/parevo/core/auth"
	"github.com/parevo/core/social"
	mockprovider "github.com/parevo/core/social/providers/mock"
	"github.com/parevo/core/storage/memory"
)

func main() {
	authSvc, err := auth.NewService(auth.Config{
		Issuer:    "parevo",
		Audience:  "parevo-api",
		SecretKey: []byte("change-this-in-production"),
	})
	if err != nil {
		panic(err)
	}

	socialSvc := social.NewService(
		&memory.SocialAccountStore{},
		auth.AccessTokenIssuer{Service: authSvc},
		mockprovider.Provider{ProviderName: "google"},
	)

	result, err := socialSvc.HandleCallback(context.Background(), "google", "user1", "http://localhost/callback", "tenant-a")
	if err != nil {
		panic(err)
	}
	fmt.Println("user:", result.UserID)
	fmt.Println("token:", result.AccessToken)
}
