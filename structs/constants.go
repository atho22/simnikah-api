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

// Define constants for DataOrangTua Hubungan
const (
	HubunganAyah = "Ayah"
	HubunganIbu  = "Ibu"
)

// Define constants for DataOrangTua Status_keberadaan
const (
	StatusKeberadaanHidup     = "Hidup"
	StatusKeberadaanMeninggal = "Meninggal"
)

// Define constants for PendaftaranNikah Status_pendaftaran
const (
	StatusPendaftaranDraft                      = "Draft"
	StatusPendaftaranMenungguVerifikasi         = "Menunggu Verifikasi"
	StatusPendaftaranMenungguPengumpulanBerkas  = "Menunggu Pengumpulan Berkas"
	StatusPendaftaranBerkasDiterima             = "Berkas Diterima"
	StatusPendaftaranMenungguPenugasan          = "Menunggu Penugasan"
	StatusPendaftaranPenghuluDitugaskan         = "Penghulu Ditugaskan"
	StatusPendaftaranMenungguVerifikasiPenghulu = "Menunggu Verifikasi Penghulu"
	StatusPendaftaranMenungguBimbingan          = "Menunggu Bimbingan"
	StatusPendaftaranSudahBimbingan             = "Sudah Bimbingan"
	StatusPendaftaranSelesai                    = "Selesai"
	StatusPendaftaranDitolak                    = "Ditolak"
)

// Define constants for PendaftaranNikah Status_bimbingan
const (
	StatusBimbinganBelum                 = "Belum"
	StatusBimbinganSudah                 = "Sudah"
	StatusBimbinganSertifikatDiterbitkan = "Sertifikat Diterbitkan"
)

// Define constants for WaliNikah Status_keberadaan
const (
	WaliStatusKeberadaanHidup     = "Hidup"
	WaliStatusKeberadaanMeninggal = "Meninggal"
)

// Define constants for WaliNikah Status_kehadiran
const (
	WaliStatusKehadiranBelumDiketahui = "Belum Diketahui"
	WaliStatusKehadiranHadir          = "Hadir"
	WaliStatusKehadiranTidakHadir     = "Tidak Hadir"
)

// Define constants for WaliNikah Hubungan (Urutan Wali Nasab sesuai Syariat Islam)
const (
	WaliHubunganAyahKandung            = "Ayah Kandung"
	WaliHubunganKakek                  = "Kakek" // Ayah dari Ayah
	WaliHubunganSaudaraLakiLakiKandung = "Saudara Laki-Laki Kandung"
	WaliHubunganSaudaraLakiLakiSeayah  = "Saudara Laki-Laki Seayah"
	WaliHubunganKeponakanLakiLaki      = "Keponakan Laki-Laki" // Anak dari saudara laki-laki kandung
	WaliHubunganPamanKandung           = "Paman Kandung"       // Saudara ayah kandung
	WaliHubunganPamanSeayah            = "Paman Seayah"        // Saudara ayah seayah
	WaliHubunganSepupuLakiLaki         = "Sepupu Laki-Laki"    // Anak dari paman kandung
	WaliHubunganWaliHakim              = "Wali Hakim"          // Jika tidak ada wali nasab
	WaliHubunganLainnya                = "Lainnya"
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

// Define constants for BimbinganPerkawinan Status
const (
	BimbinganStatusAktif      = "Aktif"
	BimbinganStatusSelesai    = "Selesai"
	BimbinganStatusDibatalkan = "Dibatalkan"
)

// Define constants for PendaftaranBimbingan Status_kehadiran
const (
	PendaftaranBimbinganKehadiranBelum      = "Belum"
	PendaftaranBimbinganKehadiranHadir      = "Hadir"
	PendaftaranBimbinganKehadiranTidakHadir = "Tidak Hadir"
)

// Define constants for PendaftaranBimbingan Status_sertifikat
const (
	PendaftaranBimbinganSertifikatBelum = "Belum"
	PendaftaranBimbinganSertifikatSudah = "Sudah"
)

// ==================== HELPER FUNCTIONS ====================

// GetUrutanWaliNasab - Mengembalikan urutan wali nasab sesuai syariat Islam
// Urutan dari yang paling berhak menjadi wali
func GetUrutanWaliNasab() []string {
	return []string{
		WaliHubunganAyahKandung,
		WaliHubunganKakek,
		WaliHubunganSaudaraLakiLakiKandung,
		WaliHubunganSaudaraLakiLakiSeayah,
		WaliHubunganKeponakanLakiLaki,
		WaliHubunganPamanKandung,
		WaliHubunganPamanSeayah,
		WaliHubunganSepupuLakiLaki,
		WaliHubunganWaliHakim,
	}
}

// IsValidWaliNikah - Validasi apakah hubungan wali valid
func IsValidWaliNikah(hubunganWali string, statusAyahCalonIstri string) bool {
	validWali := map[string]bool{
		WaliHubunganAyahKandung:            true,
		WaliHubunganKakek:                  true,
		WaliHubunganSaudaraLakiLakiKandung: true,
		WaliHubunganSaudaraLakiLakiSeayah:  true,
		WaliHubunganKeponakanLakiLaki:      true,
		WaliHubunganPamanKandung:           true,
		WaliHubunganPamanSeayah:            true,
		WaliHubunganSepupuLakiLaki:         true,
		WaliHubunganWaliHakim:              true,
		WaliHubunganLainnya:                true,
	}

	if !validWali[hubunganWali] {
		return false
	}

	// Validasi khusus: Jika ayah masih hidup, maka wali HARUS ayah kandung
	if statusAyahCalonIstri == StatusKeberadaanHidup && hubunganWali != WaliHubunganAyahKandung {
		return false
	}

	// Validasi khusus: Jika ayah meninggal/tidak ada, maka TIDAK BOLEH ayah kandung
	if statusAyahCalonIstri == StatusKeberadaanMeninggal && hubunganWali == WaliHubunganAyahKandung {
		return false
	}

	return true
}

// GetPesanValidasiWali - Memberikan pesan validasi untuk wali nikah
func GetPesanValidasiWali(statusAyahCalonIstri string) string {
	if statusAyahCalonIstri == StatusKeberadaanHidup {
		return "Jika ayah kandung masih hidup, maka wali nikah harus Ayah Kandung sesuai syariat Islam"
	}
	return "Jika ayah kandung meninggal/tidak ada, wali nikah berpindah ke nasab berikutnya: Kakek, Saudara Laki-Laki Kandung, Paman, atau Wali Hakim"
}
