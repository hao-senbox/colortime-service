package topic

type Topic struct {
	ID           string `json:"id" bson:"_id"`
	Name         string `json:"name" bson:"name"`
	MainImageUrl string `json:"main_image_url" bson:"main_image_url"`
	VideoUrl     string `json:"video_url" bson:"video_url"`
}
