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
	CreateDefaultDayColorTime(ctx context.Context, req *CreateDefaultDayColorTimeRequest, userID string) (*DefaultDayColorTimeResponse, error)
	GetDefaultDayColorTime(ctx context.Context, orgID, date, userID string, languageID *int) (*DefaultDayColorTimeResponse, error)
	GetDefaultDayColorTimesInRange(ctx context.Context, orgID, startDate, endDate, userID string, languageID *int) ([]*DefaultDayColorTimeResponse, error)
	GetAllDefaultDayColorTimes(ctx context.Context, orgID string) ([]*DefaultDayColorTimeResponse, error)
	GetBlockBySlotID(ctx context.Context, dayID, slotID string) (*BlockWithSlotResponse, error)
	UpdateDefaultColorSlot(ctx context.Context, dayID, slotID string, req *UpdateDefaultColorSlotRequest) error
	DeleteDefaultDayColorTime(ctx context.Context, id string) error
	DeleteDefaultDayColorTimeSlot(ctx context.Context, dayID, slotID string, userID string) error
	DeleteDefaultDayColorTimeBlock(ctx context.Context, dayID, blockID string, userID string) error
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

func isTimeSlotConflict(newStart, newEnd time.Time, existingSlots []*DefaultColortimeSlot, excludeSlotID *primitive.ObjectID) bool {
	for _, slot := range existingSlots {
		if excludeSlotID != nil && slot.SlotID == *excludeSlotID {
			continue
		}

		if newStart.Before(slot.EndTime) && newEnd.After(slot.StartTime) {
			return true
		}
	}
	return false
}

func (s *defaultColorTimeService) CreateDefaultDayColorTime(ctx context.Context, req *CreateDefaultDayColorTimeRequest, userID string) (*DefaultDayColorTimeResponse, error) {

	if userID == "" {
		return nil, errors.New("user id is required")
	}

	if req.OrganizationID == "" {
		return nil, errors.New("organization id is required")
	}

	if req.Date == "" {
		return nil, errors.New("date is required")
	}

	if req.StartTime == "" {
		return nil, errors.New("start_time is required")
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

	if req.ColorTimeSlotLanguage != nil {
		if req.ColorTimeSlotLanguage.LanguageID == 0 {
			return nil, errors.New("language id is required")
		}
		if req.ColorTimeSlotLanguage.Title == "" {
			return nil, errors.New("title is required")
		}
	}

	date, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		return nil, fmt.Errorf("invalid date format: %w", err)
	}

	startTime, err := time.Parse("15:04", req.StartTime)
	if err != nil {
		return nil, fmt.Errorf("invalid start_time format (use HH:MM): %w", err)
	}

	endTime := startTime.Add(time.Duration(req.Duration) * time.Second)

	repeatType := req.RepeatType
	if repeatType == "" {
		repeatType = "none"
	}

	repeatInterval := req.RepeatInterval
	if repeatInterval == 0 {
		repeatInterval = 1
	}

	var repeatUntil *time.Time
	if req.RepeatUntil != "" {
		if parsedUntil, err := time.Parse("2006-01-02", req.RepeatUntil); err == nil {
			repeatUntil = &parsedUntil
		}
	}

	dates := []time.Time{date}
	maxEnd := time.Time{}
	if repeatUntil != nil {
		maxEnd = *repeatUntil
	}

	if repeatType != "none" && !maxEnd.IsZero() {
		switch repeatType {
		case "daily":
			for d := date.AddDate(0, 0, repeatInterval); !d.After(maxEnd); d = d.AddDate(0, 0, repeatInterval) {
				dates = append(dates, d)
			}
		case "weekly":
			for d := date.AddDate(0, 0, 7*repeatInterval); !d.After(maxEnd); d = d.AddDate(0, 0, 7*repeatInterval) {
				dates = append(dates, d)
			}
		case "monthly":
			for d := date.AddDate(0, repeatInterval, 0); !d.After(maxEnd); d = d.AddDate(0, repeatInterval, 0) {
				dates = append(dates, d)
			}
		case "custom":
			for _, off := range req.RepeatDays {
				d := date.AddDate(0, 0, off)
				if !d.After(maxEnd) {
					if !sameDay(d, date) {
						dates = append(dates, d)
					}
				}
			}
		default:

		}
	}

	baseSlot := &DefaultColortimeSlot{
		SlotID:                primitive.NewObjectID(),
		Sessions:              1,
		Title:                 req.Title,
		ColorTimeSlotLanguage: []*DefaultColorTimeSlotLanguage{req.ColorTimeSlotLanguage},
		StartTime:             startTime,
		EndTime:               endTime,
		Duration:              req.Duration,
		Color:                 req.Color,
		Note:                  req.Note,
		CreatedAt:             time.Now(),
		UpdatedAt:             time.Now(),
	}

	var baseDayResult *DefaultDayColorTime
	var baseBlockID *primitive.ObjectID

	for idx, d := range dates {
		isBase := sameDay(d, date)
		existingDay, err := s.DefaultColorTimeRepository.GetDefaultDayColorTime(ctx, d, req.OrganizationID)
		if err != nil {
			return nil, err
		}

		var dayColorTime *DefaultDayColorTime
		if existingDay != nil {
			dayColorTime = existingDay
		} else {
			dayColorTime = &DefaultDayColorTime{
				ID:             primitive.NewObjectID(),
				OrganizationID: req.OrganizationID,
				Date:           d,
				TimeSlots:      []*DefaultColorBlock{},
				CreatedBy:      userID,
				CreatedAt:      time.Now(),
				UpdatedAt:      time.Now(),

				IsBaseTemplate: isBase,
				RepeatType:     repeatType,
				RepeatUntil:    repeatUntil,
				RepeatInterval: repeatInterval,
				RepeatDays:     req.RepeatDays,
			}
		}

		var targetBlock *DefaultColorBlock
		if idx == 0 {
			if req.BlockID != "" {
				id, err := primitive.ObjectIDFromHex(req.BlockID)
				if err == nil {
					baseBlockID = &id
				} else {
					newBlock := &DefaultColorBlock{
						BlockID: primitive.NewObjectID(),
						Slots:   []*DefaultColortimeSlot{},
					}
					dayColorTime.TimeSlots = append(dayColorTime.TimeSlots, newBlock)
					baseBlockID = &newBlock.BlockID
				}
			} else {
				newBlock := &DefaultColorBlock{
					BlockID: primitive.NewObjectID(),
					Slots:   []*DefaultColortimeSlot{},
				}
				dayColorTime.TimeSlots = append(dayColorTime.TimeSlots, newBlock)
				baseBlockID = &newBlock.BlockID
			}

			for _, b := range dayColorTime.TimeSlots {
				if baseBlockID != nil && b.BlockID == *baseBlockID {
					targetBlock = b
					break
				}
			}
		} else {
			for _, b := range dayColorTime.TimeSlots {
				if baseBlockID != nil && b.BlockID == *baseBlockID {
					targetBlock = b
					break
				}
			}

			if targetBlock == nil {
				if baseBlockID == nil {
					targetBlock = &DefaultColorBlock{
						BlockID: primitive.NewObjectID(),
						Slots:   []*DefaultColortimeSlot{},
					}
				} else {
					targetBlock = &DefaultColorBlock{
						BlockID: *baseBlockID,
						Slots:   []*DefaultColortimeSlot{},
					}
				}
				dayColorTime.TimeSlots = append(dayColorTime.TimeSlots, targetBlock)
			}
		}

		var slotToAdd *DefaultColortimeSlot
		if isBase {
			slotToAdd = baseSlot
			slotToAdd.Sessions = len(targetBlock.Slots) + 1
		} else {
			slotCopy := &DefaultColortimeSlot{
				SlotID:                primitive.NewObjectID(),
				Sessions:              len(targetBlock.Slots) + 1,
				Title:                 baseSlot.Title,
				ColorTimeSlotLanguage: baseSlot.ColorTimeSlotLanguage,
				StartTime:             baseSlot.StartTime,
				EndTime:               baseSlot.EndTime,
				Duration:              baseSlot.Duration,
				Color:                 baseSlot.Color,
				Note:                  baseSlot.Note,
				CreatedAt:             time.Now(),
				UpdatedAt:             time.Now(),
			}
			slotToAdd = slotCopy
		}

		// Check for conflicts with existing slots across entire day
		var allSlots []*DefaultColortimeSlot
		for _, b := range dayColorTime.TimeSlots {
			allSlots = append(allSlots, b.Slots...)
		}

		// Validate time slot conflict across entire day
		if isTimeSlotConflict(slotToAdd.StartTime, slotToAdd.EndTime, allSlots, nil) {
			return nil, errors.New("time slot conflicts with existing slots in the day")
		}

		targetBlock.Slots = append(targetBlock.Slots, slotToAdd)

		if existingDay == nil {
			if err := s.DefaultColorTimeRepository.CreateDefaultDayColorTime(ctx, dayColorTime); err != nil {
				return nil, fmt.Errorf("failed to create day %s: %w", d.Format("2006-01-02"), err)
			}
		} else {
			dayColorTime.UpdatedAt = time.Now()
			if err := s.DefaultColorTimeRepository.UpdateDefaultDayColorTime(ctx, dayColorTime.ID, dayColorTime); err != nil {
				return nil, fmt.Errorf("failed to update day %s: %w", d.Format("2006-01-02"), err)
			}
		}

		if isBase || idx == 0 {
			baseDayResult = dayColorTime
		}
	}

	if baseDayResult == nil {
		baseDayResult, err = s.DefaultColorTimeRepository.GetDefaultDayColorTime(ctx, date, req.OrganizationID)
		if err != nil {
			return nil, err
		}
		if baseDayResult == nil {
			return nil, fmt.Errorf("failed to retrieve base day after creation")
		}
	}

	result := &DefaultDayColorTimeResponse{
		ID:             baseDayResult.ID,
		OrganizationID: baseDayResult.OrganizationID,
		Date:           baseDayResult.Date,
		TimeSlots:      baseDayResult.TimeSlots,
		IsBaseTemplate: baseDayResult.IsBaseTemplate,
		RepeatType:     baseDayResult.RepeatType,
		RepeatUntil:    baseDayResult.RepeatUntil,
		RepeatInterval: baseDayResult.RepeatInterval,
		RepeatDays:     baseDayResult.RepeatDays,
		CreatedBlockID: baseBlockID,
		CreatedBy:      baseDayResult.CreatedBy,
		CreatedAt:      baseDayResult.CreatedAt,
		UpdatedAt:      baseDayResult.UpdatedAt,
	}

	return result, nil
}

func (s *defaultColorTimeService) GetDefaultDayColorTime(ctx context.Context, orgID, date, userID string, languageID *int) (*DefaultDayColorTimeResponse, error) {
	if orgID == "" {
		return nil, errors.New("organization id is required")
	}

	if date == "" {
		return nil, errors.New("date is required")
	}

	parsedDate, err := time.Parse("2006-01-02", date)
	if err != nil {
		return nil, fmt.Errorf("invalid date format: %w", err)
	}

	dayColorTime, err := s.DefaultColorTimeRepository.GetDefaultDayColorTime(ctx, parsedDate, orgID)
	if err != nil {
		return nil, err
	}

	if dayColorTime == nil {
		return nil, fmt.Errorf("default day color time not found for date: %s", date)
	}

	if languageID != nil {
		for _, block := range dayColorTime.TimeSlots {
			for _, slot := range block.Slots {
				var filteredLanguageTitle []*DefaultColorTimeSlotLanguage
				for _, lang := range slot.ColorTimeSlotLanguage {
					if lang.LanguageID == *languageID {
						filteredLanguageTitle = append(filteredLanguageTitle, lang)
					}
				}
				slot.ColorTimeSlotLanguage = filteredLanguageTitle
			}
		}
	}

	response := &DefaultDayColorTimeResponse{
		ID:             dayColorTime.ID,
		OrganizationID: dayColorTime.OrganizationID,
		Date:           dayColorTime.Date,
		TimeSlots:      dayColorTime.TimeSlots,
		IsBaseTemplate: dayColorTime.IsBaseTemplate,
		RepeatType:     dayColorTime.RepeatType,
		RepeatUntil:    dayColorTime.RepeatUntil,
		RepeatInterval: dayColorTime.RepeatInterval,
		RepeatDays:     dayColorTime.RepeatDays,
		CreatedBy:      dayColorTime.CreatedBy,
		CreatedAt:      dayColorTime.CreatedAt,
		UpdatedAt:      dayColorTime.UpdatedAt,
	}

	return response, nil
}

func (s *defaultColorTimeService) GetDefaultDayColorTimesInRange(ctx context.Context, orgID, startDate, endDate, userID string, languageID *int) ([]*DefaultDayColorTimeResponse, error) {
	if orgID == "" {
		return nil, errors.New("organization id is required")
	}

	if startDate == "" || endDate == "" {
		return nil, errors.New("start date and end date are required")
	}

	startParsed, err := time.Parse("2006-01-02", startDate)
	if err != nil {
		return nil, fmt.Errorf("invalid start date format: %w", err)
	}

	endParsed, err := time.Parse("2006-01-02", endDate)
	if err != nil {
		return nil, fmt.Errorf("invalid end date format: %w", err)
	}

	dayColorTimes, err := s.DefaultColorTimeRepository.GetDefaultDayColorTimesInRange(ctx, startParsed, endParsed, orgID)
	if err != nil {
		return nil, err
	}

	var responses []*DefaultDayColorTimeResponse
	for _, day := range dayColorTimes {
		if languageID != nil {
			for _, block := range day.TimeSlots {
				for _, slot := range block.Slots {
					var filteredLanguageTitle []*DefaultColorTimeSlotLanguage
					for _, lang := range slot.ColorTimeSlotLanguage {
						if lang.LanguageID == *languageID {
							filteredLanguageTitle = append(filteredLanguageTitle, lang)
						}
					}
					slot.ColorTimeSlotLanguage = filteredLanguageTitle
				}
			}
		}
		response := &DefaultDayColorTimeResponse{
			ID:             day.ID,
			OrganizationID: day.OrganizationID,
			Date:           day.Date,
			TimeSlots:      day.TimeSlots,
			IsBaseTemplate: day.IsBaseTemplate,
			RepeatType:     day.RepeatType,
			RepeatUntil:    day.RepeatUntil,
			RepeatInterval: day.RepeatInterval,
			RepeatDays:     day.RepeatDays,
			CreatedBy:      day.CreatedBy,
			CreatedAt:      day.CreatedAt,
			UpdatedAt:      day.UpdatedAt,
		}
		responses = append(responses, response)
	}

	return responses, nil
}

func (s *defaultColorTimeService) GetAllDefaultDayColorTimes(ctx context.Context, orgID string) ([]*DefaultDayColorTimeResponse, error) {
	if orgID == "" {
		return nil, errors.New("organization id is required")
	}

	dayColorTimes, err := s.DefaultColorTimeRepository.GetAllDefaultDayColorTimes(ctx, orgID)
	if err != nil {
		return nil, err
	}

	var responses []*DefaultDayColorTimeResponse
	for _, day := range dayColorTimes {
		response := &DefaultDayColorTimeResponse{
			ID:             day.ID,
			OrganizationID: day.OrganizationID,
			Date:           day.Date,
			TimeSlots:      day.TimeSlots,
			IsBaseTemplate: day.IsBaseTemplate,
			RepeatType:     day.RepeatType,
			RepeatUntil:    day.RepeatUntil,
			RepeatInterval: day.RepeatInterval,
			RepeatDays:     day.RepeatDays,
			CreatedBy:      day.CreatedBy,
			CreatedAt:      day.CreatedAt,
			UpdatedAt:      day.UpdatedAt,
		}
		responses = append(responses, response)
	}

	return responses, nil
}

func (s *defaultColorTimeService) DeleteDefaultDayColorTime(ctx context.Context, id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("invalid id format")
	}

	return s.DefaultColorTimeRepository.DeleteDefaultDayColorTime(ctx, objID)
}

func (s *defaultColorTimeService) GetBlockBySlotID(ctx context.Context, dayID, slotID string) (*BlockWithSlotResponse, error) {

	dayObjectID, err := primitive.ObjectIDFromHex(dayID)
	if err != nil {
		return nil, errors.New("invalid day_id format")
	}

	day, err := s.DefaultColorTimeRepository.GetDefaultDayColorTimeByID(ctx, dayObjectID)
	if err != nil {
		return nil, fmt.Errorf("failed to get day: %w", err)
	}

	if day == nil {
		return nil, fmt.Errorf("day not found")
	}

	slotObjectID, err := primitive.ObjectIDFromHex(slotID)
	if err != nil {
		return nil, errors.New("invalid slot_id format")
	}

	day, err = s.DefaultColorTimeRepository.GetDefaultDayColorTimeBySlotID(ctx, slotObjectID)
	if err != nil {
		return nil, fmt.Errorf("failed to find day containing slot: %w", err)
	}

	if day == nil {
		return nil, fmt.Errorf("slot not found")
	}

	for _, block := range day.TimeSlots {
		for _, slot := range block.Slots {
			if slot.SlotID == slotObjectID {
				return &BlockWithSlotResponse{
					Block:   block,
					Slot:    slot,
					DayID:   day.ID,
					DayDate: day.Date,
					OrgID:   day.OrganizationID,
				}, nil
			}
		}
	}

	return nil, fmt.Errorf("slot not found in day")
}

func (s *defaultColorTimeService) UpdateDefaultColorSlot(ctx context.Context, dayID, slotID string, req *UpdateDefaultColorSlotRequest) error {

	dayObjectID, err := primitive.ObjectIDFromHex(dayID)
	if err != nil {
		return errors.New("invalid day_id format")
	}

	slotObjectID, err := primitive.ObjectIDFromHex(slotID)
	if err != nil {
		return errors.New("invalid slot_id format")
	}

	day, err := s.DefaultColorTimeRepository.GetDefaultDayColorTimeByID(ctx, dayObjectID)
	if err != nil {
		return fmt.Errorf("failed to get day: %w", err)
	}

	if day == nil {
		return fmt.Errorf("day not found")
	}

	if req.ColorTimeSlotLanguage != nil {
		if req.ColorTimeSlotLanguage.LanguageID == 0 {
			return errors.New("language id is required")
		}
		if req.ColorTimeSlotLanguage.Title == "" {
			return errors.New("title is required")
		}
	}

	var slotFound bool
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

				var newStartTime *time.Time
				if req.StartTime != "" {
					parsedTime, err := time.Parse("15:04", req.StartTime)
					if err != nil {
						return fmt.Errorf("invalid start_time format (use HH:MM): %w", err)
					}
					newStartTime = &parsedTime
					slot.StartTime = time.Date(slot.StartTime.Year(), slot.StartTime.Month(), slot.StartTime.Day(),
						parsedTime.Hour(), parsedTime.Minute(), 0, 0, slot.StartTime.Location())
				}

				if req.Duration > 0 {
					slot.Duration = req.Duration
				}

				if newStartTime != nil || req.Duration > 0 {
					slot.EndTime = slot.StartTime.Add(time.Duration(slot.Duration) * time.Second)
				}

				if req.ColorTimeSlotLanguage != nil {
					languageExists := false
					for i, lang := range slot.ColorTimeSlotLanguage {
						if lang.LanguageID == req.ColorTimeSlotLanguage.LanguageID {
							slot.ColorTimeSlotLanguage[i].Title = req.ColorTimeSlotLanguage.Title
							languageExists = true
							break
						}
					}
					if !languageExists {
						slot.ColorTimeSlotLanguage = append(slot.ColorTimeSlotLanguage, req.ColorTimeSlotLanguage)
					}
				}

				slot.UpdatedAt = time.Now()
				slotFound = true
				break
			}
		}
		if slotFound {
			break
		}
	}

	if !slotFound {
		return fmt.Errorf("slot not found in day")
	}

	day.UpdatedAt = time.Now()
	return s.DefaultColorTimeRepository.UpdateDefaultDayColorTime(ctx, dayObjectID, day)
}

func (s *defaultColorTimeService) DeleteDefaultDayColorTimeSlot(ctx context.Context, dayID, slotID string, userID string) error {
	if dayID == "" {
		return errors.New("day id is required")
	}
	if slotID == "" {
		return errors.New("slot id is required")
	}
	if userID == "" {
		return errors.New("user id is required")
	}

	dayObjectID, err := primitive.ObjectIDFromHex(dayID)
	if err != nil {
		return errors.New("invalid day id format")
	}

	slotObjectID, err := primitive.ObjectIDFromHex(slotID)
	if err != nil {
		return errors.New("invalid slot id format")
	}

	day, err := s.DefaultColorTimeRepository.GetDefaultDayColorTimeByID(ctx, dayObjectID)
	if err != nil {
		return fmt.Errorf("failed to get day: %w", err)
	}

	if day == nil {
		return fmt.Errorf("day not found")
	}

	for _, block := range day.TimeSlots {
		for i, slot := range block.Slots {
			if slot.SlotID == slotObjectID {
				block.Slots = append(block.Slots[:i], block.Slots[i+1:]...)
				break
			}
		}
	}

	if err := s.DefaultColorTimeRepository.UpdateDefaultDayColorTime(ctx, dayObjectID, day); err != nil {
		return fmt.Errorf("failed to update day: %w", err)
	}

	return nil
}

func (s *defaultColorTimeService) DeleteDefaultDayColorTimeBlock(ctx context.Context, dayID, blockID string, userID string) error {
	if dayID == "" {
		return errors.New("day id is required")
	}

	if blockID == "" {
		return errors.New("block id is required")
	}

	if userID == "" {
		return errors.New("user id is required")
	}

	dayObjectID, err := primitive.ObjectIDFromHex(dayID)
	if err != nil {
		return errors.New("invalid day id format")
	}

	blockObjectID, err := primitive.ObjectIDFromHex(blockID)
	if err != nil {
		return errors.New("invalid block id format")
	}

	day, err := s.DefaultColorTimeRepository.GetDefaultDayColorTimeByID(ctx, dayObjectID)
	if err != nil {
		return fmt.Errorf("failed to get day: %w", err)
	}

	if day == nil {
		return fmt.Errorf("day not found")
	}

	var targetBlockIndex = -1
	for i, block := range day.TimeSlots {
		if block.BlockID == blockObjectID {
			targetBlockIndex = i
			break
		}
	}

	if targetBlockIndex == -1 {
		return errors.New("block not found")
	}

	day.TimeSlots = append(day.TimeSlots[:targetBlockIndex], day.TimeSlots[targetBlockIndex+1:]...)

	if err := s.DefaultColorTimeRepository.UpdateDefaultDayColorTime(ctx, dayObjectID, day); err != nil {
		return fmt.Errorf("failed to update day: %w", err)
	}

	return nil
}

func sameDay(a, b time.Time) bool {
	y1, m1, d1 := a.Date()
	y2, m2, d2 := b.Date()
	return y1 == y2 && m1 == m2 && d1 == d2
}
