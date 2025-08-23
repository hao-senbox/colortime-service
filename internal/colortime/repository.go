package colortime

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type ColorTimeRepository interface {
	CreateColorTime(ctx context.Context, colortime *ColorTime) error
	GetColorTimes(ctx context.Context, userID string, startDate, endDate, start, end time.Time) ([]*ColorTime, error)
	GetColorTime(ctx context.Context, id primitive.ObjectID) (*ColorTime, error)
	UpdateColorTime(ctx context.Context, colortime *ColorTime) error
	DeleteColorTime(ctx context.Context, id primitive.ObjectID) error
}

type colorTimeRepository struct {
	ColorTimeCollection *mongo.Collection
}

func NewColorTimeRepository(colorTimeCollection *mongo.Collection) ColorTimeRepository {
	return &colorTimeRepository{
		ColorTimeCollection: colorTimeCollection,
	}
}

func (r *colorTimeRepository) CreateColorTime(ctx context.Context, colortime *ColorTime) error {
	_, err := r.ColorTimeCollection.InsertOne(ctx, colortime)
	return err
}

func (r *colorTimeRepository) GetColorTimes(ctx context.Context, userID string, startDate, endDate, start, end time.Time) ([]*ColorTime, error) {

	filter := bson.M{
		"user_id": userID,
	}

	if start.Before(end) {
		filter["date"] = bson.M{
			"$gte": startDate,
			"$lte": endDate,
		}
		filter["time"] = bson.M{
			"$gte": start,
			"$lte": end,
		}
	} else {
		filter["$and"] = []bson.M{
			{"date": bson.M{"$gte": startDate, "$lte": endDate}},
			{"$or": []bson.M{
				{
					"time": bson.M{"$gte": start}, 
				},
				{
					"time": bson.M{"$lte": end}, 
				},
			}},
		}
	}

	cursor, err := r.ColorTimeCollection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}

	var colorTimes []*ColorTime
	if err := cursor.All(ctx, &colorTimes); err != nil {
		return nil, err
	}

	return colorTimes, nil
}

func (r *colorTimeRepository) GetColorTime(ctx context.Context, id primitive.ObjectID) (*ColorTime, error) {

	var colortime ColorTime

	if err := r.ColorTimeCollection.FindOne(ctx, bson.M{"_id": id}).Decode(&colortime); err != nil {
		return nil, err
	}

	return &colortime, nil

}

func (r *colorTimeRepository) UpdateColorTime(ctx context.Context, colortime *ColorTime) error {
	_, err := r.ColorTimeCollection.UpdateOne(ctx, bson.M{"_id": colortime.ID}, bson.M{"$set": colortime})
	return err
}

func (r *colorTimeRepository) DeleteColorTime(ctx context.Context, id primitive.ObjectID) error {
	_, err := r.ColorTimeCollection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}