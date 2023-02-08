package rss_parser

import "github.com/mmcdole/gofeed"

type Parser struct {
	parser *gofeed.Parser
}

func New() *Parser {
	return &Parser{
		parser: gofeed.NewParser(),
	}
}
