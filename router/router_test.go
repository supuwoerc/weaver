package router

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"github.com/supuwoerc/weaver/conf"
)

type RouterSuite struct {
	suite.Suite
	engine      *gin.Engine
	config      *conf.Config
	routerGroup *gin.RouterGroup
}

func TestRouterSuite(t *testing.T) {
	suite.Run(t, new(RouterSuite))
}

func (s *RouterSuite) SetupSuite() {
	gin.SetMode(gin.TestMode)
	path, err := projectRootPath()
	require.NoError(s.T(), err)
	config := &conf.Config{
		System: conf.SystemConfig{
			DefaultLang:      "cn",
			DefaultLocaleKey: "Locale",
			LocalePath: conf.LocalePath{
				Zh: filepath.Join(path, "pkg", "locale", "zh"),
				En: filepath.Join(path, "pkg", "locale", "en"),
			},
		},
	}
	s.config = config
}

func (s *RouterSuite) TearDownSuite() {
	s.engine = nil
	s.config = nil
	s.routerGroup = nil
}

func (s *RouterSuite) SetupSubTest() {
	engin := gin.New()
	routerGroup := NewRouter(engin, s.config)
	s.engine = engin
	s.routerGroup = routerGroup
}

func (s *RouterSuite) TearDownSubTest() {
	s.engine = nil
	s.routerGroup = nil
}

func projectRootPath() (string, error) {
	currentDir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	for {
		if _, err = os.Stat(filepath.Join(currentDir, "go.mod")); err == nil {
			return currentDir, nil
		}
		parent := filepath.Dir(currentDir)
		if parent == currentDir {
			return "", errors.New("fail to find work dir")
		}
		currentDir = parent
	}
}

func (s *RouterSuite) TestNewRouter_Basic() {
	s.Run("basic new router", func() {
		gin.SetMode(gin.TestMode)
		t := s.T()
		assert.NotNil(t, s.routerGroup)
		assert.Equal(t, "/api/v1", s.routerGroup.BasePath())
	})
}

func (s *RouterSuite) TestNewRouter_WithRoutes() {
	t := s.T()
	s.Run("NewRouter with routes - GET request", func() {
		s.routerGroup.GET("/test", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "get success"})
		})
		w := httptest.NewRecorder()
		req, reqErr := http.NewRequest(http.MethodGet, "/api/v1/test", nil)
		require.NoError(t, reqErr)
		s.engine.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "get success")
	})

	s.Run("NewRouter with routes - POST request", func() {
		s.routerGroup.POST("/test", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "post success"})
		})
		w := httptest.NewRecorder()
		req, reqErr := http.NewRequest(http.MethodPost, "/api/v1/test", nil)
		require.NoError(t, reqErr)
		s.engine.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "post success")
	})
}

func (s *RouterSuite) TestNewRouter_InvalidPath() {
	t := s.T()
	s.Run("invalid route path", func() {
		w := httptest.NewRecorder()
		req, err := http.NewRequest("GET", "/api/v1/invalid", nil)
		require.NoError(t, err)
		s.engine.ServeHTTP(w, req)
		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}

func (s *RouterSuite) TestNewRouter_DifferentHTTPMethods() {
	s.Run("different http methods", func() {
		s.routerGroup.GET("/methods", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"method": "GET"})
		})
		s.routerGroup.POST("/methods", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"method": "POST"})
		})
		s.routerGroup.PUT("/methods", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"method": "PUT"})
		})
		s.routerGroup.DELETE("/methods", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"method": "DELETE"})
		})
		methods := []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete}
		t := s.T()
		for _, method := range methods {
			w := httptest.NewRecorder()
			req, err := http.NewRequest(method, "/api/v1/methods", nil)
			require.NoError(t, err)
			s.engine.ServeHTTP(w, req)
			assert.Equal(t, http.StatusOK, w.Code)
			assert.Contains(t, w.Body.String(), method)
		}
	})

}

func (s *RouterSuite) TestNewRouter_SubGroups() {
	t := s.T()
	s.Run("user sub group", func() {
		userGroup := s.routerGroup.Group("/users")
		userGroup.GET("/profile", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "user profile"})
		})
		w := httptest.NewRecorder()
		req, err := http.NewRequest("GET", "/api/v1/users/profile", nil)
		require.NoError(t, err)
		s.engine.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "user profile")
	})

	s.Run("admin sub group", func() {
		adminGroup := s.routerGroup.Group("/admin")
		adminGroup.GET("/dashboard", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "admin dashboard"})
		})
		w := httptest.NewRecorder()
		req, err := http.NewRequest("GET", "/api/v1/admin/dashboard", nil)
		require.NoError(t, err)
		s.engine.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "admin dashboard")
	})
}

func (s *RouterSuite) TestNewRouter_ConfigurationValidation() {
	tests := []struct {
		name        string
		config      *conf.Config
		shouldPanic bool
		panicVal    any
	}{
		{
			name:        "valid config",
			config:      s.config,
			shouldPanic: false,
		},
		{
			name:        "nil config",
			config:      nil,
			shouldPanic: true,
		},
		{
			name: "empty DefaultLocaleKey",
			config: &conf.Config{
				System: conf.SystemConfig{
					DefaultLang:      "cn",
					DefaultLocaleKey: "",
					LocalePath:       s.config.System.LocalePath,
				},
			},
			shouldPanic: true,
			panicVal:    "miss locale key",
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			engine := gin.New()
			t := s.T()
			if tt.shouldPanic {
				if tt.panicVal == nil {
					assert.Panics(t, func() {
						NewRouter(engine, tt.config)
					})
				} else {
					assert.PanicsWithValue(t, tt.panicVal, func() {
						NewRouter(engine, tt.config)
					})
				}
			} else {
				assert.NotPanics(t, func() {
					routerGroup := NewRouter(engine, tt.config)
					assert.NotNil(t, routerGroup)
				})
			}
		})
	}
}

func (s *RouterSuite) TestNewRouter_WithI18nMiddleware() {
	t := s.T()
	s.Run("with i18n middleware", func() {
		s.routerGroup.GET("/i18n-test", func(c *gin.Context) {
			translator, exists := c.Get("i18n_translator")
			assert.True(t, exists)
			assert.NotNil(t, translator)
			validatorTranslator, exists := c.Get("validator_translator")
			assert.True(t, exists)
			assert.NotNil(t, validatorTranslator)
			c.JSON(http.StatusOK, gin.H{"message": "i18n middleware working"})
		})
		w := httptest.NewRecorder()
		req, err := http.NewRequest("GET", "/api/v1/i18n-test", nil)
		require.NoError(t, err)
		req.Header.Set("Locale", "cn")
		s.engine.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "i18n middleware working")
	})
}

func (s *RouterSuite) TestNewRouter_WithoutLocaleHeader() {
	t := s.T()
	s.Run("without locale header", func() {
		s.routerGroup.GET("/default-locale", func(c *gin.Context) {
			translator, exists := c.Get("i18n_translator")
			assert.True(t, exists)
			assert.NotNil(t, translator)
			c.JSON(http.StatusOK, gin.H{"message": "default locale working"})
		})
		w := httptest.NewRecorder()
		req, err := http.NewRequest("GET", "/api/v1/default-locale", nil)
		require.NoError(t, err)
		s.engine.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "default locale working")
	})
}

func (s *RouterSuite) TestNewRouter_EnglishLocale() {
	t := s.T()
	s.Run("en locale", func() {
		s.routerGroup.GET("/english-locale", func(c *gin.Context) {
			translator, exists := c.Get("i18n_translator")
			assert.True(t, exists)
			assert.NotNil(t, translator)
			c.JSON(http.StatusOK, gin.H{"message": "english locale working"})
		})
		w := httptest.NewRecorder()
		req, err := http.NewRequest("GET", "/api/v1/english-locale", nil)
		require.NoError(t, err)
		req.Header.Set("Locale", "en")
		s.engine.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "english locale working")
	})
}

func (s *RouterSuite) TestNewRouter_InvalidLocale() {
	t := s.T()
	s.Run("invalid locale", func() {
		s.routerGroup.GET("/invalid-locale", func(c *gin.Context) {
			translator, exists := c.Get("i18n_translator")
			assert.True(t, exists)
			assert.NotNil(t, translator)
			c.JSON(http.StatusOK, gin.H{"message": "invalid locale handled"})
		})
		w := httptest.NewRecorder()
		req, err := http.NewRequest("GET", "/api/v1/invalid-locale", nil)
		require.NoError(t, err)
		req.Header.Set("Locale", "invalid")
		s.engine.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "invalid locale handled")
	})
}
