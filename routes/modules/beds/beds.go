package beds

import (
	"github.com/gin-gonic/gin"

	"github.com/antonybholmes/go-beds"
	"github.com/antonybholmes/go-beds/bedsdbcache"
	"github.com/antonybholmes/go-dna"
	"github.com/antonybholmes/go-web"
)

type ReqBedsParams struct {
	Location string   `json:"location"`
	Beds     []string `json:"beds"`
}

type BedsParams struct {
	Location *dna.Location `json:"location"`
	Beds     []string      `json:"beds"`
}

func ParseBedParamsFromPost(c *gin.Context) (*BedsParams, error) {

	var params ReqBedsParams

	err := c.Bind(&params)

	if err != nil {
		return nil, err
	}

	location, err := dna.ParseLocation(params.Location)

	if err != nil {
		return nil, err
	}

	return &BedsParams{Location: location, Beds: params.Beds}, nil
}

func GenomeRoute(c *gin.Context) {
	platforms, err := bedsdbcache.Genomes()

	if err != nil {
		c.Error(err)
		return
	}

	web.MakeDataResp(c, "", platforms)
}

func PlatformRoute(c *gin.Context) {
	genome := c.Param("assembly")

	platforms, err := bedsdbcache.Platforms(genome)

	if err != nil {
		c.Error(err)
		return
	}

	web.MakeDataResp(c, "", platforms)
}

func SearchBedsRoute(c *gin.Context) {
	genome := c.Param("assembly")

	if genome == "" {
		web.BadReqResp(c, "must supply a genome")
	}

	query := c.Query("search")

	tracks, err := bedsdbcache.Search(genome, query)

	if err != nil {
		c.Error(err)
		return
	}

	web.MakeDataResp(c, "", tracks)
}

func BedRegionsRoute(c *gin.Context) {

	params, err := ParseBedParamsFromPost(c)

	if err != nil {
		c.Error(err)
		return
	}

	if len(params.Beds) == 0 {
		web.BadReqResp(c, "at least 1 bed id must be supplied")
	}

	ret := make([][]*beds.BedRegion, 0, len(params.Beds))

	for _, bed := range params.Beds {

		//log.Debug().Msgf("bed id %s", bed)

		reader, err := bedsdbcache.ReaderFromId(bed)

		if err != nil {
			c.Error(err)
			return
		}

		features, _ := reader.OverlappingRegions(params.Location)

		ret = append(ret, features)
	}

	web.MakeDataResp(c, "", ret)
}
