package greader_api

import "net/http"

type Client struct {
	s   IGReaderStore
	log ILogger
}

type ClientConfig struct {
	Store  IGReaderStore
	Logger ILogger
}

func New(config *ClientConfig) *Client {
	return &Client{
		s:   config.Store,
		log: NewDefaultLogger(),
	}
}

func (r *Client) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	method := request.Method
	path := request.URL.Path
	req := WrapStdHttpRequest(request)
	ctx := setContextPath(r.setAuth(request.Context(), req), path)
	r.log.Info(ctx, "[request_] method=%s, path=%s", method, path)

	if err := request.ParseForm(); err != nil {
		r.renderErr(ctx, writer, err)
		return
	}

	switch path {
	// auth
	case "/accounts/ClientLogin":
		r.Auth(ctx, req, writer)
	case "/reader/api/0/token":
		r.Token(ctx, req, writer)

	// tag
	case "/reader/api/0/tag/list":
		r.TagList(ctx, req, writer)
	case "/reader/api/0/disable-tag":
		r.TagDelete(ctx, req, writer)
	case "/reader/api/0/rename-tag":
		r.TagRename(ctx, req, writer)

	// sub
	case "/reader/api/0/subscription/list":
		r.ListSubscription(ctx, req, writer)
	case "/reader/api/0/subscription/quickadd":
		r.AddSubscription(ctx, req, writer)
	case "/reader/api/0/subscription/edit":
		r.EditSubscription(ctx, req, writer)
	case "/reader/api/0/edit-tag":
		r.EditSubscriptionStatus(ctx, req, writer)

	// content
	case "/reader/api/0/stream/items/ids":
		r.ListItemIDs(ctx, req, writer)
	case "/reader/api/0/stream/items/contents":
		r.LoadItem(ctx, req, writer)
	}
}

type HttpReader interface {
	FormList(key string) []string
	FormString(key string) string
	HeaderString(key string) string
	QueryString(key string) string
}
