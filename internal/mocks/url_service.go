package mocks

import (
	"github.com/stretchr/testify/mock"
	"github.com/yan-cerqueira-unvoid/url-shortener/internal/models"
)

type URLService struct {
	mock.Mock
}

func (m *URLService) ShortenURL(originalURL string, customCode string) (*models.URL, error) {
	args := m.Called(originalURL, customCode)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*models.URL), args.Error(1)
}

func (m *URLService) GetURL(shortCode string) (*models.URL, error) {
	args := m.Called(shortCode)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*models.URL), args.Error(1)
}
