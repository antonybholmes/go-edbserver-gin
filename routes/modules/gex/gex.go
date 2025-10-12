package gex

import (
	"strings"

	"github.com/antonybholmes/go-gex"
	"github.com/antonybholmes/go-gex/gexdbcache"
	"github.com/antonybholmes/go-web"
	"github.com/gin-gonic/gin"
)

type GexParams struct {
	Species    string        `json:"species"`
	Technology string        `json:"technology"`
	ExprType   *gex.ExprType `json:"exprType"`
	Genes      []string      `json:"genes"`
	Datasets   []string      `json:"datasets"`
}

func parseParamsFromPost(c *gin.Context) (*GexParams, error) {

	var params GexParams

	err := c.Bind(&params)

	if err != nil {
		return nil, err
	}

	return &params, nil
}

func SpeciesRoute(c *gin.Context) {

	types, err := gexdbcache.Species()

	if err != nil {
		c.Error(err)
		return
	}

	web.MakeDataResp(c, "", types)
}

func TechnologiesRoute(c *gin.Context) {

	//species := c.Param("species")

	// technologies, err := gexdbcache.AllTechnologies() //gexdbcache.Technologies() //species)

	// if err != nil {
	// 	c.Error(err)
	// 	return
	// }

	technologies := gexdbcache.Technologies() //gexdbcache.Technologies()

	web.MakeDataResp(c, "", technologies)
}

func ExprTypesRoute(c *gin.Context) {

	params, err := parseParamsFromPost(c)

	if err != nil {
		c.Error(err)
		return
	}

	exprTypes, err := gexdbcache.ExprTypes(params.Datasets)

	if err != nil {
		c.Error(err)
		return
	}

	web.MakeDataResp(c, "", exprTypes)
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

func GexDatasetsRoute(c *gin.Context) {

	species := c.Param("species")

	technology := c.Param("technology")

	datasets, err := gexdbcache.Datasets(species, technology)

	if err != nil {
		c.Error(err)
		return
	}

	web.MakeDataResp(c, "", datasets)
}

func GexGeneExprRoute(c *gin.Context) {
	params, err := parseParamsFromPost(c)

	if err != nil {
		c.Error(err)
		return
	}

	results := make([]*gex.SearchResults, 0, len(params.Datasets))

	for _, datasetId := range params.Datasets {

		if strings.Contains(datasetId, "microarray") {
			// microarray
			ret, err := gexdbcache.FindMicroarrayValues(datasetId, params.Genes)

			if err != nil {
				c.Error(err)
				return
			}

			results = append(results, ret)
		} else {
			// default to rna-seq
			ret, err := gexdbcache.FindSeqValues(datasetId, params.ExprType, params.Genes)

			if err != nil {
				c.Error(err)
				return
			}

			results = append(results, ret)
		}

		web.MakeDataResp(c, "", results)
	}
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
