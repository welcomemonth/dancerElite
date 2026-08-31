package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/welcomemonth/dancer-elite/internal/config"
	"github.com/welcomemonth/dancer-elite/internal/model"
	"github.com/welcomemonth/dancer-elite/internal/store"
)

func TestRBACAuthAllowsMenuReadonlyAccess(t *testing.T) {
	gin.SetMode(gin.TestMode)

	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	if err := db.AutoMigrate(&model.Menu{}, &model.RoleMenu{}); err != nil {
		t.Fatal(err)
	}

	menu := model.Menu{Path: "/admin/articles", Type: 2, Status: 1}
	if err := db.Create(&menu).Error; err != nil {
		t.Fatal(err)
	}
	if err := db.Create(&model.RoleMenu{RoleID: 2, MenuID: menu.ID}).Error; err != nil {
		t.Fatal(err)
	}

	st := store.NewWithDB(db, &config.Config{})

	engine := gin.New()
	engine.Use(func(c *gin.Context) {
		c.Set("role", "editor")
		c.Set("role_id", float64(2))
		c.Next()
	})
	engine.Use(RBACAuth(st))
	engine.GET("/api/admin/articles/", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})
	engine.POST("/api/admin/articles/", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	getRecorder := httptest.NewRecorder()
	engine.ServeHTTP(getRecorder, httptest.NewRequest(http.MethodGet, "/api/admin/articles/", nil))
	if getRecorder.Code != http.StatusOK {
		t.Fatalf("GET status = %d, want %d", getRecorder.Code, http.StatusOK)
	}

	postRecorder := httptest.NewRecorder()
	engine.ServeHTTP(postRecorder, httptest.NewRequest(http.MethodPost, "/api/admin/articles/", nil))
	if postRecorder.Code != http.StatusForbidden {
		t.Fatalf("POST status = %d, want %d", postRecorder.Code, http.StatusForbidden)
	}
}
