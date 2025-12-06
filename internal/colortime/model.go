package colortime

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type WeekColorTime struct {
	ID             primitive.ObjectID `bson:"_id" json:"id"`
	OrganizationID string             `bson:"organization_id" json:"organization_id"`
	Owner          *Owner             `json:"owner" bson:"owner"`
	StartDate      time.Time          `bson:"start_date" json:"start_date"`
	EndDate        time.Time          `bson:"end_date" json:"end_date"`
	TopicID        *string            `bson:"topic_id" json:"topic_id"`
	ColorTimes     []*ColorTime       `bson:"colortimes" json:"colortimes"`
	CreatedBy      string             `bson:"created_by" json:"created_by"`
	CreatedAt      time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt      time.Time          `bson:"updated_at" json:"updated_at"`
}

type ColorTime struct {
	ID        primitive.ObjectID `bson:"_id" json:"id"`
	Date      time.Time          `bson:"date" json:"date"`
	TopicID   *string            `bson:"topic_id" json:"topic_id"`
	TimeSlots []*ColorBlock      `bson:"time_slots" json:"time_slots"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at" json:"updated_at"`
}

type ColorBlock struct {
	BlockID    primitive.ObjectID  `json:"block_id" bson:"block_id"`
	BlockIDOld *primitive.ObjectID `json:"block_id_old" bson:"block_id_old"`
	Slots      []*ColortimeSlot    `json:"slots" bson:"slots"`
}

type ColortimeSlot struct {
	SlotID                primitive.ObjectID       `json:"slot_id" bson:"slot_id"`
	SlotIDOld             *primitive.ObjectID      `json:"slot_id_old" bson:"slot_id_old"`
	Sessions              int                      `json:"sessions" bson:"sessions"`
	Title                 string                   `json:"title" bson:"title"`
	ColorTimeSlotLanguage []*ColorTimeSlotLanguage `json:"color_time_slot_language" bson:"color_time_slot_language"`
	Tracking              string                   `json:"tracking" bson:"tracking"`
	UseCount              int                      `json:"use_count" bson:"use_count"`
	StartTime             time.Time                `json:"start_time" bson:"start_time"`
	EndTime               time.Time                `json:"end_time" bson:"end_time"`
	Duration              int                      `json:"duration" bson:"duration"`
	Color                 string                   `json:"color" bson:"color"`
	Note                  string                   `json:"note" bson:"note"`
	ProductID             *string                  `json:"product_id" bson:"product_id"`
	CreatedAt             time.Time                `json:"created_at" bson:"created_at"`
	UpdatedAt             time.Time                `json:"updated_at" bson:"updated_at"`
}

type ColorTimeSlotLanguage struct {
	LanguageID int    `json:"language_id" bson:"language_id"`
	Title      string `json:"title" bson:"title"`
}

type Owner struct {
	OwnerID   string `json:"owner_id" bson:"owner_id"`
	OwnerRole string `json:"owner_role" bson:"owner_role"`
}
