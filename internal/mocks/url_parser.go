package mocks

import (
	"github.com/stretchr/testify/mock"
	"github.com/yan-cerqueira-unvoid/url-shortener/internal/parser"
)

type URLParser struct {
	mock.Mock
}

func (m *URLParser) Parse(rawURL string) (*parser.URLParseResult, error) {
	args := m.Called(rawURL)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*parser.URLParseResult), args.Error(1)
}
