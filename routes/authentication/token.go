package authentication

import (
	"errors"

	"github.com/antonybholmes/go-web"
	"github.com/antonybholmes/go-web/auth"
	"github.com/antonybholmes/go-web/auth/token/tokengen"
	"github.com/antonybholmes/go-web/middleware"
	"github.com/gin-gonic/gin"
)

// func RenewTokenRoute(c *gin.Context) {
// 	user := c.Get(middleware.SESSION_USER).(*jwt.Token)
// 	claims := user.Claims.(*auth.JwtCustomClaims)

// 	// Throws unauthorized error
// 	//if username != "edb" || password != "tod4EwVHEyCRK8encuLE" {
// 	//	return echo.ErrUnauthorized
// 	//}

// 	// Set custom claims
// 	renewClaims := auth.JwtCustomClaims{
// 		UserId: claims.UserId,
// 		//Email: authUser.Email,
// 		IpAddr: claims.IpAddr,
// 		RegisteredClaims: jwt.RegisteredClaims{
// 			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * auth.JWT_TOKEN_EXPIRES_HOURS)),
// 		},
// 	}

// 	// Create token with claims
// 	token := jwt.NewWithClaims(jwt.SigningMethodHS256, renewClaims)

// 	// Generate encoded token and send it as response.
// 	t, err := token.SignedString([]byte(consts.JWT_SECRET))

// 	if err != nil {
// 		return web.ErrorReq("error signing token")
// 	}

// 	return MakeDataResp(c, "", &JwtResp{t})
// }

var (
	ErrCreatingToken = errors.New("error creating token")
)

func TokenInfoRoute(c *gin.Context) {

	// user is a jwt
	user, err := middleware.GetJwtUser(c)

	if err != nil {
		c.Error(err)
		return
	}

	if err != nil {
		c.Error(err)
		return
	}

	web.MakeDataResp(c, "", &auth.JwtInfo{
		UserId: user.Subject,
		Type:   user.Type, //.TokenTypeString(claims.Type),
		//IpAddr:  claims.IpAddr,
		Expires: user.ExpiresAt.UTC().String()})

}

func NewAccessTokenRoute(c *gin.Context) {
	middleware.NewValidator(c).CheckIsValidRefreshToken().Success(func(validator *middleware.Validator) {

		// Generate encoded token and send it as response.
		accessToken, err := tokengen.AccessTokenUsingPermissions(c,
			validator.Claims.Subject,
			validator.Claims.Permissions)

		if err != nil {
			web.BadReqResp(c, ErrCreatingToken)
		}

		web.MakeDataResp(c, "", &web.AccessTokenResp{AccessToken: accessToken})
	})

}
