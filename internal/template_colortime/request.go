package templatecolortime

type CreateTemplateColorTimeRequest struct {
	OrganizationID string `json:"organization_id" binding:"required"`
	TermID         string `json:"term_id" binding:"required"`
	Date           string `json:"date" binding:"required"`
	StartTime      string `json:"start_time" binding:"required"`
	Duration       int    `json:"duration" binding:"required"`
	Title          string `json:"title" binding:"required"`
	Color          string `json:"color" binding:"required"`
	Note           string `json:"note"`
	BlockID        string `json:"block_id"`
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
	StartTime string `json:"start_time" binding:"required"`
	Duration  int    `json:"duration" binding:"required"`
	Title     string `json:"title" binding:"required"`
	Color     string `json:"color" binding:"required"`
	Note      string `json:"note"`
	BlockID   string `json:"block_id"`
}
