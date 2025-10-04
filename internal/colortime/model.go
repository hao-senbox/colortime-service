package colortime

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ColorTime struct {
	ID             primitive.ObjectID `bson:"_id" json:"id"`
	OrganizationID string             `bson:"organization_id" json:"organization_id"`
	UserID         string             `bson:"user_id" json:"user_id"`
	Tracking       string             `bson:"tracking" json:"tracking"`
	UseCount       int                `bson:"use_count" json:"use_count"`
	Title          string             `bson:"title" json:"title"`
	Date           time.Time          `bson:"date" json:"date"`
	Time           time.Time          `bson:"time" json:"time"`
	Duration       int                `bson:"duration" json:"duration"`
	Color          string             `bson:"color" json:"color"`
	Note           string             `bson:"note" json:"note"`
	ProductID      primitive.ObjectID `bson:"product_id" json:"product_id"`
	CreatedAt      time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt      time.Time          `bson:"updated_at" json:"updated_at"`
}
