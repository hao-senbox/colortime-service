package colortime

import (
	"colortime-service/internal/product"
	"context"
	"errors"
	"log"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ColorTimeService interface {
	CreateColorTime(ctx context.Context, req CreateColorTimeRequest) error
	GetColorTimes(ctx context.Context, userID string, baseDate string, timeRange string) ([]*ColorTimeResponse, error)
	GetColorTime(ctx context.Context, id string) (*ColorTimeResponse, error)
	UpdateColorTime(ctx context.Context, req UpdateColorTimeRequest, id string) error
	DeleteColorTime(ctx context.Context, id string) error
}

type colorTimeService struct {
	ColorTimeRepository ColorTimeRepository
	ProductService      product.ProductService
}

func NewColorTimeService(colorTimeRepository ColorTimeRepository, productService product.ProductService) ColorTimeService {
	return &colorTimeService{
		ColorTimeRepository: colorTimeRepository,
		ProductService:      productService,
	}
}

func (s *colorTimeService) CreateColorTime(ctx context.Context, req CreateColorTimeRequest) error {

	if req.UserID == "" {
		return errors.New("user id is required")
	}

	if req.Title == "" {
		return errors.New("title is required")
	}

	if req.Date == "" {
		return errors.New("date is required")
	}

	dateParese, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		return errors.New("invalid date format")
	}

	if req.Time == "" {
		return errors.New("time is required")
	}

	timeParse, err := time.Parse("15:04", req.Time)
	if err != nil {
		return errors.New("invalid time format")
	}

	timeOnly := time.Date(1970, 1, 1,
		timeParse.Hour(), timeParse.Minute(), timeParse.Second(),
		0, time.UTC,
	)

	if req.Duration == 0 {
		return errors.New("duration is required")
	}

	if req.Color == "" {
		return errors.New("color is required")
	}

	if req.ProductID == "" {
		return errors.New("product id is required")
	}

	objProductID, err := primitive.ObjectIDFromHex(req.ProductID)
	if err != nil {
		return errors.New("invalid product id format")
	}

	if req.Note == "" {
		return errors.New("note is required")
	}

	colorTime := &ColorTime{
		ID:        primitive.NewObjectID(),
		UserID:    req.UserID,
		Title:     req.Title,
		Date:      dateParese,
		Time:      timeOnly,
		Duration:  req.Duration,
		Color:     req.Color,
		Note:      req.Note,
		ProductID: objProductID,
	}

	return s.ColorTimeRepository.CreateColorTime(ctx, colorTime)

}

func (s *colorTimeService) GetColorTimes(ctx context.Context, userID string, baseDate string, timeRange string) ([]*ColorTimeResponse, error) {

	var colortimesData []*ColorTimeResponse

	if userID == "" {
		return nil, errors.New("user id is required")
	}

	if baseDate == "" {
		return nil, errors.New("base date is required")
	}

	if timeRange == "" {
		return nil, errors.New("time range is required")
	}

	baseDateParse, err := time.Parse("2006-01-02", baseDate)
	if err != nil {
		return nil, errors.New("invalid base date format")
	}

	startDate, endDate := getWeekRange(baseDateParse)

	start, end, err := parseTimeRange(timeRange)
	if err != nil {
		return nil, err
	}

	colortimes, err := s.ColorTimeRepository.GetColorTimes(ctx, userID, startDate, endDate, start, end)
	if err != nil {
		return nil, err
	}

	for _, colortime := range colortimes {

		productData, err := s.ProductService.GetProductInfor(colortime.ProductID.Hex())
		if err != nil {
			log.Println(err)
		}

		var productRes *product.Product
		if productData != nil { 
			productRes = &product.Product{
				ID:                   productData.ID,
				ProductName:          productData.ProductName,
				OriginalPriceStore:   productData.OriginalPriceStore,
				OriginalPriceService: productData.OriginalPriceService,
				ProductDescription:   productData.ProductDescription,
				CategoryName:         productData.CategoryName,
				TopicName:            productData.TopicName,
			}
		} else {
			productRes = nil 
		}

		colorTimeData := &ColorTimeResponse{
			ID:       colortime.ID,
			UserID:   colortime.UserID,
			Title:    colortime.Title,
			Date:     colortime.Date,
			Time:     colortime.Time,
			Duration: colortime.Duration,
			Color:    colortime.Color,
			Note:     colortime.Note,
			Product:  productRes,
		}

		colortimesData = append(colortimesData, colorTimeData)
	}

	return colortimesData, nil

}

func (s *colorTimeService) GetColorTime(ctx context.Context, id string) (*ColorTimeResponse, error) {

	if id == "" {
		return nil, errors.New("id is required")
	}

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New("invalid id format")
	}

	colortime, err := s.ColorTimeRepository.GetColorTime(ctx, objectID)
	if err != nil {
		return nil, err
	}

	productData, err := s.ProductService.GetProductInfor(colortime.ProductID.Hex())
	if err != nil {
		log.Println(err)
	}

	colorTimeData := &ColorTimeResponse{
		ID:       colortime.ID,
		UserID:   colortime.UserID,
		Title:    colortime.Title,
		Date:     colortime.Date,
		Time:     colortime.Time,
		Duration: colortime.Duration,
		Color:    colortime.Color,
		Note:     colortime.Note,
		Product: &product.Product{
			ID:                   productData.ID,
			ProductName:          productData.ProductName,
			OriginalPriceStore:   productData.OriginalPriceStore,
			OriginalPriceService: productData.OriginalPriceService,
			ProductDescription:   productData.ProductDescription,
			CategoryName:         productData.CategoryName,
			TopicName:            productData.TopicName,
		},
	}
	return colorTimeData, nil

}

func (s *colorTimeService) UpdateColorTime(ctx context.Context, req UpdateColorTimeRequest, id string) error {

	var date time.Time
	var productID primitive.ObjectID

	if id == "" {
		return errors.New("id is required")
	}

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("invalid id format")
	}

	colorTime, err := s.ColorTimeRepository.GetColorTime(ctx, objectID)
	if err != nil {
		return err
	}

	if colorTime == nil {
		return errors.New("color time not found")
	}

	if req.Title == "" {
		req.Title = colorTime.Title
	}

	if req.Date == "" {
		date = colorTime.Date
	} else {
		dateParse, err := time.Parse("2006-01-02", req.Date)
		if err != nil {
			return errors.New("invalid date format")
		}

		date = dateParse
	}

	if req.Time == "" {
		req.Time = colorTime.Time.Format("15:04")
	}

	if req.Duration == 0 {
		req.Duration = colorTime.Duration
	}

	if req.Color == "" {
		req.Color = colorTime.Color
	}

	if req.ProductID == "" {
		productID = colorTime.ProductID
	} else {
		objProductID, err := primitive.ObjectIDFromHex(req.ProductID)
		if err != nil {
			return errors.New("invalid product id format")
		}
		productID = objProductID
	}

	if req.Note == "" {
		req.Note = colorTime.Note
	}

	timeParse, err := time.Parse("15:04", req.Time)
	if err != nil {
		return errors.New("invalid time format")
	}

	timeOnly := time.Date(1970, 1, 1,
		timeParse.Hour(), timeParse.Minute(), timeParse.Second(),
		0, time.UTC,
	)

	data := &ColorTime{
		ID:        colorTime.ID,
		UserID:    colorTime.UserID,
		Title:     req.Title,
		Date:      date,
		Time:      timeOnly,
		Duration:  req.Duration,
		Color:     req.Color,
		ProductID: productID,
		Note:      req.Note,
		UpdatedAt: time.Now(),
	}

	return s.ColorTimeRepository.UpdateColorTime(ctx, data)

}

func (s *colorTimeService) DeleteColorTime(ctx context.Context, id string) error {

	if id == "" {
		return errors.New("id is required")
	}

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("invalid id format")
	}

	return s.ColorTimeRepository.DeleteColorTime(ctx, objectID)

}

func getWeekRange(baseDate time.Time) (time.Time, time.Time) {

	weekday := int(baseDate.Weekday())
	if weekday == 0 {
		weekday = 7
	}

	startOfWeek := baseDate.AddDate(0, 0, -weekday+1)
	startOfWeek = time.Date(startOfWeek.Year(), startOfWeek.Month(), startOfWeek.Day(), 0, 0, 0, 0, baseDate.Location())

	endOfWeek := startOfWeek.AddDate(0, 0, 6)
	endOfWeek = time.Date(endOfWeek.Year(), endOfWeek.Month(), endOfWeek.Day(), 23, 59, 59, 0, baseDate.Location())

	return startOfWeek, endOfWeek

}

func parseTimeRange(tr string) (time.Time, time.Time, error) {

	parts := strings.Split(tr, "-")
	if len(parts) != 2 {
		return time.Time{}, time.Time{}, errors.New("invalid time range")
	}

	layout := "15:04"

	startT, err1 := time.Parse(layout, parts[0])
	endT, err2 := time.Parse(layout, parts[1])
	if err1 != nil || err2 != nil {
		return time.Time{}, time.Time{}, errors.New("invalid time format")
	}

	start := time.Date(1970, 1, 1, startT.Hour(), startT.Minute(), 0, 0, time.UTC)

	end := time.Date(1970, 1, 1, endT.Hour(), endT.Minute(), 59, 0, time.UTC)

	return start, end, nil

}
