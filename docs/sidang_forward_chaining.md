# Penjelasan Sistem SIPENA
## Forward Chaining untuk Distribusi Jadwal Penghulu

## 1. Ringkasan Eksekutif
SIPENA sekarang difokuskan hanya pada penjadwalan penghulu. Sistem tidak lagi meniru SIMKAH untuk validasi berkas N1-N4, KTP, atau proses administrasi nikah yang kompleks. Alur bisnis disederhanakan menjadi 4 tahap: catin mengecek slot kosong dan mengisi data esensial, engine Forward Chaining merekomendasikan penghulu, Kepala KUA menyetujui assignment, lalu penghulu melihat jadwal penugasan yang sudah disetujui.

## 2. Alur Aplikasi 4 Tahap
1. Catin mengirim tanggal, waktu, alamat, dan koordinat lokasi nikah.
2. Sistem menjalankan Forward Chaining untuk mengecek slot dan memilih penghulu yang paling layak.
3. Kepala KUA melakukan approval final dengan transaksi aman.
4. Penghulu melihat daftar penugasan yang sudah disetujui beserta alamat detail dan geolocation.

## 3. Kesesuaian Metodologi Forward Chaining
Metode Forward Chaining terlihat jelas karena sistem bergerak dari fakta awal menuju kesimpulan.

### Fakta Awal
Working memory hanya berisi fakta yang relevan untuk scheduling:
- Tanggal nikah
- Waktu nikah
- Tempat nikah
- Alamat lengkap
- Latitude dan longitude
- Status pendaftaran
- Penghulu aktif
- Beban jadwal harian dan per jam

### Rule Base
Rule yang dievaluasi hanya constraint penjadwalan:
- Status penghulu aktif
- Jadwal bentrok atau tidak
- Kapasitas harian penghulu
- Kapasitas per jam penghulu
- Kesesuaian lokasi nikah
- Jarak geolocation ke lokasi penugasan

### Inferensi Bertahap
Saat sebuah rule terpenuhi, engine menghasilkan fakta turunan seperti:
- Slot masih tersedia
- Penghulu aktif
- Kapasitas harian memadai
- Kapasitas per jam memadai
- Lokasi sesuai
- Jarak layak
- Penghulu dapat direkomendasikan

Jika salah satu rule blocking gagal, engine menghentikan proses dan memberi alasan penolakan.

## 4. Bedah Arsitektur Kode

### A. Layer Service
File inti ada di [internal/services/forward_chaining_engine.go](internal/services/forward_chaining_engine.go).

Service bertugas:
- Membentuk snapshot working memory
- Menghitung ketersediaan jadwal
- Menjalankan EvaluateRules
- Menghasilkan rekomendasi penghulu
- Menyajikan detail evaluasi dan alternatif

Pada versi baru, fakta registrasi dibatasi ke data inti scheduling saja. Struct minimal yang dipakai alur baru adalah `MarriageRegistrationFact` dan `PendaftaranJadwal`.

### B. Layer Controller
Controller yang relevan ada di:
- [internal/handlers/catin/daftar.go](internal/handlers/catin/daftar.go)
- [internal/handlers/kepala_kua/forward_chaining_handlers.go](internal/handlers/kepala_kua/forward_chaining_handlers.go)
- [internal/handlers/penghulu/penghulu.go](internal/handlers/penghulu/penghulu.go)

Tugas controller:
- Menerima request dari frontend
- Validasi URI dan payload JSON
- Menjalankan RBAC
- Menjalankan approval transaction
- Mengembalikan response JSON yang ringan dan jelas

## 5. Endpoint Utama

### Catin
Endpoint baru:
- POST `/simnikah/check-schedule`

Fungsinya:
- Mengecek slot tanggal/waktu
- Menolak slot yang fully booked
- Mengembalikan status available / conflict
- Menyediakan data rekomendasi awal jika slot masih kosong

### Kepala KUA
Endpoint forward chaining yang dipakai:
- GET `/simnikah/kepala-kua/forward-chaining/recommendation/:id`
- GET `/simnikah/kepala-kua/forward-chaining/evaluation/:id`
- POST `/simnikah/kepala-kua/forward-chaining/assign/:id`
- GET `/simnikah/kepala-kua/forward-chaining/config`

Assignment final tetap aman karena:
- memakai GORM transaction
- memakai row locking `FOR UPDATE`
- hanya bisa diakses role `Kepala KUA`

### Penghulu
Endpoint baru:
- GET `/simnikah/penghulu/jadwal-penugasan`

Fungsinya:
- Mengambil jadwal penugasan yang sudah disetujui
- Memastikan data ditampilkan berdasarkan `user_id` penghulu login
- Menampilkan alamat lengkap, latitude, dan longitude untuk kebutuhan rute/map

## 6. Alasan Desain Baru
Refactor ini sengaja menghilangkan fitur yang bukan fokus utama SIPENA agar tidak menduplikasi SIMKAH.

Yang dipertahankan hanya:
- Logika scheduling
- Distribusi beban kerja penghulu
- Approval assignment yang aman
- View jadwal untuk penghulu

Yang tidak lagi menjadi fokus:
- Validasi dokumen N1-N4
- Alur administrasi berkas nikah yang kompleks
- Field SIMKAH yang tidak dipakai penjadwalan

## 7. Selling Points untuk Sidang
- Benar-benar Forward Chaining: keputusan dibangun dari fakta awal menuju konklusi.
- Aman dari race condition: assignment memakai transaction + row locking.
- Aman secara akses: assignment hanya bisa dilakukan Kepala KUA.
- Fokus domain jelas: sistem hanya menangani jadwal penghulu.
- Lebih ringan: data yang diproses hanya field inti scheduling.
- Mudah dikembangkan: konfigurasi rule dapat dipindah ke `system_configs` kapan saja.

## 8. Kode/Struktur Lama yang Kandidat Dihapus
Berikut komponen lama yang tidak relevan dengan alur baru dan layak dihapus setelah semua route/UI lama dipastikan tidak dipakai:

### A. Kandidat fungsi legacy di controller catin
File: [internal/handlers/catin/daftar.go](internal/handlers/catin/daftar.go)
- `UpdateMarriageLocation`
- `GetUserRegistrationStatus`
- `ListRegistrations`
- `GetRegistrationDetail`
- `GetCalendarAvailability`
- `GetAvailableTimeSlots`
- `GetWeddingsByDate`
- `CreateFeedbackPernikahan`

Alasan: fungsi-fungsi ini masih membawa pola SIMKAH/administrasi atau tidak diperlukan untuk 4 tahap inti.

### B. Kandidat fungsi legacy di controller penghulu
File: [internal/handlers/penghulu/penghulu.go](internal/handlers/penghulu/penghulu.go)
- `VerifyRegistrationDocuments`
- `ListMyAssignments`
- `GetTodaySchedule`
- `CompleteMarriage`

Alasan: sebagian besar masih memakai alur lama yang bercampur verifikasi dokumen dan penyelesaian pernikahan, bukan view jadwal penugasan sederhana.

### C. Kandidat fungsi legacy di controller Kepala KUA
File: [internal/handlers/kepala_kua/kepala_kua.go](internal/handlers/kepala_kua/kepala_kua.go)
- Flow assignment lama di sekitar `AssignMarriageOfficer`
- Endpoint statistik/pengumuman/feedback yang tidak mendukung scheduling-only
- Endpoint CRUD staff/penghulu yang tidak terkait engine

Alasan: file ini masih memuat fitur operasional lain di luar fokus Forward Chaining.

### D. Kandidat field model lama
File: [internal/models/models.go](internal/models/models.go)
- `Nomor_pendaftaran`
- `Pendaftar_id`
- `Calon_suami_id`
- `Calon_istri_id`
- `Wali_nikah_id`
- `Tanggal_pendaftaran`
- `Disetujui_oleh`
- `Disetujui_pada`
- `Penghulu_assigned_by`
- `Catatan` jika tidak dipakai pada flow baru

Alasan: field-field ini bukan inti penjadwalan. Untuk arsitektur paling bersih, data yang dipakai cukup ID, tanggal, waktu, tempat, alamat, koordinat, status, dan penghulu.

### E. Kandidat struct/service lama
File: [internal/services/forward_chaining_engine.go](internal/services/forward_chaining_engine.go)
- `PenghuluFact` field yang tidak dipakai scheduling-only seperti `RejectRate`, `AverageDuration`, `BaseLocation`, `SpecializedArea`
- helper API yang masih terkait compatibility lama jika tidak dipakai UI baru
- evaluasi yang masih menampilkan nomor pendaftaran jika tidak diperlukan lagi

Alasan: makin sedikit fakta dan field, makin mudah menjaga engine tetap fokus dan mudah diaudit.

## 9. Contoh Narasi Sidang
- “Sistem ini tidak lagi meniru SIMKAH. SIPENA sekarang hanya menangani alur penjadwalan penghulu dengan Forward Chaining.”
- “Forward Chaining bekerja dari fakta awal seperti tanggal, waktu, lokasi, dan beban kerja penghulu, lalu menghasilkan konklusi rekomendasi.”
- “Approval final tetap aman karena memakai transaction dan row locking untuk mencegah bentrok data.”
- “Penghulu hanya melihat jadwal yang sudah disetujui, sehingga UI dan backend tetap sederhana.”

## 10. Kesimpulan
Refactor ini menjadikan SIPENA sebagai sistem penjadwalan penghulu yang fokus, aman, dan mudah dijelaskan secara akademik. Forward Chaining menjadi pusat logika, sementara controller hanya menangani validasi, otorisasi, dan hasil akhir yang diperlukan pengguna.
