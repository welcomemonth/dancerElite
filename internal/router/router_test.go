package router

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestSetupAdminUI(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name       string
		path       string
		accept     string
		wantStatus int
		wantBody   string
	}{
		{
			name:       "index",
			path:       "/web/",
			wantStatus: http.StatusOK,
			wantBody:   `<div id="app"></div>`,
		},
		{
			name:       "spa fallback",
			path:       "/web/admin/articles",
			accept:     "text/html",
			wantStatus: http.StatusOK,
			wantBody:   `<div id="app"></div>`,
		},
		{
			name:       "missing asset",
			path:       "/web/assets/not-found.js",
			accept:     "text/html",
			wantStatus: http.StatusNotFound,
			wantBody:   `"message":"Not found"`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			engine := gin.New()
			setupAdminUI(engine)

			req := httptest.NewRequest(http.MethodGet, tt.path, nil)
			if tt.accept != "" {
				req.Header.Set("Accept", tt.accept)
			}
			recorder := httptest.NewRecorder()

			engine.ServeHTTP(recorder, req)

			if recorder.Code != tt.wantStatus {
				t.Fatalf("status = %d, want %d", recorder.Code, tt.wantStatus)
			}
			if !strings.Contains(recorder.Body.String(), tt.wantBody) {
				t.Fatalf("body does not contain %q: %s", tt.wantBody, recorder.Body.String())
			}
		})
	}
}
