package gex

import (
	"github.com/antonybholmes/go-gex"
	"github.com/antonybholmes/go-gex/gexdbcache"
	"github.com/antonybholmes/go-web"
	"github.com/gin-gonic/gin"
)

type GexParams struct {
	Platform     *gex.ValueType    `json:"platform"`
	GexValueType *gex.GexValueType `json:"gexValueType"`
	Genes        []string          `json:"genes"`
	Datasets     []string          `json:"datasets"`
}

func parseParamsFromPost(c *gin.Context) (*GexParams, error) {

	var params GexParams

	err := c.Bind(&params)

	if err != nil {
		return nil, err
	}

	return &params, nil
}

func PlaformsRoute(c *gin.Context) {

	types, err := gexdbcache.Platforms()

	if err != nil {
		c.Error(err)
		return
	}

	web.MakeDataResp(c, "", types)
}

func GexValueTypesRoute(c *gin.Context) {

	params, err := parseParamsFromPost(c)

	if err != nil {
		c.Error(err)
		return
	}

	valueTypes, err := gexdbcache.GexValueTypes(params.Platform.Id)

	if err != nil {
		c.Error(err)
		return
	}

	web.MakeDataResp(c, "", valueTypes)
}

func GexDatasetsRoute(c *gin.Context) {

	params, err := parseParamsFromPost(c)

	if err != nil {
		c.Error(err)
		return
	}

	datasets, err := gexdbcache.Datasets(params.Platform.Id)

	if err != nil {
		c.Error(err)
		return
	}

	web.MakeDataResp(c, "", datasets)
}

func GexGeneExpRoute(c *gin.Context) {
	params, err := parseParamsFromPost(c)

	if err != nil {
		c.Error(err)
		return
	}

	// convert search genes into actual genes in the database
	gexGenes, err := gexdbcache.GetGenes(params.Genes)

	if err != nil {
		c.Error(err)
		return
	}

	if params.Platform.Id == 2 {
		// microarray
		ret, err := gexdbcache.MicroarrayValues(gexGenes,
			params.Platform,
			params.GexValueType, params.Datasets)

		if err != nil {
			c.Error(err)
			return
		}

		web.MakeDataResp(c, "", ret)
	} else {
		// default to rna-seq
		ret, err := gexdbcache.RNASeqValues(gexGenes,
			params.Platform,
			params.GexValueType,
			params.Datasets)

		if err != nil {
			c.Error(err)
			return
		}

		web.MakeDataResp(c, "", ret)
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
