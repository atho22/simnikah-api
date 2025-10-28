package validator

import (
	"fmt"
	"net/http"

	"simnikah/pkg/utils"

	"github.com/gin-gonic/gin"
)

// ValidateParentFields validates conditional fields for parents using reflection
func ValidateParentFields(parent interface{}, parentType string, c *gin.Context) {
	// Use type assertion to access fields
	switch p := parent.(type) {
	case struct {
		StatusKeberadaan   string
		Nama               string
		Nik                string
		Kewarganegaraan    string
		NegaraAsal         string
		NomorPaspor        string
		TempatLahir        string
		TanggalLahir       string
		Agama              string
		Pekerjaan          string
		DeskripsiPekerjaan string
		Alamat             string
	}:
		if p.StatusKeberadaan == "Hidup" {
			// All fields are required if parent is alive
			if p.Nama == "" {
				c.JSON(http.StatusBadRequest, gin.H{
					"success": false,
					"message": "Validasi gagal",
					"error":   fmt.Sprintf("Nama %s wajib diisi", parentType),
					"field":   "nama",
					"type":    "required",
				})
				return
			}
			if p.Kewarganegaraan == "" {
				c.JSON(http.StatusBadRequest, gin.H{
					"success": false,
					"message": "Validasi gagal",
					"error":   fmt.Sprintf("Kewarganegaraan %s wajib diisi", parentType),
					"field":   "kewarganegaraan",
					"type":    "required",
				})
				return
			}
			if p.Agama == "" {
				c.JSON(http.StatusBadRequest, gin.H{
					"success": false,
					"message": "Validasi gagal",
					"error":   fmt.Sprintf("Agama %s wajib diisi", parentType),
					"field":   "agama",
					"type":    "required",
				})
				return
			}
			if p.Pekerjaan == "" {
				c.JSON(http.StatusBadRequest, gin.H{
					"success": false,
					"message": "Validasi gagal",
					"error":   fmt.Sprintf("Pekerjaan %s wajib diisi", parentType),
					"field":   "pekerjaan",
					"type":    "required",
				})
				return
			}
			if p.Alamat == "" {
				c.JSON(http.StatusBadRequest, gin.H{
					"success": false,
					"message": "Validasi gagal",
					"error":   fmt.Sprintf("Alamat %s wajib diisi", parentType),
					"field":   "alamat",
					"type":    "required",
				})
				return
			}

			// Validate citizenship-specific fields
			if p.Kewarganegaraan == "WNA" {
				if p.NegaraAsal == "" {
					c.JSON(http.StatusBadRequest, gin.H{
						"success": false,
						"message": "Negara asal %s wajib diisi untuk WNA",
						"error":   fmt.Sprintf("Negara asal %s wajib diisi untuk WNA", parentType),
						"field":   "negara_asal",
						"type":    "required",
					})
					return
				}
				if p.NomorPaspor == "" {
					c.JSON(http.StatusBadRequest, gin.H{
						"success": false,
						"message": "Validasi gagal",
						"error":   fmt.Sprintf("Nomor paspor %s wajib diisi untuk WNA", parentType),
						"field":   "nomor_paspor",
						"type":    "required",
					})
					return
				}
			}

			// Validate occupation description if occupation is "Lainnya"
			if p.Pekerjaan == "Lainnya" && p.DeskripsiPekerjaan == "" {
				c.JSON(http.StatusBadRequest, gin.H{
					"success": false,
					"message": "Validasi gagal",
					"error":   fmt.Sprintf("Deskripsi pekerjaan %s wajib diisi jika pekerjaan adalah 'Lainnya'", parentType),
					"field":   "deskripsi_pekerjaan",
					"type":    "required",
				})
				return
			}
		}
	}
}

// ValidatePersonFields validates conditional fields for groom and bride using reflection
func ValidatePersonFields(person interface{}, personType string, c *gin.Context) {
	// Use type assertion to access fields
	switch p := person.(type) {
	case struct {
		NamaLengkap        string
		Nik                string
		Kewarganegaraan    string
		NomorPaspor        string
		TempatLahir        string
		TanggalLahir       string
		Status             string
		Agama              string
		Pendidikan         string
		Pekerjaan          string
		DeskripsiPekerjaan string
		NomorTelepon       string
		Email              string
		Alamat             string
	}:
		// Validate citizenship-specific fields
		if p.Kewarganegaraan == "WNA" {
			if p.NomorPaspor == "" {
				c.JSON(http.StatusBadRequest, gin.H{
					"success": false,
					"message": "Validasi gagal",
					"error":   fmt.Sprintf("Nomor paspor %s wajib diisi untuk WNA", personType),
					"field":   "nomor_paspor",
					"type":    "required",
				})
				return
			}
		}

		// Validate occupation description if occupation is "Lainnya"
		if p.Pekerjaan == "Lainnya" && p.DeskripsiPekerjaan == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "Validasi gagal",
				"error":   fmt.Sprintf("Deskripsi pekerjaan %s wajib diisi jika pekerjaan adalah 'Lainnya'", personType),
				"field":   "deskripsi_pekerjaan",
				"type":    "required",
			})
			return
		}

		// Validate phone number format
		if len(p.NomorTelepon) < 10 || !utils.StartsWith(p.NomorTelepon, "08") {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "Validasi gagal",
				"error":   fmt.Sprintf("Nomor telepon %s harus minimal 10 digit dan dimulai dengan '08'", personType),
				"field":   "nomor_telepon",
				"type":    "format",
			})
			return
		}

		// Validate email format (basic validation)
		if !utils.IsValidEmail(p.Email) {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "Validasi gagal",
				"error":   fmt.Sprintf("Format email %s tidak valid", personType),
				"field":   "email",
				"type":    "format",
			})
			return
		}

		// Validate minimum length requirements
		if len(p.NamaLengkap) < 3 {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "Validasi gagal",
				"error":   fmt.Sprintf("Nama lengkap %s minimal 3 karakter", personType),
				"field":   "nama_lengkap",
				"type":    "min_length",
			})
			return
		}
		if len(p.TempatLahir) < 2 {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "Validasi gagal",
				"error":   fmt.Sprintf("Tempat lahir %s minimal 2 karakter", personType),
				"field":   "tempat_lahir",
				"type":    "min_length",
			})
			return
		}
		if len(p.Alamat) < 10 {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "Validasi gagal",
				"error":   fmt.Sprintf("Alamat %s minimal 10 karakter", personType),
				"field":   "alamat",
				"type":    "min_length",
			})
			return
		}
	}
}

// CheckValidValue checks if a string is in a slice of valid values
func CheckValidValue(validValues []string, value string) bool {
	for _, validValue := range validValues {
		if validValue == value {
			return true
		}
	}
	return false
}
