package dna

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/antonybholmes/go-dna"
	"github.com/antonybholmes/go-dna/dnadbcache"
	"github.com/antonybholmes/go-web/routes"
	"github.com/gin-gonic/gin"
)

const DEFAULT_ASSEMBLY string = "grch38"
const DEFAULT_CHR string = "chr1" //"chr3"
const DEFAULT_START uint = 100000 //187728170
const DEFAULT_END uint = 100100   //187752257

type ReqLocs struct {
	Locations []string `json:"locations"`
}

type DNA struct {
	Location *dna.Location `json:"location"`
	Seq      string        `json:"seq"`
}

type DNAResp struct {
	Assembly     string `json:"assembly"`
	Format       string `json:"format"`
	IsRev        bool   `json:"isRev"`
	IsComplement bool   `json:"isComp"`
	Seqs         []*DNA `json:"seqs"`
}

type DNAQuery struct {
	Rev        bool
	Comp       bool
	Format     string
	RepeatMask string
}

func ParseLocation(c *gin.Context) (*dna.Location, error) {
	chr := DEFAULT_CHR
	start := DEFAULT_START
	end := DEFAULT_END

	var v string
	var err error
	var n uint64

	v = c.Query("chr")

	if v != "" {
		chr = v
	}

	v = c.Query("start")

	if v != "" {
		n, err = strconv.ParseUint(v, 10, 0)

		if err != nil {
			return nil, fmt.Errorf("%s is an invalid start", v)
		}

		start = uint(n)
	}

	v = c.Query("end")

	if v != "" {
		n, err = strconv.ParseUint(v, 10, 0)

		if err != nil {
			return nil, fmt.Errorf("%s is an invalid end", v)
		}

		end = uint(n)
	}

	loc := dna.NewLocation(chr, start, end)

	return loc, nil
}

func ParseLocationsFromPost(c *gin.Context) ([]*dna.Location, error) {

	var locs ReqLocs

	err := c.Bind(&locs)

	if err != nil {
		return nil, err
	}

	ret, err := dna.ParseLocations(locs.Locations)

	if err != nil {
		return nil, err
	}

	return ret, nil
}

func ParseDNAQuery(c *gin.Context) (*DNAQuery, error) {
	var err error

	rev := false
	v := c.Query("rev")

	if v != "" {
		rev, err = strconv.ParseBool(v)

		if err != nil {
			rev = false
		}
	}

	comp := false
	v = c.Query("comp")

	if v != "" {
		comp, err = strconv.ParseBool(v)

		if err != nil {
			comp = false
		}
	}

	format := ""
	v = c.Query("format")

	if v != "" {
		if strings.Contains(strings.ToLower(v), "lower") {
			format = "lower"
		} else {
			format = "upper"
		}
	}

	repeatMask := ""
	v = c.Query("mask")

	if v != "" {
		if strings.Contains(strings.ToLower(v), "lower") {
			repeatMask = "lower"
		} else {
			repeatMask = "n"
		}
	}

	return &DNAQuery{
			Rev:        rev,
			Comp:       comp,
			Format:     format,
			RepeatMask: repeatMask},
		nil
}

func GenomesRoute(c *gin.Context) {
	routes.MakeDataResp(c, "", dnadbcache.GetInstance().List())
}

func DNARoute(c *gin.Context) {
	locations, err := ParseLocationsFromPost(c)

	if err != nil {
		c.Error(err)
		return
	}

	assembly := c.Param("assembly")

	query, err := ParseDNAQuery(c)

	if err != nil {
		c.Error(err)
		return
	}

	dnadb, err := dnadbcache.Db(assembly)

	if err != nil {
		c.Error(err)
		return
	}

	seqs := make([]*DNA, 0, len(locations))

	for _, location := range locations {
		seq, err := dnadb.DNA(location,
			query.Format,
			query.RepeatMask,
			query.Rev,
			query.Comp)

		if err != nil {
			c.Error(err)
			return
		}

		seqs = append(seqs, &DNA{Location: location, Seq: seq})
	}

	routes.MakeDataResp(c,
		"",
		&DNAResp{
			Assembly:     assembly,
			Format:       query.Format,
			IsRev:        query.Rev,
			IsComplement: query.Comp,
			Seqs:         seqs})
}
