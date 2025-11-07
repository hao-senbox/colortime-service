package default_colortime

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type DefaultColorTimeRepository interface {
	CreateDefaultColorTimeWeek(ctx context.Context, colortimeWeek *DefaultWeekColorTime) error
	GetDefaultColorTimeWeek(ctx context.Context, startDate, endDate *time.Time, organizationID string) (*DefaultWeekColorTime, error)
	GetDefaultColorTimeWeekByID(ctx context.Context, id primitive.ObjectID) (*DefaultWeekColorTime, error)
	UpdateDefaultColorTimeWeek(ctx context.Context, id primitive.ObjectID, colortimeWeek *DefaultWeekColorTime) error
	DeleteDefaultColorTimeWeek(ctx context.Context, id primitive.ObjectID) error
	GetAllDefaultColorTimeWeeks(ctx context.Context, organizationID string) ([]*DefaultWeekColorTime, error)
}

type defaultColorTimeRepository struct {
	DefaultColorTimeCollection *mongo.Collection
}

func NewDefaultColorTimeRepository(defaultColorTimeCollection *mongo.Collection) DefaultColorTimeRepository {
	return &defaultColorTimeRepository{
		DefaultColorTimeCollection: defaultColorTimeCollection,
	}
}

func (r *defaultColorTimeRepository) CreateDefaultColorTimeWeek(ctx context.Context, colortimeWeek *DefaultWeekColorTime) error {
	_, err := r.DefaultColorTimeCollection.InsertOne(ctx, colortimeWeek)
	return err
}

func (r *defaultColorTimeRepository) GetDefaultColorTimeWeek(ctx context.Context, startDate, endDate *time.Time, organizationID string) (*DefaultWeekColorTime, error) {
	filter := bson.M{
		"organization_id": organizationID,
		"start_date":      bson.M{"$lte": endDate},
		"end_date":        bson.M{"$gte": startDate},
	}

	var colortimeWeek DefaultWeekColorTime

	if err := r.DefaultColorTimeCollection.FindOne(ctx, filter).Decode(&colortimeWeek); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}

	return &colortimeWeek, nil
}

func (r *defaultColorTimeRepository) GetDefaultColorTimeWeekByID(ctx context.Context, id primitive.ObjectID) (*DefaultWeekColorTime, error) {
	var colortimeWeek DefaultWeekColorTime

	if err := r.DefaultColorTimeCollection.FindOne(ctx, bson.M{"_id": id}).Decode(&colortimeWeek); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}

	return &colortimeWeek, nil
}

func (r *defaultColorTimeRepository) UpdateDefaultColorTimeWeek(ctx context.Context, id primitive.ObjectID, colortimeWeek *DefaultWeekColorTime) error {
	_, err := r.DefaultColorTimeCollection.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": colortimeWeek})
	return err
}

func (r *defaultColorTimeRepository) DeleteDefaultColorTimeWeek(ctx context.Context, id primitive.ObjectID) error {
	_, err := r.DefaultColorTimeCollection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}

func (r *defaultColorTimeRepository) GetAllDefaultColorTimeWeeks(ctx context.Context, organizationID string) ([]*DefaultWeekColorTime, error) {
	filter := bson.M{
		"organization_id": organizationID,
	}

	var colortimeWeeks []*DefaultWeekColorTime

	cursor, err := r.DefaultColorTimeCollection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	if err := cursor.All(ctx, &colortimeWeeks); err != nil {
		return nil, err
	}

	return colortimeWeeks, nil
}
