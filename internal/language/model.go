package language

type APIGateWayResponse[T any] struct {
	StatusCode int    `json:"status_code"`
	Message    string `json:"message"`
	Data       T      `json:"data"`
}

type MessageLanguageResponse struct {
	LangID   uint              `json:"language_id"`
	Contents map[string]string `json:"contents"`
}

type UploadMessageRequest struct {
	TypeID     string `json:"type_id" binding:"required"`
	Type       string `json:"type" binding:"required"`
	Key        string `json:"key" binding:"required"`
	Value      string `json:"message" binding:"required"`
	LanguageID uint   `json:"language_id" binding:"required"`
}

type UploadMessageLanguagesRequest struct {
	MessageLanguages []UploadMessageRequest `json:"message_languages" binding:"required"`
}
