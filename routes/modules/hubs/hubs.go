package gex

import (
	"github.com/antonybholmes/go-gex"
	"github.com/antonybholmes/go-hubs/hubsdbcache"
	"github.com/antonybholmes/go-web"
	"github.com/gin-gonic/gin"
)

type HubParams struct {
	Platform     *gex.ValueType    `json:"platform"`
	GexValueType *gex.GexValueType `json:"gexValueType"`
	Genes        []string          `json:"genes"`
	Datasets     []string          `json:"datasets"`
}

func parseParamsFromPost(c *gin.Context) (*HubParams, error) {

	var params HubParams

	err := c.Bind(&params)

	if err != nil {
		return nil, err
	}

	return &params, nil
}

func HubsRoute(c *gin.Context) {

	types, err := hubsdbcache.Hubs()

	if err != nil {
		c.Error(err)
		return
	}

	web.MakeDataResp(c, "", types)
}
