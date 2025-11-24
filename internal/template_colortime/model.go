package templatecolortime

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TemplateColorTime struct {
	ID             primitive.ObjectID   `bson:"_id" json:"id"`
	Date           string               `bson:"date" json:"date"`
	OrganizationID string               `bson:"organization_id" json:"organization_id"`
	TermID         string               `bson:"term_id" json:"term_id"`
	ColorTimes     []*ColorTimeTemplate `bson:"color_times" json:"color_times"`
	CreatedBy      string               `bson:"created_by" json:"created_by"`
	CreatedAt      time.Time            `bson:"created_at" json:"created_at"`
	UpdatedAt      time.Time            `bson:"updated_at" json:"updated_at"`
}

type ColorTimeTemplate struct {
	BlockID primitive.ObjectID `json:"block_id" bson:"block_id"`
	Slots   []*ColortimeSlot   `json:"slots" bson:"slots"`
}

type ColortimeSlot struct {
	SlotID    primitive.ObjectID `json:"slot_id" bson:"slot_id"`
	Sessions  int                `json:"sessions" bson:"sessions"`
	Title     string             `json:"title" bson:"title"`
	StartTime time.Time          `json:"start_time" bson:"start_time"`
	EndTime   time.Time          `json:"end_time" bson:"end_time"`
	Duration  int                `json:"duration" bson:"duration"`
	Color     string             `json:"color" bson:"color"`
	Note      string             `json:"note" bson:"note"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time          `json:"updated_at" bson:"updated_at"`
}
