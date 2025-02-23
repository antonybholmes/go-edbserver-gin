package auth0

import (
	"github.com/antonybholmes/go-edb-server-gin/routes"
	"github.com/gin-gonic/gin"
)

func ValidateAuth0TokenRoute(c *gin.Context) {

	// bytes, err := os.ReadFile("auth0.key.pub")
	// if err != nil {
	// 	log.Fatal().Msgf("%s", err)
	// }

	// key, err := jwt.ParseRSAPublicKeyFromPEM(bytes)
	// if err != nil {
	// 	log.Fatal().Msgf("%s", err)
	// }

	// h := c.Request().Header.Get("Authorization")

	// tokens := strings.SplitN(h, " ", 2)
	// token := tokens[1]

	// log.Debug().Msgf("tok: %v", h)

	// hmm, err := jwt.ParseWithClaims(token, &TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
	// 	return key, nil
	// })

	// if err != nil {
	// 	log.Debug().Msgf("%s", err)
	// }

	//user := c.Get("user").(*jwt.Token)
	//myClaims := user.Claims.(*auth.Auth0TokenClaims)

	//user := c.Get("user").(*jwt.Token)
	//claims := user.Claims.(*TokenClaims)

	//log.Debug().Msgf("auth0 claims %v", myClaims)
	//log.Debug().Msgf("auth0 claims %v", myClaims.Email)

	routes.MakeOkResp(c, "user was signed out")
}
