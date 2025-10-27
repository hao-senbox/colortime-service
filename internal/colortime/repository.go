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
	GetColorTimes(ctx context.Context, start string, end string) ([]*ColorTime, error)
	GetColorTime(ctx context.Context, id primitive.ObjectID) (*ColorTime, error)
	GetColorTimeByDate(ctx context.Context, organizationID string, date time.Time) (*ColorTime, error)
	UpdateColorTime(ctx context.Context, id primitive.ObjectID, colortime *ColorTime) error
	CreateColorTimeWeek(ctx context.Context, colortimeWeek *WeekColorTime) error
	GetColorTimeWeek(ctx context.Context, startDate, endDate *time.Time, organizationID, userID string) (*WeekColorTime, error)
	UpdateColorTimeWeek(ctx context.Context, id primitive.ObjectID, colortimeWeek *WeekColorTime) error

	CreateTemplateColorTime(ctx context.Context, template *ColorTimeTemplate) error
	GetTemplateColorTimes(ctx context.Context) ([]*ColorTimeTemplate, error)
	GetTemplateColorTime(ctx context.Context, id primitive.ObjectID) (*ColorTimeTemplate, error)
	UpdateTemplateColorTime(ctx context.Context, template *ColorTimeTemplate) error
	DeleteTemplateColorTime(ctx context.Context, id primitive.ObjectID) error
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

func (r *colorTimeRepository) CreateColorTime(ctx context.Context, colortime *ColorTime) error {
	_, err := r.ColorTimeCollection.InsertOne(ctx, colortime)
	return err
}

func (r *colorTimeRepository) GetColorTimeByDate(ctx context.Context, organizationID string, date time.Time) (*ColorTime, error) {

	var colortime ColorTime

	filter := bson.M{
		"organization_id": organizationID,
		"date":            date,
	}

	err := r.ColorTimeCollection.FindOne(ctx, filter).Decode(&colortime)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}

	return &colortime, err

}

func (r *colorTimeRepository) GetColorTimes(ctx context.Context, start string, end string) ([]*ColorTime, error) {

	filter := bson.M{}

	if start != "" && end != "" {
		startParse, err := time.Parse("2006-01-02", start)
		if err != nil {
			return nil, err
		}

		endParse, err := time.Parse("2006-01-02", end)
		if err != nil {
			return nil, err
		}

		filter = bson.M{
			"date": bson.M{
				"$gte": startParse,
				"$lte": endParse,
			},
		}
	}

	var colortimes []*ColorTime

	cursor, err := r.ColorTimeCollection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}

	if err := cursor.All(ctx, &colortimes); err != nil {
		return nil, err
	}

	return colortimes, nil

}

func (r *colorTimeRepository) GetColorTime(ctx context.Context, id primitive.ObjectID) (*ColorTime, error) {

	var colortime ColorTime

	if err := r.ColorTimeCollection.FindOne(ctx, bson.M{"_id": id}).Decode(&colortime); err != nil {
		return nil, err
	}

	return &colortime, nil

}

func (r *colorTimeRepository) UpdateColorTime(ctx context.Context, id primitive.ObjectID, colortime *ColorTime) error {

	_, err := r.ColorTimeCollection.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": colortime})
	return err

}

func (r *colorTimeRepository) CreateColorTimeWeek(ctx context.Context, colortimeWeek *WeekColorTime) error {
	_, err := r.ColorTimeCollection.InsertOne(ctx, colortimeWeek)
	return err
}

func (r *colorTimeRepository) GetColorTimeWeek(ctx context.Context, startDate, endDate *time.Time, organizationID, userID string) (*WeekColorTime, error) {

	filter := bson.M{
		"organization_id": organizationID,
		"owner.owner_id":  userID,
		"start_date":      startDate,
		"end_date":        endDate,
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

func (r *colorTimeRepository) UpdateColorTimeWeek(ctx context.Context, id primitive.ObjectID, colortimeWeek *WeekColorTime) error {
	_, err := r.ColorTimeCollection.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": colortimeWeek})
	return err
}

func (r *colorTimeRepository) CreateTemplateColorTime(ctx context.Context, template *ColorTimeTemplate) error {
	_, err := r.ColorTimeTemplateCollection.InsertOne(ctx, template)
	return err
}

func (r *colorTimeRepository) GetTemplateColorTimes(ctx context.Context) ([]*ColorTimeTemplate, error) {

	var templates []*ColorTimeTemplate

	filter := bson.M{
		"is_deleted": false,
	}

	cursor, err := r.ColorTimeTemplateCollection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}

	if err := cursor.All(ctx, &templates); err != nil {
		return nil, err
	}

	return templates, nil

}

func (r *colorTimeRepository) GetTemplateColorTime(ctx context.Context, id primitive.ObjectID) (*ColorTimeTemplate, error) {

	var template ColorTimeTemplate

	if err := r.ColorTimeTemplateCollection.FindOne(ctx, bson.M{"_id": id}).Decode(&template); err != nil {
		return nil, err
	}

	return &template, nil

}

func (r *colorTimeRepository) UpdateTemplateColorTime(ctx context.Context, template *ColorTimeTemplate) error {
	_, err := r.ColorTimeTemplateCollection.UpdateOne(ctx, bson.M{"_id": template.ID}, bson.M{"$set": template})
	return err
}

func (r *colorTimeRepository) DeleteTemplateColorTime(ctx context.Context, id primitive.ObjectID) error {

	filter := bson.M{"_id": id}
	update := bson.M{"$set": bson.M{"is_deleted": true}}

	_, err := r.ColorTimeTemplateCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	return nil
}
