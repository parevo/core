package auth

type TenantOverridePolicy interface {
	CanOverride(claims *Claims, requestedTenantID string) bool
}

type StaticTenantOverridePolicy struct {
	Allow bool
}

func (p StaticTenantOverridePolicy) CanOverride(_ *Claims, _ string) bool {
	return p.Allow
}

type RoleBasedTenantOverridePolicy struct {
	AllowedRoles map[string]struct{}
}

func (p RoleBasedTenantOverridePolicy) CanOverride(claims *Claims, _ string) bool {
	if claims == nil {
		return false
	}
	for _, role := range claims.Roles {
		if _, ok := p.AllowedRoles[role]; ok {
			return true
		}
	}
	return false
}
