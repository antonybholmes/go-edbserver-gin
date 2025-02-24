package cytobands

import (
	"github.com/antonybholmes/go-cytobands/cytobandsdbcache"
	"github.com/antonybholmes/go-web"
	"github.com/gin-gonic/gin"
)

func CytobandsRoute(c *gin.Context) {

	cytobands, _ := cytobandsdbcache.Cytobands(c.Param("assembly"), c.Param("chr"))

	// if err != nil {
	// 	return web.ErrorReq(err)
	// }

	web.MakeDataResp(c, "", cytobands)
}
