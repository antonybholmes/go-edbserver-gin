package scrna

import (
	"github.com/antonybholmes/go-scrna/scrnadbcache"
	"github.com/antonybholmes/go-web"
	"github.com/gin-gonic/gin"
)

type Scrna struct {
	Genes    []string `json:"genes"`
	Datasets []string `json:"datasets"`
}

func parseParamsFromPost(c *gin.Context) (*Scrna, error) {

	var params Scrna

	err := c.Bind(&params)

	if err != nil {
		return nil, err
	}

	return &params, nil
}

func ScrnaSpeciesRoute(c *gin.Context) {

	types, err := scrnadbcache.Species()

	if err != nil {
		c.Error(err)
		return
	}

	web.MakeDataResp(c, "", types)
}

func ScrnaAssembliesRoute(c *gin.Context) {

	species := c.Param("species")

	// technologies, err := gexdbcache.AllTechnologies() //gexdbcache.Technologies() //species)

	// if err != nil {
	// 	c.Error(err)
	// 	return
	// }

	assemblies, err := scrnadbcache.Assemblies(species) //gexdbcache.Technologies()

	if err != nil {
		c.Error(err)
		return
	}

	web.MakeDataResp(c, "", assemblies)
}

// func GexValueTypesRoute(c *gin.Context) {

// 	params, err := parseParamsFromPost(c)

// 	if err != nil {
// 		c.Error(err)
// 		return
// 	}

// 	valueTypes, err := gexdbcache.GexValueTypes(params.Platform.Id)

// 	if err != nil {
// 		c.Error(err)
// 		return
// 	}

// 	web.MakeDataResp(c, "", valueTypes)
// }

func ScrnaDatasetsRoute(c *gin.Context) {

	species := c.Param("species")
	technology := c.Param("technology")

	datasets, err := scrnadbcache.Datasets(species, technology)

	if err != nil {
		c.Error(err)
		return
	}

	web.MakeDataResp(c, "", datasets)
}

func ScrnaGeneExpRoute(c *gin.Context) {
	params, err := parseParamsFromPost(c)

	if err != nil {
		c.Error(err)
		return
	}

	// default to rna-seq
	ret, err := scrnadbcache.FindGexValues(params.Datasets, params.Genes)

	if err != nil {
		c.Error(err)
		return
	}

	web.MakeDataResp(c, "", ret)

}

// func GexRoute(c *gin.Context) {
// 	gexType := c.Param("type")

// 	params, err := ParseParamsFromPost(c)

// 	if err != nil {
// 		return web.ErrorReq(err)
// 	}

// 	search, err := gexdbcache.GetInstance().Search(gexType, params.Datasets, params.Genes)

// 	if err != nil {
// 		return web.ErrorReq(err)
// 	}

// 	web.MakeDataResp(c, "", search)

// 	//web.MakeDataResp(c, "", mutationdbcache.GetInstance().List())
// }
