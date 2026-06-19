package services

import (
	"fmt"
	"log"
	"time"

	structs "simnikah/internal/models"

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

	// Notifikasi untuk staff dan kepala KUA
	staffNotification := structs.Notifikasi{
		User_id:     "ALL_STAFF",
		Judul:       "Pendaftaran Nikah Baru",
		Pesan:       fmt.Sprintf("Pendaftaran nikah %s & %s pada %s pukul %s di %s", pendaftaran.Nama_suami, pendaftaran.Nama_istri, pendaftaran.Tanggal_nikah.Format("02 Januari 2006"), pendaftaran.Waktu_nikah, pendaftaran.Tempat_nikah),
		Tipe:        structs.NotifikasiTipeInfo,
		Status_baca: structs.NotifikasiStatusBelumDibaca,
		Link:        fmt.Sprintf("/simnikah/pendaftaran/%d", pendaftaranID),
	}

	// Kirim ke semua staff dan kepala KUA
	if err := ns.sendToRole(structs.UserRoleStaff, staffNotification); err != nil {
		log.Printf("Gagal mengirim notifikasi ke staff: %v", err)
	}
	if err := ns.sendToRole(structs.UserRoleKepalaKUA, staffNotification); err != nil {
		log.Printf("Gagal mengirim notifikasi ke kepala KUA: %v", err)
	}

	// Notifikasi untuk pendaftar
	if pendaftarID != "" {
		pendaftarNotification := structs.Notifikasi{
			User_id:     pendaftarID,
			Judul:       "Pendaftaran Nikah Berhasil",
			Pesan:       fmt.Sprintf("Pendaftaran nikah Anda pada %s pukul %s berhasil dibuat. Silakan tunggu penugasan penghulu dari KUA.", pendaftaran.Tanggal_nikah.Format("02 Januari 2006"), pendaftaran.Waktu_nikah),
			Tipe:        structs.NotifikasiTipeSuccess,
			Status_baca: structs.NotifikasiStatusBelumDibaca,
			Link:        fmt.Sprintf("/simnikah/pendaftaran/%d", pendaftaranID),
		}

		if err := ns.DB.Create(&pendaftarNotification).Error; err != nil {
			log.Printf("Gagal mengirim notifikasi ke pendaftar: %v", err)
		}
	}

	return nil
}

// SendStatusUpdateNotification mengirim notifikasi saat status pendaftaran berubah
func (ns *NotificationService) SendStatusUpdateNotification(pendaftaranID uint, statusLama, statusBaru string, pendaftarID string) error {
	// Ambil data pendaftaran
	var pendaftaran structs.PendaftaranNikah
	if err := ns.DB.Where("id = ?", pendaftaranID).First(&pendaftaran).Error; err != nil {
		return fmt.Errorf("gagal mengambil data pendaftaran: %v", err)
	}

	// Tentukan tipe notifikasi berdasarkan status
	var tipe string
	var pesan string

	switch statusBaru {
	case structs.StatusPendaftaranMenungguPenugasan:
		tipe = structs.NotifikasiTipeInfo
		pesan = fmt.Sprintf("Pendaftaran nikah Anda pada %s pukul %s sedang menunggu penugasan penghulu.", pendaftaran.Tanggal_nikah.Format("02 Januari 2006"), pendaftaran.Waktu_nikah)
	case structs.StatusPendaftaranPenghuluDitugaskan:
		tipe = structs.NotifikasiTipeSuccess
		pesan = fmt.Sprintf("Penghulu telah ditugaskan untuk nikah Anda pada %s pukul %s.", pendaftaran.Tanggal_nikah.Format("02 Januari 2006"), pendaftaran.Waktu_nikah)
	case structs.StatusPendaftaranSelesai:
		tipe = structs.NotifikasiTipeSuccess
		pesan = fmt.Sprintf("Proses nikah Anda pada %s telah selesai. Semoga menjadi keluarga yang sakinah, mawaddah, wa rahmah.", pendaftaran.Tanggal_nikah.Format("02 Januari 2006"))
	case structs.StatusPendaftaranDitolak:
		tipe = structs.NotifikasiTipeError
		pesan = fmt.Sprintf("Maaf, pendaftaran nikah Anda pada %s ditolak. Silakan hubungi KUA untuk informasi lebih lanjut.", pendaftaran.Tanggal_nikah.Format("02 Januari 2006"))
	default:
		tipe = structs.NotifikasiTipeInfo
		pesan = fmt.Sprintf("Status pendaftaran nikah Anda pada %s telah diubah menjadi %s.", pendaftaran.Tanggal_nikah.Format("02 Januari 2006"), statusBaru)
	}

	// Notifikasi untuk pendaftar
	if pendaftarID != "" {
		notification := structs.Notifikasi{
			User_id:     pendaftarID,
			Judul:       "Update Status Pendaftaran Nikah",
			Pesan:       pesan,
			Tipe:        tipe,
			Status_baca: structs.NotifikasiStatusBelumDibaca,
			Link:        fmt.Sprintf("/simnikah/pendaftaran/%d", pendaftaranID),
		}

		if err := ns.DB.Create(&notification).Error; err != nil {
			log.Printf("Gagal mengirim notifikasi ke pendaftar: %v", err)
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

	// Notifikasi untuk penghulu
	penghuluNotification := structs.Notifikasi{
		User_id: penghuluID,
		Judul:   "Penugasan Nikah Baru",
		Pesan: fmt.Sprintf("Anda ditugaskan untuk memimpin nikah %s & %s pada %s pukul %s di %s.",
			pendaftaran.Nama_suami,
			pendaftaran.Nama_istri,
			pendaftaran.Tanggal_nikah.Format("02 Januari 2006"),
			pendaftaran.Waktu_nikah,
			pendaftaran.Tempat_nikah),
		Tipe:        structs.NotifikasiTipeInfo,
		Status_baca: structs.NotifikasiStatusBelumDibaca,
		Link:        fmt.Sprintf("/simnikah/pendaftaran/%d", pendaftaranID),
	}

	if err := ns.DB.Create(&penghuluNotification).Error; err != nil {
		log.Printf("Gagal mengirim notifikasi ke penghulu: %v", err)
	}

	return nil
}

// SendReminderNotification mengirim notifikasi pengingat
func (ns *NotificationService) SendReminderNotification() error {
	// Pengingat untuk nikah yang akan datang (1 hari sebelum)
	tomorrow := time.Now().AddDate(0, 0, 1).Format("2006-01-02")

	var pendaftaranBesok []structs.PendaftaranNikah
	if err := ns.DB.Where("DATE(tanggal_nikah) = ? AND status_pendaftaran IN (?)", tomorrow, []string{structs.StatusPendaftaranMenungguPenugasan, structs.StatusPendaftaranPenghuluDitugaskan}).Find(&pendaftaranBesok).Error; err != nil {
		log.Printf("Gagal mengambil data pendaftaran besok: %v", err)
		return err
	}

	// Batch-fetch semua penghulu yang ditugaskan untuk menghindari N+1 query
	penghuluIDs := make([]uint, 0)
	for _, p := range pendaftaranBesok {
		if p.Penghulu_id != nil {
			penghuluIDs = append(penghuluIDs, *p.Penghulu_id)
		}
	}

	penghuluMap := map[uint]*structs.Penghulu{}
	if len(penghuluIDs) > 0 {
		var penghulus []structs.Penghulu
		if err := ns.DB.Where("id IN ?", penghuluIDs).Find(&penghulus).Error; err == nil {
			for i := range penghulus {
				p := penghulus[i]
				penghuluMap[p.ID] = &p
			}
		}
	}

	for _, pendaftaran := range pendaftaranBesok {
		// Notifikasi pengingat untuk penghulu yang ditugaskan
		if pendaftaran.Penghulu_id != nil {
			penghulu, ok := penghuluMap[*pendaftaran.Penghulu_id]
			if !ok {
				continue
			}

			reminderNotification := structs.Notifikasi{
				User_id: penghulu.User_id,
				Judul:   "Pengingat Nikah Besok",
				Pesan: fmt.Sprintf("Pengingat: Nikah %s & %s akan dilaksanakan besok (%s) pukul %s di %s.",
					pendaftaran.Nama_suami,
					pendaftaran.Nama_istri,
					pendaftaran.Tanggal_nikah.Format("02 Januari 2006"),
					pendaftaran.Waktu_nikah,
					pendaftaran.Tempat_nikah),
				Tipe:        structs.NotifikasiTipeWarning,
				Status_baca: structs.NotifikasiStatusBelumDibaca,
				Link:        fmt.Sprintf("/simnikah/pendaftaran/%d", pendaftaran.ID),
			}

			if err := ns.DB.Create(&reminderNotification).Error; err != nil {
				log.Printf("Gagal mengirim notifikasi pengingat ke penghulu: %v", err)
			}
		}

		// Notifikasi pengingat untuk staff dan kepala KUA
		staffReminder := structs.Notifikasi{
			Judul:   "Pengingat Nikah Besok",
			Pesan:   fmt.Sprintf("Pengingat: Nikah %s & %s akan dilaksanakan besok (%s) pukul %s di %s.", pendaftaran.Nama_suami, pendaftaran.Nama_istri, pendaftaran.Tanggal_nikah.Format("02 Januari 2006"), pendaftaran.Waktu_nikah, pendaftaran.Tempat_nikah),
			Tipe:    structs.NotifikasiTipeWarning,
			Link:    fmt.Sprintf("/simnikah/pendaftaran/%d", pendaftaran.ID),
		}

		if err := ns.sendToRole(structs.UserRoleStaff, staffReminder); err != nil {
			log.Printf("Gagal mengirim notifikasi pengingat ke staff: %v", err)
		}
		if err := ns.sendToRole(structs.UserRoleKepalaKUA, staffReminder); err != nil {
			log.Printf("Gagal mengirim notifikasi pengingat ke kepala KUA: %v", err)
		}
	}

	return nil
}

// SendNotificationToRole public method untuk mengirim notifikasi ke semua user dengan role tertentu
func (ns *NotificationService) SendNotificationToRole(role, judul, pesan, tipe, link string) error {
	notification := structs.Notifikasi{
		Judul:       judul,
		Pesan:       pesan,
		Tipe:        tipe,
		Status_baca: structs.NotifikasiStatusBelumDibaca,
		Link:        link,
	}

	return ns.sendToRole(role, notification)
}

// sendToRole helper function untuk mengirim notifikasi ke semua user dengan role tertentu
func (ns *NotificationService) sendToRole(role string, notification structs.Notifikasi) error {
	// Ambil semua user dengan role tersebut
	var users []structs.Users
	if err := ns.DB.Where("role = ? AND status = ?", role, structs.UserStatusAktif).Find(&users).Error; err != nil {
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
		Status_baca: structs.NotifikasiStatusBelumDibaca,
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
			Status_baca: structs.NotifikasiStatusBelumDibaca,
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
		Tipe:        structs.NotifikasiTipeSuccess,
		Status_baca: structs.NotifikasiStatusBelumDibaca,
		Link:        "/simnikah/dashboard",
	}

	// Simpan notifikasi untuk staff baru
	if err := ns.DB.Create(&staffNotification).Error; err != nil {
		return fmt.Errorf("gagal mengirim notifikasi ke staff baru: %v", err)
	}

	// Notifikasi untuk kepala KUA bahwa staff baru telah dibuat
	kepalaKuaNotification := structs.Notifikasi{
		User_id:     "ALL_KEPALA_KUA",
		Judul:       "Staff Baru Dibuat",
		Pesan:       fmt.Sprintf("Staff baru %s dengan jabatan %s telah berhasil dibuat dan dapat login ke sistem.", staffNama, jabatan),
		Tipe:        structs.NotifikasiTipeInfo,
		Status_baca: structs.NotifikasiStatusBelumDibaca,
		Link:        "/simnikah/internal/handlers/staff",
	}

	// Kirim ke semua kepala KUA
	if err := ns.sendToRole(structs.UserRoleKepalaKUA, kepalaKuaNotification); err != nil {
		log.Printf("Gagal mengirim notifikasi ke kepala KUA: %v", err)
	}

	return nil
}
