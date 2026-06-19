package services

import (
	"log"
	"sync"
	"time"

	"gorm.io/gorm"
)

// CronJobService untuk mengelola cron job notifikasi
type CronJobService struct {
	DB                  *gorm.DB
	NotificationService *NotificationService
	stopCh              chan struct{}
	wg                  sync.WaitGroup
}

// NewCronJobService membuat instance baru dari CronJobService
func NewCronJobService(db *gorm.DB) *CronJobService {
	return &CronJobService{
		DB:                  db,
		NotificationService: NewNotificationService(db),
		stopCh:              make(chan struct{}),
	}
}

// StartReminderCronJobWithSchedule memulai cron job dengan jadwal yang bisa dikustomisasi
func (cjs *CronJobService) StartReminderCronJobWithSchedule(hour, minute int) {
	now := time.Now()
	nextRun := time.Date(now.Year(), now.Month(), now.Day(), hour, minute, 0, 0, now.Location())

	if nextRun.Before(now) {
		nextRun = nextRun.Add(24 * time.Hour)
	}

	duration := nextRun.Sub(now)
	log.Printf("Pengingat notifikasi akan dijalankan pada %s (dalam %v)", nextRun.Format("2006-01-02 15:04:05"), duration)

	cjs.wg.Add(1)
	go func() {
		defer cjs.wg.Done()

		// Timer untuk jadwal pertama
		timer := time.NewTimer(duration)
		defer timer.Stop()

		select {
		case <-timer.C:
			log.Println("Menjalankan pengingat notifikasi sesuai jadwal...")
			if err := cjs.NotificationService.SendReminderNotification(); err != nil {
				log.Printf("Gagal mengirim notifikasi pengingat: %v", err)
			}
		case <-cjs.stopCh:
			log.Println("Cron job dihentikan sebelum jadwal pertama")
			return
		}

		// Ticker untuk jadwal berikutnya (24 jam)
		ticker := time.NewTicker(24 * time.Hour)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				log.Println("Menjalankan pengingat notifikasi harian...")
				if err := cjs.NotificationService.SendReminderNotification(); err != nil {
					log.Printf("Gagal mengirim notifikasi pengingat: %v", err)
				}
			case <-cjs.stopCh:
				log.Println("Cron job dihentikan")
				return
			}
		}
	}()

	log.Printf("Cron job pengingat notifikasi telah dimulai dengan jadwal %02d:%02d", hour, minute)
}

// StopReminderCronJob menghentikan cron job secara graceful
func (cjs *CronJobService) StopReminderCronJob() {
	log.Println("Menghentikan cron job pengingat notifikasi...")
	close(cjs.stopCh)
	cjs.wg.Wait()
	log.Println("Cron job pengingat notifikasi berhasil dihentikan")
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

	if nextRun.Before(now) {
		nextRun = nextRun.Add(24 * time.Hour)
	}

	return nextRun
}
