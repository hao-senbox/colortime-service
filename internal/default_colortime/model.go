package default_colortime

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type DefaultDayColorTime struct {
	ID             primitive.ObjectID   `bson:"_id" json:"id"`
	OrganizationID string               `bson:"organization_id" json:"organization_id"`
	Date           time.Time            `bson:"date" json:"date"`
	TimeSlots      []*DefaultColorBlock `bson:"time_slots" json:"time_slots"`
	CreatedBy      string               `bson:"created_by" json:"created_by"`
	CreatedAt      time.Time            `bson:"created_at" json:"created_at"`
	UpdatedAt      time.Time            `bson:"updated_at" json:"updated_at"`

	// Repeat configuration
	IsBaseTemplate bool                `bson:"is_base_template" json:"is_base_template"` // true if this is the base template for repeating
	RepeatType     string              `bson:"repeat_type" json:"repeat_type"`           // "none", "daily", "weekly", "monthly", "custom"
	RepeatUntil    *time.Time          `bson:"repeat_until" json:"repeat_until"`         // when to stop repeating
	RepeatInterval int                 `bson:"repeat_interval" json:"repeat_interval"`   // repeat every N days/weeks/months (default 1)
	RepeatDays     []int               `bson:"repeat_days" json:"repeat_days"`           // for weekly: [0=Sun,1=Mon,...,6=Sat], for custom dates
	BaseTemplateID *primitive.ObjectID `bson:"base_template_id" json:"base_template_id"` // reference to base template if this is a generated day
}

// // Legacy types kept for backward compatibility
// type DefaultWeekColorTime struct {
// 	ID             primitive.ObjectID  `bson:"_id" json:"id"`
// 	OrganizationID string              `bson:"organization_id" json:"organization_id"`
// 	StartDate      time.Time           `bson:"start_date" json:"start_date"`
// 	EndDate        time.Time           `bson:"end_date" json:"end_date"`
// 	ColorTimes     []*DefaultColorTime `bson:"colortimes" json:"colortimes"`
// 	CreatedBy      string              `bson:"created_by" json:"created_by"`
// 	CreatedAt      time.Time           `bson:"created_at" json:"created_at"`
// 	UpdatedAt      time.Time           `bson:"updated_at" json:"updated_at"`

// 	// Repeat configuration
// 	IsBaseTemplate bool                `bson:"is_base_template" json:"is_base_template"` // true if this is the base template for repeating
// 	RepeatType     string              `bson:"repeat_type" json:"repeat_type"`           // "none", "daily", "weekly", "monthly", "custom"
// 	RepeatUntil    *time.Time          `bson:"repeat_until" json:"repeat_until"`         // when to stop repeating
// 	RepeatInterval int                 `bson:"repeat_interval" json:"repeat_interval"`   // repeat every N days/weeks/months (default 1)
// 	RepeatDays     []int               `bson:"repeat_days" json:"repeat_days"`           // for weekly: [0=Sun,1=Mon,...,6=Sat], for custom dates
// 	BaseTemplateID *primitive.ObjectID `bson:"base_template_id" json:"base_template_id"` // reference to base template if this is a generated week
// }

// type DefaultColorTime struct {
// 	ID        primitive.ObjectID   `bson:"_id" json:"id"`
// 	Date      time.Time            `bson:"date" json:"date"`
// 	TimeSlots []*DefaultColorBlock `bson:"time_slots" json:"time_slots"`
// 	CreatedAt time.Time            `bson:"created_at" json:"created_at"`
// 	UpdatedAt time.Time            `bson:"updated_at" json:"updated_at"`
// }

type DefaultColorBlock struct {
	BlockID primitive.ObjectID      `json:"block_id" bson:"block_id"`
	Slots   []*DefaultColortimeSlot `json:"slots" bson:"slots"`
}

type DefaultColortimeSlot struct {
	SlotID                primitive.ObjectID              `json:"slot_id" bson:"slot_id"`
	Sessions              int                             `json:"sessions" bson:"sessions"`
	Title                 string                          `json:"title" bson:"title"`
	ColorTimeSlotLanguage []*DefaultColorTimeSlotLanguage `json:"color_time_slot_language" bson:"color_time_slot_language"`
	StartTime             time.Time                       `json:"start_time" bson:"start_time"`
	EndTime               time.Time                       `json:"end_time" bson:"end_time"`
	Duration              int                             `json:"duration" bson:"duration"`
	Color                 string                          `json:"color" bson:"color"`
	Note                  string                          `json:"note" bson:"note"`
	CreatedAt             time.Time                       `json:"created_at" bson:"created_at"`
	UpdatedAt             time.Time                       `json:"updated_at" bson:"updated_at"`
}

type DefaultColorTimeSlotLanguage struct {
	LanguageID int    `json:"language_id" bson:"language_id"`
	Title      string `json:"title" bson:"title"`
}
