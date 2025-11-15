# Hệ Thống Quản Lý Thời Gian Màu (Colortime)

## Tổng quan

Hệ thống Colortime là một giải pháp quản lý thời gian biểu với 3 tầng dữ liệu:

- **Template Colortime**: Mẫu thời gian biểu theo ngày trong tuần (Monday, Tuesday, etc.)
- **Default Colortime**: Áp dụng template vào các ngày cụ thể trong kỳ học
- **User Colortime**: Dữ liệu cá nhân của từng học sinh, sync từ Default

## 1. Template Colortime - Nguồn Cấu Trúc Cơ Bản

### 1.1. Mô tả
Template là mẫu thời gian biểu theo ngày trong tuần, chứa cấu trúc cơ bản của các khối thời gian và slot hoạt động.

### 1.2. Cấu trúc dữ liệu

```go
type TemplateColorTime struct {
    Date           string               `bson:"date"` // "monday", "tuesday", etc.
    OrganizationID string               `bson:"organization_id"`
    TermID         string               `bson:"term_id"`
    ColorTimes     []*ColorTimeTemplate `bson:"color_times"`
}

type ColorTimeTemplate struct {
    BlockID primitive.ObjectID `bson:"block_id"`
    Slots   []*ColortimeSlot   `bson:"slots"`
}

type ColortimeSlot struct {
    SlotID    primitive.ObjectID `bson:"slot_id"`
    Sessions  int                `bson:"sessions"`
    Title     string             `bson:"title"`
    StartTime time.Time          `bson:"start_time"`
    EndTime   time.Time          `bson:"end_time"`
    Duration  int                `bson:"duration"` // phút
    Color     string             `bson:"color"`
    Note      string             `bson:"note"`
}
```

### 1.3. Ví dụ dữ liệu Template

```json
{
  "_id": "67369e769968fdd6f5000e29",
  "date": "monday",
  "organization_id": "6ba5042b-6213-11f0-91d9-a637109e411e",
  "term_id": "68e4966d212b467510654a09",
  "color_times": [
    {
      "block_id": "67369e769968fdd6f5000e2a",
      "slots": [
        {
          "slot_id": "67369e769968fdd6f5000e2b",
          "sessions": 1,
          "title": "Đi ăn",
          "start_time": "2024-11-14T07:00:00Z",
          "end_time": "2024-11-14T07:30:00Z",
          "duration": 30,
          "color": "#4CAF50",
          "note": "Đi ăn nhé",
          "created_at": "2024-11-14T03:13:58.962Z",
          "updated_at": "2024-11-14T03:13:58.962Z"
        },
        {
          "slot_id": "67369e919968fdd6f5000e2c",
          "sessions": 2,
          "title": "Đi chơi",
          "start_time": "2024-11-14T08:00:00Z",
          "end_time": "2024-11-14T08:45:00Z",
          "duration": 45,
          "color": "#2196F3",
          "note": "Đi chơi nhé",
          "created_at": "2024-11-14T03:14:25.006Z",
          "updated_at": "2024-11-14T03:14:25.006Z"
        }
      ]
    }
  ],
  "created_by": "d8a764b0-2c8e-11f0-bf8b-1ae0d65ab5bf",
  "created_at": "2024-11-14T03:13:58.962Z",
  "updated_at": "2024-11-14T03:13:58.962Z"
}
```

### 1.4. API Endpoints

#### Tạo Template Mới
```http
POST /template-colortime
```

**Request Body:**
```json
{
  "organization_id": "6ba5042b-6213-11f0-91d9-a637109e411e",
  "term_id": "68e4966d212b467510654a09",
  "date": "monday",
  "start_time": "07:00",
  "duration": 90,
  "title": "Ăn sáng",
  "color": "#4CAF50",
  "note": "Bữa sáng healthy"
}
```

**Logic xử lý:**
1. Validate input (start_time format HH:MM, duration tính bằng giây)
2. Tính EndTime = StartTime + Duration (giây)
3. Convert Duration từ giây sang phút lưu database
4. Tìm template hiện có cho ngày đó
5. Nếu có: thêm slot vào block đầu tiên
6. Nếu chưa có: tạo template mới với block và slot đầu tiên

#### Lấy Template
```http
GET /template-colortime?organization_id=...&term_id=...&date=monday
GET /template-colortime?organization_id=...&term_id=...&date=week
```

#### Duplicate Template
```http
POST /template-colortime/duplicate
```

**Request Body:**
```json
{
  "organization_id": "6ba5042b-6213-11f0-91d9-a637109e411e",
  "term_id": "68e4966d212b467510654a09",
  "origin_date": "monday",
  "target_date": "tuesday"
}
```

**Logic:**
- Nếu target_date empty: duplicate sang tất cả ngày còn lại trong tuần (Tue-Sun)
- Nếu target_date specified: duplicate sang ngày cụ thể
- Skip nếu ngày target đã có template

#### Update Slot Template
```http
PUT /template-colortime/slot
```

**Request Body:**
```json
{
  "slot_id": "67369e769968fdd6f5000e2b",
  "start_time": "07:30",
  "duration": 120,
  "title": "Ăn sáng và tập thể dục",
  "color": "#FF9800",
  "note": "Bữa sáng + exercise"
}
```

**Logic:**
- Tìm slot theo SlotID
- Update fields: Title, StartTime, EndTime, Duration, Color, Note
- Tính lại EndTime = StartTime + Duration (giây)

## 2. Default Colortime - Áp dụng Template vào Kỳ Học

### 2.1. Mô tả
Default colortime áp dụng cấu trúc từ template vào các ngày cụ thể trong kỳ học.

### 2.2. Cấu trúc dữ liệu

```go
type DefaultDayColorTime struct {
    Date           time.Time            `bson:"date"` // Ngày cụ thể
    OrganizationID string               `bson:"organization_id"`
    TimeSlots      []*DefaultColorBlock `bson:"time_slots"`
}

type DefaultColorBlock struct {
    BlockID primitive.ObjectID      `bson:"block_id"`
    Slots   []*DefaultColortimeSlot `bson:"slots"`
}

type DefaultColortimeSlot struct {
    SlotID    primitive.ObjectID `bson:"slot_id"`
    Sessions  int                `bson:"sessions"`
    Title     string             `bson:"title"`
    StartTime time.Time          `bson:"start_time"`
    EndTime   time.Time          `bson:"end_time"`
    Duration  int                `bson:"duration"`
    Color     string             `bson:"color"`
    Note      string             `bson:"note"`
}
```

### 2.3. API Endpoints

#### Áp dụng Template vào Default
```http
POST /template-colortime/apply
```

**Request Body:**
```json
{
  "organization_id": "6ba5042b-6213-11f0-91d9-a637109e411e",
  "term_id": "68e4966d212b467510654a09",
  "start_date": "2024-11-18",
  "end_date": "2024-11-24"
}
```

**Logic xử lý:**
1. Lấy tất cả template từ Monday đến Sunday
2. Duyệt từng ngày từ start_date đến end_date
3. Xác định weekday của ngày đó (Monday, Tuesday, etc.)
4. Tìm template tương ứng với weekday
5. Nếu có template:
   - Nếu default đã tồn tại: merge template vào default
   - Nếu chưa có: tạo mới từ template

### 2.4. Logic Merge Template vào Default

**Khi Default đã tồn tại:**
```go
existingDefaultColorTime.TimeSlots = mergeTemplateIntoDefault(existingDefaultColorTime, template)
```

**Hàm mergeTemplateIntoDefault:**
1. Map các block hiện có theo BlockID
2. Duyệt template blocks:
   - Nếu BlockID match: merge slots
   - Nếu không match: thêm block mới
3. Giữ lại blocks default không có trong template

**Hàm mergeSlots:**
1. Map các slot hiện có theo SlotID
2. Duyệt template slots:
   - Nếu SlotID match: update dữ liệu từ template (giữ custom fields)
   - Nếu không match: thêm slot mới
3. Giữ lại slots default không có trong template

## 3. User Colortime - Dữ liệu Cá Nhân Học Sinh

### 3.1. Mô tả
User colortime là dữ liệu cá nhân của từng học sinh, được sync từ default và có thể có dữ liệu riêng.

### 3.2. Cấu trúc dữ liệu

```go
type WeekColorTime struct {
    ColorTimes []*ColorTime `bson:"colortimes"`
}

type ColorTime struct {
    Date      time.Time    `bson:"date"`
    TopicID   string       `bson:"topic_id"`
    TimeSlots []*ColorBlock `bson:"time_slots"`
}

type ColorBlock struct {
    BlockID    primitive.ObjectID `bson:"block_id"`
    BlockIDOld *primitive.ObjectID `bson:"block_id_old"` // Reference to default block
    Slots      []*ColortimeSlot   `bson:"slots"`
}

type ColortimeSlot struct {
    SlotID    primitive.ObjectID `bson:"slot_id"`
    SlotIDOld *primitive.ObjectID `bson:"slot_id_old"` // Reference to default slot
    Sessions  int                `bson:"sessions"`
    Title     string             `bson:"title"`
    StartTime time.Time          `bson:"start_time"`
    EndTime   time.Time          `bson:"end_time"`
    Duration  int                `bson:"duration"`
    Color     string             `bson:"color"`
    Note      string             `bson:"note"`
    // Custom fields
    Tracking  string             `bson:"tracking,omitempty"`
    Sbt       string             `bson:"sbt,omitempty"`
}
```

### 3.3. Flow Sync Default → User Colortime

#### API: Lấy Tuần Colortime
```http
GET /colortime/week?start_date=2024-11-18&end_date=2024-11-24
```

**Logic xử lý:**

1. **Lấy Default Data:**
   - Query default_colortime theo organization_id và date range
   - Clone data thành user colortime format

2. **Check Existing User Data:**
   - Tìm user week hiện có trong colortime collection
   - Nếu có: merge với default data
   - Nếu chưa có: tạo mới từ default

3. **Merge Logic:**
   ```go
   merged := mergeColorTimes(existingUserData, defaultData)
   ```

4. **Sync Fields:**
   - Luôn update: StartTime, EndTime, Duration, Title, Color, Note
   - Giữ lại: Tracking, Sbt, và các custom fields khác
   - Thêm mới: Slots từ default chưa có
   - Giữ lại: Slots riêng của user không có trong default

### 3.4. Ví dụ Sync

**Default Data:**
```json
{
  "date": "2024-11-18",
  "time_slots": [{
    "block_id": "674a1b2c3d4e5f6789abcdef1",
    "slots": [{
      "slot_id": "674a1b2c3d4e5f6789abcdef2",
      "title": "Ăn sáng",
      "start_time": "07:00",
      "duration": 30,
      "color": "#4CAF50"
    }]
  }]
}
```

**User Data hiện tại:**
```json
{
  "date": "2024-11-18",
  "time_slots": [{
    "block_id": "675b2c3d4e5f6789abcdef0",
    "block_id_old": "674a1b2c3d4e5f6789abcdef1",
    "slots": [{
      "slot_id": "675b2c3d4e5f6789abcdef1",
      "slot_id_old": "674a1b2c3d4e5f6789abcdef2",
      "title": "Ăn sáng",
      "start_time": "07:00",
      "duration": 30,
      "color": "#4CAF50",
      "tracking": "completed",
      "sbt": "nutrition"
    }]
  }]
}
```

**Sau sync:**
```json
{
  "date": "2024-11-18",
  "time_slots": [{
    "block_id": "675b2c3d4e5f6789abcdef0",
    "block_id_old": "674a1b2c3d4e5f6789abcdef1",
    "slots": [{
      "slot_id": "675b2c3d4e5f6789abcdef1",
      "slot_id_old": "674a1b2c3d4e5f6789abcdef2",
      "title": "Ăn sáng",        // Từ default
      "start_time": "07:00",      // Từ default
      "duration": 30,             // Từ default
      "color": "#4CAF50",         // Từ default
      "tracking": "completed",    // GIỮ LẠI user data
      "sbt": "nutrition"          // GIỮ LẠI user data
    }]
  }]
}
```

## 4. Workflow Hoàn Chỉnh

### 4.1. Thiết lập Ban Đầu
1. **Admin tạo Template:**
   - Tạo template cho Monday
   - Duplicate sang các ngày khác trong tuần

2. **Admin áp dụng Template:**
   - Apply template vào kỳ học (start_date → end_date)
   - Tạo default data cho từng ngày

### 4.2. Học Sinh Sử dụng
1. **Lấy Thời Khóa Biểu:**
   - API GetColorTimeWeek tự động sync default → user data
   - Học sinh thấy thời khóa biểu theo template

2. **Cập nhật Dữ liệu Cá Nhân:**
   - Học sinh thêm tracking, sbt, etc.
   - Dữ liệu này được preserve qua các lần sync

### 4.3. Cập nhật Template
1. **Admin sửa Template:**
   - Thay đổi thời gian, title, etc.

2. **Apply lại Template:**
   - Chạy API apply để update default data
   - Lần tới học sinh get data sẽ nhận update từ default

## 5. Các Rule Quan Trọng

### 5.1. ID Management
- **Template IDs:** Giữ nguyên, là source of truth
- **Default IDs:** Tạo mới khi tạo từ template
- **User IDs:** Tạo mới, có reference (SlotIDOld, BlockIDOld) tới default

### 5.2. Data Flow
- **Template → Default:** Merge thông minh, giữ custom data
- **Default → User:** Sync fields cơ bản, preserve user data
- **User:** Chỉ update custom fields, không thay đổi structure

### 5.3. Time Calculation
- **Input Duration:** Luôn tính bằng giây (90 = 1.5 phút)
- **Storage:** Duration lưu bằng phút trong database
- **EndTime:** = StartTime + Duration (giây)

### 5.4. Matching Logic
- **Template ↔ Default:** Match theo SlotID (cùng source)
- **Default ↔ User:** Match theo SlotIDOld reference
- **Fallback:** Title + StartTime nếu cần

## 6. API Reference

### Template APIs
- `POST /template-colortime` - Tạo template
- `GET /template-colortime` - Lấy template
- `POST /template-colortime/duplicate` - Duplicate template
- `PUT /template-colortime/slot` - Update slot
- `POST /template-colortime/apply` - Apply to default

### Default APIs
- Internal operations, không có public API trực tiếp

### User APIs
- `GET /colortime/week` - Lấy tuần colortime (tự động sync)
- `POST /colortime/week` - Tạo tuần colortime
- `PUT /colortime/week` - Update tuần colortime

## 7. Troubleshooting

### Lỗi Thường Gặp
1. **Template không apply:** Check term dates, template existence
2. **Data bị mất:** Check merge logic, ensure custom fields preserved
3. **Time calculation:** Verify duration input (seconds) vs storage (minutes)
4. **ID conflicts:** Ensure proper ID generation for new entities

### Debug Tips
- Log merge operations để track data flow
- Check SlotIDOld/BlockIDOld references
- Verify time calculations (seconds vs minutes)
- Test với data nhỏ trước khi production
