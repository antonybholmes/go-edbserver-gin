package utils

import (
	"bytes"
	"crypto/rand"
	b64 "encoding/base64"
	"fmt"
	"math/big"
	"strconv"

	"github.com/antonybholmes/go-sys"
	"github.com/antonybholmes/go-web"
	"github.com/antonybholmes/go-web/auth"
	"github.com/gin-gonic/gin"
)

const Alphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

var (
	ErrLengthCannotBeZero = fmt.Errorf("length cannot be zero")
)

type XlsxReq struct {
	Sheet          string `json:"sheet"`
	Headers        int    `json:"headers"`
	Indexes        int    `json:"indexes"`
	SkipRows       int    `json:"skipRows"`
	TrimWhitespace bool   `json:"trimWhitespace"`
	Xlsx           string `json:"b64xlsx"`
}

type XlsxResp struct {
	Table *sys.Table `json:"table"`
}

type XlsxSheetsResp struct {
	Sheets []string `json:"sheets"`
}

type HashResp struct {
	Password string `json:"password"`
	Hash     string `json:"hash"`
}

type KeyResp struct {
	Key    string `json:"key"`
	Length int    `json:"length"`
}

func makeXlsxReader(b64data string) (*bytes.Reader, error) {
	xlsxb, err := b64.StdEncoding.DecodeString(b64data)

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

	web.MakeDataResp(c, "", resp)
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

	table, err := sys.XlsxToJson(reader,
		req.Sheet,
		req.Indexes,
		req.Headers,
		req.SkipRows,
		req.TrimWhitespace)

	if err != nil {
		c.Error(err)
		return
	}

	resp := XlsxResp{Table: table}

	web.MakeDataResp(c, "", resp)
}

// generateRandomString generates a random string of specified length from the letters set.
func generateRandomString(length int) (string, error) {
	b := make([]byte, length)
	for i := range b {
		// Generate a random index
		index, err := rand.Int(rand.Reader, big.NewInt(int64(len(Alphabet))))

		if err != nil {
			return "", err
		}

		b[i] = Alphabet[index.Int64()]
	}
	return string(b), nil
}

func HashedPasswordRoute(c *gin.Context) {

	password := c.Query("password")

	if len(password) == 0 {
		web.BadReqResp(c, auth.ErrPasswordDoesNotMeetCriteria)
		return
	}

	hash := auth.HashPassword(password)

	ret := HashResp{Password: password, Hash: hash}

	web.MakeDataResp(c, "", ret)
}

func RandomKeyRoute(c *gin.Context) {

	l, err := strconv.Atoi(c.Query("l"))

	if err != nil || l < 1 {
		web.BadReqResp(c, ErrLengthCannotBeZero)
		return
	}

	key, err := generateRandomString(l)

	if err != nil {
		c.Error(err)
		return
	}

	ret := KeyResp{Key: key, Length: l}

	web.MakeDataResp(c, "", ret)
}

func UUIDv7Route(c *gin.Context) {

	uuid, err := sys.Uuidv7()

	if err != nil {
		web.BadReqResp(c, err)
		return
	}

	web.MakeDataResp(c, "UUID v7", uuid)
}
