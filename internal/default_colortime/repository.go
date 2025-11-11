package default_colortime

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type DefaultColorTimeRepository interface {
	CreateDefaultDayColorTime(ctx context.Context, dayColorTime *DefaultDayColorTime) error
	GetDefaultDayColorTime(ctx context.Context, date time.Time, organizationID string) (*DefaultDayColorTime, error)
	GetDefaultDayColorTimeByID(ctx context.Context, id primitive.ObjectID) (*DefaultDayColorTime, error)
	GetDefaultDayColorTimeBySlotID(ctx context.Context, slotID primitive.ObjectID) (*DefaultDayColorTime, error)
	UpdateDefaultDayColorTime(ctx context.Context, id primitive.ObjectID, dayColorTime *DefaultDayColorTime) error
	DeleteDefaultDayColorTime(ctx context.Context, id primitive.ObjectID) error
	GetDefaultDayColorTimesInRange(ctx context.Context, startDate, endDate time.Time, organizationID string) ([]*DefaultDayColorTime, error)
	GetAllDefaultDayColorTimes(ctx context.Context, organizationID string) ([]*DefaultDayColorTime, error)
}

type defaultColorTimeRepository struct {
	DefaultColorTimeCollection *mongo.Collection
}

func NewDefaultColorTimeRepository(defaultColorTimeCollection *mongo.Collection) DefaultColorTimeRepository {
	return &defaultColorTimeRepository{
		DefaultColorTimeCollection: defaultColorTimeCollection,
	}
}

func (r *defaultColorTimeRepository) CreateDefaultDayColorTime(ctx context.Context, dayColorTime *DefaultDayColorTime) error {
	_, err := r.DefaultColorTimeCollection.InsertOne(ctx, dayColorTime)
	return err
}

func (r *defaultColorTimeRepository) GetDefaultDayColorTime(ctx context.Context, date time.Time, organizationID string) (*DefaultDayColorTime, error) {
	filter := bson.M{
		"organization_id": organizationID,
		"date": bson.M{
			"$gte": time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location()),
			"$lt":  time.Date(date.Year(), date.Month(), date.Day(), 23, 59, 59, 999999999, date.Location()),
		},
	}

	var dayColorTime DefaultDayColorTime

	if err := r.DefaultColorTimeCollection.FindOne(ctx, filter).Decode(&dayColorTime); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}

	return &dayColorTime, nil
}

func (r *defaultColorTimeRepository) GetDefaultDayColorTimeByID(ctx context.Context, id primitive.ObjectID) (*DefaultDayColorTime, error) {
	var dayColorTime DefaultDayColorTime

	if err := r.DefaultColorTimeCollection.FindOne(ctx, bson.M{"_id": id}).Decode(&dayColorTime); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}

	return &dayColorTime, nil
}

func (r *defaultColorTimeRepository) GetDefaultDayColorTimeBySlotID(ctx context.Context, slotID primitive.ObjectID) (*DefaultDayColorTime, error) {

	filter := bson.M{
		"time_slots.slots.slot_id": slotID,
	}

	var dayColorTime DefaultDayColorTime
	err := r.DefaultColorTimeCollection.FindOne(ctx, filter).Decode(&dayColorTime)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}

	return &dayColorTime, nil
}

func (r *defaultColorTimeRepository) UpdateDefaultDayColorTime(ctx context.Context, id primitive.ObjectID, dayColorTime *DefaultDayColorTime) error {
	_, err := r.DefaultColorTimeCollection.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": dayColorTime})
	return err
}

func (r *defaultColorTimeRepository) DeleteDefaultDayColorTime(ctx context.Context, id primitive.ObjectID) error {
	_, err := r.DefaultColorTimeCollection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}

func (r *defaultColorTimeRepository) GetDefaultDayColorTimesInRange(ctx context.Context, startDate, endDate time.Time, organizationID string) ([]*DefaultDayColorTime, error) {
	filter := bson.M{
		"organization_id": organizationID,
		"date": bson.M{
			"$gte": time.Date(startDate.Year(), startDate.Month(), startDate.Day(), 0, 0, 0, 0, startDate.Location()),
			"$lte": time.Date(endDate.Year(), endDate.Month(), endDate.Day(), 23, 59, 59, 999999999, endDate.Location()),
		},
	}

	var dayColorTimes []*DefaultDayColorTime

	cursor, err := r.DefaultColorTimeCollection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	if err := cursor.All(ctx, &dayColorTimes); err != nil {
		return nil, err
	}

	return dayColorTimes, nil
}

func (r *defaultColorTimeRepository) GetAllDefaultDayColorTimes(ctx context.Context, organizationID string) ([]*DefaultDayColorTime, error) {
	filter := bson.M{
		"organization_id": organizationID,
	}

	var dayColorTimes []*DefaultDayColorTime

	cursor, err := r.DefaultColorTimeCollection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	if err := cursor.All(ctx, &dayColorTimes); err != nil {
		return nil, err
	}

	return dayColorTimes, nil
}
