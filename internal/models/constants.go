package structs

// ==================== PENDAFTARAN NIKAH STATUS ====================
// Flow scheduling-only: Menunggu Penugasan -> Penghulu Ditugaskan -> Selesai
// (Draft dan Disetujui dihapus karena staff verification adalah alur SIMKAH)

const (
	StatusPendaftaranMenungguPenugasan  = "Menunggu Penugasan"
	StatusPendaftaranPenghuluDitugaskan = "Penghulu Ditugaskan"
	StatusPendaftaranSelesai            = "Selesai"
	StatusPendaftaranDitolak            = "Ditolak"
)

// ==================== TEMPAT NIKAH ====================

const (
	TempatNikahDiKUA     = "Di KUA"
	TempatNikahDiLuarKUA = "Di Luar KUA"
)

// ==================== USERS ROLE ====================

const (
	UserRoleUserBiasa = "user_biasa"
	UserRolePenghulu  = "penghulu"
	UserRoleStaff     = "staff"
	UserRoleKepalaKUA = "kepala_kua"
)

// ==================== USERS STATUS ====================

const (
	UserStatusAktif    = "Aktif"
	UserStatusNonaktif = "Nonaktif"
	UserStatusBlokir   = "Blokir"
)

// ==================== STAFF KUA ====================

const (
	StaffJabatanStaff     = "Staff"
	StaffJabatanPenghulu  = "Penghulu"
	StaffJabatanKepalaKUA = "Kepala KUA"
)

const (
	StaffStatusAktif    = "Aktif"
	StaffStatusNonaktif = "Nonaktif"
)

// ==================== PENGHULU STATUS ====================

const (
	PenghuluStatusAktif    = "Aktif"
	PenghuluStatusNonaktif = "Nonaktif"
)

// ==================== NOTIFIKASI ====================

const (
	NotifikasiTipeInfo    = "Info"
	NotifikasiTipeSuccess = "Success"
	NotifikasiTipeWarning = "Warning"
	NotifikasiTipeError   = "Error"
)

const (
	NotifikasiStatusBelumDibaca = "Belum Dibaca"
	NotifikasiStatusSudahDibaca = "Sudah Dibaca"
)
