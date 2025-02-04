package utilroutes

import (
	"bytes"
	b64 "encoding/base64"
	"fmt"

	"github.com/antonybholmes/go-edb-server-gin/routes"
	"github.com/antonybholmes/go-sys"
	"github.com/gin-gonic/gin"
)

type XlsxReq struct {
	Sheet    string `json:"sheet"`
	Headers  int    `json:"headers"`
	Indexes  int    `json:"indexes"`
	SkipRows int    `json:"skipRows"`
	Xlsx     string `json:"b64xlsx"`
}

type XlsxResp struct {
	Table *sys.Table `json:"table"`
}

type XlsxSheetsResp struct {
	Sheets []string `json:"sheets"`
}

func makeXlsxReader(data string) (*bytes.Reader, error) {
	xlsxb, err := b64.StdEncoding.DecodeString(data)

	if err != nil {
		return nil, err
	}

	reader := bytes.NewReader(xlsxb)

	return reader, nil
}

func XlsxSheetsRoute(c *gin.Context) {

	var req XlsxReq

	err := c.Bind(&req)

	if err != nil {
		c.Error(err)
		return
	}

	//log.Debug().Msgf("m:%s", req.Xlsx)

	reader, err := makeXlsxReader(req.Xlsx)

	if err != nil {
		c.Error(err)
		return
	}

	sheets, err := sys.XlsxSheetNames(reader)

	if err != nil {
		c.Error(err)
		return
	}

	resp := XlsxSheetsResp{Sheets: sheets}

	routes.MakeDataResp(c, "", resp)
}

func XlsxToRoute(c *gin.Context) {

	format := c.Param("format")

	if format != "json" {
		c.Error(fmt.Errorf("unsupported format: %s", format))
		return
	}

	var req XlsxReq

	err := c.Bind(&req)

	if err != nil {
		c.Error(err)
		return
	}

	reader, err := makeXlsxReader(req.Xlsx)

	if err != nil {
		c.Error(err)
		return
	}

	table, err := sys.XlsxToJson(reader, req.Sheet, req.Indexes, req.Headers, req.SkipRows)

	if err != nil {
		c.Error(err)
		return
	}

	resp := XlsxResp{Table: table}

	routes.MakeDataResp(c, "", resp)
}
