package greader_api

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"
)

func (r *GReader) renderErr(ctx context.Context, writer http.ResponseWriter, err error) {
	if err != nil {
		writer.WriteHeader(400)
		writer.Write([]byte(err.Error()))
	}
}

func (r *GReader) renderData(ctx context.Context, writer http.ResponseWriter, data interface{}) {
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

func (r *GReader) mustJson(req HttpReader) error {
	if res := req.QueryString("output"); res != "" && res != "json" {
		panic("output must be json")
		return errors.New("output must be json")
	}
	return nil
}

// get label tag name
//
// "user/-/label/<name>" -> "<name>"
func getUserLabelName(s string) string {
	if !strings.HasPrefix(s, "user/-/label/") {
		return s
	}
	return s[strings.Index(s, "/label/")+len("/label/"):]
}

func buildUserLabelName(s string) string {
	if strings.HasPrefix(s, "user/-/label/") {
		return s
	}
	return "user/-/label/" + s
}

func getFeedID(feedID string) string {
	if strings.HasPrefix(feedID, "feed/") {
		return feedID[len("feed/"):]
	}
	return feedID
}

// get tagged item id
//
// tag:google.com,2005:reader/item/<id> -> <id>
func getTaggedItemHexID(s string) string {
	return s[strings.Index(s, "/item/")+len("/item/"):]
}

func hex16ToInt(s string) (int64, error) {
	return strconv.ParseInt(s, 16, 64)
}

func intToHex16(i int64) string {
	return strconv.FormatInt(i, 16)
}
