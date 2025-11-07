package colortime

type CreateColorTimeRequest struct {
	Tracking       string  `bson:"tracking" json:"tracking"`
	OrganizationID string  `bson:"organization_id" json:"organization_id"`
	Title          string  `json:"title"`
	Owner          *Owner  `json:"owner" bson:"owner"`
	Date           string  `json:"date"`
	StartTime      string  `json:"start_time"`
	Duration       int     `json:"duration"`
	Color          string  `json:"color"`
	Note           *string `json:"note"`
	ProductID      string  `json:"product_id"`
	LanguageID     uint    `json:"language_id"`
}

type UpdateColorTimeRequest struct {
	Tracking       string `bson:"tracking" json:"tracking"`
	UseCount       int    `bson:"use_count" json:"use_count"`
	Owner          *Owner `json:"owner" bson:"owner"`
	OrganizationID string `bson:"organization_id" json:"organization_id"`
	Title          string `json:"title"`
	Date           string `json:"date"`
	StartTime      string `json:"start_time"`
	Duration       int    `json:"duration"`
	Color          string `json:"color"`
	Note           string `json:"note"`
	ProductID      string `json:"product_id"`
	LanguageID     uint   `json:"language_id"`
}

type DeleteColorTimeRequest struct {
	OrganizationID string `bson:"organization_id" json:"organization_id"`
	Date           string `json:"date"`
}

type CreateTemplateColorTimeRequest struct {
	Name string `json:"name"`
}

type UpdateTemplateColorTimeRequest struct {
	Name string `json:"name"`
}

type AddSlotsToTemplateColorTimeRequest struct {
	Title     string `json:"title" bson:"title"`
	Tracking  string `json:"tracking" bson:"tracking"`
	StartTime string `json:"start_time" bson:"start_time"`
	Duration  int    `json:"duration" bson:"duration"`
	Color     string `json:"color" bson:"color"`
	Note      string `json:"note" bson:"note"`
	ProductID string `json:"product_id" bson:"product_id"`
}

type EditSlotsToTemplateColorTimeRequest struct {
	Title     string `json:"title" bson:"title"`
	Tracking  string `json:"tracking" bson:"tracking"`
	StartTime string `json:"start_time" bson:"start_time"`
	Duration  int    `json:"duration" bson:"duration"`
	Color     string `json:"color" bson:"color"`
	Note      string `json:"note" bson:"note"`
	ProductID string `json:"product_id" bson:"product_id"`
}

type ApplyTemplateColorTimeRequest struct {
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
	Owner     *Owner `json:"owner" bson:"owner"`
}

type AddTopicToColorTimeWeekRequest struct {
	TopicID string `json:"topic_id" bson:"topic_id"`
}

type AddTopicToColorTimeDayRequest struct {
	Date    string `json:"date" binding:"required"`
	TopicID string `json:"topic_id" binding:"required"`
}

type DeleteTopicToColorTimeDayRequest struct {
	Date string `json:"date" binding:"required"`
}

type CreateColorBlockWithSlotRequest struct {
	Date           string `json:"date" binding:"required"`
	StartTime      string `json:"start_time" binding:"required"`
	Duration       int    `json:"duration" binding:"required"`
	Title          string `json:"title" binding:"required"`
	Tracking       string `json:"tracking" binding:"required"`
	Color          string `json:"color" binding:"required"`
	Note           string `json:"note"`
	ProductID      string `json:"product_id"`
	OrganizationID string `json:"organization_id" binding:"required"`
	BlockID        string `json:"block_id"`
}

type AddSlotToColorBlockRequest struct {
	Date           string `json:"date" binding:"required"`
	BlockID        string `json:"block_id" binding:"required"`
	StartTime      string `json:"start_time" binding:"required"`
	Duration       int    `json:"duration" binding:"required"`
	Title          string `json:"title" binding:"required"`
	Tracking       string `json:"tracking" binding:"required"`
	Color          string `json:"color" binding:"required"`
	Note           string `json:"note"`
	ProductID      string `json:"product_id"`
	OrganizationID string `json:"organization_id" binding:"required"`
}
