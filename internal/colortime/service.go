package colortime

import (
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

	CreateColorBlockAndSaveSlot(ctx context.Context, weekColorTimeID string, req *CreateColorBlockWithSlotRequest, userID string) (*ColorBlock, error)
	CreateColorBlockForSession(ctx context.Context, weekColorTimeID string, req *CreateColorBlockWithSlotRequest, userID string) (*ColorBlock, error)
}

type colorTimeService struct {
	ColorTimeRepository ColorTimeRepository
	ProductService      product.ProductService
	LanguageService     language.MessageLanguageGateway
	UserService         user.UserService
	TopicService        topic.TopicService
}

func NewColorTimeService(colorTimeRepository ColorTimeRepository,
	productService product.ProductService,
	languageService language.MessageLanguageGateway,
	userService user.UserService,
	topicService topic.TopicService) ColorTimeService {
	return &colorTimeService{
		ColorTimeRepository: colorTimeRepository,
		ProductService:      productService,
		LanguageService:     languageService,
		UserService:         userService,
		TopicService:        topicService,
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

	colortimeWeek, err := s.ColorTimeRepository.GetColorTimeWeek(ctx, &startDate, &endDate, orgID, userID, role)
	if err != nil {
		return nil, err
	}

	if colortimeWeek == nil {
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
		week.ColorTimes = append(week.ColorTimes, &ColorTime{
			ID:        primitive.NewObjectID(),
			Date:      dateParse,
			TopicID:   &req.TopicID,
			TimeSlots: []*ColorBlock{},
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		})
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

func (s *colorTimeService) CreateColorBlockAndSaveSlot(
	ctx context.Context,
	weekColorTimeID string,
	req *CreateColorBlockWithSlotRequest,
	userID string,
) (*ColorBlock, error) {

	if userID == "" {
		return nil, errors.New("user id is required")
	}

	if req.OrganizationID == "" {
		return nil, errors.New("organization id is required")
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

	if req.Tracking == "" {
		return nil, fmt.Errorf("tracking is required")
	}

	// Get the WeekColorTime
	weekObjectID, err := primitive.ObjectIDFromHex(weekColorTimeID)
	if err != nil {
		return nil, errors.New("invalid week colortime ID format")
	}

	week, err := s.ColorTimeRepository.GetColorTimeWeekByID(ctx, weekObjectID)
	if err != nil {
		return nil, fmt.Errorf("failed to get week colortime: %w", err)
	}

	if week == nil {
		return nil, errors.New("week colortime not found")
	}

	// Find or create the ColorTime for the specific date
	var colorTimeForDate *ColorTime
	for _, ct := range week.ColorTimes {
		if sameDay(ct.Date, dateParse) {
			colorTimeForDate = ct
			break
		}
	}

	if colorTimeForDate == nil {
		colorTimeForDate = &ColorTime{
			ID:        primitive.NewObjectID(),
			Date:      dateParse,
			TopicID:   nil, // Assuming topic is set at week or day level, not block
			TimeSlots: []*ColorBlock{},
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		week.ColorTimes = append(week.ColorTimes, colorTimeForDate)
	}

	// Create the new slot
	endTime := startTimeParse.Add(time.Duration(req.Duration) * time.Minute)

	productID := &req.ProductID
	if req.ProductID == "" {
		productID = nil
	}

	for _, block := range colorTimeForDate.TimeSlots {
		for _, slot := range block.Slots {
			// nếu thời gian slot mới trùng khoảng với slot cũ
			if (startTimeParse.Before(slot.EndTime) && endTime.After(slot.StartTime)) ||
				(startTimeParse.Equal(slot.StartTime) || endTime.Equal(slot.EndTime)) {
				return nil, fmt.Errorf("slot time conflicts with existing slot [%s - %s]",
					slot.StartTime.Format("15:04"), slot.EndTime.Format("15:04"))
			}
		}
	}

	newSlot := &ColortimeSlot{
		SlotID:    primitive.NewObjectID(),
		Sessions:  1, // Always 1 for new slots in this flow
		Title:     req.Title,
		Tracking:  req.Tracking,
		UseCount:  0,
		StartTime: startTimeParse,
		EndTime:   endTime,
		Duration:  req.Duration,
		Color:     req.Color,
		Note:      req.Note,
		ProductID: productID,
		UpdatedAt: time.Now(),
	}

	// Create a new ColorBlock with the single slot
	newColorBlock := &ColorBlock{
		BlockID: primitive.NewObjectID(),
		Slots:   []*ColortimeSlot{newSlot},
	}

	colorTimeForDate.TimeSlots = append(colorTimeForDate.TimeSlots, newColorBlock)
	week.UpdatedAt = time.Now()

	err = s.ColorTimeRepository.UpdateColorTimeWeek(ctx, week.ID, week)
	if err != nil {
		return nil, fmt.Errorf("failed to update week colortime with new block: %w", err)
	}

	return newColorBlock, nil
}

func (s *colorTimeService) CreateColorBlockForSession(
	ctx context.Context,
	weekColorTimeID string,
	req *CreateColorBlockWithSlotRequest,
	userID string,
) (*ColorBlock, error) {

	if userID == "" {
		return nil, errors.New("user id is required")
	}

	if req.OrganizationID == "" {
		return nil, errors.New("organization id is required")
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

	if req.Tracking == "" {
		return nil, fmt.Errorf("tracking is required")
	}

	weekObjectID, err := primitive.ObjectIDFromHex(weekColorTimeID)
	if err != nil {
		return nil, errors.New("invalid week colortime ID format")
	}

	week, err := s.ColorTimeRepository.GetColorTimeWeekByID(ctx, weekObjectID)
	if err != nil {
		return nil, fmt.Errorf("failed to get week colortime: %w", err)
	}

	if week == nil {
		return nil, errors.New("week colortime not found")
	}

	var colorTimeForDate *ColorTime
	for _, ct := range week.ColorTimes {
		if sameDay(ct.Date, dateParse) {
			colorTimeForDate = ct
			break
		}
	}

	if colorTimeForDate == nil {
		colorTimeForDate = &ColorTime{
			ID:        primitive.NewObjectID(),
			Date:      dateParse,
			TimeSlots: []*ColorBlock{},
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		week.ColorTimes = append(week.ColorTimes, colorTimeForDate)
	}

	endTime := startTimeParse.Add(time.Duration(req.Duration) * time.Minute)
	var productID *string
	if req.ProductID != "" {
		productID = &req.ProductID
	}

	for _, block := range colorTimeForDate.TimeSlots {
		for _, slot := range block.Slots {
			if (startTimeParse.Before(slot.EndTime) && endTime.After(slot.StartTime)) ||
				startTimeParse.Equal(slot.StartTime) || endTime.Equal(slot.EndTime) {
				return nil, fmt.Errorf("slot time conflicts with existing slot [%s - %s]",
					slot.StartTime.Format("15:04"), slot.EndTime.Format("15:04"))
			}
		}
	}

	newSlot := &ColortimeSlot{
		SlotID:    primitive.NewObjectID(),
		Title:     req.Title,
		Tracking:  req.Tracking,
		Sessions:  1,
		StartTime: startTimeParse,
		EndTime:   endTime,
		Duration:  req.Duration,
		Color:     req.Color,
		Note:      req.Note,
		ProductID: productID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	var targetBlock *ColorBlock
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
		newBlock := &ColorBlock{
			BlockID: primitive.NewObjectID(),
			Slots:   []*ColortimeSlot{newSlot},
		}
		colorTimeForDate.TimeSlots = append(colorTimeForDate.TimeSlots, newBlock)
		targetBlock = newBlock
	}

	colorTimeForDate.UpdatedAt = time.Now()
	week.UpdatedAt = time.Now()

	if err := s.ColorTimeRepository.UpdateColorTimeWeek(ctx, week.ID, week); err != nil {
		return nil, fmt.Errorf("failed to update week colortime with new block: %w", err)
	}

	return targetBlock, nil
	
}
