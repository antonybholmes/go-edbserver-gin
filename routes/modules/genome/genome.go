package genes

import (
	"bytes"
	"encoding/csv"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/antonybholmes/go-dna"
	dnaroutes "github.com/antonybholmes/go-edbserver-gin/routes/modules/dna"
	"github.com/antonybholmes/go-genome"
	"github.com/antonybholmes/go-genome/genomedbcache"
	basemath "github.com/antonybholmes/go-math"
	"github.com/antonybholmes/go-sys/log"
	"github.com/antonybholmes/go-web"
	"github.com/gin-gonic/gin"
)

// A GeneQuery contains info from query params.
type (
	GeneQuery struct {
		Feature  genome.Feature
		Db       genome.GeneDB
		Assembly string
		GeneType string // e.g. "protein_coding", "non_coding", etc.
		// only show canonical genes
		Canonical bool
	}

	GenesResp struct {
		Location *dna.Location            `json:"location"`
		Features []*genome.GenomicFeature `json:"features"`
	}

	AnnotationResponse struct {
		Status int                      `json:"status"`
		Data   []*genome.GeneAnnotation `json:"data"`
	}
)

const (
	DefaultClosestN uint = 5
	MaxAnnotations  uint = 100
)

var (
	ErrLocationCannotBeEmpty = errors.New("location cannot be empty")
	ErrSearchTooShort        = errors.New("search too short")
)

func parseGeneQuery(c *gin.Context, assembly string) (*GeneQuery, error) {

	var feature genome.Feature = genome.GeneFeature

	switch c.Query("feature") {
	case "exon":
		feature = genome.ExonFeature
	case "transcript":
		feature = genome.TranscriptFeature
	default:
		feature = genome.GeneFeature
	}

	canonical := strings.HasPrefix(strings.ToLower(c.Query("canonical")), "t")

	geneType := c.Query("type")

	// user can specify gene type in query string, but we sanitize it
	switch {
	case strings.Contains(geneType, "protein"):
		geneType = "protein_coding"
	default:
		geneType = ""
	}

	db, err := genomedbcache.GeneDB(assembly)

	if err != nil {
		return nil, fmt.Errorf("unable to open database for assembly %s %s", assembly, err)
	}

	return &GeneQuery{
			Assembly:  assembly,
			GeneType:  geneType,
			Db:        db,
			Feature:   feature,
			Canonical: canonical},
		nil
}

// func GeneDBInfoRoute(c *gin.Context) {
// 	query, err := ParseGeneQuery(c, c.Param("assembly"))

// 	if err != nil {
// 		return web.ErrorReq(err)
// 	}

// 	info, _ := query.Db.GeneDBInfo()

// 	// if err != nil {
// 	// 	return web.ErrorReq(err)
// 	// }

// 	web.MakeDataResp(c, "", &info)
// }

func GenomesRoute(c *gin.Context) {
	infos, err := genomedbcache.GetInstance().List()

	if err != nil {
		c.Error(err)
		return
	}

	web.MakeDataResp(c, "", infos)
}

func OverlappingGenesRoute(c *gin.Context) {
	locations, err := dnaroutes.ParseLocationsFromPost(c) // dnaroutes.ParseLocationsFromPost(c)

	if err != nil {
		c.Error(err)
		return
	}

	query, err := parseGeneQuery(c, c.Param("assembly"))

	if err != nil {
		c.Error(err)
		return
	}

	if len(locations) == 0 {
		web.BadReqResp(c, ErrLocationCannotBeEmpty)
	}

	ret := make([]*GenesResp, 0, len(locations))

	for _, location := range locations {
		features, err := query.Db.OverlappingGenes(location, query.Feature, query.Canonical, query.GeneType)

		if err != nil {
			c.Error(err)
			return
		}

		ret = append(ret, &GenesResp{Location: location, Features: features})

	}

	web.MakeDataResp(c, "", &ret)
}

func SearchForGeneByNameRoute(c *gin.Context) {
	search := c.Query("search") // dnaroutes.ParseLocationsFromPost(c)

	if search == "" {
		web.BadReqResp(c, ErrSearchTooShort)
		return
	}

	fuzzyMode := c.Query("mode") == "fuzzy"

	n := web.ParseN(c, 20)

	query, err := parseGeneQuery(c, c.Param("assembly"))

	if err != nil {
		c.Error(err)
		return
	}

	canonical := strings.HasPrefix(strings.ToLower(c.Query("canonical")), "t")

	features, _ := query.Db.SearchForGeneByName(search,
		query.Feature,
		fuzzyMode,
		canonical,
		c.Query("type"),
		uint16(n))

	// if err != nil {
	// 	return web.ErrorReq(err)
	// }

	web.MakeDataResp(c, "", &features)
}

func WithinGenesRoute(c *gin.Context) {
	locations, err := dnaroutes.ParseLocationsFromPost(c) // dnaroutes.ParseLocationsFromPost(c)

	if err != nil {
		c.Error(err)
		return
	}

	query, err := parseGeneQuery(c, c.Param("assembly"))

	if err != nil {
		c.Error(err)
		return
	}

	data := make([]*genome.GenomicFeatures, len(locations))

	for li, location := range locations {
		genes, err := query.Db.WithinGenes(location, query.Feature)

		if err != nil {
			c.Error(err)
			return
		}

		data[li] = genes
	}

	web.MakeDataResp(c, "", &data)
}

// Find the n closest genes to a location
func ClosestGeneRoute(c *gin.Context) {
	locations, err := dnaroutes.ParseLocationsFromPost(c)

	if err != nil {
		c.Error(err)
		return
	}

	query, err := parseGeneQuery(c, c.Param("assembly"))

	if err != nil {
		c.Error(err)
		return
	}

	closestN := web.ParseNumParam(c, "closest", DefaultClosestN)

	data := make([]*genome.GenomicFeatures, len(locations))

	for li, location := range locations {
		genes, err := query.Db.ClosestGenes(location, uint16(closestN))

		if err != nil {
			c.Error(err)
			return
		}

		data[li] = &genome.GenomicFeatures{Location: location, Feature: genome.GeneFeature, Features: genes}
	}

	web.MakeDataResp(c, "", &data)
}

func ParsePromoterRegion(c *gin.Context) *dna.PromoterRegion {

	v := c.Query("promoter")

	if v == "" {
		return &dna.DEFAULT_PROMOTER_REGION
	}

	tokens := strings.Split(v, ",")

	s, err := strconv.ParseUint(tokens[0], 10, 0)

	if err != nil {
		return &dna.DEFAULT_PROMOTER_REGION
	}

	e, err := strconv.ParseUint(tokens[1], 10, 0)

	if err != nil {
		return &dna.DEFAULT_PROMOTER_REGION
	}

	return dna.NewPromoterRegion(uint(s), uint(e))
}

func AnnotateRoute(c *gin.Context) {
	locations, err := dnaroutes.ParseLocationsFromPost(c)

	if err != nil {
		c.Error(err)
		return
	}

	// limit amount of data returned per request to 1000 entries at a time
	locations = locations[0:basemath.Min(len(locations), int(MaxAnnotations))]

	query, err := parseGeneQuery(c, c.Param("assembly"))

	if err != nil {
		c.Error(err)
		return
	}

	closestN := web.ParseNumParam(c, "closest", DefaultClosestN)

	tssRegion := ParsePromoterRegion(c)

	output := web.ParseOutput(c)

	annotationDb := genome.NewAnnotateDb(query.Db, tssRegion, uint8(closestN))

	data := make([]*genome.GeneAnnotation, len(locations))

	for li, location := range locations {
		log.Debug().Msgf("Annotating locationsss %s", location)
		annotations, err := annotationDb.Annotate(location, genome.TranscriptFeature)

		log.Debug().Msgf("Annotated location %v", annotations)

		if err != nil {
			log.Error().Msgf("Error annotating location %s: %v", location, err)
			c.Error(err)
			return
		}

		data[li] = annotations
	}

	if output == "text" {
		tsv, err := MakeGeneTable(data, tssRegion)

		if err != nil {
			c.Error(err)
			return
		}

		c.String(http.StatusOK, tsv)
	} else {

		c.JSON(http.StatusOK, AnnotationResponse{Status: http.StatusOK, Data: data})
	}
}

func MakeGeneTable(
	data []*genome.GeneAnnotation,
	ts *dna.PromoterRegion,
) (string, error) {
	var buffer bytes.Buffer
	wtr := csv.NewWriter(&buffer)
	wtr.Comma = '\t'

	closestN := len(data[0].ClosestGenes)

	headers := make([]string, 5+4*closestN)

	headers[0] = "Location"
	headers[1] = "Gene Id"
	headers[2] = "Gene Symbol"
	headers[3] = fmt.Sprintf(
		"Relative To Gene (prom=-%d/+%dkb)",
		ts.Offset5P()/1000,
		ts.Offset3P()/1000)
	headers[4] = "TSS Distance"
	//headers[5] = "Gene Location"

	idx := 6
	for i := 1; i <= closestN; i++ {
		headers[idx] = fmt.Sprintf("#%d Closest Id", i)
		headers[idx] = fmt.Sprintf("#%d Closest Gene Symbols", i)
		headers[idx] = fmt.Sprintf(
			"#%d Relative To Closet Gene (prom=-%d/+%dkb)",
			i,
			ts.Offset5P()/1000,
			ts.Offset3P()/1000)
		headers[idx] = fmt.Sprintf("#%d TSS Closest Distance", i)
		headers[idx] = fmt.Sprintf("#%d Gene Location", i)
	}

	err := wtr.Write(headers)

	if err != nil {
		return "", err
	}

	for _, annotation := range data {
		n := len(annotation.WithinGenes)
		geneIds := make([]string, n)
		geneNames := make([]string, n)
		promLabels := make([]string, n)
		tssDists := make([]string, n)

		for i, gene := range annotation.WithinGenes {
			geneIds[i] = gene.GeneId
			geneNames[i] = gene.GeneSymbol
			promLabels[i] = gene.PromLabel
			tssDists[i] = strconv.Itoa(gene.TssDist)

		}

		row := []string{annotation.Location.String(),
			strings.Join(geneIds, genome.FeatureSeparator),
			strings.Join(geneNames, genome.FeatureSeparator),
			strings.Join(promLabels, genome.FeatureSeparator),
			strings.Join(tssDists, genome.FeatureSeparator)}

		for _, closestGene := range annotation.ClosestGenes {
			row = append(row, closestGene.GeneId)
			row = append(row, genome.GeneWithStrandLabel(closestGene.GeneSymbol, closestGene.Strand))
			row = append(row, closestGene.PromLabel)
			row = append(row, strconv.Itoa(closestGene.TssDist))
			//row = append(row, closestGene.Location.String())
		}

		err := wtr.Write(row)

		if err != nil {
			return "", err
		}
	}

	wtr.Flush()

	return buffer.String(), nil
}
