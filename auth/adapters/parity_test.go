package adapters_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/gofiber/fiber/v2"
	"github.com/labstack/echo/v4"

	"github.com/parevo/core/auth"
	"github.com/parevo/core/auth/adapters"
	echoadapter "github.com/parevo/core/auth/adapters/echo"
	fiberadapter "github.com/parevo/core/auth/adapters/fiber"
	ginadapter "github.com/parevo/core/auth/adapters/gin"
	nethttpadapter "github.com/parevo/core/auth/adapters/nethttp"
)

func testService(t *testing.T) *auth.Service {
	t.Helper()
	svc, err := auth.NewService(auth.Config{
		Issuer:    "parevo",
		Audience:  "parevo-api",
		SecretKey: []byte("super-secret-key"),
	})
	if err != nil {
		t.Fatalf("new service: %v", err)
	}
	return svc
}

func tokenFor(t *testing.T, svc *auth.Service) string {
	t.Helper()
	token, err := svc.IssueAccessToken(auth.Claims{
		UserID:      "u1",
		TenantID:    "t1",
		Permissions: []string{"users:read"},
	})
	if err != nil {
		t.Fatalf("issue token: %v", err)
	}
	return token
}

func TestAdaptersParityUnauthorized(t *testing.T) {
	svc := testService(t)
	opts := adapters.Options{}

	tests := []struct {
		name string
		run  func(t *testing.T, svc *auth.Service, opts adapters.Options) int
	}{
		{name: "nethttp", run: runNetHTTP},
		{name: "gin", run: runGin},
		{name: "echo", run: runEcho},
		{name: "fiber", run: runFiber},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			status := tc.run(t, svc, opts)
			if status != http.StatusUnauthorized {
				t.Fatalf("expected 401, got %d", status)
			}
		})
	}
}

func TestAdaptersParityAuthorized(t *testing.T) {
	svc := testService(t)
	token := tokenFor(t, svc)
	opts := adapters.Options{}

	tests := []struct {
		name string
		run  func(t *testing.T, svc *auth.Service, opts adapters.Options, token string) int
	}{
		{name: "nethttp", run: runNetHTTPWithToken},
		{name: "gin", run: runGinWithToken},
		{name: "echo", run: runEchoWithToken},
		{name: "fiber", run: runFiberWithToken},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			status := tc.run(t, svc, opts, token)
			if status != http.StatusOK {
				t.Fatalf("expected 200, got %d", status)
			}
		})
	}
}

func runNetHTTP(t *testing.T, svc *auth.Service, opts adapters.Options) int {
	t.Helper()
	handler := nethttpadapter.AuthMiddleware(svc, opts)(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	req := httptest.NewRequest(http.MethodGet, "/secure", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	return rec.Code
}

func runNetHTTPWithToken(t *testing.T, svc *auth.Service, opts adapters.Options, token string) int {
	t.Helper()
	handler := nethttpadapter.AuthMiddleware(svc, opts)(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	req := httptest.NewRequest(http.MethodGet, "/secure", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	return rec.Code
}

func runGin(t *testing.T, svc *auth.Service, opts adapters.Options) int {
	t.Helper()
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.GET("/secure", ginadapter.AuthMiddleware(svc, opts), func(c *gin.Context) {
		c.Status(http.StatusOK)
	})
	req := httptest.NewRequest(http.MethodGet, "/secure", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	return rec.Code
}

func runGinWithToken(t *testing.T, svc *auth.Service, opts adapters.Options, token string) int {
	t.Helper()
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.GET("/secure", ginadapter.AuthMiddleware(svc, opts), func(c *gin.Context) {
		c.Status(http.StatusOK)
	})
	req := httptest.NewRequest(http.MethodGet, "/secure", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	return rec.Code
}

func runEcho(t *testing.T, svc *auth.Service, opts adapters.Options) int {
	t.Helper()
	e := echo.New()
	e.GET("/secure", func(c echo.Context) error {
		return c.NoContent(http.StatusOK)
	}, echoadapter.AuthMiddleware(svc, opts))

	req := httptest.NewRequest(http.MethodGet, "/secure", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	return rec.Code
}

func runEchoWithToken(t *testing.T, svc *auth.Service, opts adapters.Options, token string) int {
	t.Helper()
	e := echo.New()
	e.GET("/secure", func(c echo.Context) error {
		return c.NoContent(http.StatusOK)
	}, echoadapter.AuthMiddleware(svc, opts))

	req := httptest.NewRequest(http.MethodGet, "/secure", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	return rec.Code
}

func runFiber(t *testing.T, svc *auth.Service, opts adapters.Options) int {
	t.Helper()
	app := fiber.New()
	app.Get("/secure", fiberadapter.AuthMiddleware(svc, opts), func(c *fiber.Ctx) error {
		return c.SendStatus(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/secure", nil)
	resp, err := app.Test(req, -1)
	if err != nil {
		t.Fatalf("fiber test request failed: %v", err)
	}
	defer resp.Body.Close()
	io.Copy(io.Discard, resp.Body)
	return resp.StatusCode
}

func runFiberWithToken(t *testing.T, svc *auth.Service, opts adapters.Options, token string) int {
	t.Helper()
	app := fiber.New()
	app.Get("/secure", fiberadapter.AuthMiddleware(svc, opts), func(c *fiber.Ctx) error {
		return c.SendStatus(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/secure", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	resp, err := app.Test(req, -1)
	if err != nil {
		t.Fatalf("fiber test request failed: %v", err)
	}
	defer resp.Body.Close()
	io.Copy(io.Discard, resp.Body)
	return resp.StatusCode
}
