package gex

import (
	"github.com/antonybholmes/go-hubs/hubsdbcache"
	"github.com/antonybholmes/go-web"
	"github.com/gin-gonic/gin"
)

func HubsRoute(c *gin.Context) {

	assembly := c.Param("assembly")

	hubs, err := hubsdbcache.Hubs(assembly)

	if err != nil {
		c.Error(err)
		return
	}

	web.MakeDataResp(c, "", hubs)
}
