package penghulu

import (
	"fmt"
	"net/http"

	structs "simnikah/internal/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// InDB struct untuk dependency injection
type InDB struct {
	DB *gorm.DB
}

// requirePenghulu memvalidasi bahwa user yang login adalah penghulu aktif.
// Mengembalikan data penghulu, user ID, dan flag ok.
func (h *InDB) requirePenghulu(c *gin.Context) (*structs.Penghulu, string, bool) {
	userIDValue, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "Unauthorized",
			"error":   "User ID tidak ditemukan",
		})
		return nil, "", false
	}

	roleValue, roleExists := c.Get("role")
	if !roleExists || roleValue.(string) != structs.UserRolePenghulu {
		c.JSON(http.StatusForbidden, gin.H{
			"success": false,
			"message": "Forbidden",
			"error":   "Hanya user dengan role Penghulu yang dapat mengakses aksi ini",
		})
		return nil, "", false
	}

	userID := fmt.Sprint(userIDValue)
	var user structs.Users
	if err := h.DB.Select("user_id, role, status").Where("user_id = ?", userID).First(&user).Error; err != nil {
		c.JSON(http.StatusForbidden, gin.H{
			"success": false,
			"message": "Forbidden",
			"error":   "Akun penghulu tidak valid atau tidak aktif",
		})
		return nil, "", false
	}

	if user.Role != structs.UserRolePenghulu || user.Status != structs.UserStatusAktif {
		c.JSON(http.StatusForbidden, gin.H{
			"success": false,
			"message": "Forbidden",
			"error":   "Role atau status akun tidak memenuhi syarat Penghulu aktif",
		})
		return nil, "", false
	}

	var penghulu structs.Penghulu
	if err := h.DB.Where("user_id = ?", userID).First(&penghulu).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "Data penghulu tidak ditemukan",
			"error":   err.Error(),
		})
		return nil, "", false
	}

	return &penghulu, userID, true
}

// GetJadwalPenugasan menampilkan daftar penugasan yang sudah disetujui untuk penghulu yang login.
// Response memunculkan AlamatLengkap, Latitude, Longitude, dan URL navigasi secara jelas
// agar penghulu bisa menjadikannya patokan rute (Google Maps, Waze, OSM).
func (h *InDB) GetJadwalPenugasan(c *gin.Context) {
	penghulu, _, ok := h.requirePenghulu(c)
	if !ok {
		return
	}

	var assignments []structs.PendaftaranJadwal
	if err := h.DB.
		Table("pendaftaran_nikahs").
		Select("id, nama_suami, umur_suami, nama_istri, umur_istri, tanggal_nikah, waktu_nikah, tempat_nikah, alamat_akad, latitude, longitude, status_pendaftaran, penghulu_id").
		Where("penghulu_id = ? AND status_pendaftaran = ?", penghulu.ID, structs.StatusPendaftaranPenghuluDitugaskan).
		Order("tanggal_nikah ASC, waktu_nikah ASC").
		Find(&assignments).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "Data jadwal tidak ditemukan",
			"error":   err.Error(),
		})
		return
	}

	jadwal := make([]gin.H, 0, len(assignments))
	for _, a := range assignments {
		item := gin.H{
			"id":                 a.ID,
			"nama_suami":         a.Nama_suami,
			"umur_suami":         a.Umur_suami,
			"nama_istri":         a.Nama_istri,
			"umur_istri":         a.Umur_istri,
			"tanggal_nikah":      a.Tanggal_nikah.Format("2006-01-02"),
			"waktu_nikah":        a.Waktu_nikah,
			"tempat_nikah":       a.Tempat_nikah,
			"alamat_lengkap":     a.Alamat_akad,
			"status_pendaftaran": a.Status_pendaftaran,
		}

		// Geolocation info untuk navigasi
		if a.Latitude != nil && a.Longitude != nil {
			item["latitude"] = *a.Latitude
			item["longitude"] = *a.Longitude
			item["has_coordinates"] = true
			item["google_maps_url"] = fmt.Sprintf("https://www.google.com/maps/search/?api=1&query=%f,%f", *a.Latitude, *a.Longitude)
			item["google_maps_directions_url"] = fmt.Sprintf("https://www.google.com/maps/dir/?api=1&destination=%f,%f", *a.Latitude, *a.Longitude)
			item["waze_url"] = fmt.Sprintf("https://www.waze.com/ul?ll=%f,%f&navigate=yes", *a.Latitude, *a.Longitude)
			item["osm_url"] = fmt.Sprintf("https://www.openstreetmap.org/?mlat=%f&mlon=%f&zoom=16", *a.Latitude, *a.Longitude)
		} else {
			item["has_coordinates"] = false
		}

		// Flag untuk nikah di luar KUA (penghulu perlu datang ke lokasi)
		item["is_outside_kua"] = a.Tempat_nikah == structs.TempatNikahDiLuarKUA

		jadwal = append(jadwal, item)
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Jadwal penugasan berhasil diambil",
		"data": gin.H{
			"penghulu_id":   penghulu.ID,
			"penghulu_nama": penghulu.Nama_lengkap,
			"total":         len(jadwal),
			"jadwal":        jadwal,
		},
	})
}
