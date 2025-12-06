package templatecolortime

type CreateTemplateColorTimeRequest struct {
	OrganizationID        string                 `json:"organization_id" binding:"required"`
	TermID                string                 `json:"term_id" binding:"required"`
	Date                  string                 `json:"date" binding:"required"`
	StartTime             string                 `json:"start_time" binding:"required"`
	Duration              int                    `json:"duration" binding:"required"`
	Title                 string                 `json:"title" binding:"required"`
	ColorTimeSlotLanguage *ColorTimeSlotLanguage `json:"color_time_slot_language" binding:"required"`
	Color                 string                 `json:"color" binding:"required"`
	Note                  string                 `json:"note"`
	BlockID               string                 `json:"block_id"`
}

type DuplicateTemplateColorTimeRequest struct {
	OrganizationID string `json:"organization_id" binding:"required"`
	TermID         string `json:"term_id" binding:"required"`
	OriginDate     string `json:"origin_date" binding:"required"`
	TargetDate     string `json:"target_date"`
}

type ApplyTemplateColorTimeRequest struct {
	OrganizationID string `json:"organization_id" binding:"required"`
	TermID         string `json:"term_id" binding:"required"`
	StartDate      string `json:"start_date" binding:"required"`
	EndDate        string `json:"end_date" binding:"required"`
}

type UpdateTemplateColorTimeSlotRequest struct {
	StartTime             string                 `json:"start_time"`
	Duration              int                    `json:"duration"`
	Title                 string                 `json:"title"`
	ColorTimeSlotLanguage *ColorTimeSlotLanguage `json:"color_time_slot_language"`
	Color                 string                 `json:"color"`
	Note                  string                 `json:"note"`
	BlockID               string                 `json:"block_id"`
}

type CopySlotToTemplateColorTimeRequest struct {
	OrganizationID string  `json:"organization_id" binding:"required"`
	TermID         string  `json:"term_id" binding:"required"`
	SlotID         *string `json:"slot_id"`
	BlockIDTarget  *string `json:"block_id_target"`
	OriginDate     string  `json:"origin_date" binding:"required"`
	TargetDate     string  `json:"target_date" binding:"required"`
	BaseHour       *int    `json:"base_hour"`
}
