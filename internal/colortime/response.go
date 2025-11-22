package colortime

import (
	"colortime-service/internal/user"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ColorTimeResponse struct {
	ID        primitive.ObjectID `bson:"_id" json:"id"`
	Date      time.Time          `bson:"date" json:"date"`
	Topic     Topic              `bson:"topic" json:"topic"`
	TimeSlots []*BlockResponse   `bson:"time_slots" json:"time_slots"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at" json:"updated_at"`
}

type TopicToColorTimeWeekResponse struct {
	ID             primitive.ObjectID   `bson:"_id" json:"id"`
	OrganizationID string               `bson:"organization_id" json:"organization_id"`
	Owner          *user.UserInfor      `json:"owner" bson:"owner"`
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

type ProductInfo struct {
	ID                   string  `json:"id"`
	ProductName          string  `json:"product_name"`
	OriginalPriceStore   float64 `json:"original_price_store"`
	OriginalPriceService float64 `json:"original_price_service"`
	ProductImage         string  `json:"product_image"`
	TopicName            string  `json:"topic_name"`
	CategoryName         string  `json:"category_name"`
}

type SlotResponse struct {
	SlotID    primitive.ObjectID `json:"slot_id"`
	SlotIDOld primitive.ObjectID `json:"slot_id_old"`
	Sessions  int                `json:"sessions"`
	Title     string             `json:"title"`
	Tracking  string             `json:"tracking"`
	UseCount  int                `json:"use_count"`
	StartTime time.Time          `json:"start_time"`
	EndTime   time.Time          `json:"end_time"`
	Duration  int                `json:"duration"`
	Color     string             `json:"color"`
	Note      string             `json:"note"`
	ProductID *string            `json:"product_id"`
	Product   *ProductInfo       `json:"product,omitempty"`
	CreatedAt time.Time          `json:"created_at"`
	UpdatedAt time.Time          `json:"updated_at"`
}

type BlockResponse struct {
	BlockID    primitive.ObjectID `json:"block_id"`
	BlockIDOld primitive.ObjectID `json:"block_id_old"`
	Slots      []*SlotResponse    `json:"slots"`
}

type ColorTimeDayResponse struct {
	
}