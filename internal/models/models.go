package structs

import "time"

// SimNikah Models
type CalonPasangan struct {
	ID                  uint      `gorm:"primaryKey" json:"id"`
	User_id             string    `gorm:"size:20;not null;unique" json:"id_pengguna"`
	NIK                 string    `gorm:"size:16;not null;unique" json:"nik"`
	Nama_lengkap        string    `gorm:"size:100;not null" json:"nama_lengkap"`
	Tempat_lahir        string    `gorm:"size:50;not null" json:"tempat_lahir"`
	Tanggal_lahir       time.Time `gorm:"not null" json:"tanggal_lahir"`
	Jenis_kelamin       string    `gorm:"type:VARCHAR(1);not null" json:"jenis_kelamin"` // postgres: gunakan VARCHAR(1) untuk L/P
	Alamat              string    `gorm:"size:200;not null" json:"alamat"`
	RT                  string    `gorm:"size:3" json:"rt"`
	RW                  string    `gorm:"size:3" json:"rw"`
	Kelurahan           string    `gorm:"size:50;not null" json:"kelurahan"`
	Kecamatan           string    `gorm:"size:50;not null" json:"kecamatan"`
	Kabupaten           string    `gorm:"size:50;not null" json:"kabupaten"`
	Provinsi            string    `gorm:"size:50;not null" json:"provinsi"`
	Agama               string    `gorm:"size:20;not null" json:"agama"`
	Status_perkawinan   string    `gorm:"size:20;not null;default:'Belum Kawin'" json:"status_perkawinan"` // Use constants from constants.go
	Pekerjaan           string    `gorm:"size:50" json:"pekerjaan"`
	Deskripsi_pekerjaan string    `gorm:"size:200" json:"deskripsi_pekerjaan"`
	Pendidikan_terakhir string    `gorm:"size:50" json:"pendidikan_terakhir"`
	No_hp               string    `gorm:"size:15" json:"nomor_telepon"`
	Email               string    `gorm:"size:100" json:"email"`
	Warga_negara        string    `gorm:"size:20;default:'WNI'" json:"warga_negara"`
	Created_at          time.Time `json:"dibuat_pada"`
	Updated_at          time.Time `json:"diperbarui_pada"`
}

// DataOrangTua untuk data orang tua calon pasangan
type DataOrangTua struct {
	ID                  uint       `gorm:"primaryKey" json:"id"`
	User_id             string     `gorm:"size:20;not null" json:"id_pengguna"`
	Jenis_kelamin_calon string     `gorm:"type:VARCHAR(1);not null" json:"jenis_kelamin_calon"` // L = Suami, P = Istri
	Hubungan            string     `gorm:"type:VARCHAR(10);not null" json:"hubungan"`           // enum -> varchar
	Nama_lengkap        string     `gorm:"size:100;not null" json:"nama_lengkap"`
	NIK                 string     `gorm:"size:16" json:"nik"`
	Warga_negara        string     `gorm:"size:20" json:"warga_negara"`
	Agama               string     `gorm:"size:20" json:"agama"`
	Tempat_lahir        string     `gorm:"size:50" json:"tempat_lahir"`
	Negara_asal         string     `gorm:"size:50" json:"negara_asal"`
	Pekerjaan           string     `gorm:"size:50" json:"pekerjaan"`
	No_paspor           string     `gorm:"size:20" json:"nomor_paspor"`
	Tanggal_lahir       *time.Time `json:"tanggal_lahir"`
	Pekerjaan_lain      string     `gorm:"size:50" json:"pekerjaan_lain"`
	Alamat              string     `gorm:"size:200" json:"alamat"`
	Status_keberadaan   string     `gorm:"size:20;not null;default:'Hidup'" json:"status_keberadaan"` // Use constants from constants.go
	Jenis_kelamin       string     `gorm:"size:1;not null" json:"jenis_kelamin"`                      // L = Ayah, P = Ibu
	Created_at          time.Time  `json:"dibuat_pada"`
	Updated_at          time.Time  `json:"diperbarui_pada"`
}

type PendaftaranNikah struct {
	ID                   uint       `gorm:"primaryKey" json:"id"`
	Nomor_pendaftaran    string     `gorm:"size:20;not null;unique" json:"nomor_pendaftaran"`
	Pendaftar_id         string     `gorm:"size:20;not null" json:"id_pendaftar"` // User ID yang mendaftar (suami atau istri)
	Calon_suami_id       string     `gorm:"size:20;not null" json:"id_calon_suami"`
	Calon_istri_id       string     `gorm:"size:20;not null" json:"id_calon_istri"`
	Tanggal_pendaftaran  time.Time  `gorm:"not null" json:"tanggal_pendaftaran"`
	Tanggal_nikah        time.Time  `gorm:"not null" json:"tanggal_nikah"`
	Waktu_nikah          string     `gorm:"size:10;not null" json:"waktu_nikah"` // format: HH:MM
	Tempat_nikah         string     `gorm:"size:100;not null" json:"tempat_nikah"`
	Nomor_dispensasi     string     `gorm:"size:50" json:"nomor_dispensasi"`
	Alamat_akad          string     `gorm:"size:200" json:"alamat_akad"`
	Latitude             *float64   `json:"latitude"`                                                   // Koordinat lintang untuk alamat nikah di luar KUA
	Longitude            *float64   `json:"longitude"`                                                  // Koordinat bujur untuk alamat nikah di luar KUA
	Status_pendaftaran   string     `gorm:"size:40;not null;default:'Draft'" json:"status_pendaftaran"` // Use constants from constants.go
	Status_bimbingan     string     `gorm:"size:30;not null;default:'Belum'" json:"status_bimbingan"`   // Use constants from constants.go
	Penghulu_id          *uint      `json:"id_penghulu"`                                                // ID penghulu yang ditugaskan
	Penghulu_assigned_by string     `gorm:"size:20" json:"penghulu_ditugaskan_oleh"`                    // ID kepala KUA yang assign
	Penghulu_assigned_at *time.Time `json:"penghulu_ditugaskan_pada"`                                   // Waktu assign penghulu
	Catatan              string     `gorm:"size:500" json:"catatan"`
	Disetujui_oleh       string     `gorm:"size:20" json:"disetujui_oleh"`
	Disetujui_pada       *time.Time `json:"disetujui_pada"`
	Created_at           time.Time  `json:"dibuat_pada"`
	Updated_at           time.Time  `json:"diperbarui_pada"`
}

type WaliNikah struct {
	ID                uint      `gorm:"primaryKey" json:"id"`
	Pendaftaran_id    uint      `gorm:"not null" json:"id_pendaftaran"`
	NIK               string    `gorm:"size:16;not null" json:"nik"`
	Nama_lengkap      string    `gorm:"size:100;not null" json:"nama_lengkap"`
	Hubungan_wali     string    `gorm:"size:50;not null" json:"hubungan_wali"`
	Alamat            string    `gorm:"size:200;not null" json:"alamat"`
	No_hp             string    `gorm:"size:15" json:"nomor_telepon"`
	Agama             string    `gorm:"size:20;not null" json:"agama"`
	Status_keberadaan  string    `gorm:"size:20;not null;default:'Hidup'" json:"status_keberadaan"` // Use constants from constants.go
	Created_at        time.Time `json:"dibuat_pada"`
	Updated_at        time.Time `json:"diperbarui_pada"`
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

// Role definitions - role tersimpan langsung di tabel Users
// Role yang tersedia:
// - user_biasa: User biasa untuk daftar nikah
// - penghulu: Penghulu untuk memimpin nikah
// - staff: Staff KUA untuk verifikasi
// - kepala_kua: Kepala KUA untuk approval

// StaffKUA untuk data staff KUA
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

// BimbinganPerkawinan model untuk sesi bimbingan perkawinan
type BimbinganPerkawinan struct {
	ID                uint      `gorm:"primaryKey" json:"id"`
	Tanggal_bimbingan time.Time `gorm:"not null" json:"tanggal_bimbingan"`
	Waktu_mulai       string    `gorm:"size:10;not null" json:"waktu_mulai"`
	Waktu_selesai     string    `gorm:"size:10;not null" json:"waktu_selesai"`
	Tempat_bimbingan  string    `gorm:"size:100;not null" json:"tempat_bimbingan"`
	Pembimbing        string    `gorm:"size:100;not null" json:"pembimbing"`
	Kapasitas         int       `gorm:"default:10" json:"kapasitas"`
	Status            string    `gorm:"size:20;not null;default:'Aktif'" json:"status"` // Use constants from constants.go
	Catatan           string    `gorm:"size:500" json:"catatan"`
	Created_at        time.Time `json:"dibuat_pada"`
	Updated_at        time.Time `json:"diperbarui_pada"`
}

// ==================== NEW MARRIAGE REGISTRATION FORM STRUCTS ====================

// DataFormPendaftaranNikah untuk input form pendaftaran nikah sesuai kontrak API baru
type DataFormPendaftaranNikah struct {
	JadwalDanLokasi struct {
		LokasiNikah     string `json:"weddingLocation" binding:"required"` // "Di KUA" atau "Di Luar KUA"
		AlamatNikah     string `json:"weddingAddress"`
		TanggalNikah    string `json:"weddingDate" binding:"required"`
		WaktuNikah      string `json:"weddingTime" binding:"required"`
		NomorDispensasi string `json:"dispensationNumber"`
	} `json:"scheduleAndLocation" binding:"required"`

	CalonSuami struct {
		NamaLengkap        string `json:"groomFullName" binding:"required"`
		Nik                string `json:"groomNik" binding:"required"`
		Kewarganegaraan    string `json:"groomCitizenship" binding:"required"`
		NomorPaspor        string `json:"groomPassportNumber"`
		TempatLahir        string `json:"groomPlaceOfBirth" binding:"required"`
		TanggalLahir       string `json:"groomDateOfBirth" binding:"required"`
		Status             string `json:"groomStatus" binding:"required"`
		Agama              string `json:"groomReligion" binding:"required"`
		Pendidikan         string `json:"groomEducation" binding:"required"`
		Pekerjaan          string `json:"groomOccupation" binding:"required"`
		DeskripsiPekerjaan string `json:"groomOccupationDescription"`
		NomorTelepon       string `json:"groomPhoneNumber" binding:"required"`
		Email              string `json:"groomEmail" binding:"required"`
		Alamat             string `json:"groomAddress" binding:"required"`
	} `json:"groom" binding:"required"`

	CalonIstri struct {
		NamaLengkap        string `json:"brideFullName" binding:"required"`
		Nik                string `json:"brideNik" binding:"required"`
		Kewarganegaraan    string `json:"brideCitizenship" binding:"required"`
		NomorPaspor        string `json:"bridePassportNumber"`
		TempatLahir        string `json:"bridePlaceOfBirth" binding:"required"`
		TanggalLahir       string `json:"brideDateOfBirth" binding:"required"`
		Status             string `json:"brideStatus" binding:"required"`
		Agama              string `json:"brideReligion" binding:"required"`
		Pendidikan         string `json:"brideEducation" binding:"required"`
		Pekerjaan          string `json:"brideOccupation" binding:"required"`
		DeskripsiPekerjaan string `json:"brideOccupationDescription"`
		NomorTelepon       string `json:"bridePhoneNumber" binding:"required"`
		Email              string `json:"brideEmail" binding:"required"`
		Alamat             string `json:"brideAddress" binding:"required"`
	} `json:"bride" binding:"required"`

	OrangTuaCalonSuami struct {
		Ayah struct {
			StatusKeberadaan   string `json:"groomFatherPresenceStatus" binding:"required"`
			Nama               string `json:"groomFatherName"`
			Nik                string `json:"groomFatherNik"`
			Kewarganegaraan    string `json:"groomFatherCitizenship"`
			NegaraAsal         string `json:"groomFatherCountryOfOrigin"`
			NomorPaspor        string `json:"groomFatherPassportNumber"`
			TempatLahir        string `json:"groomFatherPlaceOfBirth"`
			TanggalLahir       string `json:"groomFatherDateOfBirth"`
			Agama              string `json:"groomFatherReligion"`
			Pekerjaan          string `json:"groomFatherOccupation"`
			DeskripsiPekerjaan string `json:"groomFatherOccupationDescription"`
			Alamat             string `json:"groomFatherAddress"`
		} `json:"groomFather"`

		Ibu struct {
			StatusKeberadaan   string `json:"groomMotherPresenceStatus" binding:"required"`
			Nama               string `json:"groomMotherName"`
			Nik                string `json:"groomMotherNik"`
			Kewarganegaraan    string `json:"groomMotherCitizenship"`
			NegaraAsal         string `json:"groomMotherCountryOfOrigin"`
			NomorPaspor        string `json:"groomMotherPassportNumber"`
			TempatLahir        string `json:"groomMotherPlaceOfBirth"`
			TanggalLahir       string `json:"groomMotherDateOfBirth"`
			Agama              string `json:"groomMotherReligion"`
			Pekerjaan          string `json:"groomMotherOccupation"`
			DeskripsiPekerjaan string `json:"groomMotherOccupationDescription"`
			Alamat             string `json:"groomMotherAddress"`
		} `json:"groomMother"`
	} `json:"groomParents" binding:"required"`

	OrangTuaCalonIstri struct {
		Ayah struct {
			StatusKeberadaan   string `json:"brideFatherPresenceStatus" binding:"required"`
			Nama               string `json:"brideFatherName"`
			Nik                string `json:"brideFatherNik"`
			Kewarganegaraan    string `json:"brideFatherCitizenship"`
			NegaraAsal         string `json:"brideFatherCountryOfOrigin"`
			NomorPaspor        string `json:"brideFatherPassportNumber"`
			TempatLahir        string `json:"brideFatherPlaceOfBirth"`
			TanggalLahir       string `json:"brideFatherDateOfBirth"`
			Agama              string `json:"brideFatherReligion"`
			Pekerjaan          string `json:"brideFatherOccupation"`
			DeskripsiPekerjaan string `json:"brideFatherOccupationDescription"`
			Alamat             string `json:"brideFatherAddress"`
		} `json:"brideFather"`

		Ibu struct {
			StatusKeberadaan   string `json:"brideMotherPresenceStatus" binding:"required"`
			Nama               string `json:"brideMotherName"`
			Nik                string `json:"brideMotherNik"`
			Kewarganegaraan    string `json:"brideMotherCitizenship"`
			NegaraAsal         string `json:"brideMotherCountryOfOrigin"`
			NomorPaspor        string `json:"brideMotherPassportNumber"`
			TempatLahir        string `json:"brideMotherPlaceOfBirth"`
			TanggalLahir       string `json:"brideMotherDateOfBirth"`
			Agama              string `json:"brideMotherReligion"`
			Pekerjaan          string `json:"brideMotherOccupation"`
			DeskripsiPekerjaan string `json:"brideMotherOccupationDescription"`
			Alamat             string `json:"brideMotherAddress"`
		} `json:"brideMother"`
	} `json:"brideParents" binding:"required"`

	WaliNikah struct {
		NamaLengkapWali  string `json:"guardianFullName" binding:"required"`
		NikWali          string `json:"guardianNik" binding:"required"`
		HubunganWali     string `json:"guardianRelationship" binding:"required"`
		StatusWali       string `json:"guardianStatus" binding:"required"`
		AgamaWali        string `json:"guardianReligion" binding:"required"`
		AlamatWali       string `json:"guardianAddress" binding:"required"`
		NomorTeleponWali string `json:"guardianPhoneNumber" binding:"required"`
	} `json:"guardian" binding:"required"`
}

// PendaftaranBimbingan model untuk pendaftaran bimbingan perkawinan
type PendaftaranBimbingan struct {
	ID                      uint      `gorm:"primaryKey" json:"id"`
	Pendaftaran_nikah_id    uint      `gorm:"not null" json:"id_pendaftaran_nikah"`
	Bimbingan_perkawinan_id uint      `gorm:"not null" json:"id_bimbingan_perkawinan"`
	Calon_suami_id          string    `gorm:"size:20;not null" json:"id_calon_suami"`
	Calon_istri_id          string    `gorm:"size:20;not null" json:"id_calon_istri"`
	Status_kehadiran        string    `gorm:"size:20;not null;default:'Belum'" json:"status_kehadiran"`  // Use constants from constants.go
	Status_sertifikat       string    `gorm:"size:20;not null;default:'Belum'" json:"status_sertifikat"` // Use constants from constants.go
	No_sertifikat           string    `gorm:"size:50" json:"nomor_sertifikat"`
	Catatan                 string    `gorm:"size:500" json:"catatan"`
	Created_at              time.Time `json:"dibuat_pada"`
	Updated_at              time.Time `json:"diperbarui_pada"`
}
