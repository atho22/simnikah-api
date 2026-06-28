package structs

import "time"

// ==================== PENDAFTARAN NIKAH (SCHEDULING-ONLY) ====================

// PendaftaranNikah menyimpan data referensi pasangan (dari Excel Kepala KUA) + data inti penjadwalan.
// Field Nama_suami/Umur_suami/Nama_istri/Umur_istri adalah DATA REFERENSI display, BUKAN variabel keputusan FC.
// Variabel keputusan FC hanya: Tanggal_nikah, Waktu_nikah, Tempat_nikah, Latitude, Longitude.
// Field administratif SIMKAH (Nomor_pendaftaran, Wali_nikah_id, dll) sudah dihapus.
type PendaftaranNikah struct {
	ID                 uint       `gorm:"primaryKey" json:"id"`
	Nama_suami         string     `gorm:"size:100" json:"nama_suami"`                                      // data referensi dari Excel (bukan variabel FC)
	Umur_suami         int        `gorm:"default:0" json:"umur_suami"`                                     // data referensi dari Excel (bukan variabel FC)
	Nama_istri         string     `gorm:"size:100" json:"nama_istri"`                                      // data referensi dari Excel (bukan variabel FC)
	Umur_istri         int        `gorm:"default:0" json:"umur_istri"`                                     // data referensi dari Excel (bukan variabel FC)
	Tanggal_nikah      time.Time  `gorm:"not null;index:idx_pendaftaran_status_tanggal" json:"tanggal_nikah"`
	Waktu_nikah        string     `gorm:"size:10;not null" json:"waktu_nikah"`                           // format: HH:MM
	Tempat_nikah       string     `gorm:"size:100;not null" json:"tempat_nikah"`                          // "Di KUA" atau "Di Luar KUA"
	Alamat_akad        string     `gorm:"size:200" json:"alamat_akad"`                                    // alamat lengkap jika di luar KUA
	Latitude           *float64   `json:"latitude"`                                                        // koordinat lintang
	Longitude          *float64   `json:"longitude"`                                                       // koordinat bujur
	Status_pendaftaran string     `gorm:"size:40;not null;default:'Menunggu Penugasan';index:idx_pendaftaran_status_tanggal" json:"status_pendaftaran"`
	Pendaftar_id       string     `gorm:"size:20;index:idx_pendaftaran_pendaftar_id" json:"pendaftar_id"`                          // user_id catin yang mendaftar
	Penghulu_id        *uint      `gorm:"index:idx_pendaftaran_penghulu_id" json:"penghulu_id"`           // ID penghulu yang ditugaskan
	Created_at         time.Time  `json:"dibuat_pada"`
	Updated_at         time.Time  `json:"diperbarui_pada"`
}

// PendaftaranJadwal subset ringan untuk rekomendasi, pengecekan slot, dan tampilan jadwal penghulu.
type PendaftaranJadwal struct {
	ID                 uint       `json:"id"`
	Nama_suami         string     `json:"nama_suami"`
	Umur_suami         int        `json:"umur_suami"`
	Nama_istri         string     `json:"nama_istri"`
	Umur_istri         int        `json:"umur_istri"`
	Tanggal_nikah      time.Time  `json:"tanggal_nikah"`
	Waktu_nikah        string     `json:"waktu_nikah"`
	Tempat_nikah       string     `json:"tempat_nikah"`
	Alamat_akad        string     `json:"alamat_akad"`
	Latitude           *float64   `json:"latitude"`
	Longitude          *float64   `json:"longitude"`
	Status_pendaftaran string     `json:"status_pendaftaran"`
	Penghulu_id        *uint      `json:"penghulu_id"`
}

// ==================== USER MANAGEMENT & ROLES ====================

// Users untuk authentication dan role management
type Users struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	User_id       string    `gorm:"size:20;not null;unique" json:"id_pengguna"`
	Username      string    `gorm:"size:50;not null;unique" json:"nama_pengguna"`
	Email         string    `gorm:"size:100;not null;unique" json:"email"`
	Password      string    `gorm:"size:255;not null" json:"kata_sandi"`
	Role          string    `gorm:"size:20;not null" json:"peran"`
	Status        string    `gorm:"size:20;not null;default:'Aktif'" json:"status"`
	Nama          string    `gorm:"size:100;not null" json:"nama"`
	Profile_photo string    `gorm:"size:500" json:"foto_profil"`
	Created_at    time.Time `gorm:"autoCreateTime" json:"dibuat_pada"`
	Updated_at    time.Time `gorm:"autoUpdateTime" json:"diperbarui_pada"`
}

type StaffKUA struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	User_id      string    `gorm:"size:20;not null;unique" json:"id_pengguna"`
	NIP          string    `gorm:"size:30;unique" json:"nip"`
	Nama_lengkap string    `gorm:"size:100;not null" json:"nama_lengkap"`
	Jabatan      string    `gorm:"size:50;not null" json:"jabatan"`
	Bagian       string    `gorm:"size:50;not null" json:"bagian"`
	No_hp        string    `gorm:"size:15" json:"nomor_telepon"`
	Email        string    `gorm:"size:100" json:"email"`
	Alamat       string    `gorm:"size:200" json:"alamat"`
	Status       string    `gorm:"size:20;not null;default:'Aktif'" json:"status"`
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
	Latitude     *float64  `json:"latitude"`
	Longitude    *float64  `json:"longitude"`
	Status       string    `gorm:"size:20;not null;default:'Aktif'" json:"status"`
	Jumlah_nikah int       `gorm:"default:0" json:"jumlah_nikah"`
	Rating       float64   `gorm:"default:0" json:"rating"`
	Created_at   time.Time `json:"dibuat_pada"`
	Updated_at   time.Time `json:"diperbarui_pada"`
}

// ==================== NOTIFIKASI ====================

// Notifikasi untuk user
type Notifikasi struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	User_id     string    `gorm:"size:20;not null" json:"id_pengguna"`
	Judul       string    `gorm:"size:100;not null" json:"judul"`
	Pesan       string    `gorm:"size:500;not null" json:"pesan"`
	Tipe        string    `gorm:"size:10;not null;default:'Info'" json:"tipe"`
	Status_baca string    `gorm:"size:20;not null;default:'Belum Dibaca'" json:"status_dibaca"`
	Link        string    `gorm:"size:200" json:"tautan"`
	Created_at  time.Time `json:"dibuat_pada"`
	Updated_at  time.Time `json:"diperbarui_pada"`
}
