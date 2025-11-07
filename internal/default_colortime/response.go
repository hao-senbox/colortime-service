package default_colortime

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type DefaultColorTimeResponse struct {
	ID        primitive.ObjectID   `bson:"_id" json:"id"`
	Date      time.Time            `bson:"date" json:"date"`
	TimeSlots []*DefaultColorBlock `bson:"time_slots" json:"time_slots"`
	CreatedAt time.Time            `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time            `bson:"updated_at" json:"updated_at"`
}

type TopicToDefaultColorTimeWeekResponse struct {
	ID             primitive.ObjectID          `bson:"_id" json:"id"`
	OrganizationID string                      `bson:"organization_id" json:"organization_id"`
	StartDate      time.Time                   `bson:"start_date" json:"start_date"`
	EndDate        time.Time                   `bson:"end_date" json:"end_date"`
	ColorTimes     []*DefaultColorTimeResponse `bson:"colortimes" json:"colortimes"`
	CreatedBy      string                      `bson:"created_by" json:"created_by"`
	CreatedAt      time.Time                   `bson:"created_at" json:"created_at"`
	UpdatedAt      time.Time                   `bson:"updated_at" json:"updated_at"`
}

