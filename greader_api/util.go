package greader_api

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"
)

func (r *Client) renderErr(ctx context.Context, writer http.ResponseWriter, err error) {
	r.log.Error(ctx, "[response] [err] path=%s, err=%s", getContextPath(ctx), err)

	if err != nil {
		writer.WriteHeader(400)
		writer.Write([]byte(err.Error()))
	}
}

func (r *Client) renderData(ctx context.Context, writer http.ResponseWriter, data interface{}) {
	if data == nil {
		r.log.Info(ctx, "[response] [ok] path=%s, resp=nil", getContextPath(ctx))
		writer.WriteHeader(200)
		return
	}
	switch data := data.(type) {
	case string:
		r.log.Info(ctx, "[response] [ok] path=%s, resp=%s", getContextPath(ctx), data)
		writer.Write([]byte(data))
	default:
		bs, _ := json.Marshal(data)
		r.log.Info(ctx, "[response] [ok] path=%s, resp=%s", getContextPath(ctx), string(bs))
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
	if !strings.HasPrefix(s, "user/-/label/") {
		return s
	}
	return s[strings.Index(s, "/label/")+len("/label/"):]
}

func buildUserLabelName(s string) {
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
