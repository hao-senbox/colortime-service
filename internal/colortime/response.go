package colortime

import (
	"colortime-service/internal/language"
	"colortime-service/internal/product"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ColorTimeResponse struct {
	ID               primitive.ObjectID                 `bson:"_id" json:"id"`
	UserID           string                             `bson:"user_id" json:"user_id"`
	Title            string                             `bson:"title" json:"title"`
	Date             time.Time                          `bson:"date" json:"date"`
	Time             time.Time                          `bson:"time" json:"time"`
	Duration         int                                `bson:"duration" json:"duration"`
	Color            string                             `bson:"color" json:"color"`
	Note             string                             `bson:"note" json:"note"`
	Product          *product.Product                   `bson:"product" json:"product"`
	MessageLanguages []language.MessageLanguageResponse `json:"message_languages"`
	CreatedAt        time.Time                          `bson:"created_at" json:"created_at"`
	UpdatedAt        time.Time                          `bson:"updated_at" json:"updated_at"`
}
