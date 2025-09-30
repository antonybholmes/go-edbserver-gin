package authorization

import (
	"errors"

	"github.com/antonybholmes/go-edbserver-gin/consts"
	authenticationroutes "github.com/antonybholmes/go-edbserver-gin/routes/authentication"
	"github.com/antonybholmes/go-web"
	"github.com/antonybholmes/go-web/auth"
	"github.com/antonybholmes/go-web/middleware"
	"github.com/antonybholmes/go-web/tokengen"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
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
	t, err := middleware.ParseToken(c)

	if err != nil {
		c.Error(err)
		return
	}

	claims := auth.TokenClaims{}

	_, err = jwt.ParseWithClaims(t, &claims, func(token *jwt.Token) (any, error) {
		return consts.JWT_RSA_PUBLIC_KEY, nil
	})

	if err != nil {
		c.Error(err)
		return
	}

	web.MakeDataResp(c, "", &auth.JwtInfo{
		UserId: claims.UserId,
		Type:   claims.Type, //.TokenTypeString(claims.Type),
		//IpAddr:  claims.IpAddr,
		Expires: claims.ExpiresAt.UTC().String()})

}

func NewAccessTokenRoute(c *gin.Context) {
	authenticationroutes.NewValidator(c).CheckIsValidRefreshToken().Success(func(validator *authenticationroutes.Validator) {

		// Generate encoded token and send it as response.
		accessToken, err := tokengen.AccessToken(c,
			validator.Claims.UserId,
			validator.Claims.Roles)

		if err != nil {
			web.BadReqResp(c, ErrCreatingToken)
		}

		web.MakeDataResp(c, "", &web.AccessTokenResp{AccessToken: accessToken})
	})

}
