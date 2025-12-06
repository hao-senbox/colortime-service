package default_colortime

type CreateDefaultColorTimeWeekRequest struct {
	StartDate      string `json:"start_date" binding:"required"`
	EndDate        string `json:"end_date" binding:"required"`
	OrganizationID string `json:"organization_id" binding:"required"`

	// Repeat configuration
	RepeatType     string `json:"repeat_type"`     // "none", "daily", "weekly", "monthly", "custom"
	RepeatUntil    string `json:"repeat_until"`    // optional: when to stop repeating (YYYY-MM-DD)
	RepeatInterval int    `json:"repeat_interval"` // repeat every N units (default 1)
	RepeatDays     []int  `json:"repeat_days"`     // for weekly: [0=Sun,1=Mon,...,6=Sat], for custom dates
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
	Date      string `json:"date" binding:"required"`
	StartTime string `json:"start_time" binding:"required"`
	Duration  int    `json:"duration" binding:"required"`
	Title     string `json:"title" binding:"required"`
	Color     string `json:"color" binding:"required"`
	Note      string `json:"note"`
	BlockID   string `json:"block_id"`

	// Repeat configuration (optional)
	RepeatType     string `json:"repeat_type"`     // "none", "daily", "weekly", "monthly", "custom"
	RepeatUntil    string `json:"repeat_until"`    // optional: when to stop repeating (YYYY-MM-DD)
	RepeatInterval int    `json:"repeat_interval"` // repeat every N units (default 1)
	RepeatDays     []int  `json:"repeat_days"`     // for weekly: [0=Sun,1=Mon,...,6=Sat], for custom dates
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
	Title                 string                        `json:"title"`
	ColorTimeSlotLanguage *DefaultColorTimeSlotLanguage `json:"color_time_slot_language"`
	StartTime             string                        `json:"start_time"`
	Duration              int                           `json:"duration"`
	Color                 string                        `json:"color"`
	Note                  string                        `json:"note"`
}

type CreateDefaultDayColorTimeRequest struct {
	Date           string `json:"date" binding:"required"`
	OrganizationID string `json:"organization_id" binding:"required"`

	StartTime             string                        `json:"start_time" binding:"required"`
	Duration              int                           `json:"duration" binding:"required"`
	Title                 string                        `json:"title" binding:"required"`
	Color                 string                        `json:"color" binding:"required"`
	Note                  string                        `json:"note"`
	ColorTimeSlotLanguage *DefaultColorTimeSlotLanguage `json:"color_time_slot_language"`

	BlockID string `json:"block_id"`

	RepeatType     string `json:"repeat_type"`     // "none", "daily", "weekly", "monthly", "custom"
	RepeatUntil    string `json:"repeat_until"`    // optional: when to stop repeating (YYYY-MM-DD)
	RepeatInterval int    `json:"repeat_interval"` // repeat every N units (default 1)
	RepeatDays     []int  `json:"repeat_days"`     // for weekly: [0=Sun,1=Mon,...,6=Sat], for custom dates
}

type UpdateDefaultDayColorTimeRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}
