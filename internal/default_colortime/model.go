package default_colortime

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type DefaultWeekColorTime struct {
	ID             primitive.ObjectID  `bson:"_id" json:"id"`
	OrganizationID string              `bson:"organization_id" json:"organization_id"`
	StartDate      time.Time           `bson:"start_date" json:"start_date"`
	EndDate        time.Time           `bson:"end_date" json:"end_date"`
	ColorTimes     []*DefaultColorTime `bson:"colortimes" json:"colortimes"`
	CreatedBy      string              `bson:"created_by" json:"created_by"`
	CreatedAt      time.Time           `bson:"created_at" json:"created_at"`
	UpdatedAt      time.Time           `bson:"updated_at" json:"updated_at"`
}

type DefaultColorTime struct {
	ID        primitive.ObjectID   `bson:"_id" json:"id"`
	Date      time.Time            `bson:"date" json:"date"`
	TimeSlots []*DefaultColorBlock `bson:"time_slots" json:"time_slots"`
	CreatedAt time.Time            `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time            `bson:"updated_at" json:"updated_at"`
}

type DefaultColorBlock struct {
	BlockID primitive.ObjectID      `json:"block_id" bson:"block_id"`
	Slots   []*DefaultColortimeSlot `json:"slots" bson:"slots"`
}

type DefaultColortimeSlot struct {
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
