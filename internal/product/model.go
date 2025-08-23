package product

import "go.mongodb.org/mongo-driver/bson/primitive"

type Product struct {
	ID                   primitive.ObjectID `bson:"_id" json:"id"`
	ProductName          string             `bson:"product_name" json:"product_name"`
	OriginalPriceStore   float64            `bson:"original_price_store" json:"original_price_store"`
	OriginalPriceService float64            `bson:"original_price_service" json:"original_price_service"`
	ProductDescription   string             `bson:"product_description" json:"product_description"`
	ProductImage         string             `bson:"product_image" json:"product_image"`
	TopicName            string             `bson:"topic_name" json:"topic_name"`
	CategoryName         string             `bson:"category_name" json:"category_name"`
}
