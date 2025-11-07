package default_colortime

type CreateDefaultColorTimeWeekRequest struct {
	StartDate      string `json:"start_date" binding:"required"`
	EndDate        string `json:"end_date" binding:"required"`
	OrganizationID string `json:"organization_id" binding:"required"`
}

type UpdateDefaultColorTimeWeekRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type AddTopicToDefaultColorTimeWeekRequest struct {
	TopicID string `json:"topic_id" bson:"topic_id"`
}

type AddTopicToDefaultColorTimeDayRequest struct {
	Date    string `json:"date" binding:"required"`
	TopicID string `json:"topic_id" binding:"required"`
}

type DeleteTopicToDefaultColorTimeDayRequest struct {
	Date string `json:"date" binding:"required"`
}

type CreateDefaultColorBlockWithSlotRequest struct {
	Date           string `json:"date" binding:"required"`
	StartTime      string `json:"start_time" binding:"required"`
	Duration       int    `json:"duration" binding:"required"`
	Title          string `json:"title" binding:"required"`
	Color          string `json:"color" binding:"required"`
	Note           string `json:"note"`
	BlockID        string `json:"block_id"`
}

type AddSlotToDefaultColorBlockRequest struct {
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

type UpdateDefaultColorSlotRequest struct {
	Title     string `json:"title"`
	Tracking  string `json:"tracking"`
	StartTime string `json:"start_time"`
	Duration  int    `json:"duration"`
	Color     string `json:"color"`
	Note      string `json:"note"`
	ProductID string `json:"product_id"`
}
