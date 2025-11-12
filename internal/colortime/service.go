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
	"sort"
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

	var colortimeWeek *WeekColorTime

	existingWeek, err := s.ColorTimeRepository.GetColorTimeWeek(ctx, &startDate, &endDate, orgID, userID, role)
	if err != nil {
		return nil, fmt.Errorf("failed to check existing week: %w", err)
	}

	if len(defaultDayColorTimes) > 0 {
		colorTimes := cloneDefaultDayColorTimesToColorTimes(defaultDayColorTimes)

		if existingWeek != nil {
			updatedColorTimes := s.mergeColorTimes(existingWeek.ColorTimes, colorTimes)
			existingWeek.ColorTimes = updatedColorTimes
			existingWeek.UpdatedAt = time.Now()

			if err := s.ColorTimeRepository.UpdateColorTimeWeek(ctx, existingWeek.ID, existingWeek); err != nil {
				return nil, fmt.Errorf("failed to update existing week: %w", err)
			}
			colortimeWeek = existingWeek
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
				ColorTimes:     colorTimes,
				CreatedBy:      userID,
				CreatedAt:      time.Now(),
				UpdatedAt:      time.Now(),
			}

			if err := s.ColorTimeRepository.CreateColorTimeWeek(ctx, newColorTimeWeek); err != nil {
				return nil, err
			}
			colortimeWeek = newColorTimeWeek
		}
	} else {
		if existingWeek != nil {
			// Return existing week if it already exists
			colortimeWeek = existingWeek
		} else {
			// Create new empty week only if it doesn't exist
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
		blockResponses, err := s.convertBlocksWithProductInfo(ctx, day.TimeSlots)
		if err != nil {
			return nil, fmt.Errorf("failed to convert blocks for day %v: %w", day.Date, err)
		}

		colorTimeResponses = append(colorTimeResponses, &ColorTimeResponse{
			ID:        day.ID,
			Date:      day.Date,
			Topic:     dayTopic,
			TimeSlots: blockResponses,
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
					SlotID:    defaultSlot.SlotID,
					SlotIDOld: &defaultSlot.SlotID,
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
				BlockID:    defaultBlock.BlockID,
				BlockIDOld: &defaultBlock.BlockID,
				Slots:      colorSlots,
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

	oldTracking := targetSlot.Tracking
	targetTracking := req.Tracking

	targetSlot.ProductID = nil
	if req.ProductID != "" {
		targetSlot.ProductID = &req.ProductID
	}

	trackingCount, err := s.ColorTimeRepository.CountTrackingUsage(ctx, week.OrganizationID, week.Owner.OwnerID, week.Owner.OwnerRole, req.Tracking)
	if err != nil {
		return fmt.Errorf("failed to count tracking usage: %w", err)
	}

	targetSlot.Tracking = req.Tracking
	targetSlot.UseCount = trackingCount + 1

	targetSlot.UpdatedAt = time.Now()
	week.UpdatedAt = time.Now()

	if err := s.ColorTimeRepository.UpdateColorTimeWeek(ctx, week.ID, week); err != nil {
		return fmt.Errorf("failed to update week colortime: %w", err)
	}

	if oldTracking != targetTracking && oldTracking != "" {
		if err := s.normalizeTrackingGlobal(ctx, week.OrganizationID, week.Owner.OwnerID, week.Owner.OwnerRole, oldTracking); err != nil {
			return fmt.Errorf("failed to normalize tracking global: %w", err)
		}
		if err := s.normalizeTrackingGlobal(ctx, week.OrganizationID, week.Owner.OwnerID, week.Owner.OwnerRole, targetTracking); err != nil {
			return fmt.Errorf("failed to normalize tracking global: %w", err)
		}
	}

	return nil
}

func (s *colorTimeService) normalizeTrackingGlobal(ctx context.Context, organizationID, userID, role, tracking string) error {

	slots, weeks, err := s.ColorTimeRepository.GetAllSlotsByTracking(ctx, organizationID, userID, role, tracking)
	if err != nil {
		return fmt.Errorf("failed to get all slots by tracking: %w", err)
	}

	sort.Slice(slots, func(i, j int) bool {
		return slots[i].CreatedAt.Before(slots[j].CreatedAt)
	})

	for i, slot := range slots {
		slot.UseCount = i + 1
		slot.UpdatedAt = time.Now()
	}

	for _, week := range weeks {
		week.UpdatedAt = time.Now()
		if err := s.ColorTimeRepository.UpdateColorTimeWeek(ctx, week.ID, week); err != nil {
			return err
		}
	}

	return nil

}

func (s *colorTimeService) convertBlocksWithProductInfo(ctx context.Context, blocks []*ColorBlock) ([]*BlockResponse, error) {
	blockResponses := make([]*BlockResponse, 0, len(blocks))

	for _, block := range blocks {
		slotResponses := make([]*SlotResponse, 0, len(block.Slots))

		for _, slot := range block.Slots {
			slotResponse := &SlotResponse{
				SlotID:    slot.SlotID,
				SlotIDOld: *slot.SlotIDOld,
				Sessions:  slot.Sessions,
				Title:     slot.Title,
				Tracking:  slot.Tracking,
				UseCount:  slot.UseCount,
				StartTime: slot.StartTime,
				EndTime:   slot.EndTime,
				Duration:  slot.Duration,
				Color:     slot.Color,
				Note:      slot.Note,
				ProductID: slot.ProductID,
				CreatedAt: slot.CreatedAt,
				UpdatedAt: slot.UpdatedAt,
			}

			if slot.ProductID != nil && *slot.ProductID != "" {
				product, err := s.ProductService.GetProductInfor(ctx, *slot.ProductID)
				if err != nil {
					fmt.Printf("Warning: failed to get product info for ID %s: %v\n", *slot.ProductID, err)
				} else if product != nil {
					slotResponse.Product = &ProductInfo{
						ID:                   product.ID.Hex(),
						ProductName:          product.ProductName,
						OriginalPriceStore:   product.OriginalPriceStore,
						OriginalPriceService: product.OriginalPriceService,
						ProductImage:         product.ProductImage,
						TopicName:            product.TopicName,
						CategoryName:         product.CategoryName,
					}
				}
			}

			slotResponses = append(slotResponses, slotResponse)
		}

		blockResponse := &BlockResponse{
			BlockID:    block.BlockID,
			BlockIDOld: *block.BlockIDOld,
			Slots:      slotResponses,
		}

		blockResponses = append(blockResponses, blockResponse)
	}

	return blockResponses, nil
}

func (s *colorTimeService) mergeColorTimes(existingCT, defaultCT []*ColorTime) []*ColorTime {
	merged := make([]*ColorTime, 0)

	existingMap := make(map[string]*ColorTime)
	for _, e := range existingCT {
		existingMap[e.Date.Format("2006-01-02")] = e
	}

	for _, def := range defaultCT {
		dateStr := def.Date.Format("2006-01-02")
		if userCT, exists := existingMap[dateStr]; exists {
			mergedCT := s.mergeSingleColorTime(userCT, def)
			merged = append(merged, mergedCT)
			delete(existingMap, dateStr)
		} else {
			newCT := *def
			newCT.ID = primitive.NewObjectID()
			newCT.CreatedAt = time.Now()
			newCT.UpdatedAt = time.Now()
			merged = append(merged, &newCT)
		}
	}

	return merged
}

func (s *colorTimeService) mergeSingleColorTime(existingCT, defaultCT *ColorTime) *ColorTime {
	mergedCT := &ColorTime{
		ID:        existingCT.ID,
		Date:      existingCT.Date,
		TopicID:   existingCT.TopicID,
		TimeSlots: make([]*ColorBlock, 0),
		CreatedAt: existingCT.CreatedAt,
		UpdatedAt: time.Now(),
	}

	existingBlocksMap := make(map[string]*ColorBlock)
	for _, b := range existingCT.TimeSlots {
		key := ""
		if b.BlockIDOld != nil {
			key = b.BlockIDOld.Hex()
		} else {
			key = b.BlockID.Hex()
		}
		existingBlocksMap[key] = b
	}

	for _, defaultBlock := range defaultCT.TimeSlots {
		defKey := defaultBlock.BlockID.Hex()
		if userBlock, exists := existingBlocksMap[defKey]; exists {
			mergedBlock := s.mergeSingleBlock(userBlock, defaultBlock)
			mergedCT.TimeSlots = append(mergedCT.TimeSlots, mergedBlock)
			delete(existingBlocksMap, defKey)
		} else {
			newBlock := cloneBlock(defaultBlock)
			newBlock.BlockIDOld = &defaultBlock.BlockID
			mergedCT.TimeSlots = append(mergedCT.TimeSlots, newBlock)
		}
	}
	return mergedCT
}

func (s *colorTimeService) mergeSingleBlock(existingBlock, defaultBlock *ColorBlock) *ColorBlock {
	mergedBlock := &ColorBlock{
		BlockID:    existingBlock.BlockID,
		BlockIDOld: existingBlock.BlockIDOld,
		Slots:      make([]*ColortimeSlot, 0),
	}

	existingSlotsMap := make(map[string]*ColortimeSlot)
	for _, s := range existingBlock.Slots {
		key := ""
		if s.SlotIDOld != nil {
			key = s.SlotIDOld.Hex()
		} else {
			key = s.SlotID.Hex()
		}
		existingSlotsMap[key] = s
	}

	for _, defaultSlot := range defaultBlock.Slots {
		defKey := defaultSlot.SlotID.Hex()
		if userSlot, exists := existingSlotsMap[defKey]; exists {
			mergedBlock.Slots = append(mergedBlock.Slots, userSlot)
			delete(existingSlotsMap, defKey)
		} else {
			newSlot := cloneSlot(defaultSlot)
			newSlot.SlotIDOld = &defaultSlot.SlotID
			mergedBlock.Slots = append(mergedBlock.Slots, newSlot)
		}
	}
	
	return mergedBlock
}

func cloneBlock(block *ColorBlock) *ColorBlock {
	newBlock := *block
	newBlock.BlockID = primitive.NewObjectID()
	newSlots := make([]*ColortimeSlot, 0)
	for _, s := range block.Slots {
		newSlots = append(newSlots, cloneSlot(s))
	}
	newBlock.Slots = newSlots
	return &newBlock
}

func cloneSlot(s *ColortimeSlot) *ColortimeSlot {
	newSlot := *s
	newSlot.SlotID = primitive.NewObjectID()
	return &newSlot
}
