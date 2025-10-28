package colortime

import (
	"colortime-service/internal/language"
	"colortime-service/internal/product"
	"colortime-service/internal/topic"
	"colortime-service/internal/user"
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type ColorTimeService interface {
	CreateColorTime(ctx context.Context, req *CreateColorTimeRequest, userID string) error
	GetColorTimes(ctx context.Context, start string, end string) ([]*ColorTime, error)
	GetColorTime(ctx context.Context, id string) (*ColorTime, error)
	UpdateColorTime(ctx context.Context, req *UpdateColorTimeRequest, id string) error
	DeleteColorTime(ctx context.Context, req *DeleteColorTimeRequest, id string) error
	AddTopicToColorTimeWeek(ctx context.Context, id string, req *AddTopicToColorTimeWeekRequest, userID string) error
	GetColorTimeWeek(ctx context.Context, userID, role, orgID, start, end string) (*TopicToColorTimeWeekResponse, error)
	DeleteTopicToColorTimeWeek(ctx context.Context, id string) error
	AddTopicToColorTimeDay(ctx context.Context, id string, req *AddTopicToColorTimeDayRequest) error
	DeleteTopicToColorTimeDay(ctx context.Context, id string, req *DeleteTopicToColorTimeDayRequest) error

	CreateTemplateColorTime(ctx context.Context, req *CreateTemplateColorTimeRequest, userID string) error
	GetTemplateColorTimes(ctx context.Context) ([]*ColorTimeTemplate, error)
	GetTemplateColorTime(ctx context.Context, id string) (*ColorTimeTemplate, error)	
	UpdateTemplateColorTime(ctx context.Context, req *UpdateTemplateColorTimeRequest, id string) error
	DeleteTemplateColorTime(ctx context.Context, id string) error
	AddSlotsToTemplateColorTime(ctx context.Context, req *AddSlotsToTemplateColorTimeRequest, id string) error
	EditSlotsToTemplateColorTime(ctx context.Context, req *EditSlotsToTemplateColorTimeRequest, id string, slot_id string) error
	ApplyTemplateColorTime(ctx context.Context, req *ApplyTemplateColorTimeRequest, id string, userID string) error
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

func (s *colorTimeService) CreateColorTime(ctx context.Context, req *CreateColorTimeRequest, userID string) error {

	if req.Owner == nil {
		return errors.New("owner is required")
	} else {
		if req.Owner.OwnerID == "" {
			return errors.New("owner id is required")
		}

		if req.Owner.OwnerRole == "" {
			return errors.New("owner role is required")
		}
	}

	if req.Date == "" {
		return errors.New("date is required")
	}

	if req.OrganizationID == "" {
		return errors.New("organization id is required")
	}

	dateParse, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		return err
	}

	if req.Title == "" {
		return errors.New("title is required")
	}

	if req.Color == "" {
		return errors.New("color is required")
	}

	if req.Duration == 0 {
		return errors.New("duration is required")
	}

	if req.StartTime == "" {
		return errors.New("start_time is required")
	}

	if req.Tracking == "" {
		return errors.New("tracking is required")
	}

	var product_id *string
	if req.ProductID == "" {
		product_id = nil
	}

	startTimeParse, err := time.Parse("15:04", req.StartTime)
	if err != nil {
		return err
	}

	endTime := startTimeParse.Add(time.Duration(req.Duration) * time.Minute)

	if startTimeParse.Hour() != endTime.Hour() {
		return errors.New("start_time and end_time must be in the same hour")
	}

	colortime, err := s.ColorTimeRepository.GetColorTimeByDate(ctx, req.OrganizationID, dateParse)
	if err != nil {
		return err
	}

	if colortime == nil {
		data := &ColorTime{
			ID:      primitive.NewObjectID(),
			Date:    dateParse,
			TopicID: nil,
			TimeSlots: []*TemplateSlot{
				{
					SlotID:    primitive.NewObjectID(),
					Title:     req.Title,
					Tracking:  req.Tracking,
					UseCount:  0,
					StartTime: startTimeParse,
					EndTime:   endTime,
					Duration:  req.Duration,
					Color:     req.Color,
					Note:      *req.Note,
					ProductID: product_id,
				},
			},
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		return s.ColorTimeRepository.CreateColorTime(ctx, data)
	} else {
		var useCount int

		for _, slot := range colortime.TimeSlots {
			// if startTimeParse.Before(slot.EndTime) && endTime.After(slot.StartTime) {
			// 	return fmt.Errorf(
			// 		"slot overlaps with existing slot: %s (%s–%s)",
			// 		slot.Title,
			// 		slot.StartTime.Format("15:04"),
			// 		slot.EndTime.Format("15:04"),
			// 	)
			// }
			if slot.Tracking == req.Tracking {
				useCount++
			}
		}

		data := TemplateSlot{
			SlotID:    primitive.NewObjectID(),
			Title:     req.Title,
			Tracking:  req.Tracking,
			UseCount:  useCount,
			StartTime: startTimeParse,
			EndTime:   endTime,
			Duration:  req.Duration,
			Color:     req.Color,
			Note:      *req.Note,
			ProductID: product_id,
		}

		colortime.TimeSlots = append(colortime.TimeSlots, &data)

		return s.ColorTimeRepository.UpdateColorTime(ctx, colortime.ID, colortime)

	}
}

func (s *colorTimeService) GetColorTimes(ctx context.Context, start string, end string) ([]*ColorTime, error) {
	return s.ColorTimeRepository.GetColorTimes(ctx, start, end)
}

func (s *colorTimeService) GetColorTime(ctx context.Context, id string) (*ColorTime, error) {

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	return s.ColorTimeRepository.GetColorTime(ctx, objectID)

}

func (s *colorTimeService) UpdateColorTime(ctx context.Context, req *UpdateColorTimeRequest, id string) error {

	if req.Owner == nil {
		return errors.New("owner is required")
	} else {
		if req.Owner.OwnerID == "" {
			return errors.New("owner id is required")
		}

		if req.Owner.OwnerRole == "" {
			return errors.New("owner role is required")
		}
	}

	if req.Date == "" {
		return errors.New("date is required")
	}

	if req.OrganizationID == "" {
		return errors.New("organization id is required")
	}

	dateParse, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		return err
	}

	if req.Title == "" {
		return errors.New("title is required")
	}

	if req.Color == "" {
		return errors.New("color is required")
	}

	if req.Duration == 0 {
		return errors.New("duration is required")
	}

	if req.StartTime == "" {
		return errors.New("start_time is required")
	}

	if req.Tracking == "" {
		return errors.New("tracking is required")
	}

	if req.OrganizationID == "" {
		return errors.New("organization id is required")
	}

	var product_id *string
	if req.ProductID == "" {
		product_id = nil
	} else {
		product_id = &req.ProductID
	}

	startTimeParse, err := time.Parse("15:04", req.StartTime)
	if err != nil {
		return err
	}

	endTime := startTimeParse.Add(time.Duration(req.Duration) * time.Minute)

	if startTimeParse.Hour() != endTime.Hour() {
		return errors.New("start_time and end_time must be in the same hour")
	}

	colortime, err := s.ColorTimeRepository.GetColorTimeByDate(ctx, req.OrganizationID, dateParse)
	if err != nil {
		return err
	}

	if colortime == nil {
		return fmt.Errorf("color time not found")
	}

	slotObjectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("invalid id format")
	}

	var found bool
	var useCount int

	for i, slot := range colortime.TimeSlots {
		if slot.SlotID == slotObjectID {
			found = true
			for _, other := range colortime.TimeSlots {
				if other.Tracking == req.Tracking {
					useCount++
				}
			}

			colortime.TimeSlots[i].Title = req.Title
			colortime.TimeSlots[i].Color = req.Color
			colortime.TimeSlots[i].StartTime = startTimeParse
			colortime.TimeSlots[i].EndTime = endTime
			colortime.TimeSlots[i].Duration = req.Duration
			colortime.TimeSlots[i].Tracking = req.Tracking
			colortime.TimeSlots[i].UseCount = useCount
			colortime.TimeSlots[i].Note = req.Note
			colortime.TimeSlots[i].ProductID = product_id
			break
		}
	}

	if !found {
		return errors.New("slot not found")
	}

	return s.ColorTimeRepository.UpdateColorTime(ctx, colortime.ID, colortime)

}

func (s *colorTimeService) DeleteColorTime(ctx context.Context, req *DeleteColorTimeRequest, id string) error {

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("invalid id format")
	}

	if req.Date == "" {
		return errors.New("date is required")
	}

	if req.OrganizationID == "" {
		return errors.New("organization id is required")
	}

	dateParse, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		return err
	}

	colortime, err := s.ColorTimeRepository.GetColorTimeByDate(ctx, req.OrganizationID, dateParse)
	if err != nil {
		return err
	}

	if colortime == nil {
		return fmt.Errorf("color time not found")
	}

	var found bool
	for i, slot := range colortime.TimeSlots {
		if slot.SlotID == objectID {
			found = true
			colortime.TimeSlots = append(colortime.TimeSlots[:i], colortime.TimeSlots[i+1:]...)
			break
		}
	}

	if !found {
		return errors.New("slot not found")
	}

	return s.ColorTimeRepository.UpdateColorTime(ctx, colortime.ID, colortime)

}

func (s *colorTimeService) CreateTemplateColorTime(ctx context.Context, req *CreateTemplateColorTimeRequest, userID string) error {

	if req.Name == "" {
		return errors.New("name is required")
	}

	if userID == "" {
		return errors.New("user id is required")
	}

	user, err := s.UserService.GetCurrentUser(ctx)
	if err != nil {
		log.Printf("[colorTimeService] CreateTemplateColorTime: %v", err)
	}

	if user == nil {
		log.Printf("[colorTimeService] CreateTemplateColorTime: user not found")
	}

	var OrgID string
	if user != nil {
		OrgID = user.OrganizationAdmin.ID
	}

	data := &ColorTimeTemplate{
		ID:             primitive.NewObjectID(),
		OrganizationID: OrgID,
		Name:           req.Name,
		ColorTimes:     []*TemplateSlot{},
		CreatedBy:      userID,
		IsDeleted:      false,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	return s.ColorTimeRepository.CreateTemplateColorTime(ctx, data)

}

func (s *colorTimeService) GetTemplateColorTimes(ctx context.Context) ([]*ColorTimeTemplate, error) {
	return s.ColorTimeRepository.GetTemplateColorTimes(ctx)
}

func (s *colorTimeService) GetTemplateColorTime(ctx context.Context, id string) (*ColorTimeTemplate, error) {

	if id == "" {
		return nil, errors.New("id is required")
	}

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New("invalid id format")
	}

	return s.ColorTimeRepository.GetTemplateColorTime(ctx, objectID)

}

func (s *colorTimeService) UpdateTemplateColorTime(ctx context.Context, req *UpdateTemplateColorTimeRequest, id string) error {

	if id == "" {
		return errors.New("id is required")
	}

	if req.Name == "" {
		return errors.New("name is required")
	}

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("invalid id format")
	}

	template, err := s.ColorTimeRepository.GetTemplateColorTime(ctx, objectID)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return errors.New("template not found")
		}
		return err
	}

	template.Name = req.Name
	template.UpdatedAt = time.Now()

	return s.ColorTimeRepository.UpdateTemplateColorTime(ctx, template)

}

func (s *colorTimeService) DeleteTemplateColorTime(ctx context.Context, id string) error {

	if id == "" {
		return errors.New("id is required")
	}

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("invalid id format")
	}

	return s.ColorTimeRepository.DeleteTemplateColorTime(ctx, objectID)

}

func (s *colorTimeService) AddSlotsToTemplateColorTime(ctx context.Context, req *AddSlotsToTemplateColorTimeRequest, id string) error {

	if id == "" {
		return errors.New("id is required")
	}

	if req.Title == "" {
		return errors.New("title is required")
	}

	if req.Color == "" {
		return errors.New("color is required")
	}

	if req.StartTime == "" {
		return errors.New("start_time is required")
	}

	if req.Duration == 0 {
		return errors.New("duration is required")
	}

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("invalid id format")
	}

	startTimeParse, err := time.Parse("15:04", req.StartTime)
	if err != nil {
		return err
	}

	endTime := startTimeParse.Add(time.Duration(req.Duration) * time.Minute)

	if startTimeParse.Hour() != endTime.Hour() {
		return errors.New("start_time and end_time must be in the same hour")
	}

	template, err := s.ColorTimeRepository.GetTemplateColorTime(ctx, objectID)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return errors.New("template not found")
		}
		return err
	}

	if len(template.ColorTimes) <= 0 {

		data := &TemplateSlot{
			SlotID:    primitive.NewObjectID(),
			Title:     req.Title,
			Tracking:  req.Tracking,
			UseCount:  0,
			StartTime: startTimeParse,
			EndTime:   endTime,
			Duration:  req.Duration,
			Color:     req.Color,
			Note:      req.Note,
			ProductID: &req.ProductID,
		}

		template.ColorTimes = append(template.ColorTimes, data)

		return s.ColorTimeRepository.UpdateTemplateColorTime(ctx, template)
	} else {
		var useCount int

		for _, slot := range template.ColorTimes {
			// if startTimeParse.Before(slot.EndTime) && endTime.After(slot.StartTime) {
			// 	return fmt.Errorf(
			// 		"slot overlaps with existing slot: %s (%s–%s)",
			// 		slot.Title,
			// 		slot.StartTime.Format("15:04"),
			// 		slot.EndTime.Format("15:04"),
			// 	)
			// }
			if slot.Tracking == req.Tracking {
				useCount++
			}
		}

		data := &TemplateSlot{
			SlotID:    primitive.NewObjectID(),
			Title:     req.Title,
			Tracking:  req.Tracking,
			UseCount:  useCount,
			StartTime: startTimeParse,
			EndTime:   endTime,
			Duration:  req.Duration,
			Color:     req.Color,
			Note:      req.Note,
			ProductID: &req.ProductID,
		}

		template.ColorTimes = append(template.ColorTimes, data)
	}

	return s.ColorTimeRepository.UpdateTemplateColorTime(ctx, template)
}

func (s *colorTimeService) EditSlotsToTemplateColorTime(ctx context.Context, req *EditSlotsToTemplateColorTimeRequest, id string, slot_id string) error {

	if id == "" {
		return errors.New("id is required")
	}

	if req.Title == "" {
		return errors.New("title is required")
	}

	if req.Color == "" {
		return errors.New("color is required")
	}

	if req.StartTime == "" {
		return errors.New("start_time is required")
	}

	if req.Duration == 0 {
		return errors.New("duration is required")
	}

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("invalid id format")
	}

	slotOjectID, err := primitive.ObjectIDFromHex(slot_id)
	if err != nil {
		return errors.New("invalid id format")
	}

	startTimeParse, err := time.Parse("15:04", req.StartTime)
	if err != nil {
		return err
	}

	endTime := startTimeParse.Add(time.Duration(req.Duration) * time.Minute)

	if startTimeParse.Hour() != endTime.Hour() {
		return errors.New("start_time and end_time must be in the same hour")
	}

	template, err := s.ColorTimeRepository.GetTemplateColorTime(ctx, objectID)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return errors.New("template not found")
		}
		return err
	}

	if len(template.ColorTimes) <= 0 {
		return errors.New("template not found")
	}

	var found bool
	var useCount int
	for i, slot := range template.ColorTimes {
		if slot.SlotID == slotOjectID {
			found = true
			for _, other := range template.ColorTimes {
				// if startTimeParse.Before(other.EndTime) && endTime.After(other.StartTime) {
				// 	return fmt.Errorf(
				// 		"slot overlaps with existing slot: %s (%s–%s)",
				// 		other.Title,
				// 		other.StartTime.Format("15:04"),
				// 		other.EndTime.Format("15:04"),
				// 	)
				// }
				if other.Tracking == req.Tracking {
					useCount++
				}
			}

			template.ColorTimes[i].Title = req.Title
			template.ColorTimes[i].Tracking = req.Tracking
			template.ColorTimes[i].UseCount = useCount
			template.ColorTimes[i].Color = req.Color
			template.ColorTimes[i].StartTime = startTimeParse
			template.ColorTimes[i].EndTime = endTime
			template.ColorTimes[i].Duration = req.Duration
			template.ColorTimes[i].Note = req.Note
			template.ColorTimes[i].ProductID = &req.ProductID
			break
		}
	}

	if !found {
		return errors.New("slot with given tracking not found")
	}

	return s.ColorTimeRepository.UpdateTemplateColorTime(ctx, template)

}

func (s *colorTimeService) ApplyTemplateColorTime(ctx context.Context, req *ApplyTemplateColorTimeRequest, id string, userID string) error {

	if id == "" {
		return errors.New("id is required")
	}

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("invalid id format")
	}

	template, err := s.ColorTimeRepository.GetTemplateColorTime(ctx, objectID)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return errors.New("template not found")
		}
		return err
	}

	if len(template.ColorTimes) <= 0 {
		return errors.New("template not found")
	}

	if userID == "" {
		return errors.New("user id is required")
	}

	startDate, err := time.Parse("2006-01-02", req.StartDate)
	if err != nil {
		return err
	}

	endDate, err := time.Parse("2006-01-02", req.EndDate)
	if err != nil {
		return err
	}

	for startDate.Before(endDate) || startDate.Equal(endDate) {

		colortime, err := s.ColorTimeRepository.GetColorTimeByDate(ctx, template.OrganizationID, startDate)
		if err != nil {
			return err
		}
		if colortime == nil {
			newSlots := make([]*TemplateSlot, 0, len(template.ColorTimes))
			for _, slot := range template.ColorTimes {
				newSlot := *slot
				newSlots = append(newSlots, &newSlot)
			}

			newColorTime := &ColorTime{
				ID:        primitive.NewObjectID(),
				Date:      startDate,
				TopicID:   nil,
				TimeSlots: newSlots,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			}

			if err := s.ColorTimeRepository.CreateColorTime(ctx, newColorTime); err != nil {
				return err
			}
		} else {
			// existing := colortime.ColorTimes
			// for _, tmplSlot := range template.ColorTimes {
			// 	overwritten := false
			// 	for i, existSlot := range existing {
			// 		if tmplSlot.StartTime.Before(existSlot.EndTime) && tmplSlot.EndTime.After(existSlot.StartTime) {
			// 			existing[i] = tmplSlot
			// 			overwritten = true
			// 			break
			// 		}
			// 	}
			// 	if !overwritten {
			// 		existing = append(existing, tmplSlot)
			// 	}
			// }

			// colortime.ColorTimes = existing
			// colortime.UpdatedAt = time.Now()

			// if err := s.ColorTimeRepository.UpdateColorTime(ctx, colortime.ID, colortime); err != nil {
			// 	return err
			// }
			for _, tmplSlot := range template.ColorTimes {
				newSlot := *tmplSlot
				colortime.TimeSlots = append(colortime.TimeSlots, &newSlot)
			}

			colortime.UpdatedAt = time.Now()

			if err := s.ColorTimeRepository.UpdateColorTime(ctx, colortime.ID, colortime); err != nil {
				return err
			}

		}
		startDate = startDate.AddDate(0, 0, 1)
	}

	return nil
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

	result := &TopicToColorTimeWeekResponse{
		ID:             colortimeWeek.ID,
		OrganizationID: colortimeWeek.OrganizationID,
		Owner:          colortimeWeek.Owner,
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
			TimeSlots: []*TemplateSlot{},
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
