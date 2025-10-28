# 📊 Status Management Best Practices

## 🎯 TL;DR

**Gunakan VARCHAR/STRING, BUKAN ENUM!**

- ✅ **STRING**: Fleksibel, mudah maintain, type-safe di Go
- ❌ **ENUM**: Rigid, ALTER TABLE untuk perubahan, GORM compatibility issues

---

## 🤔 Mengapa STRING/VARCHAR?

### ✅ Keuntungan STRING (Recommended)

| Aspek | Penjelasan |
|-------|------------|
| **Fleksibilitas** | Tambah/ubah status tanpa ALTER TABLE |
| **GORM Support** | Perfect compatibility dengan GORM |
| **Type Safety** | Go constants memberikan compile-time safety |
| **Readable** | Status jelas di logs dan database viewer |
| **Portable** | Easy migration antar database |
| **No Schema Lock** | Tambah status tanpa database downtime |
| **Validation** | Di application level, lebih kontrol |

### ❌ Kerugian ENUM MySQL

| Aspek | Penjelasan |
|-------|------------|
| **ALTER TABLE** | Setiap perubahan = migration + downtime |
| **GORM Issues** | Mapping ENUM ke struct Go kurang reliable |
| **Not Portable** | Harus convert saat pindah database |
| **Size Limit** | Max 65,535 values (cukup tapi tetap limit) |
| **Error Handling** | Invalid values lebih susah di-handle |
| **Development Speed** | Slower iteration karena butuh migration |

---

## 📁 Struktur File

```
structs/
├── models.go       # Database models dengan VARCHAR
├── constants.go    # Status constants & validation (BARU!)
└── ...
```

---

## 🔨 Implementasi

### 1. Database Model (`structs/models.go`)

```go
type PendaftaranNikah struct {
    ID                 uint   `gorm:"primaryKey"`
    Nomor_pendaftaran  string `gorm:"size:20;not null"`
    
    // ✅ GOOD: VARCHAR dengan constants
    Status_pendaftaran string `gorm:"size:40;not null;default:'Draft'"` 
    Status_bimbingan   string `gorm:"size:30;not null;default:'Belum'"`
    
    // ❌ BAD: ENUM
    // Status_pendaftaran string `gorm:"type:ENUM('Draft','Pending','Approved')"`
}
```

### 2. Constants Definition (`structs/constants.go`)

```go
package structs

// Status Pendaftaran Nikah
const (
    StatusDraft                      = "Draft"
    StatusMenungguVerifikasi         = "Menunggu Verifikasi"
    StatusMenungguPengumpulanBerkas  = "Menunggu Pengumpulan Berkas"
    StatusBerkasDiterima             = "Berkas Diterima"
    StatusSelesai                    = "Selesai"
)

// Validation Map
var ValidStatusPendaftaran = map[string]bool{
    StatusDraft:                     true,
    StatusMenungguVerifikasi:        true,
    StatusMenungguPengumpulanBerkas: true,
    StatusBerkasDiterima:            true,
    StatusSelesai:                   true,
}

// Validation Function
func IsValidStatusPendaftaran(status string) bool {
    return ValidStatusPendaftaran[status]
}
```

### 3. Usage dalam Handler

```go
package catin

import "simnikah/structs"

func CreatePendaftaran(c *gin.Context) {
    pendaftaran := structs.PendaftaranNikah{
        Nomor_pendaftaran:  generateNomor(),
        Status_pendaftaran: structs.StatusDraft, // ✅ Type-safe!
        Status_bimbingan:   structs.StatusBimbinganBelum,
    }
    
    // Save to database
    DB.Create(&pendaftaran)
}

func UpdateStatus(c *gin.Context) {
    var input struct {
        Status string `json:"status" binding:"required"`
    }
    
    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(400, gin.H{"error": "Invalid input"})
        return
    }
    
    // ✅ Validate using helper function
    if !structs.IsValidStatusPendaftaran(input.Status) {
        c.JSON(400, gin.H{
            "error": "Status tidak valid",
            "valid_statuses": structs.GetAllStatusPendaftaran(),
        })
        return
    }
    
    // Update
    pendaftaran.Status_pendaftaran = input.Status
    DB.Save(&pendaftaran)
}
```

---

## 🛡️ Validation Strategies

### Application Level (Recommended)

```go
// Di handler
func validateStatus(status string) error {
    if !structs.IsValidStatusPendaftaran(status) {
        return fmt.Errorf("status '%s' tidak valid", status)
    }
    return nil
}
```

### Database Level (Optional - Extra Security)

```sql
-- MySQL Check Constraint (MySQL 8.0.16+)
ALTER TABLE pendaftaran_nikahs
ADD CONSTRAINT chk_status_pendaftaran
CHECK (status_pendaftaran IN (
    'Draft',
    'Menunggu Verifikasi',
    'Menunggu Pengumpulan Berkas',
    'Berkas Diterima',
    'Selesai'
));
```

**Note**: Check constraints optional. Validasi di Go level sudah cukup.

---

## 📚 Contoh Lengkap

### Status Workflow dengan Constants

```go
// Workflow progression
func GetNextStatus(currentStatus string) (string, error) {
    statusFlow := map[string]string{
        structs.StatusDraft:                     structs.StatusMenungguVerifikasi,
        structs.StatusMenungguVerifikasi:        structs.StatusMenungguPengumpulanBerkas,
        structs.StatusMenungguPengumpulanBerkas: structs.StatusBerkasDiterima,
        structs.StatusBerkasDiterima:            structs.StatusSelesai,
    }
    
    nextStatus, exists := statusFlow[currentStatus]
    if !exists {
        return "", fmt.Errorf("no next status for: %s", currentStatus)
    }
    
    return nextStatus, nil
}
```

### Switch Case dengan Type Safety

```go
func GetStatusMessage(status string) string {
    switch status {
    case structs.StatusDraft:
        return "Formulir masih draft"
    case structs.StatusMenungguVerifikasi:
        return "Menunggu verifikasi staff"
    case structs.StatusSelesai:
        return "Proses selesai"
    default:
        return "Status tidak dikenali"
    }
}
```

### Query dengan Constants

```go
// ✅ GOOD: No typo possible
func GetPendingRegistrations(db *gorm.DB) []structs.PendaftaranNikah {
    var results []structs.PendaftaranNikah
    db.Where("status_pendaftaran = ?", structs.StatusMenungguVerifikasi).
       Find(&results)
    return results
}

// ❌ BAD: Prone to typos
func GetPendingRegistrationsBad(db *gorm.DB) []structs.PendaftaranNikah {
    var results []structs.PendaftaranNikah
    db.Where("status_pendaftaran = ?", "Menunggu Verifikasi"). // Typo risk!
       Find(&results)
    return results
}
```

---

## 🔄 Migration dari ENUM ke VARCHAR

Jika sebelumnya pakai ENUM:

```sql
-- 1. Backup data
CREATE TABLE pendaftaran_nikahs_backup AS 
SELECT * FROM pendaftaran_nikahs;

-- 2. Alter table
ALTER TABLE pendaftaran_nikahs 
MODIFY COLUMN status_pendaftaran VARCHAR(40) NOT NULL DEFAULT 'Draft';

-- 3. Verify
SELECT DISTINCT status_pendaftaran FROM pendaftaran_nikahs;

-- 4. (Optional) Add check constraint
ALTER TABLE pendaftaran_nikahs
ADD CONSTRAINT chk_status_pendaftaran
CHECK (status_pendaftaran IN ('Draft', 'Menunggu Verifikasi', ...));
```

---

## 📊 Performance Comparison

| Aspect | ENUM | VARCHAR(40) | Difference |
|--------|------|-------------|------------|
| **Storage** | 1-2 bytes | 40 bytes max | ENUM wins |
| **Query Speed** | Fast (integer) | Fast (string) | ~Equal with index |
| **Index Size** | Smaller | Larger | ENUM wins |
| **Development Speed** | Slow (migration) | Fast (no migration) | **VARCHAR wins** |
| **Flexibility** | Low | High | **VARCHAR wins** |
| **Maintainability** | Hard | Easy | **VARCHAR wins** |

**Kesimpulan**: VARCHAR sedikit lebih besar storage, tapi **jauh lebih fleksibel**!

---

## ✅ Best Practices Summary

1. **✅ DO**: Use VARCHAR/STRING for all status fields
2. **✅ DO**: Define constants in `structs/constants.go`
3. **✅ DO**: Create validation functions
4. **✅ DO**: Use constants everywhere in code
5. **✅ DO**: Provide helper functions untuk dropdown/select
6. **❌ DON'T**: Use MySQL ENUM
7. **❌ DON'T**: Hardcode status strings di code
8. **❌ DON'T**: Skip validation
9. **❌ DON'T**: Mix constants dan string literals

---

## 🎓 Examples

Lihat file `examples/status_usage.go` untuk contoh lengkap:
- Create dengan constants
- Update dengan validation
- Query dengan type safety
- Workflow progression
- Dropdown/select options

---

## 🚀 Getting Started

### 1. Import Constants

```go
import "simnikah/structs"
```

### 2. Use in Code

```go
// Create
pendaftaran := structs.PendaftaranNikah{
    Status_pendaftaran: structs.StatusDraft,
}

// Validate
if !structs.IsValidStatusPendaftaran(newStatus) {
    return errors.New("invalid status")
}

// Update
pendaftaran.Status_pendaftaran = structs.StatusMenungguVerifikasi
```

### 3. API Response

```go
c.JSON(200, gin.H{
    "current_status": pendaftaran.Status_pendaftaran,
    "valid_statuses": structs.GetAllStatusPendaftaran(),
})
```

---

## 📖 References

- [GORM Data Types](https://gorm.io/docs/models.html#Fields-Tags)
- [MySQL VARCHAR vs ENUM](https://dev.mysql.com/doc/refman/8.0/en/enum.html)
- [Go Constants Best Practices](https://go.dev/blog/constants)

---

## 💡 Tips

1. **IDE Auto-complete**: Constants memberikan auto-complete di IDE
2. **Compile-time Safety**: Typo di constants = compile error
3. **Refactoring**: Rename constant = rename semua usage
4. **Testing**: Easy untuk test semua possible status
5. **Documentation**: Constants self-documenting

---

**Kesimpulan**: VARCHAR + Go Constants = Best of both worlds! 🎉

Kamu dapat type safety, flexibility, dan maintainability tinggi.

