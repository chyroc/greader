package server

import (
	"github.com/gin-gonic/gin"

	"github.com/chyroc/greader/greader_api"
	"github.com/chyroc/greader/mysql_backend"
)

type App struct {
	gin     *gin.Engine
	logger  greader_api.ILogger
	backend greader_api.IGReaderBackend
	greader *greader_api.GReader
}

func (r *App) Start(addr string) error {
	return r.gin.Run(addr)
}

func New(dsn string) (*App, error) {
	// init app
	// app := new(App)

	// init logger
	logger := greader_api.NewDefaultLogger()

	// init gin
	ginIns := gin.New()
	ginIns.Use(Log(logger))

	// backend
	backend, err := mysql_backend.New(dsn, logger)
	if err != nil {
		return nil, err
	}

	// greader
	greaderIns := greader_api.New(&greader_api.ClientConfig{
		Backend:     backend,
		Logger:      logger,
		FetchLogger: logger,
	})
	greaderIns.FetchRssBackend()

	registerAPiRoute(ginIns, greaderIns)

	return &App{
		gin:     ginIns,
		logger:  logger,
		backend: backend,
		greader: greaderIns,
	}, nil
}
