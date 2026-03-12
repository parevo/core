package auth

import "context"

type AccessTokenIssuer struct {
	Service *Service
}

func (i AccessTokenIssuer) IssueAccessToken(_ context.Context, userID string, tenantID string) (string, error) {
	return i.Service.IssueAccessToken(Claims{
		UserID:   userID,
		TenantID: tenantID,
	})
}
