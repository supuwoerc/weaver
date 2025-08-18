package router

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type RouterSuite struct {
	suite.Suite
	engine      *gin.Engine
	routerGroup *gin.RouterGroup
}

func TestRouterSuite(t *testing.T) {
	suite.Run(t, new(RouterSuite))
}

func (s *RouterSuite) SetupSuite() {
	gin.SetMode(gin.TestMode)
}

func (s *RouterSuite) TearDownSuite() {
	s.engine = nil
	s.routerGroup = nil
}

func (s *RouterSuite) SetupSubTest() {
	engin := gin.New()
	routerGroup := NewRouter(engin)
	s.engine = engin
	s.routerGroup = routerGroup
}

func (s *RouterSuite) TearDownSubTest() {
	s.engine = nil
	s.routerGroup = nil
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
