package colortime

import (
	"colortime-service/internal/language"
	"colortime-service/internal/product"
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
	AddTopicToColorTimeWeek(ctx context.Context, req *AddTopicToColorTimeWeekRequest, userID string) error
	DeleteTopicToColorTimeWeek(ctx context.Context, req *DeleteTopicToColorTimeWeekRequest, userID string) error

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
}

func NewColorTimeService(colorTimeRepository ColorTimeRepository,
	productService product.ProductService,
	languageService language.MessageLanguageGateway,
	userService user.UserService) ColorTimeService {
	return &colorTimeService{
		ColorTimeRepository: colorTimeRepository,
		ProductService:      productService,
		LanguageService:     languageService,
		UserService:         userService,
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
			CreatedBy: userID,
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
				CreatedBy: userID,
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

func (s *colorTimeService) AddTopicToColorTimeWeek(ctx context.Context, req *AddTopicToColorTimeWeekRequest, userID string) error {

	if req.TopicID == "" {
		return errors.New("topic id is required")
	}

	if req.Owner == nil {
		return errors.New("owner is required")
	}

	if req.OrganizationID == "" {
		return errors.New("organization id is required")
	}

	if req.StartDate == "" {
		return errors.New("start date is required")
	}

	if userID == "" {
		return errors.New("user id is required")
	}

	if req.EndDate == "" {
		return errors.New("end date is required")
	}

	startDate, err := time.Parse("2006-01-02", req.StartDate)
	if err != nil {
		return err
	}

	endDate, err := time.Parse("2006-01-02", req.EndDate)
	if err != nil {
		return err
	}

	colortimeWeek, err := s.ColorTimeRepository.GetColorTimeWeek(ctx, &startDate, &endDate, req.OrganizationID, req.Owner.OwnerID)
	if err != nil {
		return err
	}

	if colortimeWeek == nil {
		newColorTimeWeek := &WeekColorTime{
			ID:             primitive.NewObjectID(),
			OrganizationID: req.OrganizationID,
			Owner:          req.Owner,
			StartDate:      startDate,
			EndDate:        endDate,
			TopicID:        &req.TopicID,
			ColorTimes:     make([]*TemplateSlot, 0),
			CreatedBy:      userID,
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
		}
		if err := s.ColorTimeRepository.CreateColorTimeWeek(ctx, newColorTimeWeek); err != nil {
			return err
		}
	} else {
		colortimeWeek.TopicID = &req.TopicID
		colortimeWeek.UpdatedAt = time.Now()
		if err := s.ColorTimeRepository.UpdateColorTimeWeek(ctx, colortimeWeek.ID, colortimeWeek); err != nil {
			return err
		}
	}

	return nil

}

func (s *colorTimeService) DeleteTopicToColorTimeWeek(ctx context.Context, req *DeleteTopicToColorTimeWeekRequest, userID string) error {

	if req.Owner == nil {
		return errors.New("owner is required")
	}

	if req.OrganizationID == "" {
		return errors.New("organization id is required")
	}

	if req.StartDate == "" {
		return errors.New("start date is required")
	}

	if userID == "" {
		return errors.New("user id is required")
	}

	if req.EndDate == "" {
		return errors.New("end date is required")
	}

	startDate, err := time.Parse("2006-01-02", req.StartDate)
	if err != nil {
		return err
	}

	endDate, err := time.Parse("2006-01-02", req.EndDate)
	if err != nil {
		return err
	}

	colortimeWeek, err := s.ColorTimeRepository.GetColorTimeWeek(ctx, &startDate, &endDate, req.OrganizationID, req.Owner.OwnerID)
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
