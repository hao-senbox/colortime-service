package templatecolortime

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TemplateColorTimeResponse struct {
	ID             primitive.ObjectID   `bson:"_id" json:"id"`
	Date           string               `bson:"date" json:"date"`
	OrganizationID string               `bson:"organization_id" json:"organization_id"`
	TermID         string               `bson:"term_id" json:"term_id"`
	ColorTimes     []*ColorTimeTemplate `bson:"color_times" json:"color_times"`
	CreatedBlockID *primitive.ObjectID  `bson:"created_block_id" json:"created_block_id"`
	CreatedBy      string               `bson:"created_by" json:"created_by"`
	CreatedAt      time.Time            `bson:"created_at" json:"created_at"`
	UpdatedAt      time.Time            `bson:"updated_at" json:"updated_at"`
}
