package middleware

import (
	"encoding/json"
	"net/http"
	"strings"
	"tofash/internal/config"
	"tofash/internal/modules/user/entity"
	"tofash/internal/modules/user/service"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

type AuthMiddleware struct {
	cfg         *config.Config
	jwtService  service.JwtServiceInterface
	redisClient *config.Config // Wait, I need redis client logic?
}

// Simplified generic response, usually use shared package or module response
type DefaultResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func NewAuthMiddleware(cfg *config.Config, jwtService service.JwtServiceInterface) *AuthMiddleware {
	return &AuthMiddleware{
		cfg:        cfg,
		jwtService: jwtService,
	}
}

func (m *AuthMiddleware) CheckToken(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		respErr := DefaultResponse{Code: http.StatusUnauthorized}

		// Create Redis client on the fly? Ideally inject it.
		// For now using Config to create it as in original adapter
		// Create Redis client on the fly? Ideally inject it.
		// For now using Config to create it as in original adapter
		redisConn := m.cfg.NewRedisClient()

		authHeader := c.Request().Header.Get("Authorization")
		tokenString := ""

		if authHeader != "" {
			tokenString = strings.TrimPrefix(authHeader, "Bearer ")
		} else {
			// Fallback: Check Cookie
			cookie, err := c.Cookie("access_token")
			if err == nil {
				tokenString = cookie.Value
			}
		}

		if tokenString == "" {
			log.Errorf("[Middleware] CheckToken: %s", "missing or invalid token")
			respErr.Message = "missing or invalid token"
			return c.JSON(http.StatusUnauthorized, respErr)
		}

		_, err := m.jwtService.ValidateToken(tokenString)
		if err != nil {
			log.Errorf("[Middleware] CheckToken: %s", err.Error())
			respErr.Message = err.Error()
			return c.JSON(http.StatusUnauthorized, respErr)
		}

		var getSession string
		if redisConn == nil {
			log.Errorf("[Middleware] CheckToken: %s", "redis connection failed")
			respErr.Message = "redis connection failed"
			return c.JSON(http.StatusInternalServerError, respErr)
		}
		var errSes error
		getSession, errSes = redisConn.Get(c.Request().Context(), tokenString).Result()
		if errSes != nil || len(getSession) == 0 {
			log.Errorf("[Middleware] CheckToken: %s", "session not found/expired")
			respErr.Message = "session not found or expired"
			return c.JSON(http.StatusUnauthorized, respErr)
		}
		jwtUserData := entity.JwtUserData{}
		err = json.Unmarshal([]byte(getSession), &jwtUserData)
		if err != nil {
			log.Errorf("[Middleware] CheckToken: %v", err)
			respErr.Message = err.Error()
			respErr.Code = http.StatusInternalServerError
			return c.JSON(http.StatusInternalServerError, respErr)
		}

		// Authorization check (example)
		path := c.Request().URL.Path
		segments := strings.Split(strings.Trim(path, "/"), "/")
		if len(segments) > 0 && jwtUserData.RoleName == "Customer" && segments[0] == "admin" {
			log.Infof("[Middleware] CheckToken: %s", "customer cannot access admin routes")
			respErr.Message = "customer cannot access admin routes"
			respErr.Code = http.StatusForbidden
			return c.JSON(http.StatusForbidden, respErr)
		}

		c.Set("user", getSession)
		return next(c)
	}
}
