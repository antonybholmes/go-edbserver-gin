package cytobands

import (
	"github.com/antonybholmes/go-auth/routes"
	"github.com/gin-gonic/gin"

	"github.com/antonybholmes/go-cytobands/cytobandsdbcache"
)

func CytobandsRoute(c *gin.Context) {

	cytobands, _ := cytobandsdbcache.Cytobands(c.Param("assembly"), c.Param("chr"))

	// if err != nil {
	// 	return routes.ErrorReq(err)
	// }

	routes.MakeDataResp(c, "", cytobands)
}
