# GetColorTimeWeek Flow - Chi tiết và Ví dụ thực tế

## Tổng quan
Hàm `GetColorTimeWeek` trong `colorTimeService` chịu trách nhiệm lấy dữ liệu thời khóa biểu hàng tuần của một học sinh, với logic phức tạp để đồng bộ dữ liệu từ `default` (mẫu chung) sang `colortime` (dữ liệu cá nhân).

## Flow chính

### 1. Validation Input
```go
if userID == "" { return nil, errors.New("user id is required") }
if role == "" { return nil, errors.New("role is required") }
if orgID == "" { return nil, errors.New("organization id is required") }
if start == "" || end == "" { return nil, errors.New("start and end date are required") }
```

### 2. Parse Dates
```go
startDate, err := time.Parse("2006-01-02", start)  // "2024-01-01"
endDate, err := time.Parse("2006-01-02", end)      // "2024-01-07"
```

### 3. Lấy dữ liệu Default
```go
defaultDayColorTimes, err := s.DefaultColorTimeRepository.GetDefaultDayColorTimesInRange(ctx, startDate, endDate, orgID)
```
- Lấy tất cả dữ liệu mẫu (default) trong khoảng thời gian từ start đến end
- Ví dụ: Lấy template cho Monday, Wednesday, Friday trong tuần

### 4. Kiểm tra Existing Week
```go
existingWeek, err := s.ColorTimeRepository.GetColorTimeWeek(ctx, &startDate, &endDate, orgID, userID, role)
```
- Tìm xem học sinh đã có dữ liệu colortime cho tuần này chưa

## Logic xử lý chính

### Trường hợp A: Có dữ liệu Default (len(defaultDayColorTimes) > 0)

#### A1. Có Existing Week (existingWeek != nil)
```go
// 1. Clone dữ liệu từ default thành colorTimes
colorTimes := cloneDefaultDayColorTimesToColorTimes(defaultDayColorTimes)

// 2. Merge với existing data (giữ lại dữ liệu cá nhân đã chỉnh sửa)
updatedColorTimes := s.mergeColorTimes(existingWeek.ColorTimes, colorTimes)
existingWeek.ColorTimes = updatedColorTimes

// 3. Sync với latest default data (quan trọng!)
s.syncColorTimesWithDefault(existingWeek.ColorTimes, defaultDayColorTimes)

// 4. Update database
existingWeek.UpdatedAt = time.Now()
s.ColorTimeRepository.UpdateColorTimeWeek(ctx, existingWeek.ID, existingWeek)
```

#### A2. Không có Existing Week (existingWeek == nil)
```go
// 1. Clone dữ liệu từ default
colorTimes := cloneDefaultDayColorTimesToColorTimes(defaultDayColorTimes)

// 2. Tạo WeekColorTime mới
newColorTimeWeek := &WeekColorTime{
    ID: primitive.NewObjectID(),
    OrganizationID: orgID,
    Owner: &Owner{OwnerID: userID, OwnerRole: role},
    StartDate: startDate,
    EndDate: endDate,
    ColorTimes: colorTimes,
    CreatedBy: userID,
    CreatedAt: time.Now(),
    UpdatedAt: time.Now(),
}

// 3. Sync với default data trước khi lưu
s.syncColorTimesWithDefault(newColorTimeWeek.ColorTimes, defaultDayColorTimes)

// 4. Lưu vào database
s.ColorTimeRepository.CreateColorTimeWeek(ctx, newColorTimeWeek)
```

### Trường hợp B: Không có dữ liệu Default (len(defaultDayColorTimes) == 0)

#### B1. Có Existing Week
```go
colortimeWeek = existingWeek  // Trả về dữ liệu hiện có
```

#### B2. Không có Existing Week
```go
// Tạo WeekColorTime rỗng
newColorTimeWeek := &WeekColorTime{
    ID: primitive.NewObjectID(),
    OrganizationID: orgID,
    Owner: &Owner{OwnerID: userID, OwnerRole: role},
    StartDate: startDate,
    EndDate: endDate,
    ColorTimes: []*ColorTime{},  // Rỗng
    CreatedBy: userID,
    CreatedAt: time.Now(),
    UpdatedAt: time.Now(),
}
s.ColorTimeRepository.CreateColorTimeWeek(ctx, newColorTimeWeek)
```

## Hàm quan trọng: syncColorTimesWithDefault

### Mục đích
Đồng bộ dữ liệu từ `default` sang `colortime` - đảm bảo các trường quan trọng luôn được cập nhật từ mẫu.

### Logic
```go
func (s *colorTimeService) syncColorTimesWithDefault(colorTimes []*ColorTime, defaultDayColorTimes []*default_colortime.DefaultDayColorTime) {
    // 1. Tạo map từ date -> default day
    defaultMap := make(map[string]*default_colortime.DefaultDayColorTime)
    for _, defaultDay := range defaultDayColorTimes {
        dateStr := defaultDay.Date.Format("2006-01-02")
        defaultMap[dateStr] = defaultDay
    }

    // 2. Duyệt qua từng ngày trong colortime
    for _, colorTime := range colorTimes {
        dateStr := colorTime.Date.Format("2006-01-02")
        if defaultDay, exists := defaultMap[dateStr]; exists {
            // 3. Tạo map slot ID -> default slot
            defaultSlotMap := make(map[string]*default_colortime.DefaultColortimeSlot)
            for _, defaultBlock := range defaultDay.TimeSlots {
                for _, defaultSlot := range defaultBlock.Slots {
                    slotIDStr := defaultSlot.SlotID.Hex()
                    defaultSlotMap[slotIDStr] = defaultSlot
                }
            }

            // 4. Update colortime slots từ default
            for _, colorBlock := range colorTime.TimeSlots {
                for _, colorSlot := range colorBlock.Slots {
                    if colorSlot.SlotIDOld != nil {
                        slotIDStr := colorSlot.SlotIDOld.Hex()
                        if defaultSlot, exists := defaultSlotMap[slotIDStr]; exists {
                            // Sync các trường quan trọng
                            colorSlot.StartTime = defaultSlot.StartTime
                            colorSlot.EndTime = defaultSlot.EndTime
                            colorSlot.Duration = defaultSlot.Duration
                            colorSlot.Title = defaultSlot.Title
                            colorSlot.Color = defaultSlot.Color
                            colorSlot.Note = defaultSlot.Note
                        }
                    }
                }
            }
        }
    }
}
```

## Ví dụ thực tế

### Dữ liệu ban đầu

#### Default Templates:
```json
[
  {
    "date": "monday",
    "time_slots": [
      {
        "slots": [
          {
            "slot_id": "default_slot_1",
            "title": "Toán học",
            "start_time": "08:00",
            "duration": 60,
            "color": "#FF0000",
            "note": "Bài tập về nhà"
          }
        ]
      }
    ]
  },
  {
    "date": "wednesday",
    "time_slots": [
      {
        "slots": [
          {
            "slot_id": "default_slot_2",
            "title": "Vật lý",
            "start_time": "09:00",
            "duration": 90,
            "color": "#00FF00",
            "note": "Thí nghiệm"
          }
        ]
      }
    ]
  }
]
```

#### Colortime hiện có (existingWeek):
```json
{
  "color_times": [
    {
      "date": "2024-01-01", // Monday
      "time_slots": [
        {
          "slots": [
            {
              "slot_id": "personal_slot_1",
              "slot_id_old": "default_slot_1",
              "title": "Toán học cơ bản", // Đã chỉnh sửa
              "start_time": "08:00",
              "duration": 60,
              "color": "#FF0000",
              "note": "Bài tập về nhà"
            }
          ]
        }
      ]
    }
  ]
}
```

### Sau khi Default được cập nhật:
```json
// Default mới
{
  "date": "monday",
  "time_slots": [
    {
      "slots": [
        {
          "slot_id": "default_slot_1",
          "title": "Toán học nâng cao", // Thay đổi title
          "start_time": "08:30",        // Thay đổi time
          "duration": 75,               // Thay đổi duration
          "color": "#FF0000",
          "note": "Bài tập nâng cao"    // Thay đổi note
        }
      ]
    }
  ]
}
```

### Kết quả sau sync:
```json
{
  "color_times": [
    {
      "date": "2024-01-01",
      "time_slots": [
        {
          "slots": [
            {
              "slot_id": "personal_slot_1",
              "slot_id_old": "default_slot_1",
              "title": "Toán học nâng cao", // ✅ Cập nhật từ default
              "start_time": "08:30",        // ✅ Cập nhật từ default
              "duration": 75,               // ✅ Cập nhật từ default
              "color": "#FF0000",
              "note": "Bài tập nâng cao"    // ✅ Cập nhật từ default
            }
          ]
        }
      ]
    }
  ]
}
```

## Flow diagram

```
GET /api/v1/colortime/week
    ↓
Validate inputs
    ↓
Parse dates (startDate, endDate)
    ↓
Get defaultDayColorTimes from DB
    ↓
Check existingWeek in DB
    ↓
┌─────────────────────────────────────┐
│         Có default data?            │
│   (len(defaultDayColorTimes) > 0)   │
└─────────────────┬───────────────────┘
                  │
            ┌─────▼─────┐
            │ Có existing│
            │    week?  │
            └─────┬─────┘
                  │
        ┌─────────▼─────────┐
        │     Clone từ      │
        │   default data    │
        └─────────┬─────────┘
                  │
        ┌─────────▼─────────┐
        │     Merge với     │
        │ existing colorTimes│
        └─────────┬─────────┘
                  │
        ┌─────────▼─────────┐
        │  Sync với latest  │ ←──┐
        │   default data    │    │
        └─────────┬─────────┘    │
                  │             │
        ┌─────────▼─────────┐    │
        │   Update/Create   │    │
        │       in DB       │    │
        └─────────┬─────────┘    │
                  │             │
        ┌─────────▼─────────┐    │
        │   Build response  │    │
        │ (convert to JSON) │    │
        └─────────┬─────────┘    │
                  │             │
            ┌─────▼─────┐       │
            │ Return    │       │
            │ response  │       │
            └───────────┘       │
                               │
    ┌──────────────────────────┘
    │ Không có default data
    │
┌───▼───┐
│ Có ex- │
│ isting │
│ week?  │
└───┬───┘
    │
┌───▼───┐
│Return │
│existing│
│ week   │
└───┬───┘
    │
┌───▼───┐
│Create │
│ empty │
│ week  │
└───────┘
```

## Các hàm phụ trợ

### cloneDefaultDayColorTimesToColorTimes
- Chuyển đổi `DefaultDayColorTime` thành `ColorTime`
- Tạo `SlotIDOld` để reference về default

### mergeColorTimes
- Merge dữ liệu mới với dữ liệu hiện có
- Giữ lại dữ liệu cá nhân đã chỉnh sửa

### syncColorTimesWithDefault
- **Quan trọng nhất**: Đồng bộ dữ liệu từ default sang colortime
- Cập nhật: `start_time`, `duration`, `title`, `color`, `note`
- Dựa trên `SlotIDOld` để map đúng slot

## Lợi ích của thiết kế này

1. **Real-time sync**: Dữ liệu luôn được cập nhật từ default mỗi lần get
2. **Preserve personal changes**: Merge logic giữ lại dữ liệu cá nhân
3. **Efficient**: Chỉ sync khi có dữ liệu default
4. **Safe**: Không làm mất dữ liệu khi default thay đổi
5. **Flexible**: Hỗ trợ cả create mới và update existing

## Kết luận

Hàm `GetColorTimeWeek` không chỉ lấy dữ liệu mà còn đảm bảo tính nhất quán giữa dữ liệu mẫu (default) và dữ liệu cá nhân (colortime). Logic sync thông minh giúp tự động cập nhật các trường quan trọng từ default mà không làm mất dữ liệu cá nhân hóa của học sinh.
