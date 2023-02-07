package greader_api

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"
)

func (r *Client) renderErr(writer http.ResponseWriter, err error) {
	if err != nil {
		writer.WriteHeader(400)
		writer.Write([]byte(err.Error()))
	}
}

func (r *Client) renderData(writer http.ResponseWriter, data interface{}) {
	if data == nil {
		writer.WriteHeader(200)
		return
	}
	switch data := data.(type) {
	case string:
		writer.Write([]byte(data))
	default:
		bs, _ := json.Marshal(data)
		writer.Header().Set("Content-Type", "application/json; charset=UTF-8")
		writer.Write(bs)
	}
}

func (r *Client) mustJson(req HttpReader) error {
	if res := req.HeaderString("output"); res != "" && res != "json" {
		return errors.New("output must be json")
	}
	return nil
}

// get label tag name
//
// "user/-/label/<name>" -> "<name>"
func getUserLabelName(s string) string {
	return s[strings.Index(s, "/label/")+len("/label/"):]
}

// get tagged item id
//
// tag:google.com,2005:reader/item/<id> -> <id>
func getTaggedItemID(s string) string {
	return s[strings.Index(s, "/item/")+len("/item/"):]
}

func hex16ToInt(s string) (int64, error) {
	return strconv.ParseInt(s, 16, 64)
}
