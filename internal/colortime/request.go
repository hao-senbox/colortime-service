package colortime

type CreateColorTimeRequest struct {
	Title     string `json:"title"`
	UserID    string `json:"user_id"`
	Date      string `json:"date"`
	Time      string `json:"time"`
	Duration  int    `json:"duration"`
	Color     string `json:"color"`
	Note      string `json:"note"`
	ProductID string `json:"product_id"`
}

type UpdateColorTimeRequest struct {
	Title     string `json:"title"`
	Date      string `json:"date"`
	Time      string `json:"time"`
	Duration  int    `json:"duration"`
	Color     string `json:"color"`
	Note      string `json:"note"`
	ProductID string `json:"product_id"`
}