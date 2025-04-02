package services

import (
	"context"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type URLServiceIntegrationTestSuite struct {
	suite.Suite
	db         *mongo.Database
	urlService *URLService
}

func getMongoURI() string {
	if uri := os.Getenv("MONGO_URI"); uri != "" {
		return uri
	}

	hostname, err := os.Hostname()
	if err == nil && strings.Contains(hostname, "test") {
		return "mongodb://mongodb:27017"
	}

	return "mongodb://localhost:27018" // Match the port in docker-compose.test.yml
}

func (suite *URLServiceIntegrationTestSuite) TestShortenURL() {
	originalURL := "https://example.com/test"
	url, err := suite.urlService.ShortenURL(originalURL, "")

	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), url)
	assert.Equal(suite.T(), originalURL, url.OriginalURL)
	assert.NotEmpty(suite.T(), url.ShortCode)
	assert.Equal(suite.T(), int64(0), url.Clicks)
	assert.NotEmpty(suite.T(), url.ExpiresAt)

	url2, err := suite.urlService.ShortenURL(originalURL, "")

	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), url.ShortCode, url2.ShortCode)
}

func (suite *URLServiceIntegrationTestSuite) TestShortenURLWithCustomCode() {
	originalURL := "https://example.com/custom"
	customCode := "mycode123"
	url, err := suite.urlService.ShortenURL(originalURL, customCode)

	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), url)
	assert.Equal(suite.T(), originalURL, url.OriginalURL)
	assert.Equal(suite.T(), customCode, url.ShortCode)

	anotherURL := "https://example.com/another"
	_, err = suite.urlService.ShortenURL(anotherURL, customCode)

	assert.Error(suite.T(), err)
	assert.Contains(suite.T(), err.Error(), "already in use")
}

func (suite *URLServiceIntegrationTestSuite) TestGetURL() {
	originalURL := "https://example.com/geturl"
	customCode := "gettest123"
	createdURL, err := suite.urlService.ShortenURL(originalURL, customCode)

	assert.NoError(suite.T(), err)

	time.Sleep(300 * time.Millisecond)

	retrievedURL, err := suite.urlService.GetURL(customCode)

	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), createdURL.OriginalURL, retrievedURL.OriginalURL)
	assert.Equal(suite.T(), createdURL.ShortCode, retrievedURL.ShortCode)
	assert.Equal(suite.T(), int64(1), retrievedURL.Clicks)

	_, err = suite.urlService.GetURL("nonexistent123")

	assert.Error(suite.T(), err)
	assert.Contains(suite.T(), err.Error(), "not found")
}

func (suite *URLServiceIntegrationTestSuite) TestExpiredURL() {
	originalURL := "https://example.com/expiring"
	customCode := "expire123"

	_, err := suite.urlService.ShortenURL(originalURL, customCode)
	assert.NoError(suite.T(), err)

	time.Sleep(100 * time.Millisecond)

	updateCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err = suite.db.Collection("urls").UpdateOne(
		updateCtx,
		bson.M{"short_code": customCode},
		bson.D{
			{Key: "$set", Value: bson.D{{Key: "expires_at", Value: time.Now().Add(-1 * time.Hour)}}},
		},
	)
	assert.NoError(suite.T(), err)

	time.Sleep(100 * time.Millisecond)

	_, err = suite.urlService.GetURL(customCode)

	assert.Error(suite.T(), err)
	assert.Contains(suite.T(), err.Error(), "expired")
}

func TestURLServiceIntegrationTestSuite(t *testing.T) {
	if os.Getenv("SKIP_INTEGRATION") == "true" || testing.Short() {
		t.Skip("Skipping integration tests")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(getMongoURI()))
	if err != nil || client.Ping(ctx, nil) != nil {
		t.Skip("MongoDB not available, skipping integration tests")
	}
	client.Disconnect(ctx)

	suite.Run(t, new(URLServiceIntegrationTestSuite))
}
