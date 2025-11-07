package colortime

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type ColorTimeRepository interface {
	CreateColorTimeWeek(ctx context.Context, colortimeWeek *WeekColorTime) error
	GetColorTimeWeek(ctx context.Context, startDate, endDate *time.Time, organizationID, userID, role string) (*WeekColorTime, error)
	GetColorTimeWeekByID(ctx context.Context, id primitive.ObjectID) (*WeekColorTime, error)
	UpdateColorTimeWeek(ctx context.Context, id primitive.ObjectID, colortimeWeek *WeekColorTime) error
}

type colorTimeRepository struct {
	ColorTimeCollection         *mongo.Collection
	ColorTimeTemplateCollection *mongo.Collection
}

func NewColorTimeRepository(colorTimeCollection, colorTimeTemplateCollection *mongo.Collection) ColorTimeRepository {
	return &colorTimeRepository{
		ColorTimeCollection:         colorTimeCollection,
		ColorTimeTemplateCollection: colorTimeTemplateCollection,
	}
}

func (r *colorTimeRepository) CreateColorTimeWeek(ctx context.Context, colortimeWeek *WeekColorTime) error {
	_, err := r.ColorTimeCollection.InsertOne(ctx, colortimeWeek)
	return err
}

func (r *colorTimeRepository) GetColorTimeWeek(ctx context.Context, startDate, endDate *time.Time, organizationID, userID, role string) (*WeekColorTime, error) {

	filter := bson.M{
		"organization_id":  organizationID,
		"owner.owner_id":   userID,
		"owner.owner_role": role, 
		"start_date":       bson.M{"$lte": endDate},
		"end_date":         bson.M{"$gte": startDate},
	}

	var colortimeWeek WeekColorTime

	if err := r.ColorTimeCollection.FindOne(ctx, filter).Decode(&colortimeWeek); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}

	return &colortimeWeek, nil

}

func (r *colorTimeRepository) GetColorTimeWeekByID(ctx context.Context, id primitive.ObjectID) (*WeekColorTime, error) {

	var colortimeWeek WeekColorTime

	if err := r.ColorTimeCollection.FindOne(ctx, bson.M{"_id": id}).Decode(&colortimeWeek); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}

	return &colortimeWeek, nil

}

func (r *colorTimeRepository) UpdateColorTimeWeek(ctx context.Context, id primitive.ObjectID, colortimeWeek *WeekColorTime) error {
	_, err := r.ColorTimeCollection.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": colortimeWeek})
	return err
}
