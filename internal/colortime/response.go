package colortime

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ColorTimeResponse struct {
	ID        primitive.ObjectID `bson:"_id" json:"id"`
	Date      time.Time          `bson:"date" json:"date"`
	Topic     Topic              `bson:"topic" json:"topic"`
	TimeSlots []*TemplateSlot    `bson:"time_slots" json:"time_slots"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at" json:"updated_at"`
}

type TopicToColorTimeWeekResponse struct {
	ID             primitive.ObjectID   `bson:"_id" json:"id"`
	OrganizationID string               `bson:"organization_id" json:"organization_id"`
	Owner          *Owner               `json:"owner" bson:"owner"`
	StartDate      time.Time            `bson:"start_date" json:"start_date"`
	EndDate        time.Time            `bson:"end_date" json:"end_date"`
	Topic          Topic                `bson:"topic" json:"topic"`
	ColorTimes     []*ColorTimeResponse `bson:"colortimes" json:"colortimes"`
	CreatedBy      string               `bson:"created_by" json:"created_by"`
	CreatedAt      time.Time            `bson:"created_at" json:"created_at"`
	UpdatedAt      time.Time            `bson:"updated_at" json:"updated_at"`
}

type Topic struct {
	ID   string `json:"id" bson:"_id"`
	Name string `json:"name" bson:"name"`
}
