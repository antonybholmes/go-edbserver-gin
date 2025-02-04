package toolsroutes

import (
	"crypto/rand"
	"math/big"
	"strconv"

	"github.com/antonybholmes/go-auth"
	"github.com/antonybholmes/go-edb-server-gin/routes"
	"github.com/gin-gonic/gin"
)

type HashResp struct {
	Password string `json:"password"`
	Hash     string `json:"hash"`
}

func HashedPasswordRoute(c *gin.Context) {

	password := c.Query("password")

	if len(password) == 0 {
		routes.ErrorResp(c, "password cannot be empty")
		return
	}

	hash := auth.HashPassword(password)

	ret := HashResp{Password: password, Hash: hash}

	routes.MakeDataResp(c, "", ret)
}

type KeyResp struct {
	Key    string `json:"key"`
	Length int    `json:"length"`
}

func RandomKeyRoute(c *gin.Context) {

	l, err := strconv.Atoi(c.Query("l"))

	if err != nil || l < 1 {
		routes.ErrorResp(c, "length cannot be zero")
		return

	}

	key, err := generateRandomString(l)

	if err != nil {
		c.Error(err)
		return
	}

	ret := KeyResp{Key: key, Length: l}

	routes.MakeDataResp(c, "", ret)
}

// func generateRandomString(length int) string {
// 	randomBytes := make([]byte, length)
// 	_, err := rand.Read(randomBytes)
// 	if err != nil {
// 		panic(err) // Handle the error appropriately in your application
// 	}

// 	return base64.URLEncoding.EncodeToString(randomBytes)
// }

const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

// generateRandomString generates a random string of specified length from the letters set.
func generateRandomString(length int) (string, error) {
	b := make([]byte, length)
	for i := range b {
		// Generate a random index
		index, err := rand.Int(rand.Reader, big.NewInt(int64(len(letters))))
		if err != nil {
			return "", err
		}
		b[i] = letters[index.Int64()]
	}
	return string(b), nil
}
