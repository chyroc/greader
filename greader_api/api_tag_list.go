package greader_api

import (
	"context"
	"net/http"
)

func (r *Client) TagList(ctx context.Context, req HttpReader, writer http.ResponseWriter) {
	res, err := r.tagList(ctx, req)
	if err != nil {
		r.renderErr(writer, err)
	} else {
		r.renderData(writer, res)
	}
}

func (r *Client) tagList(ctx context.Context, req HttpReader) (*tagList, error) {
	if err := r.mustJson(req); err != nil {
		return nil, err
	}

	tagNames, err := r.s.ListTag(ctx)
	if err != nil {
		return nil, err
	}

	return &tagList{
		Tags: buildTads(tagNames),
	}, nil
}
