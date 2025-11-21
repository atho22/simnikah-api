package structs

import "time"

// SimNikah Models
// CalonPasangan - Model minimal untuk form sederhana
// Hanya menyimpan data yang diinput dari form sederhana: nama, pendidikan, umur (via tanggal_lahir)
type CalonPasangan struct {
	ID                  uint      `gorm:"primaryKey" json:"id"`
	User_id             string    `gorm:"size:20;not null;unique" json:"id_pengguna"`
	NIK                 string    `gorm:"size:16;unique" json:"nik"` // Optional, temporary untuk search
	Nama_lengkap        string    `gorm:"size:100;not null" json:"nama_lengkap"` // Dari form sederhana
	Tanggal_lahir       time.Time `gorm:"not null" json:"tanggal_lahir"` // Diperlukan untuk umur, dihitung dari form
	Jenis_kelamin       string    `gorm:"type:VARCHAR(1);not null" json:"jenis_kelamin"` // L/P (diperlukan untuk identifikasi)
	Pendidikan_terakhir string    `gorm:"size:50" json:"pendidikan_terakhir"` // Dari form sederhana
	// Field lain dihapus - tidak digunakan di form sederhana
	// Jika diperlukan field lain (alamat, tempat_lahir, dll) bisa ditambahkan kembali untuk scalability
	Created_at          time.Time `json:"dibuat_pada"`
	Updated_at          time.Time `json:"diperbarui_pada"`
}


// ==================== FEEDBACK PERNIKAHAN ====================

// FeedbackPernikahan untuk feedback, saran, kritik, atau laporan dari catin setelah pernikahan selesai
type FeedbackPernikahan struct {
	ID                uint      `gorm:"primaryKey" json:"id"`
	Pendaftaran_id    uint      `gorm:"not null" json:"id_pendaftaran"` // ID pendaftaran nikah
	User_id           string    `gorm:"size:20;not null" json:"id_pengguna"` // User ID yang memberikan feedback (catin)
	Jenis_feedback    string    `gorm:"size:20;not null" json:"jenis_feedback"` // "Rating", "Saran", "Kritik", "Laporan"
	Rating            *int      `json:"rating"` // 1-5, hanya untuk jenis "Rating"
	Judul             string    `gorm:"size:200;not null" json:"judul"`
	Pesan             string    `gorm:"type:TEXT;not null" json:"pesan"` // Isi feedback/saran/kritik/laporan
	Status_baca       string    `gorm:"size:20;not null;default:'Belum Dibaca'" json:"status_baca"` // "Belum Dibaca", "Sudah Dibaca"
	Dibaca_oleh       string    `gorm:"size:20" json:"dibaca_oleh"` // ID kepala KUA yang membaca
	Dibaca_pada       *time.Time `json:"dibaca_pada"`
	Created_at        time.Time `json:"dibuat_pada"`
	Updated_at        time.Time `json:"diperbarui_pada"`
}

// WaliNikah untuk data wali nikah (wali untuk calon pengantin perempuan)
// Data sederhana: hanya nama wali (dengan bin) dan hubungan nasab
type WaliNikah struct {
	ID                    uint      `gorm:"primaryKey" json:"id"`
	Pendaftaran_id        uint      `gorm:"not null" json:"id_pendaftaran"` // Foreign key ke PendaftaranNikah
	Nama_dan_bin          string    `gorm:"size:100;not null" json:"nama_dan_bin"` // Contoh: "Abdullah bin Muhammad"
	Hubungan_wali         string    `gorm:"size:50;not null" json:"hubungan_wali"` // Ayah Kandung, Kakek, Saudara Laki-Laki Kandung, dll
	Created_at            time.Time `json:"dibuat_pada"`
	Updated_at            time.Time `json:"diperbarui_pada"`
}

type PendaftaranNikah struct {
	ID                   uint       `gorm:"primaryKey" json:"id"`
	Nomor_pendaftaran    string     `gorm:"size:20;not null;unique" json:"nomor_pendaftaran"`
	Pendaftar_id         string     `gorm:"size:20;not null" json:"id_pendaftar"` // User ID yang mendaftar (suami atau istri)
	Calon_suami_id       string     `gorm:"size:20;not null" json:"id_calon_suami"`
	Calon_istri_id       string     `gorm:"size:20;not null" json:"id_calon_istri"`
	Wali_nikah_id        *uint      `json:"id_wali_nikah"` // Foreign key ke WaliNikah
	Tanggal_pendaftaran  time.Time  `gorm:"not null" json:"tanggal_pendaftaran"`
	Tanggal_nikah        time.Time  `gorm:"not null" json:"tanggal_nikah"`
	Waktu_nikah          string     `gorm:"size:10;not null" json:"waktu_nikah"` // format: HH:MM
	Tempat_nikah         string     `gorm:"size:100;not null" json:"tempat_nikah"`
	// Nomor_dispensasi dihapus - tidak digunakan di form sederhana (bisa ditambahkan kembali jika diperlukan)
	Alamat_akad          string     `gorm:"size:200" json:"alamat_akad"`
	Latitude             *float64   `json:"latitude"`                                                   // Koordinat lintang untuk alamat nikah di luar KUA
	Longitude            *float64   `json:"longitude"`                                                  // Koordinat bujur untuk alamat nikah di luar KUA
	Status_pendaftaran   string     `gorm:"size:40;not null;default:'Draft'" json:"status_pendaftaran"` // Use constants from constants.go
	// Status_bimbingan dihapus - tidak digunakan di form sederhana (bisa ditambahkan kembali jika diperlukan)
	Penghulu_id          *uint      `json:"id_penghulu"`                                                // ID penghulu yang ditugaskan
	Penghulu_assigned_by string     `gorm:"size:20" json:"penghulu_ditugaskan_oleh"`                    // ID kepala KUA yang assign
	Penghulu_assigned_at *time.Time `json:"penghulu_ditugaskan_pada"`                                   // Waktu assign penghulu
	Catatan              string     `gorm:"size:500" json:"catatan"`
	Disetujui_oleh       string     `gorm:"size:20" json:"disetujui_oleh"`
	Disetujui_pada       *time.Time `json:"disetujui_pada"`
	Created_at           time.Time  `json:"dibuat_pada"`
	Updated_at           time.Time  `json:"diperbarui_pada"`
}


// ==================== USER MANAGEMENT & ROLES ====================

// Users untuk user authentication dan role management
type Users struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	User_id    string    `gorm:"size:20;not null;unique" json:"id_pengguna"`
	Username   string    `gorm:"size:50;not null;unique" json:"nama_pengguna"`
	Email      string    `gorm:"size:100;not null;unique" json:"email"`
	Password   string    `gorm:"size:255;not null" json:"kata_sandi"`            // hashed with bcrypt
	Role       string    `gorm:"size:20;not null" json:"peran"`                  // user_biasa, penghulu, staff, kepala_kua
	Status     string    `gorm:"size:20;not null;default:'Aktif'" json:"status"` // Use constants from constants.go
	Nama       string    `gorm:"size:100;not null" json:"nama"`                  // Nama lengkap user
	Created_at time.Time `gorm:"autoCreateTime" json:"dibuat_pada"`
	Updated_at time.Time `gorm:"autoUpdateTime" json:"diperbarui_pada"`
}


type StaffKUA struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	User_id      string    `gorm:"size:20;not null;unique" json:"id_pengguna"`
	NIP          string    `gorm:"size:30;unique" json:"nip"`
	Nama_lengkap string    `gorm:"size:100;not null" json:"nama_lengkap"`
	Jabatan      string    `gorm:"size:50;not null" json:"jabatan"` // Staff, Penghulu, Kepala KUA
	Bagian       string    `gorm:"size:50;not null" json:"bagian"`  // Pendaftaran, Verifikasi, dll
	No_hp        string    `gorm:"size:15" json:"nomor_telepon"`
	Email        string    `gorm:"size:100" json:"email"`
	Alamat       string    `gorm:"size:200" json:"alamat"`
	Status       string    `gorm:"size:20;not null;default:'Aktif'" json:"status"` // Use constants from constants.go
	Created_at   time.Time `json:"dibuat_pada"`
	Updated_at   time.Time `json:"diperbarui_pada"`
}

// Penghulu untuk data penghulu
type Penghulu struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	User_id      string    `gorm:"size:20;not null;unique" json:"id_pengguna"`
	NIP          string    `gorm:"size:30;unique" json:"nip"`
	Nama_lengkap string    `gorm:"size:100;not null" json:"nama_lengkap"`
	No_hp        string    `gorm:"size:15" json:"nomor_telepon"`
	Email        string    `gorm:"size:100" json:"email"`
	Alamat       string    `gorm:"size:200" json:"alamat"`
	Status       string    `gorm:"size:20;not null;default:'Aktif'" json:"status"` // Use constants from constants.go
	Jumlah_nikah int       `gorm:"default:0" json:"jumlah_nikah"`
	Rating       float64   `gorm:"default:0" json:"rating"`
	Created_at   time.Time `json:"dibuat_pada"`
	Updated_at   time.Time `json:"diperbarui_pada"`
}

// ==================== ADDITIONAL SIMNIKAH MODELS ====================

// Notifikasi untuk user
type Notifikasi struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	User_id     string    `gorm:"size:20;not null" json:"id_pengguna"`
	Judul       string    `gorm:"size:100;not null" json:"judul"`
	Pesan       string    `gorm:"size:500;not null" json:"pesan"`
	Tipe        string    `gorm:"size:10;not null;default:'Info'" json:"tipe"`                  // Use constants from constants.go
	Status_baca string    `gorm:"size:20;not null;default:'Belum Dibaca'" json:"status_dibaca"` // Use constants from constants.go
	Link        string    `gorm:"size:200" json:"tautan"`
	Created_at  time.Time `json:"dibuat_pada"`
	Updated_at  time.Time `json:"diperbarui_pada"`
}


// ==================== FORM PENDAFTARAN NIKAH SEDERHANA ====================

// DataFormPendaftaranSederhana untuk form pendaftaran yang disederhanakan
// Hanya berisi field-field penting untuk memudahkan catin
type DataFormPendaftaranSederhana struct {
	// Calon Laki-laki
	CalonLakiLaki struct {
		NamaDanBin       string `json:"nama_dan_bin" binding:"required"` // Contoh: "Ahmad bin Abdullah"
		PendidikanAkhir  string `json:"pendidikan_akhir" binding:"required"`
		Umur             int    `json:"umur" binding:"required,min=1"`
	} `json:"calon_laki_laki" binding:"required"`

	// Calon Perempuan
	CalonPerempuan struct {
		NamaDanBinti     string `json:"nama_dan_binti" binding:"required"` // Contoh: "Siti binti Abdullah"
		PendidikanAkhir  string `json:"pendidikan_akhir" binding:"required"`
		Umur             int    `json:"umur" binding:"required,min=1"`
	} `json:"calon_perempuan" binding:"required"`

	// Lokasi Nikah
	LokasiNikah struct {
		TempatNikah      string `json:"tempat_nikah" binding:"required"` // "Di KUA" atau "Di Luar KUA"
		TanggalNikah     string `json:"tanggal_nikah" binding:"required"` // Format: YYYY-MM-DD
		WaktuNikah       string `json:"waktu_nikah" binding:"required"`   // Format: HH:MM
		
		// Hanya jika tempat_nikah = "Di Luar KUA"
		AlamatNikah      string `json:"alamat_nikah"`              // Alamat lengkap
		DetailAlamat     string `json:"detail_alamat"`             // Detail alamat (RT/RW/dll)
		Kelurahan        string `json:"kelurahan"`                 // Nama kelurahan (hanya lingkup Banjarmasin Utara)
	} `json:"lokasi_nikah" binding:"required"`

	// Wali Nikah (untuk calon pengantin perempuan)
	// Data sederhana: hanya nama wali (dengan bin) dan hubungan nasab
	WaliNikah struct {
		NamaDanBin      string `json:"nama_dan_bin" binding:"required"` // Contoh: "Abdullah bin Muhammad"
		HubunganWali    string `json:"hubungan_wali" binding:"required"` // Ayah Kandung, Kakek, Saudara Laki-Laki Kandung, dll
	} `json:"wali_nikah" binding:"required"`
}

// KelurahanBanjarmasinUtara - Daftar kelurahan di Kecamatan Banjarmasin Utara
var KelurahanBanjarmasinUtara = []string{
	"Sungai Miai",
	"Sungai Andai",
	"Surgi Mufti",
	"Pangeran",
	"Kuin Utara",
	"Antasan Kecil Timur",
	"Alalak Utara",
	"Alalak Tengah",
	"Alalak Selatan",
}

