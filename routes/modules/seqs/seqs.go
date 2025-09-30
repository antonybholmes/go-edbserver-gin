package seqs

import (
	"errors"

	"github.com/antonybholmes/go-dna"
	seq "github.com/antonybholmes/go-seqs"
	"github.com/antonybholmes/go-seqs/seqsdbcache"
	"github.com/antonybholmes/go-web"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

var (
	ErrNoGenomeSupplied = errors.New("must supply a genome")
)

type ReqSeqParams struct {
	Locations []string `json:"locations"`
	Scale     float64  `json:"scale"`
	BinSizes  []uint   `json:"binSizes"`
	Tracks    []string `json:"tracks"`
}

type SeqParams struct {
	Locations []*dna.Location
	Scale     float64
	BinSizes  []uint
	Tracks    []string
}

type SeqResp struct {
	Location *dna.Location         `json:"location"`
	Tracks   []*seq.TrackBinCounts `json:"tracks"`
}

func ParseSeqParamsFromPost(c *gin.Context) (*SeqParams, error) {

	var params ReqSeqParams

	err := c.Bind(&params)

	if err != nil {
		return nil, err
	}

	locations := make([]*dna.Location, 0, len(params.Locations))

	for _, loc := range params.Locations {
		location, err := dna.ParseLocation(loc)

		if err != nil {
			return nil, err
		}

		locations = append(locations, location)
	}

	return &SeqParams{
			Locations: locations,
			BinSizes:  params.BinSizes,
			Tracks:    params.Tracks,
			Scale:     params.Scale},
		nil
}

func GenomeRoute(c *gin.Context) {
	platforms, err := seqsdbcache.Genomes()

	if err != nil {
		c.Error(err)
		return
	}

	web.MakeDataResp(c, "", platforms)
}

func PlatformRoute(c *gin.Context) {
	genome := c.Param("assembly")

	platforms, err := seqsdbcache.Platforms(genome)

	if err != nil {
		c.Error(err)
		return
	}

	web.MakeDataResp(c, "", platforms)
}

func TracksRoute(c *gin.Context) {
	platform := c.Param("platform")
	genome := c.Param("assembly")

	tracks, err := seqsdbcache.Tracks(platform, genome)

	if err != nil {
		c.Error(err)
		return
	}

	web.MakeDataResp(c, "", tracks)
}

func SearchSeqRoute(c *gin.Context) {
	genome := c.Param("assembly")

	if genome == "" {
		web.BadReqResp(c, ErrNoGenomeSupplied)
		return
	}

	query := c.Query("search")

	tracks, err := seqsdbcache.Search(genome, query)

	if err != nil {
		c.Error(err)
		return
	}

	web.MakeDataResp(c, "", tracks)
}

func BinsRoute(c *gin.Context) {

	params, err := ParseSeqParamsFromPost(c)

	if err != nil {
		log.Debug().Msgf("err %s", err)
		c.Error(err)
		return
	}

	//log.Debug().Msgf("bin %v %v", params.Locations, params.BinSizes)

	ret := make([]*SeqResp, 0, len(params.Locations)) //make([]*seq.BinCounts, 0, len(params.Tracks))

	for li, location := range params.Locations {
		resp := SeqResp{Location: location, Tracks: make([]*seq.TrackBinCounts, 0, len(params.Tracks))}

		for _, track := range params.Tracks {
			reader, err := seqsdbcache.ReaderFromId(track,
				params.BinSizes[li],
				params.Scale)

			if err != nil {
				//log.Debug().Msgf("stupid err %s", err)
				c.Error(err)
				return
			}

			// guarantees something is returned even with error
			// so we can ignore the errors for now to make the api
			// more robus
			trackBinCounts, _ := reader.TrackBinCounts(location)

			// if err != nil {
			// 	return web.ErrorReq(err)
			// }

			resp.Tracks = append(resp.Tracks, trackBinCounts)
		}

		ret = append(ret, &resp)
	}

	//log.Debug().Msgf("ret %v", len(ret))

	web.MakeDataResp(c, "", ret)
}
