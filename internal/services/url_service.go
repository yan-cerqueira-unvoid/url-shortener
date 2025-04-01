package services

import (
	"context"
	"crypto/md5"
	"encoding/base64"
	"errors"
	"fmt"
	"math/rand"
	"time"

	"github.com/yan-cerqueira-unvoid/url-shortener/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type URLService struct {
	db         *mongo.Database
	ctx        *context.Context
	collection *mongo.Collection
}

func NewURLService(db *mongo.Database, ctx *context.Context) *URLService {
	collection := db.Collection("urls")
	indexModel := mongo.IndexModel{
		Keys:    bson.M{"short_code": 1},
		Options: options.Index().SetUnique(true),
	}
	_, err := collection.Indexes().CreateOne(*ctx, indexModel)
	if err != nil {
		panic(fmt.Sprintf("Failed to create index: %v", err))
	}

	return &URLService{
		db:         db,
		ctx:        ctx,
		collection: collection,
	}
}

func (service *URLService) ShortenURL(originalURL string, customCode string) (*models.URL, error) {

	var existingURL models.URL

	err := service.collection.FindOne(*service.ctx, bson.M{"original_url": originalURL}).Decode(&existingURL)

	if err == nil {
		return &existingURL, nil
	} else if err != mongo.ErrNoDocuments {
		return nil, err
	}

	var shortCode string

	if customCode != "" {
		var existingCustom models.URL

		err := service.collection.FindOne(*service.ctx, bson.M{"short_code": customCode}).Decode(&existingCustom)
		if err == nil {
			return nil, errors.New("custom short code already in use")
		} else if err != mongo.ErrNoDocuments {
			return nil, err
		}

		shortCode = customCode
	} else {
		shortCode = service.generateShortCode(originalURL)

		for {
			var existingCode models.URL

			err := service.collection.FindOne(*service.ctx, bson.M{"short_code": shortCode}).Decode(&existingCode)
			if err == mongo.ErrNoDocuments {
				break
			}

			shortCode = service.generateShortCode(originalURL + fmt.Sprint(rand.Intn(1000)))
		}
	}

	now := time.Now()
	url := models.URL{
		OriginalURL: originalURL,
		ShortCode:   shortCode,
		Clicks:      0,
		ExpiresAt:   now.AddDate(1, 0, 0), // Expires in 1 year
		CreatedBy:   "anonymous",          // Would be set from auth in a real app
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	_, err = service.collection.InsertOne(*service.ctx, url)
	if err != nil {
		return nil, err
	}

	return &url, nil
}

func (service *URLService) GetURL(shortCode string) (*models.URL, error) {
	var url models.URL

	err := service.collection.FindOne(*service.ctx, bson.M{"short_code": shortCode}).Decode(&url)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("short URL not found")
		}

		return nil, err
	}

	if time.Now().After(url.ExpiresAt) {
		return nil, errors.New("URL has expired")
	}

	_, err = service.collection.UpdateOne(
		*service.ctx,
		bson.M{"short_code": shortCode},
		bson.M{"$inc": bson.M{"clicks": 1}, "$set": bson.M{"updated_at": time.Now()}},
	)

	if err != nil {
		return nil, err
	}

	url.Clicks++

	return &url, nil
}

func (service *URLService) generateShortCode(url string) string {
	hasher := md5.New()

	hasher.Write([]byte(url + time.Now().String()))
	hash := base64.URLEncoding.EncodeToString(hasher.Sum(nil))

	return hash[:6]
}
