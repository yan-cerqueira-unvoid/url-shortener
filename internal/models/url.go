package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type URL struct {
	ID          primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	OriginalURL string             `json:"original_url" bson:"original_url"`
	ShortCode   string             `json:"short_code" bson:"short_code"`
	Clicks      int64              `json:"clicks" bson:"clicks"`
	ExpiresAt   time.Time          `json:"expires_at" bson:"expires_at"`
	CreatedBy   string             `json:"created_by" bson:"created_by"`
	CreatedAt   time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt   time.Time          `json:"updated_at" bson:"updated_at"`
}
