# ⚡ Quick Start - SimNikah API

Pilih cara deploy yang kamu mau:

---

## 🚂 Railway (Production - Recommended)

**Waktu:** ~30 menit | **Biaya:** GRATIS ($5 credit/bulan)

### Step-by-step:
1. **Push ke GitHub**
   ```bash
   git remote add origin https://github.com/YOUR_USERNAME/simnikah-backend.git
   git push -u origin main
   ```

2. **Deploy di Railway**
   - Login di [railway.app](https://railway.app) dengan GitHub
   - New Project → Deploy from GitHub
   - Pilih repository `simnikah-backend`

3. **Add MySQL Database**
   - Click "+ New" → Database → MySQL
   - Tunggu MySQL ready

4. **Set Environment Variables**
   ```bash
   DB_HOST=${{MySQL.MYSQL_HOST}}
   DB_PORT=${{MySQL.MYSQL_PORT}}
   DB_USER=${{MySQL.MYSQL_USER}}
   DB_PASSWORD=${{MySQL.MYSQL_PASSWORD}}
   DB_NAME=${{MySQL.MYSQL_DATABASE}}
   JWT_KEY=<hasil dari: openssl rand -base64 32>
   PORT=8080
   GIN_MODE=release
   ALLOWED_ORIGINS=http://localhost:3000,http://localhost:5173
   ```

5. **Generate Domain & Test**
   - Settings → Generate Domain
   - Test: `curl https://YOUR_APP.up.railway.app/health`

**✅ DONE! API sudah LIVE!**

📖 **Tutorial lengkap:** [TUTORIAL_DEPLOY_RAILWAY.md](TUTORIAL_DEPLOY_RAILWAY.md)

---

## 🐳 Docker Compose (Local Development)

**Waktu:** ~5 menit | **Biaya:** GRATIS

### Setup:

```bash
# Clone project
git clone https://github.com/YOUR_USERNAME/simnikah-backend.git
cd simnikah-backend

# Copy environment
cp env.example .env

# Edit .env jika perlu
nano .env

# Run dengan Docker
docker-compose up -d

# Check logs
docker-compose logs -f simnikah-api

# Test
curl http://localhost:8080/health
```

**✅ DONE! API berjalan di `http://localhost:8080`**

### Stop:
```bash
docker-compose down
```

### Restart:
```bash
docker-compose restart
```

---

## 🖥️ Local (Development)

**Waktu:** ~10 menit | **Biaya:** GRATIS

### Prerequisites:
- Go 1.23.6+
- MySQL 8.0+

### Setup MySQL:

```bash
# Install MySQL
sudo apt install mysql-server  # Ubuntu/Debian
# atau
brew install mysql             # macOS

# Start MySQL
sudo systemctl start mysql     # Linux
# atau
brew services start mysql      # macOS

# Login
sudo mysql -u root -p

# Create database & user
CREATE DATABASE simnikah;
CREATE USER 'simnikah_user'@'localhost' IDENTIFIED BY 'simnikah123';
GRANT ALL PRIVILEGES ON simnikah.* TO 'simnikah_user'@'localhost';
FLUSH PRIVILEGES;
EXIT;
```

### Setup Project:

```bash
# Clone
git clone https://github.com/YOUR_USERNAME/simnikah-backend.git
cd simnikah-backend

# Install dependencies
go mod download

# Copy & edit .env
cp env.example .env
nano .env

# Update di .env:
# DB_HOST=127.0.0.1
# DB_PORT=3306
# DB_USER=simnikah_user
# DB_PASSWORD=simnikah123
# DB_NAME=simnikah

# Run
go run main.go

# Or build & run
go build -o main .
./main
```

**✅ DONE! API berjalan di `http://localhost:8080`**

### Test:
```bash
curl http://localhost:8080/health
```

---

## 📊 Perbandingan

| Method | Pros | Cons | Use Case |
|--------|------|------|----------|
| **Railway** | ✅ Auto SSL<br>✅ Auto backup<br>✅ Easy scaling<br>✅ Free tier | ⚠️ Need GitHub<br>⚠️ Internet required | Production, Demo, Staging |
| **Docker** | ✅ Quick setup<br>✅ Isolated environment<br>✅ Reproducible | ⚠️ Need Docker installed<br>⚠️ Resource intensive | Local dev, Testing |
| **Local** | ✅ Full control<br>✅ Fast rebuild<br>✅ Easy debug | ⚠️ Manual MySQL setup<br>⚠️ OS dependent | Active development |

---

## 🧪 Test API

### Health Check
```bash
curl http://localhost:8080/health
```

Expected:
```json
{
  "status": "healthy",
  "message": "SimNikah API is running"
}
```

### Register User
```bash
curl -X POST http://localhost:8080/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "email": "test@example.com",
    "password": "test123",
    "nama": "Test User",
    "no_hp": "08123456789",
    "role": "user_biasa"
  }'
```

### Login
```bash
curl -X POST http://localhost:8080/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "password": "test123"
  }'
```

Copy `token` dari response untuk authenticated requests.

### Test Authenticated Endpoint
```bash
curl http://localhost:8080/simnikah/profile \
  -H "Authorization: Bearer YOUR_TOKEN_HERE"
```

---

## 📚 Next Steps

1. **Explore API:** [API_DOCUMENTATION.md](API_DOCUMENTATION.md)
2. **Test Endpoints:** [API_TESTING_DOCUMENTATION.md](API_TESTING_DOCUMENTATION.md)
3. **Setup CORS:** [CORS_SETUP.md](CORS_SETUP.md)
4. **Map Integration:** [MAP_INTEGRATION.md](MAP_INTEGRATION.md)

---

## 🐛 Troubleshooting

### Port already in use
```bash
# Find process
lsof -i :8080

# Kill process
kill -9 <PID>

# Or change port in .env
PORT=8081
```

### Database connection failed
```bash
# Check MySQL running
sudo systemctl status mysql

# Test connection
mysql -u simnikah_user -p simnikah
```

### Module not found
```bash
go mod tidy
go mod download
```

---

## 💡 Tips

### Development Workflow
```bash
# Auto reload saat edit code (install air)
go install github.com/cosmtrek/air@latest

# Run dengan air
air
```

### Database GUI
- **phpMyAdmin**: Web-based MySQL admin
- **MySQL Workbench**: Desktop MySQL client
- **TablePlus**: Modern database client (Mac/Windows/Linux)

### API Testing Tools
- **Postman**: GUI API testing
- **Insomnia**: Lightweight API client
- **curl**: Command line testing
- **Thunder Client**: VS Code extension

---

## 🔗 Resources

- **Railway Dashboard**: [railway.app/dashboard](https://railway.app/dashboard)
- **Go Documentation**: [go.dev/doc](https://go.dev/doc)
- **Gin Framework**: [gin-gonic.com](https://gin-gonic.com)
- **GORM**: [gorm.io](https://gorm.io)

---

**Selamat coding! 🚀**

Need help? Check dokumentasi lengkap atau buat issue di GitHub.


