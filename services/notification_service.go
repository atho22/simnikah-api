package services

import (
	"fmt"
	"log"
	"time"

	"simnikah/structs"

	"gorm.io/gorm"
)

// NotificationService untuk mengelola notifikasi otomatis
type NotificationService struct {
	DB *gorm.DB
}

// NewNotificationService membuat instance baru dari NotificationService
func NewNotificationService(db *gorm.DB) *NotificationService {
	return &NotificationService{DB: db}
}

// SendPendaftaranNotification mengirim notifikasi saat ada pendaftaran baru
func (ns *NotificationService) SendPendaftaranNotification(pendaftaranID uint, pendaftarID string) error {
	// Ambil data pendaftaran
	var pendaftaran structs.PendaftaranNikah
	if err := ns.DB.Where("id = ?", pendaftaranID).First(&pendaftaran).Error; err != nil {
		return fmt.Errorf("gagal mengambil data pendaftaran: %v", err)
	}

	// Ambil data calon suami dan istri
	var calonSuami, calonIstri structs.CalonPasangan
	if err := ns.DB.Where("id = ?", pendaftaran.Calon_suami_id).First(&calonSuami).Error; err != nil {
		return fmt.Errorf("gagal mengambil data calon suami: %v", err)
	}
	if err := ns.DB.Where("id = ?", pendaftaran.Calon_istri_id).First(&calonIstri).Error; err != nil {
		return fmt.Errorf("gagal mengambil data calon istri: %v", err)
	}

	// Notifikasi untuk staff dan kepala KUA
	staffNotification := structs.Notifikasi{
		User_id:     "ALL_STAFF", // Akan dipecah menjadi notifikasi individual
		Judul:       "Pendaftaran Nikah Baru",
		Pesan:       fmt.Sprintf("Pendaftaran nikah baru dari %s dan %s dengan nomor pendaftaran %s", calonSuami.Nama_lengkap, calonIstri.Nama_lengkap, pendaftaran.Nomor_pendaftaran),
		Tipe:        "Info",
		Status_baca: "Belum Dibaca",
		Link:        fmt.Sprintf("/simnikah/pendaftaran/%d", pendaftaranID),
	}

	// Kirim ke semua staff dan kepala KUA
	if err := ns.sendToRole("staff", staffNotification); err != nil {
		log.Printf("Gagal mengirim notifikasi ke staff: %v", err)
	}
	if err := ns.sendToRole("kepala_kua", staffNotification); err != nil {
		log.Printf("Gagal mengirim notifikasi ke kepala KUA: %v", err)
	}

	// Notifikasi untuk pendaftar (calon suami)
	pendaftarNotification := structs.Notifikasi{
		User_id:     pendaftarID,
		Judul:       "Pendaftaran Nikah Berhasil",
		Pesan:       fmt.Sprintf("Pendaftaran nikah Anda dengan %s berhasil dibuat dengan nomor pendaftaran %s. Silakan tunggu proses verifikasi dari KUA.", calonIstri.Nama_lengkap, pendaftaran.Nomor_pendaftaran),
		Tipe:        "Success",
		Status_baca: "Belum Dibaca",
		Link:        fmt.Sprintf("/simnikah/pendaftaran/%d", pendaftaranID),
	}

	if err := ns.DB.Create(&pendaftarNotification).Error; err != nil {
		log.Printf("Gagal mengirim notifikasi ke pendaftar: %v", err)
	}

	return nil
}

// SendStatusUpdateNotification mengirim notifikasi saat status pendaftaran berubah
func (ns *NotificationService) SendStatusUpdateNotification(pendaftaranID uint, statusLama, statusBaru string, updatedBy string) error {
	// Ambil data pendaftaran
	var pendaftaran structs.PendaftaranNikah
	if err := ns.DB.Where("id = ?", pendaftaranID).First(&pendaftaran).Error; err != nil {
		return fmt.Errorf("gagal mengambil data pendaftaran: %v", err)
	}

	// Ambil data calon suami dan istri
	var calonSuami, calonIstri structs.CalonPasangan
	if err := ns.DB.Where("user_id = ?", pendaftaran.Calon_suami_id).First(&calonSuami).Error; err != nil {
		return fmt.Errorf("gagal mengambil data calon suami: %v", err)
	}
	if err := ns.DB.Where("user_id = ?", pendaftaran.Calon_istri_id).First(&calonIstri).Error; err != nil {
		return fmt.Errorf("gagal mengambil data calon istri: %v", err)
	}

	// Tentukan tipe notifikasi berdasarkan status
	var tipe string
	var pesan string

	switch statusBaru {
	case "Menunggu Verifikasi":
		tipe = "Info"
		pesan = fmt.Sprintf("Pendaftaran nikah Anda dengan %s sedang menunggu verifikasi dari KUA.", calonIstri.Nama_lengkap)
	case "Disetujui":
		tipe = "Success"
		pesan = fmt.Sprintf("Selamat! Pendaftaran nikah Anda dengan %s telah disetujui oleh KUA.", calonIstri.Nama_lengkap)
	case "Surat Diterbitkan":
		tipe = "Success"
		pesan = fmt.Sprintf("Surat undangan nikah Anda dengan %s telah diterbitkan. Silakan cek detail surat undangan.", calonIstri.Nama_lengkap)
	case "Sudah Bimbingan":
		tipe = "Info"
		pesan = fmt.Sprintf("Anda telah menyelesaikan bimbingan perkawinan. Pendaftaran nikah Anda dengan %s siap untuk dilaksanakan.", calonIstri.Nama_lengkap)
	case "Selesai":
		tipe = "Success"
		pesan = fmt.Sprintf("Selamat! Proses nikah Anda dengan %s telah selesai. Semoga menjadi keluarga yang sakinah, mawaddah, wa rahmah.", calonIstri.Nama_lengkap)
	case "Ditolak":
		tipe = "Error"
		pesan = fmt.Sprintf("Maaf, pendaftaran nikah Anda dengan %s ditolak. Silakan hubungi KUA untuk informasi lebih lanjut.", calonIstri.Nama_lengkap)
	default:
		tipe = "Info"
		pesan = fmt.Sprintf("Status pendaftaran nikah Anda dengan %s telah diubah menjadi %s.", calonIstri.Nama_lengkap, statusBaru)
	}

	// Notifikasi untuk calon suami
	suamiNotification := structs.Notifikasi{
		User_id:     calonSuami.User_id,
		Judul:       "Update Status Pendaftaran Nikah",
		Pesan:       pesan,
		Tipe:        tipe,
		Status_baca: "Belum Dibaca",
		Link:        fmt.Sprintf("/simnikah/pendaftaran/%d", pendaftaranID),
	}

	if err := ns.DB.Create(&suamiNotification).Error; err != nil {
		log.Printf("Gagal mengirim notifikasi ke calon suami: %v", err)
	}

	// Notifikasi untuk calon istri (jika berbeda dengan suami)
	if pendaftaran.Calon_istri_id != pendaftaran.Calon_suami_id {
		istriNotification := structs.Notifikasi{
			User_id:     calonIstri.User_id,
			Judul:       "Update Status Pendaftaran Nikah",
			Pesan:       pesan,
			Tipe:        tipe,
			Status_baca: "Belum Dibaca",
			Link:        fmt.Sprintf("/simnikah/pendaftaran/%d", pendaftaranID),
		}

		if err := ns.DB.Create(&istriNotification).Error; err != nil {
			log.Printf("Gagal mengirim notifikasi ke calon istri: %v", err)
		}
	}

	return nil
}

// SendBimbinganNotification mengirim notifikasi terkait bimbingan perkawinan
func (ns *NotificationService) SendBimbinganNotification(bimbinganID uint, action string) error {
	// Ambil data bimbingan
	var bimbingan structs.BimbinganPerkawinan
	if err := ns.DB.Where("id = ?", bimbinganID).First(&bimbingan).Error; err != nil {
		return fmt.Errorf("gagal mengambil data bimbingan: %v", err)
	}

	// Ambil semua pendaftar bimbingan
	var pendaftarBimbingan []structs.PendaftaranBimbingan
	if err := ns.DB.Where("bimbingan_perkawinan_id = ?", bimbinganID).Find(&pendaftarBimbingan).Error; err != nil {
		return fmt.Errorf("gagal mengambil data pendaftar bimbingan: %v", err)
	}

	var tipe string
	var judul, pesan string

	switch action {
	case "created":
		tipe = "Info"
		judul = "Bimbingan Perkawinan Baru"
		pesan = fmt.Sprintf("Bimbingan perkawinan baru telah dijadwalkan pada %s pukul %s - %s di %s.",
			bimbingan.Tanggal_bimbingan.Format("02 Januari 2006"),
			bimbingan.Waktu_mulai,
			bimbingan.Waktu_selesai,
			bimbingan.Tempat_bimbingan)
	case "updated":
		tipe = "Warning"
		judul = "Update Jadwal Bimbingan Perkawinan"
		pesan = fmt.Sprintf("Jadwal bimbingan perkawinan telah diubah. Tanggal: %s, Waktu: %s - %s, Tempat: %s.",
			bimbingan.Tanggal_bimbingan.Format("02 Januari 2006"),
			bimbingan.Waktu_mulai,
			bimbingan.Waktu_selesai,
			bimbingan.Tempat_bimbingan)
	case "cancelled":
		tipe = "Error"
		judul = "Bimbingan Perkawinan Dibatalkan"
		pesan = fmt.Sprintf("Bimbingan perkawinan pada %s telah dibatalkan. Silakan hubungi KUA untuk informasi lebih lanjut.",
			bimbingan.Tanggal_bimbingan.Format("02 Januari 2006"))
	}

	// Kirim notifikasi ke semua peserta
	for _, peserta := range pendaftarBimbingan {
		// Ambil data calon pasangan untuk mendapatkan user_id
		var calonSuami, calonIstri structs.CalonPasangan
		ns.DB.First(&calonSuami, peserta.Calon_suami_id)
		ns.DB.First(&calonIstri, peserta.Calon_istri_id)

		notification := structs.Notifikasi{
			User_id:     calonSuami.User_id,
			Judul:       judul,
			Pesan:       pesan,
			Tipe:        tipe,
			Status_baca: "Belum Dibaca",
			Link:        fmt.Sprintf("/simnikah/bimbingan/%d", bimbinganID),
		}

		if err := ns.DB.Create(&notification).Error; err != nil {
			log.Printf("Gagal mengirim notifikasi ke calon suami %s: %v", peserta.Calon_suami_id, err)
		}

		// Notifikasi untuk calon istri juga
		if peserta.Calon_istri_id != peserta.Calon_suami_id {
			notification.User_id = calonIstri.User_id
			if err := ns.DB.Create(&notification).Error; err != nil {
				log.Printf("Gagal mengirim notifikasi ke calon istri %s: %v", peserta.Calon_istri_id, err)
			}
		}
	}

	return nil
}

// SendPenghuluAssignmentNotification mengirim notifikasi saat penghulu ditugaskan
func (ns *NotificationService) SendPenghuluAssignmentNotification(pendaftaranID uint, penghuluID string) error {
	// Ambil data pendaftaran
	var pendaftaran structs.PendaftaranNikah
	if err := ns.DB.Where("id = ?", pendaftaranID).First(&pendaftaran).Error; err != nil {
		return fmt.Errorf("gagal mengambil data pendaftaran: %v", err)
	}

	// Ambil data penghulu
	var penghulu structs.Penghulu
	if err := ns.DB.Where("user_id = ?", penghuluID).First(&penghulu).Error; err != nil {
		return fmt.Errorf("gagal mengambil data penghulu: %v", err)
	}

	// Ambil data calon suami dan istri
	var calonSuami, calonIstri structs.CalonPasangan
	if err := ns.DB.Where("user_id = ?", pendaftaran.Calon_suami_id).First(&calonSuami).Error; err != nil {
		return fmt.Errorf("gagal mengambil data calon suami: %v", err)
	}
	if err := ns.DB.Where("user_id = ?", pendaftaran.Calon_istri_id).First(&calonIstri).Error; err != nil {
		return fmt.Errorf("gagal mengambil data calon istri: %v", err)
	}

	// Notifikasi untuk penghulu
	penghuluNotification := structs.Notifikasi{
		User_id: penghuluID,
		Judul:   "Penugasan Nikah Baru",
		Pesan: fmt.Sprintf("Anda ditugaskan untuk memimpin nikah %s dan %s pada %s pukul %s di %s.",
			calonSuami.Nama_lengkap,
			calonIstri.Nama_lengkap,
			pendaftaran.Tanggal_nikah.Format("02 Januari 2006"),
			pendaftaran.Waktu_nikah,
			pendaftaran.Tempat_nikah),
		Tipe:        "Info",
		Status_baca: "Belum Dibaca",
		Link:        fmt.Sprintf("/simnikah/pendaftaran/%d", pendaftaranID),
	}

	if err := ns.DB.Create(&penghuluNotification).Error; err != nil {
		log.Printf("Gagal mengirim notifikasi ke penghulu: %v", err)
	}

	// Notifikasi untuk calon pasangan
	calonNotification := structs.Notifikasi{
		User_id:     calonSuami.User_id,
		Judul:       "Penghulu Ditugaskan",
		Pesan:       fmt.Sprintf("Penghulu %s telah ditugaskan untuk memimpin nikah Anda dengan %s.", penghulu.Nama_lengkap, calonIstri.Nama_lengkap),
		Tipe:        "Success",
		Status_baca: "Belum Dibaca",
		Link:        fmt.Sprintf("/simnikah/pendaftaran/%d", pendaftaranID),
	}

	if err := ns.DB.Create(&calonNotification).Error; err != nil {
		log.Printf("Gagal mengirim notifikasi ke calon suami: %v", err)
	}

	// Notifikasi untuk calon istri juga
	if pendaftaran.Calon_istri_id != pendaftaran.Calon_suami_id {
		calonNotification.User_id = calonIstri.User_id
		if err := ns.DB.Create(&calonNotification).Error; err != nil {
			log.Printf("Gagal mengirim notifikasi ke calon istri: %v", err)
		}
	}

	return nil
}

// SendReminderNotification mengirim notifikasi pengingat
func (ns *NotificationService) SendReminderNotification() error {
	// Pengingat untuk nikah yang akan datang (1 hari sebelum)
	tomorrow := time.Now().AddDate(0, 0, 1).Format("2006-01-02")

	var pendaftaranBesok []structs.PendaftaranNikah
	if err := ns.DB.Where("DATE(tanggal_nikah) = ? AND status_pendaftaran IN (?)", tomorrow, []string{"Disetujui", "Surat Diterbitkan", "Sudah Bimbingan"}).Find(&pendaftaranBesok).Error; err != nil {
		log.Printf("Gagal mengambil data pendaftaran besok: %v", err)
		return err
	}

	for _, pendaftaran := range pendaftaranBesok {
		// Ambil data calon suami dan istri
		var calonSuami, calonIstri structs.CalonPasangan
		if err := ns.DB.Where("id = ?", pendaftaran.Calon_suami_id).First(&calonSuami).Error; err != nil {
			continue
		}
		if err := ns.DB.Where("id = ?", pendaftaran.Calon_istri_id).First(&calonIstri).Error; err != nil {
			continue
		}

		// Notifikasi pengingat untuk calon suami
		reminderNotification := structs.Notifikasi{
			User_id: calonSuami.User_id,
			Judul:   "Pengingat Nikah Besok",
			Pesan: fmt.Sprintf("Pengingat: Nikah Anda dengan %s akan dilaksanakan besok (%s) pukul %s di %s. Pastikan semua persiapan sudah siap!",
				calonIstri.Nama_lengkap,
				pendaftaran.Tanggal_nikah.Format("02 Januari 2006"),
				pendaftaran.Waktu_nikah,
				pendaftaran.Tempat_nikah),
			Tipe:        "Warning",
			Status_baca: "Belum Dibaca",
			Link:        fmt.Sprintf("/simnikah/pendaftaran/%d", pendaftaran.ID),
		}

		if err := ns.DB.Create(&reminderNotification).Error; err != nil {
			log.Printf("Gagal mengirim notifikasi pengingat ke calon suami: %v", err)
		}

		// Notifikasi pengingat untuk calon istri juga
		if pendaftaran.Calon_istri_id != pendaftaran.Calon_suami_id {
			reminderNotification.User_id = calonIstri.User_id
			if err := ns.DB.Create(&reminderNotification).Error; err != nil {
				log.Printf("Gagal mengirim notifikasi pengingat ke calon istri: %v", err)
			}
		}
	}

	// Pengingat untuk bimbingan yang akan datang (1 hari sebelum)
	var bimbinganBesok []structs.BimbinganPerkawinan
	if err := ns.DB.Where("DATE(tanggal_bimbingan) = ? AND status = ?", tomorrow, "Aktif").Find(&bimbinganBesok).Error; err != nil {
		log.Printf("Gagal mengambil data bimbingan besok: %v", err)
		return err
	}

	for _, bimbingan := range bimbinganBesok {
		// Ambil semua peserta bimbingan
		var pesertaBimbingan []structs.PendaftaranBimbingan
		if err := ns.DB.Where("bimbingan_perkawinan_id = ?", bimbingan.ID).Find(&pesertaBimbingan).Error; err != nil {
			continue
		}

		for _, peserta := range pesertaBimbingan {
			// Notifikasi pengingat untuk calon suami
			reminderNotification := structs.Notifikasi{
				User_id: peserta.Calon_suami_id,
				Judul:   "Pengingat Bimbingan Perkawinan Besok",
				Pesan: fmt.Sprintf("Pengingat: Bimbingan perkawinan akan dilaksanakan besok (%s) pukul %s - %s di %s. Pastikan Anda hadir tepat waktu!",
					bimbingan.Tanggal_bimbingan.Format("02 Januari 2006"),
					bimbingan.Waktu_mulai,
					bimbingan.Waktu_selesai,
					bimbingan.Tempat_bimbingan),
				Tipe:        "Warning",
				Status_baca: "Belum Dibaca",
				Link:        fmt.Sprintf("/simnikah/bimbingan/%d", bimbingan.ID),
			}

			if err := ns.DB.Create(&reminderNotification).Error; err != nil {
				log.Printf("Gagal mengirim notifikasi pengingat bimbingan ke calon suami: %v", err)
			}

			// Notifikasi pengingat untuk calon istri juga
			if peserta.Calon_istri_id != peserta.Calon_suami_id {
				reminderNotification.User_id = peserta.Calon_istri_id
				if err := ns.DB.Create(&reminderNotification).Error; err != nil {
					log.Printf("Gagal mengirim notifikasi pengingat bimbingan ke calon istri: %v", err)
				}
			}
		}
	}

	return nil
}

// sendToRole helper function untuk mengirim notifikasi ke semua user dengan role tertentu
func (ns *NotificationService) sendToRole(role string, notification structs.Notifikasi) error {
	// Ambil semua user dengan role tersebut
	var users []structs.Users
	if err := ns.DB.Where("role = ? AND status = ?", role, "Aktif").Find(&users).Error; err != nil {
		return err
	}

	// Buat notifikasi untuk setiap user
	var notifications []structs.Notifikasi
	for _, user := range users {
		notif := notification
		notif.User_id = user.User_id
		notifications = append(notifications, notif)
	}

	// Simpan semua notifikasi
	return ns.DB.Create(&notifications).Error
}

// SendSystemNotification mengirim notifikasi sistem
func (ns *NotificationService) SendSystemNotification(userID, judul, pesan, tipe, link string) error {
	notification := structs.Notifikasi{
		User_id:     userID,
		Judul:       judul,
		Pesan:       pesan,
		Tipe:        tipe,
		Status_baca: "Belum Dibaca",
		Link:        link,
	}

	return ns.DB.Create(&notification).Error
}

// SendBulkNotification mengirim notifikasi ke multiple users
func (ns *NotificationService) SendBulkNotification(userIDs []string, judul, pesan, tipe, link string) error {
	var notifications []structs.Notifikasi

	for _, userID := range userIDs {
		notification := structs.Notifikasi{
			User_id:     userID,
			Judul:       judul,
			Pesan:       pesan,
			Tipe:        tipe,
			Status_baca: "Belum Dibaca",
			Link:        link,
		}
		notifications = append(notifications, notification)
	}

	return ns.DB.Create(&notifications).Error
}

// SendStaffCreatedNotification mengirim notifikasi saat staff baru dibuat
func (ns *NotificationService) SendStaffCreatedNotification(staffUserID, staffNama, jabatan string) error {
	// Notifikasi untuk staff yang baru dibuat
	staffNotification := structs.Notifikasi{
		User_id:     staffUserID,
		Judul:       "Selamat Datang di SimNikah",
		Pesan:       fmt.Sprintf("Selamat datang %s! Akun Anda sebagai %s telah berhasil dibuat. Silakan login untuk mengakses sistem.", staffNama, jabatan),
		Tipe:        "Success",
		Status_baca: "Belum Dibaca",
		Link:        "/simnikah/dashboard",
	}

	// Simpan notifikasi untuk staff baru
	if err := ns.DB.Create(&staffNotification).Error; err != nil {
		return fmt.Errorf("gagal mengirim notifikasi ke staff baru: %v", err)
	}

	// Notifikasi untuk kepala KUA bahwa staff baru telah dibuat
	kepalaKuaNotification := structs.Notifikasi{
		User_id:     "ALL_KEPALA_KUA", // Akan dipecah menjadi notifikasi individual
		Judul:       "Staff Baru Dibuat",
		Pesan:       fmt.Sprintf("Staff baru %s dengan jabatan %s telah berhasil dibuat dan dapat login ke sistem.", staffNama, jabatan),
		Tipe:        "Info",
		Status_baca: "Belum Dibaca",
		Link:        "/simnikah/staff",
	}

	// Kirim ke semua kepala KUA
	if err := ns.sendToRole("kepala_kua", kepalaKuaNotification); err != nil {
		log.Printf("Gagal mengirim notifikasi ke kepala KUA: %v", err)
	}

	return nil
}
