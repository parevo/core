package main

import (
	"context"
	"fmt"

	"github.com/parevo/core/auth"
	"github.com/parevo/core/storage/memory"
)

func main() {
	svc, err := auth.NewServiceWithModules(auth.Config{
		Issuer:    "parevo",
		Audience:  "parevo-api",
		SecretKey: []byte("change-this-in-production"),
	}, auth.Modules{
		SessionStore: &memory.SessionStore{},
		RefreshStore: &memory.RefreshStore{},
	})
	if err != nil {
		panic(err)
	}

	pair, err := svc.IssueTokenPair(context.Background(), auth.Claims{
		UserID:    "u1",
		TenantID:  "tenant-a",
		SessionID: "session-1",
	})
	if err != nil {
		panic(err)
	}
	fmt.Println("issued access:", pair.AccessToken != "")
	fmt.Println("issued refresh:", pair.RefreshToken != "")

	rotated, err := svc.RotateRefreshToken(context.Background(), pair.RefreshToken)
	if err != nil {
		panic(err)
	}
	fmt.Println("rotated access:", rotated.AccessToken != "")
	fmt.Println("rotated refresh:", rotated.RefreshToken != "")
}
