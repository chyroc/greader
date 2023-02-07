package greader_api

import "net/http"

type Client struct {
	s IStore
}

type ClientConfig struct {
	Store IStore
}

func New(config *ClientConfig) *Client {
	return &Client{
		s: config.Store,
	}
}

func (r *Client) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	path := request.URL.Path
	req := WrapStdHttpRequest(request)

	switch path {

	// auth
	case "/accounts/ClientLogin":
		r.Auth(request.Context(), req, writer)
	case "/reader/api/0/token":
		r.Token(request.Context(), req, writer)

	// tag
	case "/reader/api/0/tag/list":
		r.TagList(request.Context(), req, writer)
	case "/reader/api/0/disable-tag":
		r.TagDelete(request.Context(), req, writer)
	case "/reader/api/0/rename-tag":
		r.TagRename(request.Context(), req, writer)

	// sub
	case "/reader/api/0/subscription/list":
		r.ListSubscription(request.Context(), req, writer)
	case "/reader/api/0/subscription/quickadd":
		r.AddSubscription(request.Context(), req, writer)
	case "/reader/api/0/subscription/edit":
		r.EditSubscription(request.Context(), req, writer)

	// content
	case "/reader/api/0/stream/items/ids":
		r.ListItemIDs(request.Context(), req, writer)
	case "/reader/api/0/stream/items/contents":
		r.LoadItem(request.Context(), req, writer)
	}
}

type HttpReader interface {
	FormList(key string) []string
	FormString(key string) string
	HeaderString(key string) string
	QueryString(key string) string
}
