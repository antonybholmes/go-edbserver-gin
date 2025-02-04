package motifroutes

import (
	"github.com/antonybholmes/go-edb-server-gin/routes"
	"github.com/antonybholmes/go-motifs"
	"github.com/antonybholmes/go-motifs/motifsdb"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

const MIN_SEARCH_LEN = 3

type ReqParams struct {
	Search     string `json:"search"`
	Exact      bool   `json:"exact"`
	Reverse    bool   `json:"reverse"`
	Complement bool   `json:"complement"`
}

type MotifRes struct {
	Search     string          `json:"search"`
	Motifs     []*motifs.Motif `json:"motifs"`
	Reverse    bool            `json:"reverse"`
	Complement bool            `json:"complement"`
}

func ParseParamsFromPost(c *gin.Context) (*ReqParams, error) {

	var params ReqParams

	err := c.Bind(&params)

	if err != nil {
		return nil, err
	}

	return &params, nil
}

func DatasetsRoute(c *gin.Context) {

	// Don't care about the errors, just plug empty list into failures
	datasets, err := motifsdb.Datasets()

	if err != nil {
		c.Error(err)
		return
	}

	routes.MakeDataResp(c, "", datasets)

	//routes.MakeDataResp(c, "", mutationdbcache.GetInstance().List())
}

func SearchRoute(c *gin.Context) {

	params, err := ParseParamsFromPost(c)

	if err != nil {
		c.Error(err)
		return
	}

	search := params.Search

	if len(search) < MIN_SEARCH_LEN {
		routes.ErrorResp(c, "Search too short")
		return
	}

	log.Debug().Msgf("motif %v", params)

	// Don't care about the errors, just plug empty list into failures
	motifs, err := motifsdb.Search(search, params.Reverse, params.Complement)

	if err != nil {
		log.Debug().Msgf("motif %s", err)
		c.Error(err)
		return
	}

	routes.MakeDataResp(c, "",
		MotifRes{
			Search:     search,
			Motifs:     motifs,
			Reverse:    params.Reverse,
			Complement: params.Complement,
		})

	//routes.MakeDataResp(c, "", mutationdbcache.GetInstance().List())
}
