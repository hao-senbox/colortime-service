package colortime

type CreateColorTimeRequest struct {
	Tracking       string `bson:"tracking" json:"tracking"`
	UseCount       int    `bson:"use_count" json:"use_count"`
	OrganizationID string `bson:"organization_id" json:"organization_id"`
	Title          string `json:"title"`
	UserID         string `json:"user_id"`
	Date           string `json:"date"`
	Time           string `json:"time"`
	Duration       int    `json:"duration"`
	Color          string `json:"color"`
	Note           string `json:"note"`
	ProductID      string `json:"product_id"`
	LanguageID     uint   `json:"language_id"`
}

type UpdateColorTimeRequest struct {
	Tracking       string `bson:"tracking" json:"tracking"`
	UseCount       int    `bson:"use_count" json:"use_count"`
	OrganizationID string `bson:"organization_id" json:"organization_id"`
	Title          string `json:"title"`
	Date           string `json:"date"`
	Time           string `json:"time"`
	Duration       int    `json:"duration"`
	Color          string `json:"color"`
	Note           string `json:"note"`
	ProductID      string `json:"product_id"`
	LanguageID     uint   `json:"language_id"`
}
