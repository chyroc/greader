package greader_api

import (
	"context"
	"net/http"
)

type Client struct {
	s           IGReaderStore
	log         ILogger
	fetchLogger ILogger
}

type ClientConfig struct {
	Store       IGReaderStore
	Logger      ILogger
	FetchLogger ILogger
}

func New(config *ClientConfig) *Client {
	return &Client{
		s:           config.Store,
		log:         config.Logger,
		fetchLogger: config.FetchLogger,
	}
}

func (r *Client) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	method := request.Method
	path := request.URL.Path
	req := WrapStdHttpRequest(request)
	ctx := setContextPath(request.Context(), path)

	r.log.Info(ctx, "[request_] method=%s, path=%s", method, path)

	for methodPath, handler := range r.Routers() {
		if methodPath[1] == path {
			handler(ctx, req, writer)
			return
		}
	}
}

func (r *Client) Routers() map[[2]string]func(ctx context.Context, reader HttpReader, writer http.ResponseWriter) {
	return map[[2]string]func(ctx context.Context, reader HttpReader, writer http.ResponseWriter){
		// auth
		{"POST", "/accounts/ClientLogin"}: r.Auth, // POST
		{"GET", "/reader/api/0/token"}:    r.Token,
		// tag
		{"GET", "/reader/api/0/tag/list"}:     r.TagList,   // GET
		{"POST", "/reader/api/0/disable-tag"}: r.TagDelete, // POST
		{"POST", "/reader/api/0/rename-tag"}:  r.TagRename, // POST
		// sub
		{"GET", "/reader/api/0/subscription/list"}:      r.ListSubscription,       // GET
		{"POST", "/reader/api/0/subscription/quickadd"}: r.AddSubscription,        // POST
		{"POST", "/reader/api/0/subscription/edit"}:     r.EditSubscription,       // POST
		{"POST", "/reader/api/0/edit-tag"}:              r.EditSubscriptionStatus, // POST
		// content
		{"GET", "/reader/api/0/stream/items/ids"}:       r.ListItemIDs, // GET
		{"POST", "/reader/api/0/stream/items/contents"}: r.LoadItem,    // POST
	}
}

type HttpReader interface {
	FormList(key string) []string
	FormString(key string) string
	HeaderString(key string) string
	QueryString(key string) string
}
