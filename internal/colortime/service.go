package colortime

import (
	"colortime-service/internal/default_colortime"
	"colortime-service/internal/language"
	"colortime-service/internal/product"
	"colortime-service/internal/topic"
	"colortime-service/internal/user"
	"context"
	"errors"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ColorTimeService interface {
	AddTopicToColorTimeWeek(ctx context.Context, id string, req *AddTopicToColorTimeWeekRequest, userID string) error
	GetColorTimeWeek(ctx context.Context, userID, role, orgID, start, end string) (*TopicToColorTimeWeekResponse, error)
	DeleteTopicToColorTimeWeek(ctx context.Context, id string) error
	AddTopicToColorTimeDay(ctx context.Context, id string, req *AddTopicToColorTimeDayRequest) error
	DeleteTopicToColorTimeDay(ctx context.Context, id string, req *DeleteTopicToColorTimeDayRequest) error

	UpdateColorSlot(ctx context.Context, weekColorTimeID, slotID string, req *UpdateColorSlotRequest, userID string) error
}

type colorTimeService struct {
	ColorTimeRepository        ColorTimeRepository
	DefaultColorTimeRepository default_colortime.DefaultColorTimeRepository
	ProductService             product.ProductService
	LanguageService            language.MessageLanguageGateway
	UserService                user.UserService
	TopicService               topic.TopicService
}

func NewColorTimeService(colorTimeRepository ColorTimeRepository,
	defaultColorTimeRepository default_colortime.DefaultColorTimeRepository,
	productService product.ProductService,
	languageService language.MessageLanguageGateway,
	userService user.UserService,
	topicService topic.TopicService) ColorTimeService {
	return &colorTimeService{
		ColorTimeRepository:        colorTimeRepository,
		DefaultColorTimeRepository: defaultColorTimeRepository,
		ProductService:             productService,
		LanguageService:            languageService,
		UserService:                userService,
		TopicService:               topicService,
	}
}

func (s *colorTimeService) AddTopicToColorTimeWeek(ctx context.Context, id string, req *AddTopicToColorTimeWeekRequest, userID string) error {

	if userID == "" {
		return errors.New("user id is required")
	}

	if req.TopicID == "" {
		return errors.New("topic id is required")
	}

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("invalid id format")
	}

	colortimeWeek, err := s.ColorTimeRepository.GetColorTimeWeekByID(ctx, objectID)
	if err != nil {
		return err
	}

	if colortimeWeek == nil {
		return fmt.Errorf("color time week not found")
	}

	colortimeWeek.TopicID = &req.TopicID
	colortimeWeek.UpdatedAt = time.Now()

	if err := s.ColorTimeRepository.UpdateColorTimeWeek(ctx, colortimeWeek.ID, colortimeWeek); err != nil {
		return err
	}

	return nil

}

func (s *colorTimeService) GetColorTimeWeek(ctx context.Context, userID, role, orgID, start, end string) (*TopicToColorTimeWeekResponse, error) {

	if userID == "" {
		return nil, errors.New("user id is required")
	}

	if role == "" {
		return nil, errors.New("role is required")
	}

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

	defaultDayColorTimes, err := s.DefaultColorTimeRepository.GetDefaultDayColorTimesInRange(ctx, startDate, endDate, orgID)
	if err != nil {
		return nil, fmt.Errorf("failed to get default day colortimes: %w", err)
	}
	fmt.Println("defaultDayColorTimes", defaultDayColorTimes)
	var colortimeWeek *WeekColorTime

	if len(defaultDayColorTimes) > 0 {
		colorTimes := cloneDefaultDayColorTimesToColorTimes(defaultDayColorTimes)
		owner := &Owner{
			OwnerID:   userID,
			OwnerRole: role,
		}
		newColorTimeWeek := &WeekColorTime{
			ID:             primitive.NewObjectID(),
			OrganizationID: orgID,
			Owner:          owner,
			StartDate:      startDate,
			EndDate:        endDate,
			TopicID:        nil,
			ColorTimes:     colorTimes,
			CreatedBy:      userID,
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
		}

		if err := s.ColorTimeRepository.CreateColorTimeWeek(ctx, newColorTimeWeek); err != nil {
			return nil, err
		}
		colortimeWeek = newColorTimeWeek
	} else {
		owner := &Owner{
			OwnerID:   userID,
			OwnerRole: role,
		}
		newColorTimeWeek := &WeekColorTime{
			ID:             primitive.NewObjectID(),
			OrganizationID: orgID,
			Owner:          owner,
			StartDate:      startDate,
			EndDate:        endDate,
			TopicID:        nil,
			ColorTimes:     []*ColorTime{}, 
			CreatedBy:      userID,
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
		}

		if err := s.ColorTimeRepository.CreateColorTimeWeek(ctx, newColorTimeWeek); err != nil {
			return nil, err
		}
		colortimeWeek = newColorTimeWeek
	}

	var weekTopic Topic
	if colortimeWeek.TopicID != nil && *colortimeWeek.TopicID != "" {
		topic, err := s.TopicService.GetTopicInfor(ctx, *colortimeWeek.TopicID)
		if err != nil {
			return nil, err
		}
		if topic != nil {
			weekTopic = Topic{
				ID:   topic.ID,
				Name: topic.Name,
			}
		}
	}

	colorTimeResponses := make([]*ColorTimeResponse, 0, len(colortimeWeek.ColorTimes))

	for _, day := range colortimeWeek.ColorTimes {
		var dayTopic Topic
		if day.TopicID != nil && *day.TopicID != "" {
			topic, err := s.TopicService.GetTopicInfor(ctx, *day.TopicID)
			if err != nil {
				return nil, fmt.Errorf("failed to fetch topic for day %v: %w", day.Date, err)
			}
			if topic != nil {
				dayTopic = Topic{
					ID:   topic.ID,
					Name: topic.Name,
				}
			}
		}

		colorTimeResponses = append(colorTimeResponses, &ColorTimeResponse{
			ID:        day.ID,
			Date:      day.Date,
			Topic:     dayTopic,
			TimeSlots: day.TimeSlots,
			CreatedAt: day.CreatedAt,
			UpdatedAt: day.UpdatedAt,
		})
	}

	var studentInfor *user.UserInfor
	if colortimeWeek.Owner != nil {
		student, err := s.UserService.GetStudentInfor(ctx, colortimeWeek.Owner.OwnerID)
		if err != nil {
			return nil, err
		}
		if student != nil {
			studentInfor = student
		} else {
			studentInfor = &user.UserInfor{}
		}
	}

	result := &TopicToColorTimeWeekResponse{
		ID:             colortimeWeek.ID,
		OrganizationID: colortimeWeek.OrganizationID,
		Owner:          studentInfor,
		StartDate:      colortimeWeek.StartDate,
		EndDate:        colortimeWeek.EndDate,
		Topic:          weekTopic,
		ColorTimes:     colorTimeResponses,
		CreatedBy:      colortimeWeek.CreatedBy,
		CreatedAt:      colortimeWeek.CreatedAt,
		UpdatedAt:      colortimeWeek.UpdatedAt,
	}

	return result, nil
}

func (s *colorTimeService) DeleteTopicToColorTimeWeek(ctx context.Context, id string) error {

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("invalid id format")
	}

	colortimeWeek, err := s.ColorTimeRepository.GetColorTimeWeekByID(ctx, objectID)
	if err != nil {
		return err
	}

	if colortimeWeek == nil {
		return fmt.Errorf("color time week not found")
	} else {
		colortimeWeek.TopicID = nil
		colortimeWeek.UpdatedAt = time.Now()
		if err := s.ColorTimeRepository.UpdateColorTimeWeek(ctx, colortimeWeek.ID, colortimeWeek); err != nil {
			return err
		}
	}

	return nil

}

func (s *colorTimeService) AddTopicToColorTimeDay(ctx context.Context, id string, req *AddTopicToColorTimeDayRequest) error {

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("invalid id format")
	}

	if req.Date == "" {
		return fmt.Errorf("date is required")
	}

	dateParse, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		return err
	}

	if req.TopicID == "" {
		return fmt.Errorf("topic id is required")
	}

	week, err := s.ColorTimeRepository.GetColorTimeWeekByID(ctx, objectID)
	if err != nil {
		return err
	}

	if week == nil {
		return fmt.Errorf("color time day not found")
	}

	found := false
	for _, day := range week.ColorTimes {
		if sameDay(day.Date, dateParse) {
			day.TopicID = &req.TopicID
			day.UpdatedAt = time.Now()
			found = true
			break
		}
	}

	if !found {
		return fmt.Errorf("color time day not found")
	}

	week.UpdatedAt = time.Now()
	return s.ColorTimeRepository.UpdateColorTimeWeek(ctx, week.ID, week)

}

func (s *colorTimeService) DeleteTopicToColorTimeDay(ctx context.Context, id string, req *DeleteTopicToColorTimeDayRequest) error {

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("invalid id format")
	}

	if req.Date == "" {
		return fmt.Errorf("date is required")
	}

	dateParse, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		return fmt.Errorf("invalid date format: %w", err)
	}

	week, err := s.ColorTimeRepository.GetColorTimeWeekByID(ctx, objectID)
	if err != nil {
		return err
	}

	if week == nil {
		return fmt.Errorf("week not found")
	}

	found := false
	for i := range week.ColorTimes {
		if sameDay(week.ColorTimes[i].Date, dateParse) {
			week.ColorTimes[i].TopicID = nil
			week.ColorTimes[i].UpdatedAt = time.Now()
			found = true
			break
		}
	}

	if !found {
		return fmt.Errorf("color time day not found")
	}

	week.UpdatedAt = time.Now()
	return s.ColorTimeRepository.UpdateColorTimeWeek(ctx, week.ID, week)
}

func sameDay(a, b time.Time) bool {
	y1, m1, d1 := a.Date()
	y2, m2, d2 := b.Date()
	return y1 == y2 && m1 == m2 && d1 == d2
}

func cloneDefaultDayColorTimesToColorTimes(defaultDayColorTimes []*default_colortime.DefaultDayColorTime) []*ColorTime {
	colorTimes := make([]*ColorTime, 0, len(defaultDayColorTimes))
	for _, defaultDay := range defaultDayColorTimes {
		colorBlocks := make([]*ColorBlock, 0, len(defaultDay.TimeSlots))

		for _, defaultBlock := range defaultDay.TimeSlots {
			colorSlots := make([]*ColortimeSlot, 0, len(defaultBlock.Slots))

			for _, defaultSlot := range defaultBlock.Slots {
				colorSlot := &ColortimeSlot{
					SlotID:    primitive.NewObjectID(),
					Sessions:  defaultSlot.Sessions,
					Title:     defaultSlot.Title,
					Tracking:  "",
					UseCount:  0,
					StartTime: defaultSlot.StartTime,
					EndTime:   defaultSlot.EndTime,
					Duration:  defaultSlot.Duration,
					Color:     defaultSlot.Color,
					Note:      defaultSlot.Note,
					ProductID: nil,
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				}
				colorSlots = append(colorSlots, colorSlot)
			}

			colorBlock := &ColorBlock{
				BlockID: primitive.NewObjectID(),
				Slots:   colorSlots,
			}
			colorBlocks = append(colorBlocks, colorBlock)
		}

		colorTime := &ColorTime{
			ID:        primitive.NewObjectID(),
			Date:      defaultDay.Date,
			TopicID:   nil,
			TimeSlots: colorBlocks,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		colorTimes = append(colorTimes, colorTime)
	}

	return colorTimes
}

func (s *colorTimeService) UpdateColorSlot(ctx context.Context, weekColorTimeID, slotID string, req *UpdateColorSlotRequest, userID string) error {

	if userID == "" {
		return errors.New("user id is required")
	}

	if req.Tracking == "" {
		return errors.New("tracking is required")
	}

	weekObjectID, err := primitive.ObjectIDFromHex(weekColorTimeID)
	if err != nil {
		return errors.New("invalid week colortime ID format")
	}

	slotObjectID, err := primitive.ObjectIDFromHex(slotID)
	if err != nil {
		return errors.New("invalid slot ID format")
	}

	week, err := s.ColorTimeRepository.GetColorTimeWeekByID(ctx, weekObjectID)
	if err != nil {
		return fmt.Errorf("failed to get week colortime: %w", err)
	}

	if week == nil {
		return errors.New("week colortime not found")
	}

	// Find target slot in all blocks and days
	var targetSlot *ColortimeSlot
	for _, colorTime := range week.ColorTimes {
		for _, block := range colorTime.TimeSlots {
			for _, slot := range block.Slots {
				if slot.SlotID == slotObjectID {
					targetSlot = slot
					break
				}
			}
			if targetSlot != nil {
				break
			}
		}
		if targetSlot != nil {
			break
		}
	}

	if targetSlot == nil {
		return errors.New("slot not found")
	}

	// Update product_id
	targetSlot.ProductID = nil
	if req.ProductID != "" {
		targetSlot.ProductID = &req.ProductID
	}

	// Count how many slots have the same tracking as the new tracking
	trackingCount := 0
	for _, colorTime := range week.ColorTimes {
		for _, block := range colorTime.TimeSlots {
			for _, slot := range block.Slots {
				if slot.Tracking == req.Tracking {
					trackingCount++
				}
			}
		}
	}

	// Set tracking and use count (number of times this tracking appears)
	targetSlot.Tracking = req.Tracking
	targetSlot.UseCount = trackingCount

	targetSlot.UpdatedAt = time.Now()
	week.UpdatedAt = time.Now()

	if err := s.ColorTimeRepository.UpdateColorTimeWeek(ctx, week.ID, week); err != nil {
		return fmt.Errorf("failed to update week colortime: %w", err)
	}

	return nil
}
