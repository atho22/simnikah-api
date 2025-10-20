package examples

import (
	"log"
	"simnikah/services"

	"gorm.io/gorm"
)

// NotificationIntegrationExample menunjukkan cara mengintegrasikan notification service
// dengan handler yang sudah ada
type NotificationIntegrationExample struct {
	DB                  *gorm.DB
	NotificationService *services.NotificationService
}

// NewNotificationIntegrationExample membuat instance baru
func NewNotificationIntegrationExample(db *gorm.DB) *NotificationIntegrationExample {
	return &NotificationIntegrationExample{
		DB:                  db,
		NotificationService: services.NewNotificationService(db),
	}
}

// ExampleCreatePendaftaranWithNotification contoh integrasi saat membuat pendaftaran
func (ni *NotificationIntegrationExample) ExampleCreatePendaftaranWithNotification(pendaftaranID uint, pendaftarID string) error {
	// 1. Buat pendaftaran (kode ini sudah ada di handler yang ada)
	// ... kode untuk membuat pendaftaran ...

	// 2. Kirim notifikasi otomatis
	err := ni.NotificationService.SendPendaftaranNotification(pendaftaranID, pendaftarID)
	if err != nil {
		log.Printf("Gagal mengirim notifikasi pendaftaran: %v", err)
		// Jangan return error, karena pendaftaran sudah berhasil dibuat
	}

	return nil
}

// ExampleUpdateStatusWithNotification contoh integrasi saat update status
func (ni *NotificationIntegrationExample) ExampleUpdateStatusWithNotification(pendaftaranID uint, statusLama, statusBaru, updatedBy string) error {
	// 1. Update status pendaftaran (kode ini sudah ada di handler yang ada)
	// ... kode untuk update status ...

	// 2. Kirim notifikasi otomatis
	err := ni.NotificationService.SendStatusUpdateNotification(pendaftaranID, statusLama, statusBaru, updatedBy)
	if err != nil {
		log.Printf("Gagal mengirim notifikasi update status: %v", err)
		// Jangan return error, karena update status sudah berhasil
	}

	return nil
}

// ExampleAssignPenghuluWithNotification contoh integrasi saat assign penghulu
func (ni *NotificationIntegrationExample) ExampleAssignPenghuluWithNotification(pendaftaranID uint, penghuluID string) error {
	// 1. Assign penghulu (kode ini sudah ada di handler yang ada)
	// ... kode untuk assign penghulu ...

	// 2. Kirim notifikasi otomatis
	err := ni.NotificationService.SendPenghuluAssignmentNotification(pendaftaranID, penghuluID)
	if err != nil {
		log.Printf("Gagal mengirim notifikasi penugasan penghulu: %v", err)
		// Jangan return error, karena assign penghulu sudah berhasil
	}

	return nil
}

// ExampleCreateBimbinganWithNotification contoh integrasi saat membuat bimbingan
func (ni *NotificationIntegrationExample) ExampleCreateBimbinganWithNotification(bimbinganID uint) error {
	// 1. Buat bimbingan (kode ini sudah ada di handler yang ada)
	// ... kode untuk membuat bimbingan ...

	// 2. Kirim notifikasi otomatis
	err := ni.NotificationService.SendBimbinganNotification(bimbinganID, "created")
	if err != nil {
		log.Printf("Gagal mengirim notifikasi bimbingan: %v", err)
		// Jangan return error, karena bimbingan sudah berhasil dibuat
	}

	return nil
}

// ExampleSendCustomNotification contoh mengirim notifikasi custom
func (ni *NotificationIntegrationExample) ExampleSendCustomNotification(userID, judul, pesan, tipe, link string) error {
	err := ni.NotificationService.SendSystemNotification(userID, judul, pesan, tipe, link)
	if err != nil {
		log.Printf("Gagal mengirim notifikasi custom: %v", err)
		return err
	}

	return nil
}

// ExampleSendBulkNotification contoh mengirim notifikasi ke multiple users
func (ni *NotificationIntegrationExample) ExampleSendBulkNotification(userIDs []string, judul, pesan, tipe, link string) error {
	err := ni.NotificationService.SendBulkNotification(userIDs, judul, pesan, tipe, link)
	if err != nil {
		log.Printf("Gagal mengirim notifikasi bulk: %v", err)
		return err
	}

	return nil
}

// ExampleSendReminderNotification contoh mengirim notifikasi pengingat
func (ni *NotificationIntegrationExample) ExampleSendReminderNotification() error {
	err := ni.NotificationService.SendReminderNotification()
	if err != nil {
		log.Printf("Gagal mengirim notifikasi pengingat: %v", err)
		return err
	}

	return nil
}

// ExampleIntegrationInHandler menunjukkan cara mengintegrasikan di handler yang sudah ada
func ExampleIntegrationInHandler() {
	// Contoh integrasi di handler CreatePendaftaranNikah
	/*
		func (idb *InDB) CreatePendaftaranNikah(c *gin.Context) {
			// ... kode existing untuk membuat pendaftaran ...

			// Setelah pendaftaran berhasil dibuat, tambahkan notifikasi
			notificationService := services.NewNotificationService(idb.DB)
			err := notificationService.SendPendaftaranNotification(pendaftaran.ID, pendaftarID)
			if err != nil {
				log.Printf("Gagal mengirim notifikasi: %v", err)
				// Jangan return error, karena pendaftaran sudah berhasil
			}

			c.JSON(http.StatusCreated, gin.H{
				"message": "Pendaftaran berhasil dibuat",
				"data":    pendaftaran,
			})
		}
	*/

	// Contoh integrasi di handler UpdateStatusPendaftaran
	/*
		func (idb *InDB) UpdateStatusPendaftaran(c *gin.Context) {
			// ... kode existing untuk update status ...

			// Setelah status berhasil diupdate, tambahkan notifikasi
			notificationService := services.NewNotificationService(idb.DB)
			err := notificationService.SendStatusUpdateNotification(pendaftaran.ID, statusLama, statusBaru, updatedBy)
			if err != nil {
				log.Printf("Gagal mengirim notifikasi: %v", err)
				// Jangan return error, karena update sudah berhasil
			}

			c.JSON(http.StatusOK, gin.H{
				"message": "Status berhasil diupdate",
				"data":    pendaftaran,
			})
		}
	*/
}

// ExampleCronJobReminder contoh implementasi cron job untuk pengingat
func ExampleCronJobReminder() {
	// Contoh implementasi cron job (bisa menggunakan library seperti github.com/robfig/cron/v3)
	/*
		func StartReminderCronJob(db *gorm.DB) {
			c := cron.New()

			// Jalankan setiap hari jam 08:00
			c.AddFunc("0 8 * * *", func() {
				notificationService := services.NewNotificationService(db)
				err := notificationService.SendReminderNotification()
				if err != nil {
					log.Printf("Gagal mengirim notifikasi pengingat: %v", err)
				}
			})

			c.Start()
		}
	*/
}

// ExampleFrontendIntegration contoh integrasi dengan frontend
func ExampleFrontendIntegration() {
	// Contoh JavaScript untuk frontend
	/*
		// Mengambil notifikasi user
		async function getUserNotifications(userId, page = 1, limit = 10) {
			try {
				const response = await fetch(`/simnikah/notifikasi/user/${userId}?page=${page}&limit=${limit}`, {
					headers: {
						'Authorization': `Bearer ${localStorage.getItem('token')}`
					}
				});

				if (!response.ok) {
					throw new Error('Failed to fetch notifications');
				}

				const data = await response.json();
				return data;
			} catch (error) {
				console.error('Error fetching notifications:', error);
				throw error;
			}
		}

		// Tandai notifikasi sebagai dibaca
		async function markNotificationAsRead(notificationId) {
			try {
				const response = await fetch(`/simnikah/notifikasi/${notificationId}/status`, {
					method: 'PUT',
					headers: {
						'Authorization': `Bearer ${localStorage.getItem('token')}`,
						'Content-Type': 'application/json'
					},
					body: JSON.stringify({
						status_baca: 'Sudah Dibaca'
					})
				});

				if (!response.ok) {
					throw new Error('Failed to mark notification as read');
				}

				const data = await response.json();
				return data;
			} catch (error) {
				console.error('Error marking notification as read:', error);
				throw error;
			}
		}

		// Tandai semua notifikasi sebagai dibaca
		async function markAllNotificationsAsRead(userId) {
			try {
				const response = await fetch(`/simnikah/notifikasi/user/${userId}/mark-all-read`, {
					method: 'PUT',
					headers: {
						'Authorization': `Bearer ${localStorage.getItem('token')}`
					}
				});

				if (!response.ok) {
					throw new Error('Failed to mark all notifications as read');
				}

				const data = await response.json();
				return data;
			} catch (error) {
				console.error('Error marking all notifications as read:', error);
				throw error;
			}
		}

		// Mengambil statistik notifikasi
		async function getNotificationStats(userId) {
			try {
				const response = await fetch(`/simnikah/notifikasi/user/${userId}/stats`, {
					headers: {
						'Authorization': `Bearer ${localStorage.getItem('token')}`
					}
				});

				if (!response.ok) {
					throw new Error('Failed to fetch notification stats');
				}

				const data = await response.json();
				return data;
			} catch (error) {
				console.error('Error fetching notification stats:', error);
				throw error;
			}
		}
	*/
}
