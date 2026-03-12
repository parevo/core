package admin

import (
	"embed"
	"html/template"
	"net/http"

	"github.com/parevo/core/auth"
	"github.com/parevo/core/storage"
)

//go:embed templates/*.html
var templateFS embed.FS

var sectionConfig = map[string]struct {
	title  string
	page   string
	loader func(*adminHandler, *pageData, *http.Request)
}{
	"users":       {"Users", "templates/users.html", loadUsers},
	"tenants":     {"Tenants", "templates/tenants.html", loadTenants},
	"permissions": {"Permissions", "templates/permissions.html", loadPermissions},
	"sessions":    {"Sessions", "templates/sessions.html", loadSessions},
	"dashboard":   {"Dashboard", "templates/dashboard.html", nil},
}

type pageData struct {
	Title        string
	Section      string
	BasePath     string
	UserID       string
	TenantID     string
	LogoutPath   string
	Tenants      []storage.TenantInfo
	Users        []storage.UserInfo
	Permissions  []string
	Sessions     []string
	SubjectID    string
	TenantIDForm string
	UserIDForm   string
}

func (h *adminHandler) serveUI(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	section := r.URL.Query().Get("section")
	if section == "" {
		section = "dashboard"
	}
	cfg, ok := sectionConfig[section]
	if !ok {
		cfg = sectionConfig["dashboard"]
		section = "dashboard"
	}

	claims, _ := auth.ClaimsFromContext(r.Context())
	logoutPath := h.opts.LogoutPath
	if logoutPath == "" {
		logoutPath = "/logout"
	}

	data := pageData{
		Title:      cfg.title,
		Section:    section,
		BasePath:   h.opts.basePath(),
		LogoutPath: logoutPath,
	}
	if claims != nil {
		data.UserID = claims.UserID
		data.TenantID = claims.TenantID
	}
	if cfg.loader != nil {
		cfg.loader(h, &data, r)
	}

	tmpl, err := template.ParseFS(templateFS, "templates/base.html", cfg.page)
	if err != nil {
		http.Error(w, "template error", http.StatusInternalServerError)
		return
	}
	tmplName := "mainFragment"
	if r.Header.Get("HX-Request") != "true" {
		tmplName = "base"
	}
	_ = tmpl.ExecuteTemplate(w, tmplName, data)
}

func loadUsers(h *adminHandler, d *pageData, r *http.Request) {
	if h.opts.UserList != nil {
		d.Users, _ = h.opts.UserList.ListUsers(r.Context())
	}
}

func loadTenants(h *adminHandler, d *pageData, r *http.Request) {
	if h.opts.TenantList != nil {
		d.Tenants, _ = h.opts.TenantList.ListTenants(r.Context())
	}
}

func loadPermissions(h *adminHandler, d *pageData, r *http.Request) {
	d.SubjectID = r.URL.Query().Get("subject_id")
	d.TenantIDForm = r.URL.Query().Get("tenant_id")
	if d.SubjectID == "" {
		d.SubjectID = "user-1"
	}
	if d.TenantIDForm == "" {
		d.TenantIDForm = "tenant-a"
	}
	if h.opts.PermissionGrant != nil && d.SubjectID != "" && d.TenantIDForm != "" {
		d.Permissions, _ = h.opts.PermissionGrant.ListGrants(r.Context(), d.SubjectID, d.TenantIDForm)
	}
}

func loadSessions(h *adminHandler, d *pageData, r *http.Request) {
	d.UserIDForm = r.URL.Query().Get("user_id")
	if d.UserIDForm == "" {
		d.UserIDForm = "user-1"
	}
	if h.opts.SessionList != nil && d.UserIDForm != "" {
		d.Sessions, _ = h.opts.SessionList.ListSessionsByUser(r.Context(), d.UserIDForm)
	}
}
