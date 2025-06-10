package scrna

import (
	"fmt"
	"strconv"

	"github.com/antonybholmes/go-scrna/scrnadbcache"
	"github.com/antonybholmes/go-sys"
	"github.com/antonybholmes/go-web"
	"github.com/gin-gonic/gin"
)

const DEFAULT_LIMIT uint16 = 20

type ScrnaParams struct {
	Genes []string `json:"genes"`
}

func parseParamsFromPost(c *gin.Context) (*ScrnaParams, error) {

	var params ScrnaParams

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
	assembly := c.Param("assembly")

	datasets, err := scrnadbcache.Datasets(species, assembly)

	if err != nil {
		c.Error(err)
		return
	}

	web.MakeDataResp(c, "", datasets)
}

// Gets expression data from a given dataset
func ScrnaGexRoute(c *gin.Context) {
	publicId := c.Param("id")

	if publicId == "" {
		c.Error(fmt.Errorf("missing id"))
		return
	}

	params, err := parseParamsFromPost(c)

	if err != nil {
		c.Error(err)
		return
	}

	// default to rna-seq
	ret, err := scrnadbcache.Gex(publicId, params.Genes)

	if err != nil {
		c.Error(err)
		return
	}

	web.MakeDataResp(c, "", ret)
}

// func ScrnaMetadataRoute(c *gin.Context) {
// 	publicId := c.Param("id")

// 	if publicId == "" {
// 		c.Error(fmt.Errorf("missing id"))
// 		return
// 	}

// 	ret, err := scrnadbcache.Metadata(publicId)

// 	if err != nil {
// 		c.Error(err)
// 		return
// 	}

// 	web.MakeDataResp(c, "", ret)
// }

func ScrnaMetadataRoute(c *gin.Context) {
	publicId := c.Param("id")

	if publicId == "" {
		c.Error(fmt.Errorf("missing id"))
		return
	}

	ret, err := scrnadbcache.Metadata(publicId)

	if err != nil {
		c.Error(err)
		return
	}

	web.MakeDataResp(c, "", ret)
}

func ScrnaGenesRoute(c *gin.Context) {
	publicId := c.Param("id")

	if publicId == "" {
		c.Error(fmt.Errorf("missing id"))
		return
	}

	ret, err := scrnadbcache.Genes(publicId)

	if err != nil {
		c.Error(err)
		return
	}

	web.MakeDataResp(c, "", ret)
}

func ScrnaSearchGenesRoute(c *gin.Context) {
	publicId := c.Param("id")

	if publicId == "" {
		c.Error(fmt.Errorf("id missing"))
		return
	}

	query := c.Query("q")

	if query == "" {
		c.Error(fmt.Errorf("query missing"))
		return
	}

	limit := DEFAULT_LIMIT

	if c.Query("limit") != "" {
		v, err := strconv.ParseUint(c.Query("limit"), 10, 16)

		if err == nil {
			limit = uint16(v)
		}
	}

	safeQuery := sys.SanitizeQuery(query)

	ret, err := scrnadbcache.SearchGenes(publicId, safeQuery, limit)

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
