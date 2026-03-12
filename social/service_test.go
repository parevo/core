package social

import (
	"context"
	"testing"

	mockprovider "github.com/parevo/core/social/providers/mock"
	"github.com/parevo/core/storage/memory"
)

type testIssuer struct{}

func (testIssuer) IssueAccessToken(_ context.Context, userID string, tenantID string) (string, error) {
	return "token-for-" + userID + "-" + tenantID, nil
}

func TestHandleCallback(t *testing.T) {
	store := &memory.SocialAccountStore{}
	svc := NewService(store, testIssuer{}, mockprovider.Provider{ProviderName: "google"})

	result, err := svc.HandleCallback(context.Background(), "google", "u1", "http://localhost/callback", "tenant-a")
	if err != nil {
		t.Fatalf("callback failed: %v", err)
	}
	if result.UserID == "" || result.AccessToken == "" {
		t.Fatalf("invalid result: %+v", result)
	}
}
