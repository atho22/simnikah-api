# 🛠️ SimNikah Developer Reference

## 1. How to: Add a New API Endpoint

### Step 1: Define the Handler
```go
// File: internal/handlers/yourmodule/yourhandler.go
package yourmodule

import (
    "net/http"
    structs "simnikah/internal/models"
    "github.com/gin-gonic/gin"
    "gorm.io/gorm"
)

type InDB struct {
    DB *gorm.DB
}

// NewHandler creates a new resource
func (h *InDB) NewHandler(c *gin.Context) {
    // 1. Get user_id from context (if needed)
    userID, exists := c.Get("user_id")
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{
            "success": false,
            "message": "Unauthorized",
            "error": "User ID tidak ditemukan",
        })
        return
    }

    // 2. Bind request body
    var req struct {
        Field1 string `json:"field1" binding:"required"`
        Field2 int    `json:"field2" binding:"required,min=1"`
    }
    
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "success": false,
            "message": "Format data tidak valid",
            "error": err.Error(),
        })
        return
    }

    // 3. Validate input
    if req.Field1 == "" {
        c.JSON(http.StatusBadRequest, gin.H{
            "success": false,
            "message": "Field1 tidak boleh kosong",
        })
        return
    }

    // 4. Database operation
    data := structs.YourModel{
        Field1: req.Field1,
        Field2: req.Field2,
    }
    
    if err := h.DB.Create(&data).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "success": false,
            "message": "Gagal menyimpan data",
            "error": err.Error(),
        })
        return
    }

    // 5. Return success response
    c.JSON(http.StatusCreated, gin.H{
        "success": true,
        "message": "Data berhasil dibuat",
        "data": data,
    })
}
```

### Step 2: Register Route in main.go
```go
// File: cmd/api/main.go
// Initialize handler
yourHandler := &yourmodule.InDB{DB: DB}

// Register route
simnikahRoutes := r.Group("/simnikah")
{
    // POST with auth middleware
    simnikahRoutes.POST("/your-endpoint", 
        middleware.AuthMiddleware(), 
        yourHandler.NewHandler)
    
    // GET with role middleware
    simnikahRoutes.GET("/your-endpoint", 
        middleware.AuthMiddleware(), 
        middleware.RoleMiddleware("staff"), 
        yourHandler.NewHandler)
    
    // POST with multiple roles
    simnikahRoutes.POST("/your-endpoint", 
        middleware.AuthMiddleware(), 
        middleware.MultiRoleMiddleware("staff", "kepala_kua"), 
        yourHandler.NewHandler)
}
```

### Step 3: Test the Endpoint
```bash
# Test with curl or Postman
curl -X POST http://localhost:8080/simnikah/your-endpoint \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "field1": "value1",
    "field2": 123
  }'
```

---

## 2. How to: Add a New Database Model

### Step 1: Define the Model
```go
// File: internal/models/models.go
package structs

import "time"

type YourModel struct {
    ID          uint      `gorm:"primaryKey" json:"id"`
    User_id     string    `gorm:"size:20;not null" json:"id_pengguna"`
    Field1      string    `gorm:"size:100;not null" json:"field1"`
    Field2      int       `gorm:"default:0" json:"field2"`
    Field3      *string   `gorm:"size:200" json:"field3"` // Optional field
    Created_at  time.Time `json:"dibuat_pada"`
    Updated_at  time.Time `json:"diperbarui_pada"`
}
```

### Step 2: Add to AutoMigrate
```go
// File: cmd/api/main.go, in main() function
if err := DB.AutoMigrate(
    &structs.Users{},
    &structs.YourModel{},  // Add this line
    // ... other models
); err != nil {
    log.Fatal("Database migration failed:", err)
}
```

### Step 3: Use in Handler
```go
// Query
var record structs.YourModel
h.DB.Where("user_id = ?", userID).First(&record)

// Create
newRecord := structs.YourModel{
    User_id: userID,
    Field1: "value",
}
h.DB.Create(&newRecord)

// Update
h.DB.Model(&record).Update("field1", "new_value")

// Delete
h.DB.Delete(&record)
```

---

## 3. How to: Add Authentication to Endpoint

### Option 1: Basic Auth (User must be logged in)
```go
simnikahRoutes.GET("/protected-endpoint", 
    middleware.AuthMiddleware(),
    handler.ProtectedHandler)
```

### Option 2: Role-Based Auth
```go
// Single role
simnikahRoutes.POST("/staff-only", 
    middleware.AuthMiddleware(), 
    middleware.RoleMiddleware("staff"),
    handler.StaffOnlyHandler)

// Multiple roles
simnikahRoutes.GET("/staff-or-kepala", 
    middleware.AuthMiddleware(), 
    middleware.MultiRoleMiddleware("staff", "kepala_kua"),
    handler.StaffOrKepalaHandler)
```

### Option 3: No Auth (Public endpoint)
```go
r.GET("/public-endpoint", handler.PublicHandler)
```

---

## 4. How to: Handle File Uploads (Images)

### Using ImgBB Storage
```go
import (
    "simnikah/pkg/storage"
    "net/http"
)

func (h *InDB) UploadProfilePhoto(c *gin.Context) {
    // Get file from request
    file, err := c.FormFile("photo")
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "success": false,
            "message": "File tidak ditemukan",
        })
        return
    }

    // Open file
    src, err := file.Open()
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "success": false,
            "message": "Gagal membaca file",
        })
        return
    }
    defer src.Close()

    // Upload to ImgBB
    photoURL, err := storage.UploadFileFromMultipart(src, file.Filename)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "success": false,
            "message": "Gagal upload foto",
            "error": err.Error(),
        })
        return
    }

    // Save URL to database
    userID, _ := c.Get("user_id")
    user := structs.Users{}
    h.DB.Where("user_id = ?", userID).First(&user)
    h.DB.Model(&user).Update("profile_photo", photoURL)

    c.JSON(http.StatusOK, gin.H{
        "success": true,
        "message": "Foto berhasil diupload",
        "data": gin.H{
            "url": photoURL,
        },
    })
}
```

---

## 5. How to: Add Validation

### Built-in Validation
```go
type RegisterRequest struct {
    Username string `json:"username" binding:"required,min=3,max=50"`
    Email    string `json:"email" binding:"required,email"`
    Password string `json:"password" binding:"required,min=6"`
    Umur     int    `json:"umur" binding:"required,min=18,max=100"`
}
```

### Custom Validation
```go
func (h *InDB) CreateRegistration(c *gin.Context) {
    var req structs.DataFormPendaftaranSederhana
    
    if err := c.ShouldBindJSON(&req); err != nil {
        // ... handle error
    }

    // Custom validation
    if req.CalonLakiLaki.Umur < 19 {
        c.JSON(http.StatusBadRequest, gin.H{
            "success": false,
            "message": "Umur calon suami minimal 19 tahun",
        })
        return
    }

    if req.CalonPerempuan.Umur < 16 {
        c.JSON(http.StatusBadRequest, gin.H{
            "success": false,
            "message": "Umur calon istri minimal 16 tahun",
        })
        return
    }
}
```

---

## 6. How to: Generate Unique IDs

### User ID (Timestamp-based)
```go
import "time"

userID := fmt.Sprintf("USR%d", time.Now().Unix())
// Result: USR1701929380
```

### Registration Number
```go
import "crypto/md5"

registrationNum := fmt.Sprintf("NIKAH-%s-%d", 
    time.Now().Format("20060102"), 
    rand.Intn(10000))
// Result: NIKAH-20240105-4532
```

### Hash-based ID
```go
import (
    "crypto/md5"
    "fmt"
)

data := fmt.Sprintf("%s%s%d", userID, time.Now().String(), rand.Intn(1000))
hash := md5.Sum([]byte(data))
id := fmt.Sprintf("%x", hash)[:12]
```

---

## 7. How to: Send Notifications

### Create Notification
```go
notification := structs.Notifikasi{
    User_id: targetUserID,
    Judul: "Pendaftaran Disetujui",
    Pesan: "Pendaftaran nikah Anda telah disetujui oleh staff",
    Tipe: structs.NotifikasiTipeSuccess,
    Status_baca: structs.NotifikasiStatusBelumDibaca,
    Link: "/simnikah/pendaftaran/status",
}

if err := h.DB.Create(&notification).Error; err != nil {
    log.Printf("Failed to create notification: %v", err)
}
```

### Mark as Read
```go
h.DB.Model(&structs.Notifikasi{}).
    Where("id = ?", notifikasiID).
    Updates(map[string]interface{}{
        "status_baca": structs.NotifikasiStatusSudahDibaca,
    })
```

---

## 8. How to: Query with Conditions

### Single Condition
```go
var user structs.Users
h.DB.Where("username = ?", username).First(&user)

var registrations []structs.PendaftaranNikah
h.DB.Where("status_pendaftaran = ?", "Draft").Find(&registrations)
```

### Multiple Conditions (AND)
```go
var registrations []structs.PendaftaranNikah
h.DB.Where("status_pendaftaran = ? AND tanggal_nikah >= ?", 
    "Disetujui", 
    time.Now()).Find(&registrations)
```

### OR Condition
```go
var users []structs.Users
h.DB.Where("role = ? OR role = ?", "staff", "penghulu").Find(&users)
```

### IN Condition
```go
var registrations []structs.PendaftaranNikah
statuses := []string{"Draft", "Disetujui"}
h.DB.Where("status_pendaftaran IN ?", statuses).Find(&registrations)
```

### LIKE Query
```go
var users []structs.Users
h.DB.Where("nama LIKE ?", "%Ahmad%").Find(&users)
```

### Order & Limit
```go
var registrations []structs.PendaftaranNikah
h.DB.Where("status_pendaftaran = ?", "Draft").
    Order("created_at DESC").
    Limit(10).
    Find(&registrations)
```

---

## 9. How to: Handle Transactions

```go
func (h *InDB) ComplexOperation(c *gin.Context) {
    // Start transaction
    tx := h.DB.BeginTx(c, nil)
    
    // Operation 1
    if err := tx.Create(&record1).Error; err != nil {
        tx.Rollback()
        c.JSON(http.StatusInternalServerError, gin.H{
            "success": false,
            "message": "Gagal membuat record1",
        })
        return
    }
    
    // Operation 2
    if err := tx.Update(&record2).Error; err != nil {
        tx.Rollback()
        c.JSON(http.StatusInternalServerError, gin.H{
            "success": false,
            "message": "Gagal update record2",
        })
        return
    }
    
    // Commit if all succeed
    tx.Commit()
    
    c.JSON(http.StatusOK, gin.H{
        "success": true,
        "message": "Operasi berhasil",
    })
}
```

---

## 10. How to: Debug Issues

### Check Database Connection
```go
// In main.go after ConnectDB()
if err := DB.Exec("SELECT 1").Error; err != nil {
    log.Fatal("Database connection failed:", err)
}
```

### Debug SQL Queries
```go
// In handlers, use:
h.DB.Debug().Where(...).First(&record)
// This will print the SQL query to console
```

### Check Middleware Order
```go
// Middleware order matters!
// Rate limiter should be global
r.Use(middleware.RateLimiter())

// Auth should be per route
simnikahRoutes.GET("/protected",
    middleware.AuthMiddleware(),  // First check auth
    middleware.RoleMiddleware("staff"),  // Then check role
    handler.Handler)
```

### Print Context Values
```go
func (h *InDB) DebugHandler(c *gin.Context) {
    userID, exists := c.Get("user_id")
    log.Printf("User ID: %v, Exists: %v", userID, exists)
    
    role, _ := c.Get("role")
    log.Printf("Role: %v", role)
}
```

---

## 11. How to: Test Locally

### Start Development Server
```bash
make dev
# or
GIN_MODE=debug go run cmd/api/main.go
```

### Using Postman/Insomnia
```
1. Register: POST http://localhost:8080/register
   Body: {"username": "test", "email": "test@test.com", "password": "123456", "nama": "Test User"}

2. Login: POST http://localhost:8080/login
   Body: {"username": "test", "password": "123456"}
   Save the token from response

3. Protected endpoint: GET http://localhost:8080/profile
   Header: Authorization: Bearer <TOKEN>
```

### Using curl
```bash
# Register
curl -X POST http://localhost:8080/register \
  -H "Content-Type: application/json" \
  -d '{"username":"test","email":"test@test.com","password":"123456","nama":"Test"}'

# Login
curl -X POST http://localhost:8080/login \
  -H "Content-Type: application/json" \
  -d '{"username":"test","password":"123456"}'

# Use token
curl -X GET http://localhost:8080/profile \
  -H "Authorization: Bearer YOUR_TOKEN"
```

---

## 12. Environment Variables Checklist

```bash
# Database Configuration
DB_HOST=127.0.0.1           # or your MySQL server IP
DB_PORT=3306                # MySQL port
DB_USER=simnikah_user       # Database user
DB_PASSWORD=your_password   # Database password
DB_NAME=simnikah            # Database name

# Server
PORT=8080                   # Server port
GIN_MODE=debug              # debug or release

# Image Storage
IMGBB_API_KEY=your_api_key  # ImgBB API key (optional, anonymous if not set)

# CORS
ALLOWED_ORIGINS=http://localhost:3000,http://localhost:5173,https://example.com
```

---

## 13. Common Error Fixes

| Error | Cause | Fix |
|-------|-------|-----|
| `pq: password authentication failed` | Wrong DB_PASSWORD | Check DB_PASSWORD env var |
| `CORS policy: No 'Access-Control-Allow-Origin'` | Origin not allowed | Add origin to ALLOWED_ORIGINS |
| `404 Not Found` | Wrong endpoint path | Check route in main.go |
| `401 Unauthorized` | Invalid/missing JWT token | Include Authorization header |
| `403 Forbidden` | Wrong role | Check RoleMiddleware(role) |
| `400 Bad Request` | Invalid JSON/validation error | Check request body format |
| `500 Internal Server Error` | Database error | Check database connection, model |
| `file: no such file or directory` | Image upload failed | Check IMGBB_API_KEY |

---

## 14. Performance Tips

- Use `h.DB.Debug()` only in development, remove in production
- Add database indexes for frequently queried fields
- Use pagination for list endpoints: `.Limit(10).Offset((page-1)*10)`
- Use select to fetch only needed columns: `.Select("id", "nama", "email")`
- Cache static data like HariLiburNasional in memory
- Use connection pooling (already configured in config.go)

---

## 15. Security Checklist

- [ ] Always use parameterized queries (GORM does this)
- [ ] Hash passwords with bcrypt (crypto.HashPassword)
- [ ] Validate all user input before database operations
- [ ] Check authorization before sensitive operations
- [ ] Use HTTPS in production
- [ ] Set secure cookie flags (HttpOnly, Secure, SameSite)
- [ ] Implement rate limiting (already done)
- [ ] Sanitize error messages (don't expose internal details)
- [ ] Keep dependencies updated
- [ ] Use environment variables for secrets (not hardcoded)

---

**Last Updated: December 6, 2025**
