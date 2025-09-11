package mutations

import (
	"github.com/antonybholmes/go-dna"
	authenticationroutes "github.com/antonybholmes/go-edbserver-gin/routes/authentication"
	"github.com/antonybholmes/go-mutations"
	"github.com/antonybholmes/go-mutations/mutationdbcache"
	"github.com/antonybholmes/go-web"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

type MutationParams struct {
	Locations []*dna.Location
	Datasets  []string
}

type ReqMutationParams struct {
	Locations []string `json:"locations"`
	Datasets  []string `json:"datasets"`
}

func ParseParamsFromPost(c *gin.Context) (*MutationParams, error) {

	var locs ReqMutationParams

	err := c.Bind(&locs)

	if err != nil {
		return nil, err
	}

	locations, err := dna.ParseLocations(locs.Locations)

	if err != nil {
		return nil, err
	}

	return &MutationParams{locations, locs.Datasets}, nil
}

func MutationDatasetsRoute(c *gin.Context) {

	assembly := c.Param("assembly")

	datasets, err := mutationdbcache.List(assembly)

	if err != nil {
		c.Error(err)
		return
	}

	web.MakeDataResp(c, "", datasets)
}

func MutationsRoute(c *gin.Context) {
	assembly := c.Param("assembly")

	params, err := ParseParamsFromPost(c)

	if err != nil {
		c.Error(err)
		return
	}

	location := params.Locations[0]

	search, err := mutationdbcache.GetInstance().Search(assembly,
		location,
		params.Datasets)

	if err != nil {
		c.Error(err)
		return
	}

	// if err != nil {
	// 	return web.ErrorReq(err)
	// }

	// ret := make([]*mutations.SearchResults, len(locations))

	// for i, location := range locations {
	// 	mutations, err := db.FindMutations(location)

	// 	if err != nil {
	// 		return web.ErrorReq(err)
	// 	}

	// 	ret[i] = mutations
	// }

	web.MakeDataResp(c, "", search)

	//web.MakeDataResp(c, "", mutationdbcache.GetInstance().List())
}

type MafResp struct {
	Location *dna.Location `json:"location"`
	//Info      []*mutations.MutationDB `json:"info"`
	Samples   int        `json:"samples"`
	Mutations [][]string `json:"mutations"`
}

// func MafRoute(c *gin.Context) {
// 	 NewValidator(c).Success(func(validator *Validator) {
// 		assembly := c.Param("assembly")

// 		params, err := ParseParamsFromPost(c)

// 		if err != nil {
// 			return web.ErrorReq(err)
// 		}

// 		location := params.Locations[0]

// 		//assembly := c.Param("assembly")
// 		//name := c.Param("name")

// 		ret := MafResp{Location: location,
// 			//Info:      make([]*mutations.MutationDBInfo, len(params.Databases)),
// 			Mutations: make([][]string, location.Len())}

// 		for i := range location.Len() {
// 			ret.Mutations[i] = make([]string, 0, 10)
// 		}

// 		sampleMap := make([]map[int]struct{}, location.Len())

// 		for i := range location.Len() {
// 			sampleMap[i] = make(map[int]struct{})
// 		}

// 		for _, id := range params.Datasets {
// 			dataset, err := mutationdbcache.GetDataset(assembly, id)

// 			if err != nil {
// 				return web.ErrorReq(err)
// 			}

// 			results, err := dataset.Search(location)

// 			if err != nil {
// 				return web.ErrorReq(err)
// 			}

// 			// sum the total number of samples involved
// 			ret.Samples += len(dataset.Samples)

// 			for _, mutation := range results.Mutations {
// 				offset := mutation.Start - location.Start
// 				sample := mutation.Sample

// 				_, ok := sampleMap[offset][sample]

// 				if !ok {
// 					sampleMap[offset][sample] = struct{}{}
// 				}

// 			}

// 			//ret.Info[dbi] = db.Info
// 		}

// 		// sort each pileup
// 		for ci := range location.Len() {
// 			if len(sampleMap[ci]) > 0 {
// 				samples := make([]int, 0, len(sampleMap[ci]))

// 				for sample := range sampleMap[ci] {
// 					samples = append(samples, sample)
// 				}

// 				// sort the samples for ease of use
// 				sort.Ints(samples)

// 				ret.Mutations[ci] = append(ret.Mutations[ci], samples...)
// 			}

// 		}

// 		web.MakeDataResp(c, "", ret)
// 	})

// 	//web.MakeDataResp(c, "", mutationdbcache.GetInstance().List())
// }

type PileupResp struct {
	Location *dna.Location `json:"location"`
	//Info      []*mutations.MutationDBInfo `json:"info"`
	Samples   int                     `json:"samples"`
	Mutations [][]*mutations.Mutation `json:"mutations"`
}

func PileupRoute(c *gin.Context) {
	authenticationroutes.NewValidator(c).Success(func(validator *authenticationroutes.Validator) {

		assembly := c.Param("assembly")

		params, err := ParseParamsFromPost(c)

		if err != nil {
			c.Error(err)
			return
		}

		log.Debug().Msgf("pileup: %v", params)

		location := params.Locations[0]

		//assembly := c.Param("assembly")
		//name := c.Param("name")

		// ret := PileupResp{Location: location,
		// 	// one metadata file for each database requested
		// 	//Info:      make([]*mutations.MutationDBInfo, len(params.Databases)),
		// 	Mutations: make([][]*mutations.Mutation, location.Len())}

		// for i := range location.Len() {
		// 	ret.Mutations[i] = make([]*mutations.Mutation, 0, 10)
		// }

		search, err := mutationdbcache.GetInstance().Search(assembly,
			location,
			params.Datasets)

		if err != nil {
			log.Debug().Msgf("here 1 %s", err)
			c.Error(err)
			return
		}

		pileup, err := mutations.GetPileup(search)

		if err != nil {
			c.Error(err)
			return
		}

		web.MakeDataResp(c, "", pileup)
	})

	//web.MakeDataResp(c, "", mutationdbcache.GetInstance().List())
}
