package permission

import "strings"

func MatchPermission(granted, requested string) bool {
	if granted == requested {
		return true
	}
	if granted == "*" || granted == "*:*" {
		return true
	}
	parts := strings.Split(granted, ":")
	reqParts := strings.Split(requested, ":")
	if len(parts) != 2 || len(reqParts) != 2 {
		return false
	}
	if parts[0] == "*" && parts[1] == "*" {
		return true
	}
	if parts[0] == reqParts[0] && parts[1] == "*" {
		return true
	}
	if parts[0] == "*" && parts[1] == reqParts[1] {
		return true
	}
	return false
}
