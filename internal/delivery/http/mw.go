package http

import (
	"context"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/auth0/go-jwt-middleware/v2/jwks"
	"github.com/auth0/go-jwt-middleware/v2/validator"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type (
	Auth0Config struct {
		// Skipper defines a function to skip middleware.
		Skipper            middleware.Skipper
		Issuer             string        `json:"issuer"`
		Audience           []string      `json:"audience"`
		Permissions        []string      `json:"permissions"`
		SignatureAlgorithm string        `json:"signature_algorithm"`
		CacheDuration      time.Duration `json:"cache_duration"`
	}

	jwtCustomClaim struct {
		Permissions []string `json:"permissions"`
	}
)

func (j *jwtCustomClaim) Validate(context.Context) error {
	return nil
}

var (
	// DefaultAuth0Config is the default Auth0 middleware config.
	DefaultAuth0Config = Auth0Config{
		Skipper:            middleware.DefaultSkipper,
		Issuer:             "",
		Audience:           []string{},
		SignatureAlgorithm: "RS256",
		CacheDuration:      5 * time.Minute,
	}
)

func Auth0() echo.MiddlewareFunc {
	return Auth0WithConfig(DefaultAuth0Config)
}

func Auth0WithConfig(config Auth0Config) echo.MiddlewareFunc {
	// Defaults
	if config.Skipper == nil {
		config.Skipper = DefaultAuth0Config.Skipper
	}

	if config.Issuer == "" {
		config.Issuer = DefaultAuth0Config.Issuer
	}

	if len(config.Audience) == 0 {
		config.Audience = DefaultAuth0Config.Audience
	}

	if config.SignatureAlgorithm == "" {
		config.SignatureAlgorithm = DefaultAuth0Config.SignatureAlgorithm
	}

	if config.CacheDuration == 0 {
		config.CacheDuration = DefaultAuth0Config.CacheDuration
	}

	issuerURL, err := url.Parse(config.Issuer)
	if err != nil {
		log.Fatalf("failed to parse the issuer url: %v", err)
	}

	provider := jwks.NewCachingProvider(issuerURL, config.CacheDuration)

	// Set up the validator.
	jwtValidator, err := validator.New(
		provider.KeyFunc,
		validator.SignatureAlgorithm(config.SignatureAlgorithm),
		issuerURL.String(),
		config.Audience,
		validator.WithCustomClaims(func() validator.CustomClaims {
			return &jwtCustomClaim{}
		}),
	)
	if err != nil {
		log.Fatalf("failed to set up the validator: %v", err)
	}

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authorization := c.Request().Header.Get("Authorization")
			if authorization == "" {
				return echo.NewHTTPError(http.StatusUnauthorized, "No Authorization Header")
			}

			// check if authorization header has bearer prefix
			if !strings.HasPrefix(authorization, "Bearer ") {
				return echo.NewHTTPError(http.StatusUnauthorized, "Invalid Authorization Header")
			}

			// get token from header
			token := strings.TrimPrefix(authorization, "Bearer ")

			// Get the JWT token from the request header.
			claims, err := jwtValidator.ValidateToken(c.Request().Context(), token)
			if err != nil {
				return echo.NewHTTPError(http.StatusUnauthorized, "Invalid Token")
			}

			userPermissions := claims.(*validator.ValidatedClaims).CustomClaims.(*jwtCustomClaim).Permissions

			//checking permissions
			if !checkPermissions(config.Permissions, userPermissions) {
				return echo.NewHTTPError(http.StatusForbidden, "Access Denied")
			}

			// Set the claims in the context.
			c.Set("claims", claims.(*validator.ValidatedClaims))

			return next(c)
		}
	}
}

func checkPermissions(requirements []string, permissions []string) bool {
	for _, req := range requirements {
		var hasPermissions = false
		for _, perm := range permissions {
			if req == perm {
				hasPermissions = true
				break
			}
		}

		if !hasPermissions {
			return false
		}
	}

	return true
}
