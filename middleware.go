package main

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/antonybholmes/go-auth"
	"github.com/antonybholmes/go-edb-server-gin/consts"
	"github.com/antonybholmes/go-edb-server-gin/routes"
	authenticationroutes "github.com/antonybholmes/go-edb-server-gin/routes/authentication"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"

	"github.com/golang-jwt/jwt/v5"
)

type APIError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func ErrorHandlerMiddleware() gin.HandlerFunc {
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
				// ok indicates cast worked
				customStatus, ok := err.Meta.(int)

				if ok {
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

func parseToken(c *gin.Context) (string, error) {
	// Get the token from the "Authorization" header
	authHeader := c.GetHeader("Authorization")

	if authHeader == "" {
		return "", fmt.Errorf("authorization header missing")
	}

	// Split the token (format: "Bearer <token>")
	tokenString := strings.TrimPrefix(authHeader, "Bearer ")

	if tokenString == authHeader {
		return "", fmt.Errorf("malformed token")
	}

	return tokenString, nil
}

func JwtMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		tokenString, err := parseToken(c)

		if err != nil {
			c.Error(err)
			c.Abort()
			return
		}

		claims := auth.TokenClaims{}

		// Parse the JWT
		_, err = jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (interface{}, error) {
			// Return the secret key for verifying the token
			return consts.JWT_RSA_PUBLIC_KEY, nil
		})

		if err != nil {
			c.Error(err)
			c.Abort()
			return
		}

		// use pointer to token
		c.Set("user", &claims)

		// Continue processing the request
		c.Next()

	}
}

func JwtAuth0Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		tokenString, err := parseToken(c)

		if err != nil {
			c.Error(err)
			c.Abort()
			return
		}

		claims := auth.Auth0TokenClaims{}

		log.Debug().Msgf("tok %s", tokenString)

		// Parse the JWT
		_, err = jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (interface{}, error) {
			// Return the secret key for verifying the token
			return consts.JWT_AUTH0_RSA_PUBLIC_KEY, nil
		})

		if err != nil {
			c.Error(err)
			c.Abort()
			return
		}

		// use pointer to token
		c.Set("user", &claims)

		// Continue processing the request
		c.Next()
	}
}

func JwtIsRefreshTokenMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		user, ok := c.Get("user")

		if !ok {
			routes.AuthErrorResp(c, "no user")

			return
		}

		claims := user.(*auth.TokenClaims)

		if claims.Type != auth.REFRESH_TOKEN {
			routes.AuthErrorResp(c, "not a refresh token")

			return
		}

		c.Next()
	}
}

func JwtIsAccessTokenMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		user, ok := c.Get("user")

		if !ok {
			routes.UserDoesNotExistResp(c)

			return
		}

		claims := user.(*auth.TokenClaims)

		if claims.Type != auth.ACCESS_TOKEN {
			routes.AuthErrorResp(c, "not an access token")

			return
		}

		c.Next()
	}
}

func JwtHasAdminPermissionMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		user, ok := c.Get("user")

		if !ok {
			routes.UserDoesNotExistResp(c)
			return
		}

		claims := user.(*auth.TokenClaims)

		if !auth.IsAdmin((claims.Roles)) {
			routes.AuthErrorResp(c, "user is not an admin")

			return
		}

		c.Next()
	}
}

func JwtHasLoginPermissionMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		user, ok := c.Get("user")

		if !ok {
			routes.UserDoesNotExistResp(c)

			return
		}

		claims := user.(*auth.TokenClaims)

		if !auth.CanSignin((claims.Roles)) {
			routes.AuthErrorResp(c, "user is not allowed to login")
			return
		}

		c.Next()
	}
}

// basic check that session exists and seems to be populated with the user
func SessionIsValidMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		sessData, err := authenticationroutes.ReadSessionInfo(c)

		if err != nil {
			routes.AuthErrorResp(c, "cannot get user id from session")

			return
		}

		c.Set("authUser", sessData.AuthUser)

		c.Next()
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

		user, ok := c.Get("user")

		if !ok {
			routes.AuthErrorResp(c, "no user claims")

			return
		}

		// if user == nil {
		// 	return routes.AuthErrorReq("no jwt available")
		// }

		claims := user.(*auth.TokenClaims)

		// shortcut for admin, as we allow this for everything
		if !auth.IsAdmin(claims.Roles) {
			routes.NotAdminResp(c)

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
			routes.ErrorResp(c, "invalid role")
			return
		}

		c.Next()
	}

}

func RDFMiddleware() gin.HandlerFunc {
	return JwtRoleMiddleware("RDF")
}
