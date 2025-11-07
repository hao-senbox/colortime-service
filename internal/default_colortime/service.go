package default_colortime

import (
	"colortime-service/internal/product"
	"colortime-service/internal/topic"
	"context"
	"errors"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type DefaultColorTimeService interface {
	GetDefaultColorTimeWeek(ctx context.Context, orgID, start, end, userID string) (*TopicToDefaultColorTimeWeekResponse, error)
	GetAllDefaultColorTimeWeeks(ctx context.Context, orgID string) ([]*TopicToDefaultColorTimeWeekResponse, error)
	UpdateDefaultColorTimeWeek(ctx context.Context, id string, req *UpdateDefaultColorTimeWeekRequest) error
	DeleteDefaultColorTimeWeek(ctx context.Context, id string) error

	CreateDefaultColorBlockAndSaveSlot(ctx context.Context, weekColorTimeID string, req *CreateDefaultColorBlockWithSlotRequest, userID string) (*DefaultColorBlock, error)
	CreateDefaultColorBlockForSession(ctx context.Context, weekColorTimeID string, req *CreateDefaultColorBlockWithSlotRequest, userID string) (*DefaultColorBlock, error)

	UpdateDefaultColorSlot(ctx context.Context, weekID, slotID string, req *UpdateDefaultColorSlotRequest) error
	DeleteDefaultColorSlot(ctx context.Context, weekID, slotID string) error
}

type defaultColorTimeService struct {
	DefaultColorTimeRepository DefaultColorTimeRepository
	ProductService             product.ProductService
	TopicService               topic.TopicService
}

func NewDefaultColorTimeService(
	defaultColorTimeRepository DefaultColorTimeRepository,
	productService product.ProductService,
	topicService topic.TopicService,
) DefaultColorTimeService {
	return &defaultColorTimeService{
		DefaultColorTimeRepository: defaultColorTimeRepository,
		ProductService:             productService,
		TopicService:               topicService,
	}
}

func (s *defaultColorTimeService) GetDefaultColorTimeWeek(ctx context.Context, orgID, start, end, userID string) (*TopicToDefaultColorTimeWeekResponse, error) {
	if orgID == "" {
		return nil, errors.New("organization id is required")
	}

	if start == "" || end == "" {
		return nil, errors.New("start and end date are required")
	}

	startDate, err := time.Parse("2006-01-02", start)
	if err != nil {
		return nil, fmt.Errorf("invalid start date: %w", err)
	}

	endDate, err := time.Parse("2006-01-02", end)
	if err != nil {
		return nil, fmt.Errorf("invalid end date: %w", err)
	}

	colortimeWeek, err := s.DefaultColorTimeRepository.GetDefaultColorTimeWeek(ctx, &startDate, &endDate, orgID)
	if err != nil {
		return nil, err
	}

	if colortimeWeek == nil {
		weekDefault := &DefaultWeekColorTime{
			ID:             primitive.NewObjectID(),
			OrganizationID: orgID,
			StartDate:      startDate,
			EndDate:        endDate,
			ColorTimes:     []*DefaultColorTime{},
			CreatedBy:      userID,
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
		}
		err = s.DefaultColorTimeRepository.CreateDefaultColorTimeWeek(ctx, weekDefault)
		if err != nil {
			return nil, err
		}
		colortimeWeek = weekDefault
	}

	colorTimeResponses := make([]*DefaultColorTimeResponse, 0, len(colortimeWeek.ColorTimes))

	for _, day := range colortimeWeek.ColorTimes {
		colorTimeResponses = append(colorTimeResponses, &DefaultColorTimeResponse{
			ID:        day.ID,
			Date:      day.Date,
			TimeSlots: day.TimeSlots,
			CreatedAt: day.CreatedAt,
			UpdatedAt: day.UpdatedAt,
		})
	}

	result := &TopicToDefaultColorTimeWeekResponse{
		ID:             colortimeWeek.ID,
		OrganizationID: colortimeWeek.OrganizationID,
		StartDate:      colortimeWeek.StartDate,
		EndDate:        colortimeWeek.EndDate,
		ColorTimes:     colorTimeResponses,
		CreatedBy:      colortimeWeek.CreatedBy,
		CreatedAt:      colortimeWeek.CreatedAt,
		UpdatedAt:      colortimeWeek.UpdatedAt,
	}

	return result, nil
}

func (s *defaultColorTimeService) GetAllDefaultColorTimeWeeks(ctx context.Context, orgID string) ([]*TopicToDefaultColorTimeWeekResponse, error) {
	if orgID == "" {
		return nil, errors.New("organization id is required")
	}

	colortimeWeeks, err := s.DefaultColorTimeRepository.GetAllDefaultColorTimeWeeks(ctx, orgID)
	if err != nil {
		return nil, err
	}

	results := make([]*TopicToDefaultColorTimeWeekResponse, 0, len(colortimeWeeks))

	for _, colortimeWeek := range colortimeWeeks {
		colorTimeResponses := make([]*DefaultColorTimeResponse, 0, len(colortimeWeek.ColorTimes))

		for _, day := range colortimeWeek.ColorTimes {
			colorTimeResponses = append(colorTimeResponses, &DefaultColorTimeResponse{
				ID:        day.ID,
				Date:      day.Date,
				TimeSlots: day.TimeSlots,
				CreatedAt: day.CreatedAt,
				UpdatedAt: day.UpdatedAt,
			})
		}

		results = append(results, &TopicToDefaultColorTimeWeekResponse{
			ID:             colortimeWeek.ID,
			OrganizationID: colortimeWeek.OrganizationID,
			StartDate:      colortimeWeek.StartDate,
			EndDate:        colortimeWeek.EndDate,
			ColorTimes:     colorTimeResponses,
			CreatedBy:      colortimeWeek.CreatedBy,
			CreatedAt:      colortimeWeek.CreatedAt,
			UpdatedAt:      colortimeWeek.UpdatedAt,
		})
	}

	return results, nil
}

func (s *defaultColorTimeService) UpdateDefaultColorTimeWeek(ctx context.Context, id string, req *UpdateDefaultColorTimeWeekRequest) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("invalid id format")
	}

	colortimeWeek, err := s.DefaultColorTimeRepository.GetDefaultColorTimeWeekByID(ctx, objectID)
	if err != nil {
		return err
	}

	if colortimeWeek == nil {
		return fmt.Errorf("default colortime week not found")
	}

	colortimeWeek.UpdatedAt = time.Now()

	return s.DefaultColorTimeRepository.UpdateDefaultColorTimeWeek(ctx, colortimeWeek.ID, colortimeWeek)
}

func (s *defaultColorTimeService) DeleteDefaultColorTimeWeek(ctx context.Context, id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("invalid id format")
	}

	return s.DefaultColorTimeRepository.DeleteDefaultColorTimeWeek(ctx, objectID)
}

func sameDay(a, b time.Time) bool {
	y1, m1, d1 := a.Date()
	y2, m2, d2 := b.Date()
	return y1 == y2 && m1 == m2 && d1 == d2
}

func (s *defaultColorTimeService) CreateDefaultColorBlockAndSaveSlot(
	ctx context.Context,
	weekColorTimeID string,
	req *CreateDefaultColorBlockWithSlotRequest,
	userID string,
) (*DefaultColorBlock, error) {

	if userID == "" {
		return nil, errors.New("user id is required")
	}

	if req.Date == "" {
		return nil, errors.New("date is required")
	}

	dateParse, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		return nil, fmt.Errorf("invalid date format: %w", err)
	}

	if req.StartTime == "" {
		return nil, errors.New("start time is required")
	}

	startTimeParse, err := time.Parse("15:04", req.StartTime)
	if err != nil {
		return nil, fmt.Errorf("invalid start time format: %w", err)
	}

	if req.Duration == 0 {
		return nil, errors.New("duration is required")
	}

	if req.Title == "" {
		return nil, errors.New("title is required")
	}

	if req.Color == "" {
		return nil, errors.New("color is required")
	}

	weekObjectID, err := primitive.ObjectIDFromHex(weekColorTimeID)
	if err != nil {
		return nil, errors.New("invalid week colortime ID format")
	}

	week, err := s.DefaultColorTimeRepository.GetDefaultColorTimeWeekByID(ctx, weekObjectID)
	if err != nil {
		return nil, fmt.Errorf("failed to get week colortime: %w", err)
	}

	if week == nil {
		return nil, errors.New("week colortime not found")
	}

	var colorTimeForDate *DefaultColorTime
	for _, ct := range week.ColorTimes {
		if sameDay(ct.Date, dateParse) {
			colorTimeForDate = ct
			break
		}
	}

	if colorTimeForDate == nil {
		colorTimeForDate = &DefaultColorTime{
			ID:        primitive.NewObjectID(),
			Date:      dateParse,
			TimeSlots: []*DefaultColorBlock{},
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		week.ColorTimes = append(week.ColorTimes, colorTimeForDate)
	}

	endTime := startTimeParse.Add(time.Duration(req.Duration) * time.Minute)

	for _, block := range colorTimeForDate.TimeSlots {
		for _, slot := range block.Slots {
			if (startTimeParse.Before(slot.EndTime) && endTime.After(slot.StartTime)) ||
				(startTimeParse.Equal(slot.StartTime) || endTime.Equal(slot.EndTime)) {
				return nil, fmt.Errorf("slot time conflicts with existing slot [%s - %s]",
					slot.StartTime.Format("15:04"), slot.EndTime.Format("15:04"))
			}
		}
	}

	newSlot := &DefaultColortimeSlot{
		SlotID:    primitive.NewObjectID(),
		Sessions:  1,
		Title:     req.Title,
		StartTime: startTimeParse,
		EndTime:   endTime,
		Duration:  req.Duration,
		Color:     req.Color,
		Note:      req.Note,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	newColorBlock := &DefaultColorBlock{
		BlockID: primitive.NewObjectID(),
		Slots:   []*DefaultColortimeSlot{newSlot},
	}

	colorTimeForDate.TimeSlots = append(colorTimeForDate.TimeSlots, newColorBlock)
	week.UpdatedAt = time.Now()

	err = s.DefaultColorTimeRepository.UpdateDefaultColorTimeWeek(ctx, week.ID, week)
	if err != nil {
		return nil, fmt.Errorf("failed to update week colortime with new block: %w", err)
	}

	return newColorBlock, nil
}

func (s *defaultColorTimeService) CreateDefaultColorBlockForSession(
	ctx context.Context,
	weekColorTimeID string,
	req *CreateDefaultColorBlockWithSlotRequest,
	userID string,
) (*DefaultColorBlock, error) {

	if userID == "" {
		return nil, errors.New("user id is required")
	}

	if req.Date == "" {
		return nil, errors.New("date is required")
	}

	dateParse, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		return nil, fmt.Errorf("invalid date format: %w", err)
	}

	if req.StartTime == "" {
		return nil, errors.New("start time is required")
	}

	startTimeParse, err := time.Parse("15:04", req.StartTime)
	if err != nil {
		return nil, fmt.Errorf("invalid start time format: %w", err)
	}

	if req.Duration <= 0 {
		return nil, errors.New("duration must be greater than 0")
	}

	if req.Title == "" {
		return nil, errors.New("title is required")
	}

	if req.Color == "" {
		return nil, errors.New("color is required")
	}

	weekObjectID, err := primitive.ObjectIDFromHex(weekColorTimeID)
	if err != nil {
		return nil, errors.New("invalid week colortime ID format")
	}

	week, err := s.DefaultColorTimeRepository.GetDefaultColorTimeWeekByID(ctx, weekObjectID)
	if err != nil {
		return nil, fmt.Errorf("failed to get week colortime: %w", err)
	}

	if week == nil {
		return nil, errors.New("week colortime not found")
	}

	var colorTimeForDate *DefaultColorTime
	for _, ct := range week.ColorTimes {
		if sameDay(ct.Date, dateParse) {
			colorTimeForDate = ct
			break
		}
	}

	if colorTimeForDate == nil {
		colorTimeForDate = &DefaultColorTime{
			ID:        primitive.NewObjectID(),
			Date:      dateParse,
			TimeSlots: []*DefaultColorBlock{},
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		week.ColorTimes = append(week.ColorTimes, colorTimeForDate)
	}

	endTime := startTimeParse.Add(time.Duration(req.Duration) * time.Minute)

	for _, block := range colorTimeForDate.TimeSlots {
		for _, slot := range block.Slots {
			if (startTimeParse.Before(slot.EndTime) && endTime.After(slot.StartTime)) ||
				startTimeParse.Equal(slot.StartTime) || endTime.Equal(slot.EndTime) {
				return nil, fmt.Errorf("slot time conflicts with existing slot [%s - %s]",
					slot.StartTime.Format("15:04"), slot.EndTime.Format("15:04"))
			}
		}
	}

	newSlot := &DefaultColortimeSlot{
		SlotID:    primitive.NewObjectID(),
		Title:     req.Title,
		Sessions:  1,
		StartTime: startTimeParse,
		EndTime:   endTime,
		Duration:  req.Duration,
		Color:     req.Color,
		Note:      req.Note,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	var targetBlock *DefaultColorBlock
	if req.BlockID != "" {
		blockObjID, err := primitive.ObjectIDFromHex(req.BlockID)
		if err != nil {
			return nil, fmt.Errorf("invalid block_id format: %w", err)
		}

		for _, b := range colorTimeForDate.TimeSlots {
			if b.BlockID == blockObjID {
				targetBlock = b
				break
			}
		}

		if targetBlock == nil {
			return nil, fmt.Errorf("block not found for given block_id")
		}

		newSlot.Sessions = len(targetBlock.Slots) + 1
		targetBlock.Slots = append(targetBlock.Slots, newSlot)

	} else {
		newBlock := &DefaultColorBlock{
			BlockID: primitive.NewObjectID(),
			Slots:   []*DefaultColortimeSlot{newSlot},
		}
		colorTimeForDate.TimeSlots = append(colorTimeForDate.TimeSlots, newBlock)
		targetBlock = newBlock
	}

	colorTimeForDate.UpdatedAt = time.Now()
	week.UpdatedAt = time.Now()

	if err := s.DefaultColorTimeRepository.UpdateDefaultColorTimeWeek(ctx, week.ID, week); err != nil {
		return nil, fmt.Errorf("failed to update week colortime with new block: %w", err)
	}

	return targetBlock, nil
}

func (s *defaultColorTimeService) UpdateDefaultColorSlot(ctx context.Context, weekID, slotID string, req *UpdateDefaultColorSlotRequest) error {
	weekObjectID, err := primitive.ObjectIDFromHex(weekID)
	if err != nil {
		return errors.New("invalid week id format")
	}

	slotObjectID, err := primitive.ObjectIDFromHex(slotID)
	if err != nil {
		return errors.New("invalid slot id format")
	}

	week, err := s.DefaultColorTimeRepository.GetDefaultColorTimeWeekByID(ctx, weekObjectID)
	if err != nil {
		return err
	}

	if week == nil {
		return fmt.Errorf("week not found")
	}

	var found bool
	for _, day := range week.ColorTimes {
		for _, block := range day.TimeSlots {
			for _, slot := range block.Slots {
				if slot.SlotID == slotObjectID {
					if req.Title != "" {
						slot.Title = req.Title
					}
					if req.Color != "" {
						slot.Color = req.Color
					}
					if req.Note != "" {
						slot.Note = req.Note
					}
					if req.StartTime != "" {
						startTimeParse, err := time.Parse("15:04", req.StartTime)
						if err != nil {
							return fmt.Errorf("invalid start time format: %w", err)
						}
						slot.StartTime = startTimeParse
					}
					if req.Duration > 0 {
						slot.Duration = req.Duration
						slot.EndTime = slot.StartTime.Add(time.Duration(req.Duration) * time.Minute)
					}
					slot.UpdatedAt = time.Now()
					found = true
					break
				}
			}
			if found {
				break
			}
		}
		if found {
			break
		}
	}

	if !found {
		return fmt.Errorf("slot not found")
	}

	week.UpdatedAt = time.Now()
	return s.DefaultColorTimeRepository.UpdateDefaultColorTimeWeek(ctx, week.ID, week)
}

func (s *defaultColorTimeService) DeleteDefaultColorSlot(ctx context.Context, weekID, slotID string) error {
	weekObjectID, err := primitive.ObjectIDFromHex(weekID)
	if err != nil {
		return errors.New("invalid week id format")
	}

	slotObjectID, err := primitive.ObjectIDFromHex(slotID)
	if err != nil {
		return errors.New("invalid slot id format")
	}

	week, err := s.DefaultColorTimeRepository.GetDefaultColorTimeWeekByID(ctx, weekObjectID)
	if err != nil {
		return err
	}

	if week == nil {
		return fmt.Errorf("week not found")
	}

	var found bool
	for _, day := range week.ColorTimes {
		for _, block := range day.TimeSlots {
			for i, slot := range block.Slots {
				if slot.SlotID == slotObjectID {
					block.Slots = append(block.Slots[:i], block.Slots[i+1:]...)
					found = true
					break
				}
			}
			if found {
				break
			}
		}
		if found {
			break
		}
	}

	if !found {
		return fmt.Errorf("slot not found")
	}

	week.UpdatedAt = time.Now()
	return s.DefaultColorTimeRepository.UpdateDefaultColorTimeWeek(ctx, week.ID, week)
}
