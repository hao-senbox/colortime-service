package colortime

import (
	"colortime-service/internal/language"
	"colortime-service/pkg/constants"
)

func BuildColortimeMessagesUpdate(colortimeID string, req UpdateColorTimeRequest) language.UploadMessageLanguagesRequest {
	return language.UploadMessageLanguagesRequest{
		MessageLanguages: []language.UploadMessageRequest{
			{
				TypeID:     colortimeID,
				Type:       "colortime",
				Key:        string(constants.ColortimeNoteKey),
				Value:      req.Note,
				LanguageID: req.LanguageID,
			},
			{
				TypeID:     colortimeID,
				Type:       "colortime",
				Key:        string(constants.ColortimeTitleKey),
				Value:      req.Title,
				LanguageID: req.LanguageID,
			},
		},
	}
}
