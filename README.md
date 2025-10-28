# 🏛️ SimNikah - Sistem Manajemen Pendaftaran Nikah KUA

[![Go Version](https://img.shields.io/badge/Go-1.23.6-blue.svg)](https://golang.org/)
[![Gin Framework](https://img.shields.io/badge/Gin-v1.10.0-green.svg)](https://github.com/gin-gonic/gin)
[![MySQL](https://img.shields.io/badge/MySQL-8.0%2B-blue.svg)](https://www.mysql.com/)
[![Railway](https://img.shields.io/badge/Deploy-Railway-purple.svg)](https://railway.app)
[![License](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)

REST API untuk mengelola pendaftaran dan administrasi pernikahan di Kantor Urusan Agama (KUA).

---

## 🚀 Quick Start

### Prerequisites
- Go 1.23.6+
- MySQL 8.0+
- Git

### Installation

```bash
# Clone repository
git clone https://github.com/your-username/simpadu.git
cd simpadu

# Install dependencies
go mod download

# Copy environment variables
cp env.example .env

# Edit .env dengan database credentials Anda
nano .env

# Run application
go run main.go
```

Server akan berjalan di `http://localhost:8080`

---

## 📚 Documentation

- **[API Documentation](API_DOCUMENTATION.md)** - Complete API reference (50+ endpoints)
- **[API Testing](API_TESTING_DOCUMENTATION.md)** - Testing guide
- **[Map Integration](MAP_INTEGRATION.md)** - 🆕 Map & location integration (100% FREE)
- **[Wali Nikah Validation](WALI_NIKAH_VALIDATION.md)** - 🆕 Sharia-compliant guardian validation
- **[Status Management](STATUS_MANAGEMENT.md)** - Status constants and best practices
- **[CORS Setup](CORS_SETUP.md)** - CORS configuration guide
- **[Railway Deployment](RAILWAY_DEPLOYMENT.md)** - Railway deployment guide
- **[Migration Guide](MIGRATION_GUIDE.md)** - LeapCell to Railway migration

---

## ⚡ Features

- ✅ **Authentication & Authorization** - JWT-based with role management
- ✅ **Marriage Registration** - Complete registration workflow
- ✅ **Staff Management** - CRUD for KUA staff and penghulu
- ✅ **Multi-level Verification** - Staff → Penghulu → Approval
- ✅ **Calendar System** - Marriage schedule & availability
- ✅ **Counseling System** - Pre-marriage counseling management
- ✅ **Notification System** - Auto notifications with cron jobs
- ✅ **CORS Configured** - Ready for frontend integration
- ✅ **Map Integration (100% FREE)** 🆕
  - Geocoding & Reverse Geocoding (OpenStreetMap Nominatim API)
  - Address autocomplete search
  - Location display for penghulu
  - Navigation links (Google Maps, Waze, OpenStreetMap)
- ✅ **Wali Nikah Validation** 🆕 - Sharia-compliant guardian validation
- ✅ **Docker Ready** - Development & production

---

## 👥 User Roles

| Role | Description | Access Level |
|------|-------------|--------------|
| `user_biasa` | Calon pengantin | Register marriage, view status |
| `staff` | Staff KUA | Verify forms & documents |
| `penghulu` | Penghulu | Verify documents, view schedule |
| `kepala_kua` | Kepala KUA | Full access, assign penghulu |

---

## 🔄 Registration Workflow

```
1. Draft
   ↓
2. Menunggu Verifikasi (Staff verify form)
   ↓
3. Menunggu Pengumpulan Berkas (User bring documents)
   ↓
4. Berkas Diterima (Staff verify documents)
   ↓
5. Menunggu Penugasan (Ready for penghulu assignment)
   ↓
6. Penghulu Ditugaskan (Head assigns penghulu)
   ↓
7. Menunggu Verifikasi Penghulu (Penghulu checks)
   ↓
8. Menunggu Bimbingan (Pre-marriage counseling)
   ↓
9. Sudah Bimbingan (Counseling completed)
   ↓
10. Selesai (Marriage completed)
```

---

## 🌐 CORS Configuration

### Development
CORS sudah dikonfigurasi untuk development. Default allowed origins:
- `http://localhost:3000` (React)
- `http://localhost:5173` (Vite)
- `http://localhost:8080`

### Production
Set environment variable:
```bash
ALLOWED_ORIGINS=https://your-frontend.com,https://app.your-domain.com
```

📖 **[Read full CORS guide →](CORS_SETUP.md)**

---

## 🐳 Docker

### Development
```bash
docker-compose up -d
```

### Production
```bash
docker build -t simnikah-api .
docker run -p 8080:8080 \
  -e DB_HOST=your-db-host \
  -e DB_PASSWORD=your-password \
  -e JWT_KEY=your-jwt-key \
  -e ALLOWED_ORIGINS=https://your-frontend.com \
  simnikah-api
```

---

## 🚀 Deploy to Railway

### Quick Steps:
1. Push to GitHub
2. Create new project in Railway
3. Add MySQL database
4. Set environment variables (including `ALLOWED_ORIGINS`)
5. Deploy automatically!

📖 **[Read Railway deployment guide →](RAILWAY_DEPLOYMENT.md)**

### Required Environment Variables:
```bash
# Database (Railway auto-provides via references)
DB_HOST=${{MYSQL.MYSQL_HOST}}
DB_PORT=${{MYSQL.MYSQL_PORT}}
DB_USER=${{MYSQL.MYSQL_USER}}
DB_PASSWORD=${{MYSQL.MYSQL_PASSWORD}}
DB_NAME=${{MYSQL.MYSQL_DATABASE}}

# Manual configuration
JWT_KEY=your-32-char-secret-key
ALLOWED_ORIGINS=http://localhost:3000,https://your-frontend.vercel.app
PORT=8080
GIN_MODE=release
```

---

## 📡 API Endpoints Overview

### Authentication (3)
- `POST /register` - Register user
- `POST /login` - Login
- `GET /profile` - Get profile

### Marriage Registration (7)
- `POST /simnikah/pendaftaran/form-baru` - Create registration
- `GET /simnikah/pendaftaran/status` - Check status
- `GET /simnikah/pendaftaran` - Get all registrations
- ...and more

### Staff Management (6)
- `POST /simnikah/staff` - Create staff
- `POST /simnikah/staff/verify-formulir/:id` - Verify form
- `POST /simnikah/staff/verify-berkas/:id` - Verify documents
- ...and more

### Penghulu Management (9)
- `POST /simnikah/penghulu` - Create penghulu
- `POST /simnikah/pendaftaran/:id/assign-penghulu` - Assign
- ...and more

### Calendar & Schedule (5)
- `GET /simnikah/kalender-ketersediaan` - Availability calendar
- `GET /simnikah/penghulu-jadwal/:tanggal` - Penghulu schedule
- ...and more

### Counseling (8)
- `POST /simnikah/bimbingan` - Create session
- `POST /simnikah/bimbingan/:id/daftar` - Register participant
- ...and more

### Notifications (8)
- `POST /simnikah/notifikasi` - Create notification
- `GET /simnikah/notifikasi/user/:user_id` - Get user notifications
- ...and more

📖 **[View complete API docs →](API_DOCUMENTATION.md)**

---

## 🧪 Testing

### Health Check
```bash
curl http://localhost:8080/health
```

### Register User
```bash
curl -X POST http://localhost:8080/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "test",
    "email": "test@example.com",
    "password": "test123",
    "nama": "Test User",
    "role": "user_biasa"
  }'
```

### Login
```bash
curl -X POST http://localhost:8080/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "test",
    "password": "test123"
  }'
```

📖 **[View complete testing guide →](API_TESTING_DOCUMENTATION.md)**

---

## 🗄️ Database Models

- **Users** - Authentication & role management
- **CalonPasangan** - Bride/Groom data
- **DataOrangTua** - Parent data
- **PendaftaranNikah** - Marriage registration
- **WaliNikah** - Guardian data
- **StaffKUA** - KUA staff data
- **Penghulu** - Penghulu data
- **BimbinganPerkawinan** - Counseling sessions
- **PendaftaranBimbingan** - Counseling participants
- **Notifikasi** - Notifications

---

## 🔒 Security

- ✅ JWT Authentication (24h validity)
- ✅ bcrypt Password Hashing (10 rounds)
- ✅ Role-based Authorization
- ✅ SQL Injection Protection (GORM)
- ✅ CORS Protection
- ✅ Environment Variables
- ✅ HTTPS (Production)

---

## 🏗️ Tech Stack

- **Language**: Go 1.23.6
- **Framework**: Gin v1.10.0
- **ORM**: GORM v1.26.1
- **Database**: MySQL 8.0+
- **Authentication**: JWT (golang-jwt/jwt v5)
- **Password**: bcrypt (golang.org/x/crypto)
- **CORS**: gin-contrib/cors v1.7.2
- **Deployment**: Railway

---

## 📁 Project Structure

```
simpadu/
├── main.go              # Entry point & routing
├── config/              # Database configuration
│   ├── config.go        # DB connection
│   └── db.go            # DB struct
├── structs/             # Data models
│   └── models.go        # All database models
├── helper/              # Helper functions
│   ├── bcrypt.go        # Password hashing
│   ├── jwt.go           # JWT utilities
│   └── validation.go    # Input validation
├── services/            # Business logic
│   ├── cron_job.go      # Scheduled tasks
│   └── notification_service.go
├── catin/               # Marriage registration
├── staff/               # Staff management
├── penghulu/            # Penghulu management
├── notification/        # Notifications
├── Dockerfile           # Docker configuration
├── docker-compose.yml   # Local development
├── railway.json         # Railway deployment config
├── nixpacks.toml        # Railway build config
└── docs/                # Documentation
```

---

## 🔧 Environment Variables

Create `.env` file:
```bash
# Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_password
DB_NAME=simnikah

# JWT
JWT_KEY=your-super-secret-jwt-key-minimum-32-characters-long

# Server
PORT=8080
GIN_MODE=debug

# CORS
ALLOWED_ORIGINS=http://localhost:3000,http://localhost:5173
```

---

## 📊 Business Rules

- ✅ Max 9 marriages per day at KUA
- ✅ Unlimited marriages outside KUA
- ✅ Max 3 marriages per penghulu per day
- ✅ Counseling only on Wednesdays
- ✅ Max 10 couples per counseling session
- ✅ Auto database migration on startup
- ✅ Daily reminder notifications at 08:00

---

## 🐛 Troubleshooting

### CORS Error
```bash
# Check ALLOWED_ORIGINS is set correctly
# View CORS_SETUP.md for detailed guide
```

### Database Connection Failed
```bash
# Check DB credentials in .env
# Ensure PostgreSQL is running
# Verify port 5432 is accessible
```

### Build Failed
```bash
# Update dependencies
go mod tidy

# Clear build cache
go clean -cache

# Rebuild
go build
```

---

## 📈 Roadmap

- [x] Core API functionality
- [x] Authentication & Authorization
- [x] CORS configuration
- [x] Docker support
- [x] LeapCell deployment
- [ ] Unit tests
- [ ] API rate limiting
- [ ] File upload (documents)
- [ ] PDF generation (certificates)
- [ ] Email notifications
- [ ] SMS notifications
- [ ] Mobile app API

---

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/AmazingFeature`)
3. Commit your changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

---

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

---

## 👨‍💻 Author

**SimNikah Development Team**

---

## 📞 Support

- 📧 Email: support@simnikah.id
- 📚 Documentation: [Full Docs](API_DOCUMENTATION.md)
- 🐛 Issues: [GitHub Issues](https://github.com/your-username/simpadu/issues)

---

## 🙏 Acknowledgments

- Gin Web Framework
- GORM
- LeapCell Platform
- OpenStreetMap (Geocoding)

---

**Made with ❤️ for Indonesian KUA (Kantor Urusan Agama)**

---

## Quick Links

- 📖 [API Documentation](API_DOCUMENTATION.md)
- 🌐 [CORS Setup Guide](CORS_SETUP.md)
- 🚂 [Railway Deployment](RAILWAY_DEPLOYMENT.md)
- 🔄 [Migration Guide](MIGRATION_GUIDE.md)
- 🧪 [Testing Guide](API_TESTING_DOCUMENTATION.md)
- 🐳 [Docker Setup](docker-compose.yml)




