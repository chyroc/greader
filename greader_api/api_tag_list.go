package greader_api

import (
	"context"
	"net/http"
)

func (r *Client) TagList(ctx context.Context, req HttpReader, writer http.ResponseWriter) {
	res, err := r.tagList(ctx, req)
	if err != nil {
		r.renderErr(ctx, writer, err)
	} else {
		r.renderData(ctx, writer, res)
	}
}

func (r *Client) tagList(ctx context.Context, req HttpReader) (*tagList, error) {
	if err := r.mustJson(req); err != nil {
		return nil, err
	}
	username := getContextUsername(ctx)
	r.log.Info(ctx, "[TagList] username=%s", username)

	tagNames, err := r.s.ListTag(ctx, username)
	if err != nil {
		return nil, err
	}

	if len(tagNames) == 0 {
		tagNames = []string{"default"}
	}

	return &tagList{
		Tags: buildTads(tagNames),
	}, nil
}
