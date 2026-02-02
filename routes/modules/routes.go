package modules

import (
	bedroutes "github.com/antonybholmes/go-beds/routes"
	cytobandroutes "github.com/antonybholmes/go-cytobands/routes"
	dnaroutes "github.com/antonybholmes/go-dna/routes"
	geneconvroutes "github.com/antonybholmes/go-geneconv/routes"
	genomeroutes "github.com/antonybholmes/go-genome/routes"
	gexroutes "github.com/antonybholmes/go-gex/routes"
	hubroutes "github.com/antonybholmes/go-hubs/routes"
	motifroutes "github.com/antonybholmes/go-motifs/routes"
	mutationroutes "github.com/antonybholmes/go-mutations/routes"
	pathwayroutes "github.com/antonybholmes/go-pathway/routes"
	scrnaroutes "github.com/antonybholmes/go-scrna/routes"
	seqroutes "github.com/antonybholmes/go-seqs/routes"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine, rulesMiddleware gin.HandlerFunc) {
	//
	// module groups: start
	//

	moduleGroup := r.Group("/modules")
	//moduleGroup.Use(jwtMiddleWare,JwtIsAccessTokenMiddleware)

	dnaGroup := moduleGroup.Group("/dna")
	dnaGroup.POST("/:assembly", dnaroutes.DNARoute)
	dnaGroup.GET("/genomes", dnaroutes.GenomesRoute)

	genomeGroup := moduleGroup.Group("/genome")
	genomeGroup.GET("/genomes", genomeroutes.GenomesRoute)
	genomeGroup.POST("/within/:assembly", genomeroutes.WithinGenesRoute)
	genomeGroup.POST("/closest/:assembly", genomeroutes.ClosestGeneRoute)
	genomeGroup.POST("/annotate/:assembly", genomeroutes.AnnotateRoute)
	genomeGroup.POST("/overlap/:assembly", genomeroutes.OverlappingGenesRoute)
	genomeGroup.GET("/info/:assembly", genomeroutes.SearchForGeneByNameRoute)

	// mutationsGroup := moduleGroup.Group("/mutations",
	// 	jwtMiddleWare,
	// 	JwtIsAccessTokenMiddleware,
	// 	NewJwtPermissionsMiddleware("rdf"))

	mutationsGroup := moduleGroup.Group("/mutations")

	mutationsGroup.POST("/:assembly/:name",
		mutationroutes.MutationsRoute)
	mutationsGroup.POST("/maf/:assembly",
		mutationroutes.PileupRoute)

	mutationsProtectedGroup := mutationsGroup.Group("",
		rulesMiddleware,
	)

	mutationsProtectedGroup.GET("/datasets/:assembly",
		mutationroutes.MutationDatasetsRoute)

	mutationsProtectedGroup.POST("/pileup/:assembly",
		mutationroutes.PileupRoute,
	)

	gexGroup := moduleGroup.Group("/gex")
	gexGroup.GET("/species", gexroutes.SpeciesRoute)
	gexGroup.GET("/technologies", gexroutes.TechnologiesRoute)

	//gexGroup.GET("/types", gexroutes.GexValueTypesRoute)

	// protected routes
	gexProtectedGroup := gexGroup.Group("",
		rulesMiddleware,
		//jwtUserMiddleWare,
		//accessTokenMiddleware,
		//rdfRoleMiddleware
	)

	gexProtectedGroup.GET("/datasets/:species/:technology",
		gexroutes.GexDatasetsRoute)

	gexProtectedGroup.POST("/expr/types",
		gexroutes.ExprTypesRoute,
	)

	gexProtectedGroup.POST("/expr", gexroutes.GexGeneExprRoute)

	scrnaGroup := moduleGroup.Group("/scrna")
	scrnaGroup.GET("/species", scrnaroutes.ScrnaSpeciesRoute)
	scrnaGroup.GET("/assemblies/:species", scrnaroutes.ScrnaAssembliesRoute)
	//gexGroup.GET("/types", gexroutes.GexValueTypesRoute)

	scrnaProtectedGroup := scrnaGroup.Group("",
		rulesMiddleware,
		//jwtUserMiddleWare,
	//accessTokenMiddleware,
	//rdfRoleMiddleware
	)

	scrnaProtectedGroup.GET("/datasets/:species/:assembly",
		scrnaroutes.ScrnaDatasetsRoute)

	// scrnaGroup.GET("/clusters/:id",
	// 	jwtUserMiddleWare,
	// 	accessTokenMiddleware,
	// 	rdfRoleMiddleware,
	// 	scrnaroutes.ScrnaClustersRoute,
	// )

	scrnaProtectedGroup.GET("/metadata/:id",
		scrnaroutes.ScrnaMetadataRoute)

	scrnaProtectedGroup.GET("/genes/:id",
		scrnaroutes.ScrnaGenesRoute)

	scrnaProtectedGroup.GET("/genes/search/:id",
		scrnaroutes.ScrnaSearchGenesRoute)

	scrnaProtectedGroup.POST("/gex/:id",
		scrnaroutes.ScrnaGexRoute)

	hubsGroup := moduleGroup.Group("/hubs")
	hubsGroup.GET("/:assembly",
		rulesMiddleware,
		//jwtUserMiddleWare,
		//accessTokenMiddleware,
		//rdfRoleMiddleware,
		hubroutes.HubsRoute,
	)

	geneConvGroup := moduleGroup.Group("/geneconv")
	geneConvGroup.POST("/convert/:from/:to", geneconvroutes.ConvertRoute)

	// geneConvGroup.POST("/:species", func(c *gin.Context) {
	// 	return geneconvroutes.GeneInfoRoute(c, "")
	// })

	motifsGroup := moduleGroup.Group("/motifs")
	motifsGroup.GET("/datasets", motifroutes.DatasetsRoute)
	motifsGroup.POST("/search", motifroutes.SearchRoute)

	pathwayGroup := moduleGroup.Group("/pathway")
	pathwayGroup.GET("/genes", pathwayroutes.GenesRoute)
	pathwayGroup.POST("/dataset", pathwayroutes.DatasetRoute)
	pathwayGroup.GET("/datasets", pathwayroutes.DatasetsRoute)
	pathwayGroup.POST("/overlap", pathwayroutes.PathwayOverlapRoute)

	seqsGroup := moduleGroup.Group("/seqs",
		rulesMiddleware,
		//jwtUserMiddleWare,
		//accessTokenMiddleware,
		//rdfRoleMiddleware
	)

	//seqsGroup.GET("/genomes", seqroutes.GenomeRoute)
	//seqsGroup.GET("/platforms/:assembly", seqroutes.PlatformRoute)
	seqsGroup.GET("/platforms", seqroutes.PlatformsRoute)
	//tracksGroup.GET("/:platform/:assembly/tracks", seqroutes.TracksRoute)
	seqsGroup.GET("/search/:assembly", seqroutes.SearchSeqRoute)
	seqsGroup.POST("/bins", seqroutes.BinsRoute)

	cytobandsGroup := moduleGroup.Group("/cytobands")
	cytobandsGroup.GET("/:assembly/:chr", cytobandroutes.CytobandsRoute)

	bedsGroup := moduleGroup.Group("/beds",
		rulesMiddleware,
	//jwtUserMiddleWare,
	//accessTokenMiddleware,
	//rdfRoleMiddleware
	)
	//bedsGroup.GET("/genomes", bedroutes.GenomeRoute)
	bedsGroup.GET("/platforms/:assembly", bedroutes.PlatformsRoute)
	bedsGroup.GET("/search/:assembly", bedroutes.SearchBedsRoute)
	bedsGroup.POST("/regions", bedroutes.BedRegionsRoute)
}
