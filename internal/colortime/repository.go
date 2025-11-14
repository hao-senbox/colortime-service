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
	CountTrackingUsage(ctx context.Context, organizationID, userID, role, tracking string) (int, error)
	GetAllSlotsByTracking(ctx context.Context, organizationID, userID, role, tracking string) ([]*ColortimeSlot, []*WeekColorTime, error)
}

type colorTimeRepository struct {
	ColorTimeCollection *mongo.Collection
}

func NewColorTimeRepository(colorTimeCollection *mongo.Collection) ColorTimeRepository {
	return &colorTimeRepository{
		ColorTimeCollection: colorTimeCollection,
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

func (r *colorTimeRepository) CountTrackingUsage(ctx context.Context, organizationID, userID, role, tracking string) (int, error) {
	filter := bson.M{
		"organization_id":  organizationID,
		"owner.owner_id":   userID,
		"owner.owner_role": role,
	}

	cursor, err := r.ColorTimeCollection.Find(ctx, filter)
	if err != nil {
		return 0, err
	}
	defer cursor.Close(ctx)

	totalCount := 0
	for cursor.Next(ctx) {
		var week WeekColorTime
		if err := cursor.Decode(&week); err != nil {
			return 0, err
		}

		// Count tracking in this week
		for _, colorTime := range week.ColorTimes {
			for _, block := range colorTime.TimeSlots {
				for _, slot := range block.Slots {
					if slot.Tracking == tracking {
						totalCount++
					}
				}
			}
		}
	}

	if err := cursor.Err(); err != nil {
		return 0, err
	}

	return totalCount, nil
}

func (r *colorTimeRepository) GetAllSlotsByTracking(ctx context.Context, organizationID, userID, role, tracking string) ([]*ColortimeSlot, []*WeekColorTime, error) {

	filter := bson.M{
		"organization_id":  organizationID,
		"owner.owner_id":   userID,
		"owner.owner_role": role,
	}

	cursor, err := r.ColorTimeCollection.Find(ctx, filter)
	if err != nil {
		return nil, nil, err
	}
	defer cursor.Close(ctx)

	var slots []*ColortimeSlot
	var weeks []*WeekColorTime

	for cursor.Next(ctx) {

		var week WeekColorTime
		if err := cursor.Decode(&week); err != nil {
			return nil, nil, err
		}

		weeks = append(weeks, &week)

		for _, colorTime := range week.ColorTimes {
			for _, block := range colorTime.TimeSlots {
				for _, slot := range block.Slots {
					if slot.Tracking == tracking {
						slots = append(slots, slot)
					}
				}
			}
		}

	}

	return slots, weeks, nil
}
