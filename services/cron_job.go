package services

import (
	"log"
	"time"

	"gorm.io/gorm"
)

// CronJobService untuk mengelola cron job notifikasi
type CronJobService struct {
	DB                  *gorm.DB
	NotificationService *NotificationService
}

// NewCronJobService membuat instance baru dari CronJobService
func NewCronJobService(db *gorm.DB) *CronJobService {
	return &CronJobService{
		DB:                  db,
		NotificationService: NewNotificationService(db),
	}
}

// StartReminderCronJob memulai cron job untuk pengingat harian
func (cjs *CronJobService) StartReminderCronJob() {
	// Jalankan pengingat setiap hari jam 08:00
	ticker := time.NewTicker(24 * time.Hour)

	// Jalankan segera untuk testing (bisa dihapus di production)
	go func() {
		log.Println("Menjalankan pengingat notifikasi...")
		if err := cjs.NotificationService.SendReminderNotification(); err != nil {
			log.Printf("Gagal mengirim notifikasi pengingat: %v", err)
		}
	}()

	// Jalankan setiap hari
	go func() {
		for {
			select {
			case <-ticker.C:
				// Cek apakah sekarang jam 08:00
				now := time.Now()
				if now.Hour() == 8 && now.Minute() == 0 {
					log.Println("Menjalankan pengingat notifikasi harian...")
					if err := cjs.NotificationService.SendReminderNotification(); err != nil {
						log.Printf("Gagal mengirim notifikasi pengingat: %v", err)
					}
				}
			}
		}
	}()

	log.Println("Cron job pengingat notifikasi telah dimulai")
}

// StartReminderCronJobWithSchedule memulai cron job dengan jadwal yang bisa dikustomisasi
func (cjs *CronJobService) StartReminderCronJobWithSchedule(hour, minute int) {
	// Hitung waktu sampai jadwal berikutnya
	now := time.Now()
	nextRun := time.Date(now.Year(), now.Month(), now.Day(), hour, minute, 0, 0, now.Location())

	// Jika waktu sudah lewat hari ini, jadwalkan untuk besok
	if nextRun.Before(now) {
		nextRun = nextRun.Add(24 * time.Hour)
	}

	// Hitung durasi sampai jadwal berikutnya
	duration := nextRun.Sub(now)

	log.Printf("Pengingat notifikasi akan dijalankan pada %s (dalam %v)", nextRun.Format("2006-01-02 15:04:05"), duration)

	// Timer untuk jadwal pertama
	timer := time.NewTimer(duration)

	go func() {
		<-timer.C

		// Jalankan pengingat
		log.Println("Menjalankan pengingat notifikasi sesuai jadwal...")
		if err := cjs.NotificationService.SendReminderNotification(); err != nil {
			log.Printf("Gagal mengirim notifikasi pengingat: %v", err)
		}

		// Set timer untuk jadwal berikutnya (24 jam kemudian)
		ticker := time.NewTicker(24 * time.Hour)
		for {
			select {
			case <-ticker.C:
				log.Println("Menjalankan pengingat notifikasi harian...")
				if err := cjs.NotificationService.SendReminderNotification(); err != nil {
					log.Printf("Gagal mengirim notifikasi pengingat: %v", err)
				}
			}
		}
	}()

	log.Printf("Cron job pengingat notifikasi telah dimulai dengan jadwal %02d:%02d", hour, minute)
}

// StopReminderCronJob menghentikan cron job (untuk testing atau maintenance)
func (cjs *CronJobService) StopReminderCronJob() {
	log.Println("Cron job pengingat notifikasi dihentikan")
}

// RunReminderNow menjalankan pengingat sekarang (untuk testing)
func (cjs *CronJobService) RunReminderNow() error {
	log.Println("Menjalankan pengingat notifikasi sekarang...")
	return cjs.NotificationService.SendReminderNotification()
}

// GetNextReminderTime mendapatkan waktu pengingat berikutnya
func (cjs *CronJobService) GetNextReminderTime(hour, minute int) time.Time {
	now := time.Now()
	nextRun := time.Date(now.Year(), now.Month(), now.Day(), hour, minute, 0, 0, now.Location())

	// Jika waktu sudah lewat hari ini, jadwalkan untuk besok
	if nextRun.Before(now) {
		nextRun = nextRun.Add(24 * time.Hour)
	}

	return nextRun
}
