package templatecolortime

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type TemplateColorTimeRepository interface {
	CreateTemplateColorTime(ctx context.Context, colortimeTemplate *TemplateColorTime) error
	GetTemplateColorTime(ctx context.Context, organizationID, termID, date string) (*TemplateColorTime, error)
	GetTemplateColorTimeByID(ctx context.Context, id primitive.ObjectID) (*TemplateColorTime, error)
	UpdateTemplateColorTime(ctx context.Context, id primitive.ObjectID, colortimeTemplate *TemplateColorTime) error
	DeleteTemplateColorTime(ctx context.Context, id primitive.ObjectID) error
}

type templateColorTimeRepository struct {
	TemplateColorTimeCollection *mongo.Collection
}

func NewTemplateColorTimeRepository(templateColorTimeCollection *mongo.Collection) TemplateColorTimeRepository {
	return &templateColorTimeRepository{
		TemplateColorTimeCollection: templateColorTimeCollection,
	}
}

func (r *templateColorTimeRepository) CreateTemplateColorTime(ctx context.Context, colortimeTemplate *TemplateColorTime) error {
	_, err := r.TemplateColorTimeCollection.InsertOne(ctx, colortimeTemplate)
	return err
}

func (r *templateColorTimeRepository) GetTemplateColorTime(ctx context.Context, organizationID, termID, date string) (*TemplateColorTime, error) {
	filter := bson.M{
		"organization_id": organizationID,
		"term_id":         termID,
		"date":            date,
	}

	var colortimeTemplate TemplateColorTime

	if err := r.TemplateColorTimeCollection.FindOne(ctx, filter).Decode(&colortimeTemplate); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}

	return &colortimeTemplate, nil
}

func (r *templateColorTimeRepository) GetTemplateColorTimeByID(ctx context.Context, id primitive.ObjectID) (*TemplateColorTime, error) {
	var colortimeTemplate TemplateColorTime

	if err := r.TemplateColorTimeCollection.FindOne(ctx, bson.M{"_id": id}).Decode(&colortimeTemplate); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}

	return &colortimeTemplate, nil
}

func (r *templateColorTimeRepository) UpdateTemplateColorTime(ctx context.Context, id primitive.ObjectID, colortimeTemplate *TemplateColorTime) error {
	_, err := r.TemplateColorTimeCollection.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": colortimeTemplate})
	return err
}

func (r *templateColorTimeRepository) DeleteTemplateColorTime(ctx context.Context, id primitive.ObjectID) error {
	_, err := r.TemplateColorTimeCollection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}
