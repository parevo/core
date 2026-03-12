package main

import (
	"context"
	"fmt"
	"os"

	"github.com/parevo/core/auth"
	"github.com/parevo/core/social"
	googleprovider "github.com/parevo/core/social/providers/google"
	"github.com/parevo/core/storage/memory"
)

func main() {
	clientID := os.Getenv("GOOGLE_CLIENT_ID")
	clientSecret := os.Getenv("GOOGLE_CLIENT_SECRET")
	redirectURL := os.Getenv("GOOGLE_REDIRECT_URL")
	code := os.Getenv("GOOGLE_AUTH_CODE")
	if clientID == "" || clientSecret == "" || redirectURL == "" || code == "" {
		fmt.Println("set GOOGLE_CLIENT_ID, GOOGLE_CLIENT_SECRET, GOOGLE_REDIRECT_URL, GOOGLE_AUTH_CODE")
		return
	}

	authSvc, err := auth.NewService(auth.Config{
		Issuer:    "parevo",
		Audience:  "parevo-api",
		SecretKey: []byte("change-this-in-production"),
	})
	if err != nil {
		panic(err)
	}

	google := googleprovider.New(clientID, clientSecret, redirectURL)
	socialSvc := social.NewService(&memory.SocialAccountStore{}, auth.AccessTokenIssuer{Service: authSvc}, google)

	result, err := socialSvc.HandleCallback(context.Background(), "google", code, redirectURL, "tenant-a")
	if err != nil {
		panic(err)
	}
	fmt.Println("user:", result.UserID)
	fmt.Println("access token issued:", result.AccessToken != "")
}
