package main

import (
	"crypto/rsa"
	"fmt"
	"net/http"
	"strings"

	"github.com/antonybholmes/go-auth"
	"github.com/antonybholmes/go-edb-server-gin/consts"
	"github.com/antonybholmes/go-edb-server-gin/routes"
	authenticationroutes "github.com/antonybholmes/go-edb-server-gin/routes/authentication"
	"github.com/gin-gonic/gin"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

type APIError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Create a defer function that will be called after the handler finishes
		defer func() {
			if err := recover(); err != nil {
				// Handle panic errors with 500 status code
				c.JSON(http.StatusInternalServerError, APIError{
					Code:    http.StatusInternalServerError,
					Message: fmt.Sprintf("Internal Server Error: %v", err),
				})
			}
		}()
		// Continue processing the request
		c.Next()

		// Check if there are any errors
		if len(c.Errors) > 0 {
			// Get the last error (or you can choose how to handle multiple errors)
			err := c.Errors.Last()

			// Set a custom status code based on the error
			// If no custom status code is set, use the error's default status or fallback to 400
			statusCode := http.StatusBadRequest

			if err.Meta != nil {
				if customStatus, ok := err.Meta.(int); ok {
					statusCode = customStatus
				}
			}

			// Send the error response with custom status code
			c.JSON(statusCode, APIError{
				Code:    statusCode,
				Message: err.Error(),
			})
		}
	}
}

func parseToken(c *gin.Context, key *rsa.PublicKey) error {
	// Get the token from the "Authorization" header
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.Error(fmt.Errorf("Authorization header missing"))
		c.Abort()
		return fmt.Errorf("Authorization header missing")
	}

	// Split the token (format: "Bearer <token>")
	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	if tokenString == authHeader {
		c.Error(fmt.Errorf("Malformed token"))

		return fmt.Errorf("Malformed token")
	}

	// Parse the JWT
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Return the secret key for verifying the token
		return key, nil
	})

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token", "message": err.Error()})
		c.Error(err)
		return
	}

	c.Set("user", *token)

	return nil
}

func JWTMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		err := parseToken(c, consts.JWT_RSA_PUBLIC_KEY)

		if err != nil {
			return
		}

		// Continue processing the request
		c.Next()

	}
}

func JWTAuth0Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		err := parseToken(c, consts.JWT_AUTH0_RSA_PUBLIC_KEY)

		if err != nil {
			return
		}

		// Continue processing the request
		c.Next()

	}
}

func JwtIsRefreshTokenMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		user := c.MustGet("user").(*jwt.Token)
		claims := user.Claims.(*auth.TokenClaims)

		if claims.Type != auth.REFRESH_TOKEN {
			routes.AuthErrorReq(c, "not a refresh token")
		}

		c.Next()
	}
}

func JwtIsAccessTokenMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c *gin.Context) {
		user := c.Get("user").(*jwt.Token)
		claims := user.Claims.(*auth.TokenClaims)

		if claims.Type != auth.ACCESS_TOKEN {
			routes.AuthErrorReq("not an access token")
		}

		return next(c)
	}
}

func JwtHasAdminPermissionMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c *gin.Context) {
		user := c.Get("user").(*jwt.Token)
		claims := user.Claims.(*auth.TokenClaims)

		if !auth.IsAdmin((claims.Roles)) {
			return routes.AuthErrorReq("user is not an admin")
		}

		return next(c)
	}
}

func JwtHasLoginPermissionMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c *gin.Context) {
		user := c.Get("user").(*jwt.Token)
		claims := user.Claims.(*auth.TokenClaims)

		if !auth.CanSignin((claims.Roles)) {
			return routes.AuthErrorReq("user is not allowed to login")
		}

		return next(c)
	}
}

// basic check that session exists and seems to be populated with the user
func SessionIsValidMiddleware(next echo.HandlerFunc) echo.HandlerFunc {

	return func(c *gin.Context) {
		sessData, err := authenticationroutes.ReadSessionInfo(c)

		if err != nil {
			return routes.AuthErrorReq("cannot get user id from session")
		}

		c.Set("authUser", sessData.AuthUser)

		return next(c)
	}
}

// func ValidateJwtMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
// 	return func(c *gin.Context) {
// 		authorizationHeader := c.Request().Header.Get("authorization")

// 		if len(authorizationHeader) == 0 {
// 			return routes.AuthErrorReq("missing Authentication header")

// 		}

// 		log.Debug().Msgf("parsing authentication header")

// 		authPair := strings.SplitN(authorizationHeader, " ", 2)

// 		if len(authPair) != 2 {
// 			return routes.AuthErrorReq("wrong Authentication header definiton")
// 		}

// 		headerAuthScheme := authPair[0]
// 		headerAuthToken := authPair[1]

// 		if headerAuthScheme != "Bearer" {
// 			return routes.AuthErrorReq("wrong Authentication header definiton")
// 		}

// 		log.Debug().Msgf("validating JWT token")

// 		token, err := validateJwtToken(headerAuthToken)

// 		if err != nil {
// 			return routes.AuthErrorReq(err)
// 		}

// 		log.Debug().Msgf("JWT token is valid")
// 		c.Set("user", token)
// 		return next(c)

// 	}
// }

// Create a permissions middleware to verify jwt permissions on a token
func JwtRoleMiddleware(validRoles ...string) gin.HandlerFunc {

	return func(c *gin.Context) {

		user := c.MustGet("user").(*jwt.Token)

		// if user == nil {
		// 	return routes.AuthErrorReq("no jwt available")
		// }

		claims := user.Claims.(*auth.TokenClaims)

		// shortcut for admin, as we allow this for everything
		if !auth.IsAdmin(claims.Roles) {
			//log.Debug().Msgf("is admin")
			return
		}

		isValidRole := false

		for _, validRole := range validRoles {

			// if we find a permission, stop and move on
			if strings.Contains(claims.Roles, validRole) {
				isValidRole = true
			}

		}

		if !isValidRole {
			return
		}

		c.Next()
	}

}

func RDFMiddleware() gin.HandlerFunc {
	return JwtRoleMiddleware("RDF")
}
