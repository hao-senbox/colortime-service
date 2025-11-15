package templatecolortime

import (
	"colortime-service/internal/default_colortime"
	"colortime-service/internal/term"
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TemplateColorTimeService interface {
	CreateTemplateColorTime(ctx context.Context, req CreateTemplateColorTimeRequest, userID string) (*TemplateColorTimeResponse, error)
	GetTemplateColorTime(ctx context.Context, organizationID, termID, date string) ([]*TemplateColorTime, error)
	UpdateTemplateColorTimeSlot(ctx context.Context, templateColorTimeID, slotID string, req *UpdateTemplateColorTimeSlotRequest, userID string) error
	DuplicateTemplateColorTime(ctx context.Context, req DuplicateTemplateColorTimeRequest, userID string) error
	ApplyTemplateColorTime(ctx context.Context, req ApplyTemplateColorTimeRequest, userID string) error
}

type templateColorTimeService struct {
	TemplateColorTimeRepository TemplateColorTimeRepository
	TermService                 term.TermService
	DefaultColorTimeRepository  default_colortime.DefaultColorTimeRepository
}

func NewTemplateColorTimeService(
	templateColorTimeRepository TemplateColorTimeRepository,
	termService term.TermService,
	defaultColorTimeRepository default_colortime.DefaultColorTimeRepository,
) TemplateColorTimeService {
	return &templateColorTimeService{
		TemplateColorTimeRepository: templateColorTimeRepository,
		TermService:                 termService,
		DefaultColorTimeRepository:  defaultColorTimeRepository,
	}
}

func (s *templateColorTimeService) CreateTemplateColorTime(ctx context.Context, req CreateTemplateColorTimeRequest, userID string) (*TemplateColorTimeResponse, error) {
	if req.OrganizationID == "" {
		return nil, errors.New("organization id is required")
	}

	if req.TermID == "" {
		return nil, errors.New("term id is required")
	}

	if req.Date == "" {
		return nil, errors.New("date is required")
	}

	if req.StartTime == "" {
		return nil, errors.New("start time is required")
	}

	if req.Duration <= 0 {
		return nil, errors.New("duration must be greater than 0")
	}

	if req.Color == "" {
		return nil, errors.New("color is required")
	}

	startTime, err := time.Parse("15:04", req.StartTime)
	if err != nil {
		return nil, errors.New("invalid start time format (use HH:MM)")
	}

	// Add duration in seconds directly to start time
	endTime := startTime.Add(time.Duration(req.Duration) * time.Second)

	if endTime.After(time.Date(2025, 12, 31, 23, 59, 59, 0, time.UTC)) {
		return nil, errors.New("end time must be before or equal to 2025-12-31 23:59:59")
	}

	if endTime.Before(startTime) {
		return nil, errors.New("end time must be after start time")
	}

	colortimeTemplateData, err := s.TemplateColorTimeRepository.GetTemplateColorTime(ctx, req.OrganizationID, req.TermID, req.Date)
	if err != nil {
		return nil, errors.New("failed to get template color time")
	}

	var baseBlockID *primitive.ObjectID
	
	if colortimeTemplateData == nil {
		colortimeTemplateData = &TemplateColorTime{
			ID:             primitive.NewObjectID(),
			OrganizationID: req.OrganizationID,
			TermID:         req.TermID,
			Date:           req.Date,
			ColorTimes:     []*ColorTimeTemplate{},
			CreatedBy:      userID,
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
		}

		block := &ColorTimeTemplate{
			BlockID: primitive.NewObjectID(),
			Slots:   []*ColortimeSlot{},
		}
		colortimeTemplateData.ColorTimes = append(colortimeTemplateData.ColorTimes, block)

		slot := &ColortimeSlot{
			SlotID:    primitive.NewObjectID(),
			Sessions:  1,
			Title:     req.Title,
			StartTime: startTime,
			EndTime:   endTime,
			Duration:  req.Duration / 60, // Convert to minutes for storage
			Color:     req.Color,
			Note:      req.Note,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		block.Slots = append(block.Slots, slot)

		baseBlockID = &block.BlockID

		if err := s.TemplateColorTimeRepository.CreateTemplateColorTime(ctx, colortimeTemplateData); err != nil {
			return nil, errors.New("failed to create template color time")
		}

		return &TemplateColorTimeResponse{
			ID:             colortimeTemplateData.ID,
			Date:           colortimeTemplateData.Date,
			OrganizationID: colortimeTemplateData.OrganizationID,
			TermID:         colortimeTemplateData.TermID,
			ColorTimes:     colortimeTemplateData.ColorTimes,
			CreatedBlockID: baseBlockID,
			CreatedBy:      colortimeTemplateData.CreatedBy,
			CreatedAt:      colortimeTemplateData.CreatedAt,
			UpdatedAt:      colortimeTemplateData.UpdatedAt,
		}, nil
	} else {
		if req.BlockID != "" {
			id, err := primitive.ObjectIDFromHex(req.BlockID)
			if err == nil {
				baseBlockID = &id
			}

			// Convert duration from seconds to minutes for database storage
			durationMinutes := req.Duration / 60
			baseSlot := &ColortimeSlot{
				SlotID:    primitive.NewObjectID(),
				Title:     req.Title,
				StartTime: startTime,
				EndTime:   startTime.Add(time.Duration(req.Duration) * time.Second), // Add seconds directly
				Duration:  durationMinutes,                                          // Store as minutes
				Color:     req.Color,
				Note:      req.Note,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			}

			for _, block := range colortimeTemplateData.ColorTimes {
				if baseBlockID != nil && block.BlockID == *baseBlockID {
					baseSlot.Sessions = len(block.Slots) + 1
					block.Slots = append(block.Slots, baseSlot)
					break
				}
			}
		} else {
			newBlock := &ColorTimeTemplate{
				BlockID: primitive.NewObjectID(),
				Slots:   []*ColortimeSlot{},
			}
			colortimeTemplateData.ColorTimes = append(colortimeTemplateData.ColorTimes, newBlock)
			baseBlockID = &newBlock.BlockID

			// Convert duration from seconds to minutes for database storage
			durationMinutes := req.Duration / 60
			baseSlot := &ColortimeSlot{
				SlotID:    primitive.NewObjectID(),
				Sessions:  1,
				Title:     req.Title,
				StartTime: startTime,
				EndTime:   startTime.Add(time.Duration(req.Duration) * time.Second), // Add seconds directly
				Duration:  durationMinutes,                                          // Store as minutes
				Color:     req.Color,
				Note:      req.Note,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			}

			newBlock.Slots = append(newBlock.Slots, baseSlot)
		}

		if err := s.TemplateColorTimeRepository.UpdateTemplateColorTime(ctx, colortimeTemplateData.ID, colortimeTemplateData); err != nil {
			return nil, errors.New("failed to update template color time")
		}
	}

	return &TemplateColorTimeResponse{
		ID:             colortimeTemplateData.ID,
		Date:           colortimeTemplateData.Date,
		OrganizationID: colortimeTemplateData.OrganizationID,
		TermID:         colortimeTemplateData.TermID,
		ColorTimes:     colortimeTemplateData.ColorTimes,
		CreatedBlockID: baseBlockID,
		CreatedBy:      colortimeTemplateData.CreatedBy,
		CreatedAt:      colortimeTemplateData.CreatedAt,
		UpdatedAt:      colortimeTemplateData.UpdatedAt,
	}, nil

}

func (s *templateColorTimeService) GetTemplateColorTime(ctx context.Context, organizationID, termID, date string) ([]*TemplateColorTime, error) {
	if organizationID == "" {
		return nil, errors.New("organization id is required")
	}

	if termID == "" {
		return nil, errors.New("term id is required")
	}

	var result []*TemplateColorTime

	if strings.ToLower(date) == "week" {
		weekdays := []string{"monday", "tuesday", "wednesday", "thursday", "friday", "saturday", "sunday"}

		for _, weekday := range weekdays {
			template, err := s.TemplateColorTimeRepository.GetTemplateColorTime(ctx, organizationID, termID, weekday)
			if err != nil {
				return nil, errors.New("failed to get template color time for " + weekday)
			}
			if template != nil {
				result = append(result, template)
			}
		}
	} else {
		template, err := s.TemplateColorTimeRepository.GetTemplateColorTime(ctx, organizationID, termID, date)
		if err != nil {
			return nil, errors.New("failed to get template color time")
		}
		if template != nil {
			result = append(result, template)
		}
	}

	return result, nil
}

func (s *templateColorTimeService) UpdateTemplateColorTimeSlot(ctx context.Context, templateColorTimeID, slotID string, req *UpdateTemplateColorTimeSlotRequest, userID string) error {
	if templateColorTimeID == "" {
		return errors.New("template color time id is required")
	}

	if slotID == "" {
		return errors.New("slot id is required")
	}

	if userID == "" {
		return errors.New("user id is required")
	}

	templateColorTimeObjectID, err := primitive.ObjectIDFromHex(templateColorTimeID)
	if err != nil {
		return errors.New("invalid template color time id format")
	}

	templateColorTime, err := s.TemplateColorTimeRepository.GetTemplateColorTimeByID(ctx, templateColorTimeObjectID)
	if err != nil {
		return errors.New("failed to get template color time")
	}

	if templateColorTime == nil {
		return errors.New("template color time not found")
	}

	slotObjectID, err := primitive.ObjectIDFromHex(slotID)
	if err != nil {
		return errors.New("invalid slot id format")
	}

	var targetSlot *ColortimeSlot

	for _, block := range templateColorTime.ColorTimes {
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

	if targetSlot == nil {
		return errors.New("slot not found")
	}
	var targetBlock *ColorTimeTemplate
	if req.BlockID != "" {
		blockObjectID, err := primitive.ObjectIDFromHex(req.BlockID)
		if err != nil {
			return errors.New("invalid block id format")
		}

		for _, block := range templateColorTime.ColorTimes {
			if block.BlockID == blockObjectID {
				targetBlock = block
				break
			}
		}

		if targetBlock == nil {
			return errors.New("block not found")
		}

		targetSlot.Sessions = len(targetBlock.Slots) + 1
		targetBlock.Slots = append(targetBlock.Slots, targetSlot)
	} else {
		if req.StartTime != "" {
			startTime, err := time.Parse("15:04", req.StartTime)
			if err != nil {
				return errors.New("invalid start time format (use HH:MM)")
			}
			targetSlot.StartTime = startTime
			targetSlot.EndTime = startTime.Add(time.Duration(req.Duration) * time.Second)
		}
		if req.Duration > 0 {
			// Convert duration from seconds to minutes for storage
			durationMinutes := req.Duration / 60
			targetSlot.Duration = durationMinutes
			targetSlot.EndTime = targetSlot.StartTime.Add(time.Duration(req.Duration) * time.Second)
		}
		if req.Title != "" {
			targetSlot.Title = req.Title
		}
		if req.Color != "" {
			targetSlot.Color = req.Color
		}
		if req.Note != "" {
			targetSlot.Note = req.Note
		}
		targetSlot.UpdatedAt = time.Now()
	}

	if err := s.TemplateColorTimeRepository.UpdateTemplateColorTime(ctx, templateColorTimeObjectID, templateColorTime); err != nil {
		return errors.New("failed to update template color time")
	}

	return nil
}

func (s *templateColorTimeService) DuplicateTemplateColorTime(ctx context.Context, req DuplicateTemplateColorTimeRequest, userID string) error {
	if req.OrganizationID == "" {
		return errors.New("organization id is required")
	}

	if req.TermID == "" {
		return errors.New("term id is required")
	}

	if req.OriginDate == "" {
		return errors.New("origin date is required")
	}

	templateColorTime, err := s.TemplateColorTimeRepository.GetTemplateColorTime(ctx, req.OrganizationID, req.TermID, req.OriginDate)
	if err != nil {
		return errors.New("failed to get template color time")
	}

	if templateColorTime == nil {
		return errors.New("template color time not found")
	}

	var targetDates []string

	if req.TargetDate != "" {
		targetDates = []string{req.TargetDate}
	} else {
		targetDates = getRemainingWeekdays(req.OriginDate)
	}

	for _, targetDate := range targetDates {
		if strings.EqualFold(targetDate, req.OriginDate) {
			continue
		}

		existingTarget, err := s.TemplateColorTimeRepository.GetTemplateColorTime(ctx, req.OrganizationID, req.TermID, targetDate)
		if err != nil {
			return errors.New("failed to check existing template color time")
		}

		if existingTarget != nil {
			continue
		}

		duplicateTemplate := s.createDuplicateTemplate(templateColorTime, targetDate, userID)

		if err := s.TemplateColorTimeRepository.CreateTemplateColorTime(ctx, duplicateTemplate); err != nil {
			return errors.New("failed to create template color time")
		}
	}

	return nil
}

func (s *templateColorTimeService) ApplyTemplateColorTime(ctx context.Context, req ApplyTemplateColorTimeRequest, userID string) error {
	if req.OrganizationID == "" {
		return errors.New("organization id is required")
	}

	if req.TermID == "" {
		return errors.New("term id is required")
	}

	var result []*TemplateColorTime

	weekdays := []string{"monday", "tuesday", "wednesday", "thursday", "friday", "saturday", "sunday"}

	for _, weekday := range weekdays {
		template, err := s.TemplateColorTimeRepository.GetTemplateColorTime(ctx, req.OrganizationID, req.TermID, weekday)
		if err != nil {
			return errors.New("failed to get template color time for " + weekday)
		}
		if template == nil {
			continue
		} else {
			result = append(result, template)
		}

	}

	term, err := s.TermService.GetTermByID(ctx, req.TermID)
	if err != nil {
		return errors.New("failed to get term")
	}

	if term == nil {
		return errors.New("term not found")
	}

	startDate, err := time.Parse("2006-01-02", req.StartDate)
	if err != nil {
		return errors.New("failed to parse start date")
	}

	endDate, err := time.Parse("2006-01-02", req.EndDate)
	if err != nil {
		return errors.New("failed to parse end date")
	}

	// Create a map of weekday to template for quick lookup
	templateMap := make(map[string]*TemplateColorTime)
	for _, template := range result {
		templateMap[strings.ToLower(template.Date)] = template
	}

	// Loop through all dates from start to end
	for currentDate := startDate; !currentDate.After(endDate); currentDate = currentDate.AddDate(0, 0, 1) {
		// Get weekday name in lowercase (convert Go's weekday to our format)
		weekdayName := strings.ToLower(currentDate.Weekday().String())

		// Get template for this weekday
		template, exists := templateMap[weekdayName]
		if !exists || template == nil {
			continue // Skip if no template for this weekday
		}

		// Check if default colortime already exists for this date
		existingDefaultColorTime, err := s.DefaultColorTimeRepository.GetDefaultDayColorTime(ctx, currentDate, req.OrganizationID)
		if err != nil {
			return errors.New("failed to get existing default color time")
		}

		if existingDefaultColorTime != nil {
			existingDefaultColorTime.TimeSlots = mergeTemplateIntoDefault(existingDefaultColorTime, template)
			existingDefaultColorTime.UpdatedAt = time.Now()
			if err := s.DefaultColorTimeRepository.UpdateDefaultDayColorTime(ctx, existingDefaultColorTime.ID, existingDefaultColorTime); err != nil {
				return errors.New("failed to update default color time")
			}
		} else {
			// Create new default colortime by copying template structure
			defaultColorTime := &default_colortime.DefaultDayColorTime{
				ID:             primitive.NewObjectID(),
				OrganizationID: req.OrganizationID,
				Date:           currentDate,
				TimeSlots:      []*default_colortime.DefaultColorBlock{},
				CreatedBy:      userID,
				CreatedAt:      time.Now(),
				UpdatedAt:      time.Now(),
				IsBaseTemplate: false,
				RepeatType:     "none",
			}

			// Copy all blocks and slots from template to default
			for _, templateBlock := range template.ColorTimes {
				colorBlock := &default_colortime.DefaultColorBlock{
					BlockID: templateBlock.BlockID, // Generate new block ID
					Slots:   []*default_colortime.DefaultColortimeSlot{},
				}

				for _, templateSlot := range templateBlock.Slots {
					colorSlot := &default_colortime.DefaultColortimeSlot{
						SlotID:    templateSlot.SlotID, // Generate new slot ID
						Sessions:  templateSlot.Sessions,
						Title:     templateSlot.Title,
						StartTime: templateSlot.StartTime,
						EndTime:   templateSlot.EndTime,
						Duration:  templateSlot.Duration,
						Color:     templateSlot.Color,
						Note:      templateSlot.Note,
						CreatedAt: time.Now(),
						UpdatedAt: time.Now(),
					}
					colorBlock.Slots = append(colorBlock.Slots, colorSlot)
				}

				defaultColorTime.TimeSlots = append(defaultColorTime.TimeSlots, colorBlock)
			}

			// Create the new default colortime document
			if err := s.DefaultColorTimeRepository.CreateDefaultDayColorTime(ctx, defaultColorTime); err != nil {
				continue // Continue to next date if creation fails
			}
		}
	}

	return nil
}

func getRemainingWeekdays(originWeekday string) []string {
	weekdays := []string{"monday", "tuesday", "wednesday", "thursday", "friday", "saturday", "sunday"}

	var remainingDays []string
	originIndex := -1

	for i, day := range weekdays {
		if strings.EqualFold(day, originWeekday) {
			originIndex = i
			break
		}
	}

	if originIndex == -1 {
		return remainingDays
	}

	for i := originIndex + 1; i < len(weekdays); i++ {
		remainingDays = append(remainingDays, weekdays[i])
	}

	return remainingDays
}

func (s *templateColorTimeService) createDuplicateTemplate(sourceTemplate *TemplateColorTime, targetDate string, userID string) *TemplateColorTime {
	duplicate := &TemplateColorTime{
		ID:             primitive.NewObjectID(),
		OrganizationID: sourceTemplate.OrganizationID,
		TermID:         sourceTemplate.TermID,
		Date:           targetDate,
		ColorTimes:     []*ColorTimeTemplate{},
		CreatedBy:      userID,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	for _, block := range sourceTemplate.ColorTimes {
		newBlock := &ColorTimeTemplate{
			BlockID: primitive.NewObjectID(),
			Slots:   []*ColortimeSlot{},
		}
		duplicate.ColorTimes = append(duplicate.ColorTimes, newBlock)

		for _, slot := range block.Slots {
			newSlot := &ColortimeSlot{
				SlotID:    primitive.NewObjectID(),
				Sessions:  slot.Sessions,
				Title:     slot.Title,
				StartTime: slot.StartTime,
				EndTime:   slot.EndTime,
				Duration:  slot.Duration,
				Color:     slot.Color,
				Note:      slot.Note,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			}
			newBlock.Slots = append(newBlock.Slots, newSlot)
		}
	}

	return duplicate
}

func mergeTemplateIntoDefault(existingDefaultColorTime *default_colortime.DefaultDayColorTime, template *TemplateColorTime) []*default_colortime.DefaultColorBlock {
	merged := []*default_colortime.DefaultColorBlock{}

	defaultMap := make(map[string]*default_colortime.DefaultColorBlock)
	for _, block := range existingDefaultColorTime.TimeSlots {
		defaultMap[block.BlockID.Hex()] = block
	}

	for _, block := range template.ColorTimes {
		defKey := block.BlockID.Hex()
		if db, exists := defaultMap[defKey]; exists {
			db.Slots = mergeSlots(db.Slots, block.Slots)
			merged = append(merged, db)
			delete(defaultMap, defKey)

		} else {
			newBlock := cloneTemplateBlockToDefault(block)
			merged = append(merged, newBlock)
		}
	}

for _, remaining := range defaultMap {
	fmt.Printf("→ Keeping existing block (not in template): %s\n", remaining.BlockID.Hex())
	merged = append(merged, remaining)
}

	return merged
}

func mergeSlots(existingSlots []*default_colortime.DefaultColortimeSlot, templateSlots []*ColortimeSlot) []*default_colortime.DefaultColortimeSlot {
	merged := []*default_colortime.DefaultColortimeSlot{}

	defaultMap := make(map[string]*default_colortime.DefaultColortimeSlot)
	for _, s := range existingSlots {
		key := s.SlotID.Hex()
		defaultMap[key] = s
	}

	for _, ts := range templateSlots {

		key := ts.SlotID.Hex()

		if ds, exists := defaultMap[key]; exists {

			// Update template fields but keep user custom fields
			ds.Title = ts.Title
			ds.StartTime = ts.StartTime
			ds.EndTime = ts.EndTime
			ds.Duration = ts.Duration
			ds.Color = ts.Color
			ds.Sessions = ts.Sessions
			ds.UpdatedAt = time.Now()

			merged = append(merged, ds)
			delete(defaultMap, key)

		} else {
			// new slot from template
			newSlot := cloneTemplateSlotToDefault(ts)
			merged = append(merged, newSlot)
		}
	}

	return merged
}

func cloneTemplateBlockToDefault(templateBlock *ColorTimeTemplate) *default_colortime.DefaultColorBlock {
	block := &default_colortime.DefaultColorBlock{
		BlockID: primitive.NewObjectID(), // Generate new ID for default
		Slots:   []*default_colortime.DefaultColortimeSlot{},
	}

	for _, templateSlot := range templateBlock.Slots {
		slot := cloneTemplateSlotToDefault(templateSlot)
		block.Slots = append(block.Slots, slot)
	}

	return block
}

func cloneTemplateSlotToDefault(templateSlot *ColortimeSlot) *default_colortime.DefaultColortimeSlot {
	return &default_colortime.DefaultColortimeSlot{
		SlotID:    primitive.NewObjectID(), // Generate new ID for default
		Sessions:  templateSlot.Sessions,
		Title:     templateSlot.Title,
		StartTime: templateSlot.StartTime,
		EndTime:   templateSlot.EndTime,
		Duration:  templateSlot.Duration,
		Color:     templateSlot.Color,
		Note:      templateSlot.Note,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}
