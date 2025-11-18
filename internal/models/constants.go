package structs

// ==================== STATUS CONSTANTS ====================
// Menggunakan constants untuk type safety dan centralized management

// Define constants for CalonPasangan Status_perkawinan
const (
	StatusPerkawinanBelumKawin = "Belum Kawin"
	StatusPerkawinanKawin      = "Kawin"
	StatusPerkawinanCeraiMati  = "Cerai Mati"
	StatusPerkawinanCeraiHidup = "Cerai Hidup"
)


// Define constants for PendaftaranNikah Status_pendaftaran
// Flow sederhana: Draft → Disetujui → Menunggu Penugasan → Penghulu Ditugaskan → Selesai
const (
	StatusPendaftaranDraft              = "Draft"                // Catin daftar (form sederhana)
	StatusPendaftaranDisetujui          = "Disetujui"            // Staff menyetujui setelah verifikasi di belakang layar
	StatusPendaftaranMenungguPenugasan  = "Menunggu Penugasan"   // Menunggu kepala KUA menentukan penghulu
	StatusPendaftaranPenghuluDitugaskan = "Penghulu Ditugaskan"  // Kepala KUA sudah assign penghulu
	StatusPendaftaranSelesai            = "Selesai"              // Penghulu sudah melaksanakan nikah
	StatusPendaftaranDitolak            = "Ditolak"              // Ditolak oleh staff
)



// Define constants for Users Role
const (
	UserRoleUserBiasa = "user_biasa"
	UserRolePenghulu  = "penghulu"
	UserRoleStaff     = "staff"
	UserRoleKepalaKUA = "kepala_kua"
)

// Define constants for Users Status
const (
	UserStatusAktif    = "Aktif"
	UserStatusNonaktif = "Nonaktif"
	UserStatusBlokir   = "Blokir"
)

// Define constants for StaffKUA Jabatan
const (
	StaffJabatanStaff     = "Staff"
	StaffJabatanPenghulu  = "Penghulu"
	StaffJabatanKepalaKUA = "Kepala KUA"
)

// Define constants for StaffKUA Status
const (
	StaffStatusAktif    = "Aktif"
	StaffStatusNonaktif = "Nonaktif"
)

// Define constants for Penghulu Status
const (
	PenghuluStatusAktif    = "Aktif"
	PenghuluStatusNonaktif = "Nonaktif"
)

// Define constants for FeedbackPernikahan Jenis_feedback
const (
	FeedbackJenisRating  = "Rating"
	FeedbackJenisSaran   = "Saran"
	FeedbackJenisKritik  = "Kritik"
	FeedbackJenisLaporan = "Laporan"
)

// Define constants for FeedbackPernikahan Status_baca
const (
	FeedbackStatusBelumDibaca = "Belum Dibaca"
	FeedbackStatusSudahDibaca = "Sudah Dibaca"
)

// Define constants for Notifikasi Tipe
const (
	NotifikasiTipeInfo    = "Info"
	NotifikasiTipeSuccess = "Success"
	NotifikasiTipeWarning = "Warning"
	NotifikasiTipeError   = "Error"
)

// Define constants for Notifikasi Status_baca
const (
	NotifikasiStatusBelumDibaca = "Belum Dibaca"
	NotifikasiStatusSudahDibaca = "Sudah Dibaca"
)


