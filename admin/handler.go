package admin

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"strings"
)

type adminHandler struct {
	opts Options
}

func newAdminHandler(opts Options) *adminHandler {
	return &adminHandler{opts: opts}
}

func (h *adminHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if h.opts.AuthService != nil && h.opts.RequiredPermission != "" {
		if err := h.opts.AuthService.AuthorizePermission(r.Context(), h.opts.RequiredPermission); err != nil {
			http.Error(w, "forbidden", http.StatusForbidden)
			return
		}
	}

	path := strings.TrimPrefix(r.URL.Path, h.opts.basePath())
	path = strings.Trim(path, "/")
	parts := strings.Split(path, "/")

	switch r.Method {
	case http.MethodGet:
		if len(parts) > 0 && parts[0] == "api" {
			h.serveAPI(w, r, parts[1:])
			return
		}
		// UI: /admin, /admin?section=X, /admin/tenants
		if path != "" && path != "ui" {
			q := r.URL.Query()
			q.Set("section", path)
			r.URL.RawQuery = q.Encode()
		}
		h.serveUI(w, r)
	case http.MethodPost, http.MethodPut, http.MethodDelete:
		if parts[0] == "api" {
			h.serveAPI(w, r, parts[1:])
			return
		}
		http.NotFound(w, r)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *adminHandler) serveAPI(w http.ResponseWriter, r *http.Request, parts []string) {
	w.Header().Set("Content-Type", "application/json")
	ctx := r.Context()

	if len(parts) < 1 {
		writeJSON(w, map[string]string{"error": "missing resource"}, http.StatusBadRequest)
		return
	}

	switch parts[0] {
	case "tenants":
		h.handleTenants(w, r, ctx, parts[1:])
	case "permissions":
		h.handlePermissions(w, r, ctx, parts[1:])
	case "sessions":
		h.handleSessions(w, r, ctx, parts[1:])
	default:
		writeJSON(w, map[string]string{"error": "unknown resource"}, http.StatusNotFound)
	}
}

func (h *adminHandler) handleTenants(w http.ResponseWriter, r *http.Request, ctx context.Context, parts []string) {
	if h.opts.TenantLifecycle == nil {
		writeJSON(w, map[string]string{"error": "tenant management not configured"}, http.StatusNotImplemented)
		return
	}

	switch r.Method {
	case http.MethodGet:
		if h.opts.TenantList == nil {
			writeJSON(w, map[string]string{"error": "tenant list not configured"}, http.StatusNotImplemented)
			return
		}
		list, err := h.opts.TenantList.ListTenants(ctx)
		if err != nil {
			writeJSON(w, map[string]string{"error": err.Error()}, http.StatusInternalServerError)
			return
		}
		writeJSON(w, map[string]interface{}{"tenants": list}, http.StatusOK)
	case http.MethodPost:
		if len(parts) >= 2 {
			// POST /api/tenants/:id?action=suspend|resume|delete
			tenantID := parts[1]
			action := r.URL.Query().Get("action")
			fromForm := r.Header.Get("Content-Type") == "application/x-www-form-urlencoded"
			switch action {
			case "suspend":
				_ = h.opts.TenantLifecycle.Suspend(ctx, tenantID)
				if fromForm {
					http.Redirect(w, r, h.opts.basePath()+"/tenants", http.StatusFound)
					return
				}
				writeJSON(w, map[string]string{"ok": "suspended"}, http.StatusOK)
			case "resume":
				_ = h.opts.TenantLifecycle.Resume(ctx, tenantID)
				if fromForm {
					http.Redirect(w, r, h.opts.basePath()+"/tenants", http.StatusFound)
					return
				}
				writeJSON(w, map[string]string{"ok": "resumed"}, http.StatusOK)
			case "delete":
				_ = h.opts.TenantLifecycle.Delete(ctx, tenantID)
				if fromForm {
					http.Redirect(w, r, h.opts.basePath()+"/tenants", http.StatusFound)
					return
				}
				writeJSON(w, map[string]string{"ok": "deleted"}, http.StatusOK)
			default:
				writeJSON(w, map[string]string{"error": "action required"}, http.StatusBadRequest)
			}
			return
		}
		// POST /api/tenants = create
		var req struct {
			TenantID string `json:"tenant_id"`
			Name     string `json:"name"`
			OwnerID  string `json:"owner_id"`
		}
		if r.Header.Get("Content-Type") == "application/x-www-form-urlencoded" {
			_ = r.ParseForm()
			req.TenantID = r.FormValue("tenant_id")
			req.Name = r.FormValue("name")
			req.OwnerID = r.FormValue("owner_id")
		} else if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeJSON(w, map[string]string{"error": "invalid body"}, http.StatusBadRequest)
			return
		}
		if req.TenantID == "" {
			writeJSON(w, map[string]string{"error": "tenant_id required"}, http.StatusBadRequest)
			return
		}
		if err := h.opts.TenantLifecycle.Create(ctx, req.TenantID, req.Name, req.OwnerID); err != nil {
			writeJSON(w, map[string]string{"error": err.Error()}, http.StatusInternalServerError)
			return
		}
		if r.Header.Get("Content-Type") == "application/x-www-form-urlencoded" {
			http.Redirect(w, r, h.opts.basePath()+"/tenants", http.StatusFound)
			return
		}
		writeJSON(w, map[string]string{"ok": "created"}, http.StatusOK)
	case http.MethodPut:
		if len(parts) < 2 {
			writeJSON(w, map[string]string{"error": "tenant_id required"}, http.StatusBadRequest)
			return
		}
		tenantID := parts[1]
		action := r.URL.Query().Get("action")
		switch action {
		case "suspend":
			_ = h.opts.TenantLifecycle.Suspend(ctx, tenantID)
			writeJSON(w, map[string]string{"ok": "suspended"}, http.StatusOK)
		case "resume":
			_ = h.opts.TenantLifecycle.Resume(ctx, tenantID)
			writeJSON(w, map[string]string{"ok": "resumed"}, http.StatusOK)
		default:
			writeJSON(w, map[string]string{"error": "action must be suspend or resume"}, http.StatusBadRequest)
		}
	case http.MethodDelete:
		if len(parts) < 2 {
			writeJSON(w, map[string]string{"error": "tenant_id required"}, http.StatusBadRequest)
			return
		}
		_ = h.opts.TenantLifecycle.Delete(ctx, parts[1])
		writeJSON(w, map[string]string{"ok": "deleted"}, http.StatusOK)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *adminHandler) handlePermissions(w http.ResponseWriter, r *http.Request, ctx context.Context, parts []string) {
	if h.opts.PermissionGrant == nil {
		writeJSON(w, map[string]string{"error": "permission management not configured"}, http.StatusNotImplemented)
		return
	}

	switch r.Method {
	case http.MethodGet:
		subjectID := r.URL.Query().Get("subject_id")
		tenantID := r.URL.Query().Get("tenant_id")
		if subjectID == "" || tenantID == "" {
			writeJSON(w, map[string]string{"error": "subject_id and tenant_id required"}, http.StatusBadRequest)
			return
		}
		list, err := h.opts.PermissionGrant.ListGrants(ctx, subjectID, tenantID)
		if err != nil {
			writeJSON(w, map[string]string{"error": err.Error()}, http.StatusInternalServerError)
			return
		}
		writeJSON(w, map[string]interface{}{"permissions": list}, http.StatusOK)
	case http.MethodPost:
		if r.URL.Query().Get("action") == "revoke" {
			subjectID := r.URL.Query().Get("subject_id")
			tenantID := r.URL.Query().Get("tenant_id")
			permission := r.URL.Query().Get("permission")
			if subjectID == "" || tenantID == "" || permission == "" {
				writeJSON(w, map[string]string{"error": "subject_id, tenant_id, permission required"}, http.StatusBadRequest)
				return
			}
			_ = h.opts.PermissionGrant.Revoke(ctx, subjectID, tenantID, permission)
			if r.Header.Get("Content-Type") == "application/x-www-form-urlencoded" {
				u := h.opts.basePath() + "/permissions?" + url.Values{"subject_id": {subjectID}, "tenant_id": {tenantID}}.Encode()
				http.Redirect(w, r, u, http.StatusFound)
				return
			}
			writeJSON(w, map[string]string{"ok": "revoked"}, http.StatusOK)
			return
		}
		var req struct {
			SubjectID  string `json:"subject_id"`
			TenantID   string `json:"tenant_id"`
			Permission string `json:"permission"`
		}
		if r.Header.Get("Content-Type") == "application/x-www-form-urlencoded" {
			_ = r.ParseForm()
			req.SubjectID = r.FormValue("subject_id")
			req.TenantID = r.FormValue("tenant_id")
			req.Permission = r.FormValue("permission")
		} else if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeJSON(w, map[string]string{"error": "invalid body"}, http.StatusBadRequest)
			return
		}
		if req.SubjectID == "" || req.TenantID == "" || req.Permission == "" {
			writeJSON(w, map[string]string{"error": "subject_id, tenant_id, permission required"}, http.StatusBadRequest)
			return
		}
		if err := h.opts.PermissionGrant.Grant(ctx, req.SubjectID, req.TenantID, req.Permission); err != nil {
			writeJSON(w, map[string]string{"error": err.Error()}, http.StatusInternalServerError)
			return
		}
		if r.Header.Get("Content-Type") == "application/x-www-form-urlencoded" {
			u := h.opts.basePath() + "/permissions?" + url.Values{"subject_id": {req.SubjectID}, "tenant_id": {req.TenantID}}.Encode()
			http.Redirect(w, r, u, http.StatusFound)
			return
		}
		writeJSON(w, map[string]string{"ok": "granted"}, http.StatusOK)
	case http.MethodDelete:
		subjectID := r.URL.Query().Get("subject_id")
		tenantID := r.URL.Query().Get("tenant_id")
		permission := r.URL.Query().Get("permission")
		if subjectID == "" || tenantID == "" || permission == "" {
			writeJSON(w, map[string]string{"error": "subject_id, tenant_id, permission required"}, http.StatusBadRequest)
			return
		}
		_ = h.opts.PermissionGrant.Revoke(ctx, subjectID, tenantID, permission)
		writeJSON(w, map[string]string{"ok": "revoked"}, http.StatusOK)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *adminHandler) handleSessions(w http.ResponseWriter, r *http.Request, ctx context.Context, parts []string) {
	if h.opts.SessionStore == nil {
		writeJSON(w, map[string]string{"error": "session management not configured"}, http.StatusNotImplemented)
		return
	}

	switch r.Method {
	case http.MethodGet:
		userID := r.URL.Query().Get("user_id")
		if userID == "" {
			writeJSON(w, map[string]string{"error": "user_id required"}, http.StatusBadRequest)
			return
		}
		if h.opts.SessionList == nil {
			writeJSON(w, map[string]string{"error": "session list not configured"}, http.StatusNotImplemented)
			return
		}
		list, err := h.opts.SessionList.ListSessionsByUser(ctx, userID)
		if err != nil {
			writeJSON(w, map[string]string{"error": err.Error()}, http.StatusInternalServerError)
			return
		}
		writeJSON(w, map[string]interface{}{"sessions": list}, http.StatusOK)
	case http.MethodDelete, http.MethodPost:
		fromForm := r.Header.Get("Content-Type") == "application/x-www-form-urlencoded"
		if len(parts) >= 2 && parts[1] != "" {
			_ = h.opts.SessionStore.RevokeSession(ctx, parts[1])
			if fromForm {
				userID := r.URL.Query().Get("user_id")
				if userID != "" {
					http.Redirect(w, r, h.opts.basePath()+"/sessions?user_id="+url.QueryEscape(userID), http.StatusFound)
				} else {
					http.Redirect(w, r, h.opts.basePath()+"/sessions", http.StatusFound)
				}
				return
			}
			writeJSON(w, map[string]string{"ok": "revoked"}, http.StatusOK)
			return
		}
		userID := r.URL.Query().Get("user_id")
		if userID != "" {
			if us, ok := h.opts.SessionStore.(interface {
				RevokeAllSessionsByUser(ctx context.Context, userID string) error
			}); ok {
				_ = us.RevokeAllSessionsByUser(ctx, userID)
			}
			if fromForm {
				http.Redirect(w, r, h.opts.basePath()+"/sessions?user_id="+url.QueryEscape(userID), http.StatusFound)
				return
			}
			writeJSON(w, map[string]string{"ok": "revoked_all"}, http.StatusOK)
			return
		}
		writeJSON(w, map[string]string{"error": "session_id or user_id required"}, http.StatusBadRequest)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}


func writeJSON(w http.ResponseWriter, v interface{}, status int) {
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}
